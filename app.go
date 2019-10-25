package main

import (
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/kafka-client-go/kafka"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
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

		logger.Infof(map[string]interface{}{"event": consumerStartedEvent}, "Starting queue consumer: %v", *producerTopic)
		messageConsumer, _ = kafka.NewPerseverantConsumer(*zookeeperAddress, *consumerGroup, []string{*consumerTopic}, kafka.DefaultConsumerConfig(), time.Minute, nil)

		logger.Infof(map[string]interface{}{"event": consumerStartedEvent}, "Starting queue producer: %s", *consumerTopic)
		messageProducer, _ = kafka.NewPerseverantProducer(*brokerAddress, *producerTopic, nil, 0, time.Minute)

		go startKafkaConsumer(messageConsumer)
		go startServer(messageConsumer, messageProducer)

		waitForSignal()
		messageConsumer.Shutdown()
	}

	app.Run(os.Args)
}

func startServer(messageConsumer kafka.Consumer, messageProducer kafka.Producer) {
	hc := NewHealthCheck(messageConsumer, messageProducer)
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

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
