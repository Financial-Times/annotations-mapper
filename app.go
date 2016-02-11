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

	"github.com/Financial-Times/message-queue-go-producer/producer"
	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/Financial-Times/v1-suggestor/model"
	"github.com/Financial-Times/v1-suggestor/service"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/kr/pretty"
	"github.com/twinj/uuid"
)

var messageProducer producer.MessageProducer
var taxonomyHandlers = make(map[string]service.TaxonomyService)

func main() {
	app := cli.App("V1 suggestor", "A service to read V1 metadata publish event, filter it and output UP-specific metadata to the destination queue.")
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
		Desc:   "The topic to write the concept suggestion to",
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
		infoLogger.Printf("Using source configuration: %# v", pretty.Formatter(srcConf))
		infoLogger.Printf("Using dest configuration: %# v", pretty.Formatter(destConf))

		setupTaxonomyHandlers()
		go enableHealthChecks(srcConf, destConf)

		initializeProducer(destConf)
		readMessages(srcConf)
	}

	app.Run(os.Args)
}

func setupTaxonomyHandlers() {
	taxonomyHandlers["subjects"] = service.SubjectService{HandledTaxonomy: "subjects"}
}

func enableHealthChecks(srcConf consumer.QueueConfig, destConf producer.MessageProducerConfig) {
	healthCheck := &Healthcheck{
		client:   http.Client{},
		srcConf:  srcConf,
		destConf: destConf}
	router := mux.NewRouter()
	router.HandleFunc("/__health", healthCheck.checkHealth())
	router.HandleFunc("/__gtg", healthCheck.gtg)
	http.Handle("/", router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		errorLogger.Panicf("Couldn't set up HTTP listener: %v\n", err)
	}
}

func initializeProducer(config producer.MessageProducerConfig) {
	messageProducer = producer.NewMessageProducer(config)
	infoLogger.Printf("Producer: %# v", pretty.Formatter(messageProducer))
}

func readMessages(config consumer.QueueConfig) {
	messageConsumer := consumer.NewConsumer(config, handleMessage, http.Client{})
	infoLogger.Printf("Consumer: %# v", pretty.Formatter(messageConsumer))

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
	var metadataPublishEvent model.MetadataPublishEvent
	err := json.Unmarshal([]byte(msg.Body), &metadataPublishEvent)
	if err != nil {
		errorLogger.Printf("Cannot unmarshal message body:[%v]", err.Error())
		return
	}

	infoLogger.Printf("Processing metadata publish event for uuid [%s]", metadataPublishEvent.UUID)

	metadataXML, err := base64.StdEncoding.DecodeString(metadataPublishEvent.Value)
	if err != nil {
		errorLogger.Printf("Error decoding body for uuid:  [%s]", err.Error())
		return
	}

	metadata := model.ContentRef{}
	err = xml.Unmarshal(metadataXML, &metadata)
	if err != nil {
		errorLogger.Printf("Error unmarshalling metadata XML: [%v]", err.Error())
		return
	}

	var suggestions []model.Suggestion
	for key, value := range taxonomyHandlers {
		infoLogger.Printf("Processing taxonomy [%s]", key)
		suggestions = append(suggestions, value.BuildSuggestions(metadata)...)
	}

	conceptSuggestion := model.ConceptSuggestion{UUID: metadataPublishEvent.UUID, Suggestions: suggestions}

	marshalledSuggestions, err := json.Marshal(conceptSuggestion)
	if err != nil {
		errorLogger.Printf("Error marshalling the concept suggestions: [%v]", err.Error())
		return
	}

	var headers = buildConceptSuggestionsHeader(msg.Headers)
	message := producer.Message{Headers: headers, Body: string(marshalledSuggestions)}
	err = messageProducer.SendMessage(conceptSuggestion.UUID, message)
	if err != nil {
		errorLogger.Printf("Error sending concept suggestion to queue: [%v]", err.Error())
	}

	infoLogger.Printf("Sent suggestion message with message ID [%s]to queue.", headers["Message-Id"])
}

func buildConceptSuggestionsHeader(publishEventHeaders map[string]string) map[string]string {
	return map[string]string{
		"Message-Id":        uuid.NewV4().String(),
		"Message-Type":      "concept-suggestions",
		"Content-Type":      publishEventHeaders["Content-Type"],
		"X-Request-Id":      publishEventHeaders["X-Request-Id"],
		"Origin-System-Id":  publishEventHeaders["Origin-System-Id"],
		"Message-Timestamp": publishEventHeaders["Message-Timestamp"],
	}
}