package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	"unicode/utf8"

	"net/http"

	"github.com/Financial-Times/kafka-client-go/kafka"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/twinj/uuid"
)

const messageTimestampDateFormat = "2006-01-02T15:04:05.000Z"

var (
	messageConsumer  kafka.Consumer
	messageProducer  kafka.Producer
	logger           *AppLogger
	taxonomyHandlers map[string]TaxonomyService
	whitelist        *regexp.Regexp
)

func init() {
	logger = NewAppLogger()
	taxonomyHandlers = map[string]TaxonomyService{
		"subjects":         SubjectService{HandledTaxonomy: "subjects"},
		"sections":         SectionService{HandledTaxonomy: "sections"},
		"topics":           TopicService{HandledTaxonomy: "topics"},
		"locations":        LocationService{HandledTaxonomy: "gl"},
		"genres":           GenreService{HandledTaxonomy: "genres"},
		"specialReports":   SpecialReportService{HandledTaxonomy: "specialReports"},
		"alphavilleSeries": AlphavilleSeriesService{HandledTaxonomy: "alphavilleSeriesClassification"},
		"organisations":    OrganisationService{HandledTaxonomy: "ON"},
		"people":           PeopleService{HandledTaxonomy: "PN"},
		"authors":          AuthorService{HandledTaxonomy: "Authors"},
		"brands":           BrandService{HandledTaxonomy: "Brands"},
	}
}

func main() {
	app := cli.App("annotations-mapper", "A service to read V1 metadata publish event, filter it and output UPP-specific metadata to the destination queue.")
	zookeeperAddress := app.String(cli.StringOpt{
		Name:   "zookeeperAddress",
		Value:  "localhost:2181",
		Desc:   "Addresses used by the queue consumer to connect to the queue",
		EnvVar: "ZOOKEEPER_ADDRESS",
	})
	consumerGroup := app.String(cli.StringOpt{
		Name:   "consumerGroup",
		Desc:   "Group used to read the messages from the queue",
		EnvVar: "CONSUMER_GROUP",
	})
	consumerTopic := app.String(cli.StringOpt{
		Name:   "consumerTopic",
		Desc:   "The topic to read the meassages from",
		EnvVar: "CONSUMER_TOPIC",
	})
	brokerAddress := app.String(cli.StringOpt{
		Name:   "brokerAddress",
		Desc:   "Address used by the producer to connect to the queue",
		EnvVar: "BROKER_ADDRESS",
	})
	producerTopic := app.String(cli.StringOpt{
		Name:   "producerTopic",
		Desc:   "The topic to write the concept annotation to",
		EnvVar: "PRODUCER_TOPIC",
	})
	whitelistRegex := app.String(cli.StringOpt{
		Name:   "whitelistRegex",
		Desc:   "The regex to use to filter messages based on Origin-System-Id.",
		EnvVar: "WHITELIST_REGEX",
		Value:  "http://cmdb\\.ft\\.com/systems/methode-web-pub",
	})

	app.Action = func() {
		var err error
		whitelist, err = regexp.Compile(*whitelistRegex)
		if err != nil {
			logger.Fatal("Please specify a valid whitelist ", err)
		}

		messageProducer, err = kafka.NewProducer(*brokerAddress, *producerTopic)
		if err != nil {
			logger.Fatal("Cannot start message producer", err)
		}
		logger.QueueProducerStarted(*producerTopic)

		messageConsumer, err = kafka.NewConsumer(*zookeeperAddress, *consumerGroup, []string{*consumerTopic}, kafka.DefaultConsumerConfig())
		if err != nil {
			logger.Fatal("Cannot start message consumer", err)
		}
		logger.QueueConsumerStarted(*consumerTopic)
		messageConsumer.StartListening(handleMessage)

		go enableHealthChecks(messageConsumer)

		waitForSignal()
		logger.Info("Shutting down Kafka consumer", "", "")
		messageConsumer.Shutdown()
	}

	app.Run(os.Args)
}

func enableHealthChecks(messageConsumer kafka.Consumer) {
	hc := NewHealthCheck(messageConsumer)
	router := mux.NewRouter()
	router.HandleFunc("/__health", hc.Health())
	router.HandleFunc("/__gtg", status.NewGoodToGoHandler(hc.GTG))
	router.HandleFunc(status.PingPath, status.PingHandler)
	router.HandleFunc(status.PingPathDW, status.PingHandler)
	router.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
	router.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)
	http.Handle("/", router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatal("Couldn't set up HTTP listener", err)
	}
}

func handleMessage(msg kafka.FTMessage) error {
	tid := msg.Headers["X-Request-Id"]
	systemCode := msg.Headers["Origin-System-Id"]
	if !whitelist.MatchString(systemCode) {
		logger.Info(fmt.Sprintf("Skipping annotations published with Origin-System-Id \"%v\". It does not match the configured whitelist.", systemCode), tid, "-")
		return nil
	}

	var metadataPublishEvent MetadataPublishEvent
	err := json.Unmarshal([]byte(msg.Body), &metadataPublishEvent)
	if err != nil {
		logger.Error("Cannot unmarshal message body", tid, "", err)
		return err
	}

	logger.Info("Processing metadata publish event", tid, metadataPublishEvent.UUID)

	metadataXML, err := base64.StdEncoding.DecodeString(metadataPublishEvent.Value)
	if err != nil {
		logger.Error("Error decoding body", tid, metadataPublishEvent.UUID, err)
		return err
	}

	metadata, err, hadInvalidChars := unmarshalMetadata(metadataXML)

	if err != nil {
		logger.Error("Error unmarshalling metadata XML", tid, metadataPublishEvent.UUID, err)
		if hadInvalidChars {
			logger.Info("Metadata XML had invalid UTF8 characters.", tid, metadataPublishEvent.UUID)
		}
		return err
	}

	annotations := []annotation{}
	for _, value := range taxonomyHandlers {
		annotations = append(annotations, value.buildAnnotations(metadata)...)
	}

	conceptAnnotations := ConceptAnnotations{UUID: metadataPublishEvent.UUID, Annotations: annotations}

	marshalledAnnotations, err := json.Marshal(conceptAnnotations)
	if err != nil {
		logger.Error("Error marshalling the concept annotations", tid, metadataPublishEvent.UUID, err)
		return err
	}

	var headers = buildConceptAnnotationsHeader(msg.Headers)
	message := kafka.FTMessage{Headers: headers, Body: string(marshalledAnnotations)}
	err = messageProducer.SendMessage(message)
	if err != nil {
		logger.Error("Error sending concept annotation to queue", tid, metadataPublishEvent.UUID, err)
		return err
	}
	logger.Info("Sent annotation message to queue", tid, metadataPublishEvent.UUID)
	return nil
}

func unmarshalMetadata(metadataXML []byte) (ContentRef, error, bool) {
	metadata := ContentRef{}
	err := xml.Unmarshal(metadataXML, &metadata)
	if err == nil {
		return metadata, nil, false
	}
	return metadata, err, !utf8.Valid(metadataXML)
}

func buildConceptAnnotationsHeader(publishEventHeaders map[string]string) map[string]string {
	return map[string]string{
		"Message-Id":        uuid.NewV4().String(),
		"Message-Type":      "concept-annotation",
		"Content-Type":      publishEventHeaders["Content-Type"],
		"X-Request-Id":      publishEventHeaders["X-Request-Id"],
		"Origin-System-Id":  publishEventHeaders["Origin-System-Id"],
		"Message-Timestamp": time.Now().Format(messageTimestampDateFormat),
	}
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
