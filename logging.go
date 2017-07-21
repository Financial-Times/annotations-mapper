package main

import "github.com/Financial-Times/go-logger"

const (
	consumerStartedEvent = "consume_queue"
	producerStartedEvent = "produce_queue"
)

type AppLogger struct {
	*go_logger.AppLogger
}

func NewAppLogger(serviceName string) *AppLogger {
	return &AppLogger{go_logger.NewLogger(serviceName)}
}

func (appLogger *AppLogger) QueueConsumerStarted(queueTopic string) {
	event := map[string]interface{}{
		"event": consumerStartedEvent,
	}
	appLogger.Infof(event,"Starting queue consumer: %v", queueTopic)
}

func (appLogger *AppLogger) QueueProducerStarted(queueTopic string) {
	event := map[string]interface{}{
		"event": producerStartedEvent,
	}
	appLogger.Infof(event,"Starting queue producer: %v", queueTopic)
}