package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"encoding/base64"
	"encoding/json"
	"encoding/xml"

	"time"

	"unicode/utf8"

	"github.com/Financial-Times/message-queue-go-producer/producer"
	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/kr/pretty"
	"github.com/twinj/uuid"
)

var messageProducer producer.MessageProducer

var taxonomyHandlers = map[string]TaxonomyService{
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

const messageTimestampDateFormat = "2006-01-02T15:04:05.000Z"

func main() {
	app := cli.App("annotations-mapper", "A service to read V1 metadata publish event, filter it and output UP-specific metadata to the destination queue.")
	sourceAddresses := app.Strings(cli.StringsOpt{
		Name:   "source-addresses",
		Value:  []string{},
		Desc:   "Addresses used by the queue consumer to connect to the queue",
		EnvVar: "SRC_ADDR",
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
	sourceQueue := app.String(cli.StringOpt{
		Name:   "source-queue",
		Value:  "",
		Desc:   "Thew queue to read the messages from",
		EnvVar: "SRC_QUEUE",
	})
	sourceConcurrentProcessing := app.Bool(cli.BoolOpt{
		Name:   "source-concurrent-processing",
		Value:  false,
		Desc:   "Whether the consumer uses concurrent processing for the messages",
		EnvVar: "SRC_CONCURRENT_PROCESSING",
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
	destinationQueue := app.String(cli.StringOpt{
		Name:   "destination-queue",
		Value:  "",
		Desc:   "The queue used by the producer",
		EnvVar: "DEST_QUEUE",
	})

	app.Action = func() {
		srcConf := consumer.QueueConfig{
			Addrs:                *sourceAddresses,
			Group:                *sourceGroup,
			Topic:                *sourceTopic,
			Queue:                *sourceQueue,
			ConcurrentProcessing: *sourceConcurrentProcessing,
		}

		destConf := producer.MessageProducerConfig{
			Addr:  *destinationAddress,
			Topic: *destinationTopic,
			Queue: *destinationQueue,
		}

		initLogs(os.Stdout, os.Stdout, os.Stderr)
		infoLogger.Printf("[Startup] Using source configuration: %# v", pretty.Formatter(srcConf))
		infoLogger.Printf("[Startup] Using dest configuration: %# v", pretty.Formatter(destConf))

		infoLogger.Printf("[Startup] Handling taxonomies:")
		for key := range taxonomyHandlers {
			infoLogger.Printf("\t %v", key)
		}

		go enableHealthChecks(srcConf, destConf)

		initializeProducer(destConf)
		readMessages(srcConf)
	}

	app.Run(os.Args)
}

func enableHealthChecks(srcConf consumer.QueueConfig, destConf producer.MessageProducerConfig) {
	healthCheck := &Healthcheck{
		client:   http.Client{},
		srcConf:  srcConf,
		destConf: destConf}
	router := mux.NewRouter()
	router.HandleFunc("/__health", healthCheck.checkHealth())
	router.HandleFunc("/__gtg", healthCheck.gtg)
	router.HandleFunc(status.PingPath, status.PingHandler)
	router.HandleFunc(status.PingPathDW, status.PingHandler)
	router.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
	router.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)
	http.Handle("/", router)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		errorLogger.Panicf("Couldn't set up HTTP listener: %v\n", err)
	}
}

func initializeProducer(config producer.MessageProducerConfig) {
	messageProducer = producer.NewMessageProducer(config)
	infoLogger.Printf("[Startup] Producer: %# v", pretty.Formatter(messageProducer))
}

func readMessages(config consumer.QueueConfig) {
	messageConsumer := consumer.NewConsumer(config, handleMessage, &http.Client{})
	infoLogger.Printf("[Startup] Consumer: %# v", pretty.Formatter(messageConsumer))

	var consumerWaitGroup sync.WaitGroup
	consumerWaitGroup.Add(1)

	go func() {
		messageConsumer.Start()
		consumerWaitGroup.Done()
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	messageConsumer.Stop()
	consumerWaitGroup.Wait()
}

func handleMessage(msg consumer.Message) {
	tid := msg.Headers["X-Request-Id"]

	var metadataPublishEvent MetadataPublishEvent
	err := json.Unmarshal([]byte(msg.Body), &metadataPublishEvent)
	if err != nil {
		errorLogger.Printf("[%s] Cannot unmarshal message body:[%v]", tid, err.Error())
		return
	}

	infoLogger.Printf("[%s] Processing metadata publish event for uuid [%s]", tid, metadataPublishEvent.UUID)

	metadataXML, err := base64.StdEncoding.DecodeString(metadataPublishEvent.Value)
	if err != nil {
		errorLogger.Printf("[%s] Error decoding body for uuid:  [%s]", tid, err.Error())
		return
	}

	metadata, err, hadInvalidChars := unmarshalMetadata(metadataXML)

	if err != nil {
		errorLogger.Printf("[%s] Error unmarshalling metadata XML for UUID [%v]: [%v]", tid, metadataPublishEvent.UUID, err.Error())
		if hadInvalidChars {
			infoLogger.Printf("[%s] Metadata XML for UUID [%s] had invalid UTF8 characters.", tid, metadataPublishEvent.UUID)
		}
		return
	}

	annotations := []annotation{}
	for key, value := range taxonomyHandlers {
		infoLogger.Printf("[%s] Processing taxonomy [%s]", tid, key)
		annotations = append(annotations, value.buildAnnotations(metadata)...)
	}

	conceptAnnotations := ConceptAnnotations{UUID: metadataPublishEvent.UUID, Annotations: annotations}

	marshalledAnnotations, err := json.Marshal(conceptAnnotations)
	if err != nil {
		errorLogger.Printf("[%s] Error marshalling the concept annotations for UUID [%v]: [%v]", tid, metadataPublishEvent.UUID, err.Error())
		return
	}

	var headers = buildConceptAnnotationsHeader(msg.Headers)
	message := producer.Message{Headers: headers, Body: string(marshalledAnnotations)}
	err = messageProducer.SendMessage(conceptAnnotations.UUID, message)
	if err != nil {
		errorLogger.Printf("[%s] Error sending concept annotation to queue for UUID [%v]: [%v]", tid, metadataPublishEvent.UUID, err.Error())
	}

	infoLogger.Printf("[%s] Sent annotation message for [%s] with message ID [%s] to queue.", tid, metadataPublishEvent.UUID, headers["Message-Id"])
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
