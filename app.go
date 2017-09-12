package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/kafka-client-go/kafka"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/twinj/uuid"
)

const (
	messageTimestampDateFormat = "2006-01-02T15:04:05.000Z"
	serviceName                = "annotations-mapper"
	contentType                = "Annotations"
	consumerStartedEvent       = "consume_queue"
	producerStartedEvent       = "produce_queue"
	mapperEvent                = "Map"
)

var (
	messageConsumer  kafka.Consumer
	messageProducer  kafka.Producer
	taxonomyHandlers map[string]TaxonomyService
	whitelist        *regexp.Regexp
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
			logger.Fatalf(nil, err, "Please specify a valid whitelist")
		}

		messageProducer, err = kafka.NewProducer(*brokerAddress, *producerTopic)
		if err != nil {
			logger.Fatalf(nil, err, "Cannot start message producer")
		}
		logger.Infof(map[string]interface{}{"event": consumerStartedEvent}, "Starting queue consumer: %v", *producerTopic)

		messageConsumer, err = kafka.NewConsumer(*zookeeperAddress, *consumerGroup, []string{*consumerTopic}, kafka.DefaultConsumerConfig())
		if err != nil {
			logger.Fatalf(nil, err, "Cannot start message consumer")
		}
		logger.Infof(map[string]interface{}{"event": consumerStartedEvent}, "Starting queue producer: %s", *consumerTopic)
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
		logger.Fatalf(nil, err, "Couldn't set up HTTP listener")
	}
}

func handleMessage(msg kafka.FTMessage) error {
	tid := msg.Headers["X-Request-Id"]
	systemCode := msg.Headers["Origin-System-Id"]
	if !whitelist.MatchString(systemCode) {
		logger.NewEntry(tid).Infof("Skipping annotations published with Origin-System-Id \"%v\". It does not match the configured whitelist.", systemCode)
		return nil
	}

	// There is no proper validation in place for annotations. Everything that is parsable will be considered as being valid, everything that is not as being invalid.
	// This behaviour can be changed when a proper validation will be introduced.

	// Consider the message as invalid - logging all the error messages for this transaction as monitoring events
	msgIsValid := false

	var metadataPublishEvent MetadataPublishEvent
	err := json.Unmarshal([]byte(msg.Body), &metadataPublishEvent)
	if err != nil {
		logger.NewMonitoringEntry(mapperEvent, tid, contentType).WithValidFlag(msgIsValid).WithError(err).Error("Cannot unmarshal message body")
		return err
	}

	logger.NewEntry(tid).WithUUID(metadataPublishEvent.UUID).Info("Processing metadata publish event")

	metadataXML, err := base64.StdEncoding.DecodeString(metadataPublishEvent.Value)
	if err != nil {
		logger.NewMonitoringEntry(mapperEvent, tid, contentType).WithValidFlag(msgIsValid).WithUUID(metadataPublishEvent.UUID).WithError(err).Error("Error decoding body")
		return err
	}

	metadata, err, hadInvalidChars := unmarshalMetadata(metadataXML)
	if err != nil {
		errMsg := "Error unmarshalling metadata XML"
		if hadInvalidChars {
			logger.NewEntry(tid).WithUUID(metadataPublishEvent.UUID).WithError(err).Errorf("%s Metadata XML had invalid UTF8 characters.", errMsg)
		} else {
			logger.NewEntry(tid).WithUUID(metadataPublishEvent.UUID).WithError(err).Errorf("%s", errMsg)
		}

		// Log validation error as a monitoring event
		entry := logger.NewMonitoringEntry(mapperEvent, tid, contentType).WithValidFlag(msgIsValid)
		if metadataPublishEvent.UUID != "" {
			entry = entry.WithUUID(metadataPublishEvent.UUID)
		}
		entry.WithError(err).Error("Message is not valid due to parsing issues. ")
		return err
	}

	// if the message had no parsing errors: consider it as valid
	msgIsValid = true
	annotations := []annotation{}
	for _, value := range taxonomyHandlers {
		annotations = append(annotations, value.buildAnnotations(metadata)...)
	}

	conceptAnnotations := ConceptAnnotations{UUID: metadataPublishEvent.UUID, Annotations: annotations}

	marshalledAnnotations, err := json.Marshal(conceptAnnotations)
	if err != nil {
		logger.NewMonitoringEntry(mapperEvent, tid, contentType).WithUUID(metadataPublishEvent.UUID).WithValidFlag(msgIsValid).WithError(err).Error("Error marshalling concept annotations")
		return err
	}

	var headers = buildConceptAnnotationsHeader(msg.Headers)
	message := kafka.FTMessage{Headers: headers, Body: string(marshalledAnnotations)}
	err = messageProducer.SendMessage(message)
	if err != nil {
		logger.NewMonitoringEntry(mapperEvent, tid, contentType).WithUUID(metadataPublishEvent.UUID).WithValidFlag(msgIsValid).WithError(err).Error("Error sending concept annotations to queue")
		return err
	}

	logger.NewMonitoringEntry(mapperEvent, tid, contentType).WithUUID(metadataPublishEvent.UUID).
		WithValidFlag(msgIsValid).Info("Successfully mapped")
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
