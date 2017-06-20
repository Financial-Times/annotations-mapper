package main

import (
	"os"

	"encoding/base64"
	"encoding/json"
	"encoding/xml"

	"time"

	"unicode/utf8"

	"github.com/Financial-Times/kafka-client-go/kafka"
	"github.com/jawher/mow.cli"
	"github.com/twinj/uuid"
	"os/signal"
	"syscall"
)

const messageTimestampDateFormat = "2006-01-02T15:04:05.000Z"

var (
	messageConsumer  kafka.Consumer
	messageProducer  kafka.Producer
	logger           *AppLogger
	taxonomyHandlers map[string]TaxonomyService
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
	app := cli.App("annotations-mapper", "A service to read V1 metadata publish event, filter it and output UP-specific metadata to the destination queue.")
	zookeeperAddress := app.String(cli.StringOpt{
		Name:   "zookeeperAddress",
		Value:  "localhost:2181",
		Desc:   "Addresses used by the queue consumer to connect to the queue",
		EnvVar: "ZOOKEEPER_ADDRESS",
	})
	sourceGroup := app.String(cli.StringOpt{
		Name:   "source-group",
		Value:  "",
		Desc:   "Group used to read the messages from the queue",
		EnvVar: "SRC_GROUP",
	})
	sourceTopic := app.String(cli.StringOpt{
		Name:   "source-topic",
		Value:  "",
		Desc:   "The topic to read the meassages from",
		EnvVar: "SRC_TOPIC",
	})
	destinationAddress := app.String(cli.StringOpt{
		Name:   "destination-address",
		Value:  "",
		Desc:   "Address used by the producer to connect to the queue",
		EnvVar: "DEST_ADDRESS",
	})
	destinationTopic := app.String(cli.StringOpt{
		Name:   "destination-topic",
		Value:  "",
		Desc:   "The topic to write the concept annotation to",
		EnvVar: "DEST_TOPIC",
	})

	app.Action = func() {
		var err error
		messageProducer, err = kafka.NewProducer(*destinationAddress, *destinationTopic)
		if err != nil {
			logger.Fatal("Cannot start message producer", err)
		}
		logger.QueueProducerStarted(*destinationTopic)

		messageConsumer, err = kafka.NewConsumer(*zookeeperAddress, *sourceGroup, []string{*sourceTopic}, kafka.DefaultConsumerConfig())
		if err != nil {
			logger.Fatal("Cannot start message consumer", err)
		}
		logger.QueueConsumerStarted(*sourceTopic)
		messageConsumer.StartListening(handleMessage)

		waitForSignal()
		logger.Info("Shutting down Kafka consumer", "", "")
		messageConsumer.Shutdown()
	}

	app.Run(os.Args)
}

func handleMessage(msg kafka.FTMessage) error {
	tid := msg.Headers["X-Request-Id"]

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
