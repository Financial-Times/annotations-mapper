package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/kafka-client-go/kafka"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/twinj/uuid"
	"net/http"
)

const (
	messageTimestampDateFormat = "2006-01-02T15:04:05.000Z"
	serviceName                = "annotations-mapper"
	contentType                = "Annotations"
	consumerStartedEvent       = "consume_queue"
	producerStartedEvent       = "produce_queue"
)

var (
	messageConsumer  kafka.Consumer
	messageProducer  kafka.Producer
	taxonomyHandlers map[string]TaxonomyService
)

func init() {
	logger.InitDefaultLogger(serviceName)
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
	app := cli.App(serviceName, "A service to read V1 metadata publish event, filter it and output UPP-specific metadata to the destination queue.")
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

	app.Action = func() {
		var err error
		messageProducer, err = kafka.NewProducer(*brokerAddress, *producerTopic)
		if err != nil {
			logger.FatalEvent("Cannot start message producer", err)
		}
		logger.Infof(map[string]interface{}{"event": consumerStartedEvent}, "Starting queue consumer: %v", *producerTopic)

		messageConsumer, err = kafka.NewConsumer(*zookeeperAddress, *consumerGroup, []string{*consumerTopic}, kafka.DefaultConsumerConfig())
		if err != nil {
			logger.FatalEvent("Cannot start message consumer", err)
		}
		logger.Infof(map[string]interface{}{"event": consumerStartedEvent}, "Starting queue producer: %v", *consumerTopic)
		messageConsumer.StartListening(handleMessage)

		go enableHealthChecks(messageConsumer)

		waitForSignal()
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
		logger.FatalEvent("Couldn't set up HTTP listener", err)
	}
}

func handleMessage(msg kafka.FTMessage) error {
	tid := msg.Headers["X-Request-Id"]

	var metadataPublishEvent MetadataPublishEvent
	err := json.Unmarshal([]byte(msg.Body), &metadataPublishEvent)
	if err != nil {
		logger.ErrorEvent(tid, "Cannot unmarshal message body", err)
		return err
	}

	logger.InfoEventWithUUID(tid, metadataPublishEvent.UUID, "Processing metadata publish event")

	metadataXML, err := base64.StdEncoding.DecodeString(metadataPublishEvent.Value)
	if err != nil {
		logger.ErrorEventWithUUID(tid, metadataPublishEvent.UUID, "Error decoding body", err)
		return err
	}

	metadata, err, hadInvalidChars := unmarshalMetadata(metadataXML)
	if err != nil {
		errMsg := "Error unmarshalling metadata XML"
		if hadInvalidChars {
			logger.ErrorEventWithUUID(tid, metadataPublishEvent.UUID, errMsg+"Metadata XML had invalid UTF8 characters.", err)
		} else {
			logger.ErrorEventWithUUID(tid, metadataPublishEvent.UUID, errMsg, err)
		}
		logger.MonitoringValidationEvent(tid, metadataPublishEvent.UUID, contentType, err.Error(), false)
		return err
	}

	annotations := []annotation{}
	for _, value := range taxonomyHandlers {
		annotations = append(annotations, value.buildAnnotations(metadata)...)
	}

	conceptAnnotations := ConceptAnnotations{UUID: metadataPublishEvent.UUID, Annotations: annotations}

	marshalledAnnotations, err := json.Marshal(conceptAnnotations)
	if err != nil {
		logger.ErrorEventWithUUID(tid, metadataPublishEvent.UUID, "Error marshalling concept annotations", err)
		return err
	}

	var headers = buildConceptAnnotationsHeader(msg.Headers)
	message := kafka.FTMessage{Headers: headers, Body: string(marshalledAnnotations)}
	err = messageProducer.SendMessage(message)
	if err != nil {
		logger.ErrorEventWithUUID(tid, metadataPublishEvent.UUID, "Error sending concept annotations to queue", err)
		return err
	}

	logger.MonitoringValidationEvent(tid, metadataPublishEvent.UUID, contentType, "Successfully mapped", true)
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
