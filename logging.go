package main

import (
	"github.com/Sirupsen/logrus"
)

const (
	serviceStartedEvent  = "service_started"
	consumerStartedEvent = "consume_queue"
	producerStartedEvent = "produce_queue"
	mappingEvent         = "mapping"
	errorEvent           = "error"
)

type AppLogger struct {
	Log *logrus.Logger
}

func NewAppLogger() *AppLogger {
	logrus.SetLevel(logrus.InfoLevel)
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)

	return &AppLogger{Log: log}
}

func (appLogger *AppLogger) ServiceStartedEvent(port int) {
	event := map[string]interface{}{
		"event": serviceStartedEvent,
	}
	appLogger.Log.WithFields(event).Infof("Starting to listen on port [%d]", port)
}

func (appLogger *AppLogger) QueueConsumerStarted(queueTopic string) {
	event := map[string]interface{}{
		"event": consumerStartedEvent,
	}
	appLogger.Log.WithFields(event).Infof("Starting queue consumer: %v", queueTopic)
}

func (appLogger *AppLogger) QueueProducerStarted(queueTopic string) {
	event := map[string]interface{}{
		"event": producerStartedEvent,
	}
	appLogger.Log.WithFields(event).Infof("Starting queue producer: %v", queueTopic)
}

func (appLogger *AppLogger) Info(message string, transactionID string, contentUUID string) {
	event := map[string]interface{}{
		"event":          mappingEvent,
		"transaction_id": transactionID,
	}
	if contentUUID != "" {
		event["uuid"] = contentUUID
	}

	appLogger.Log.WithFields(event).Info(message)
}

func (appLogger *AppLogger) Warn(message string, transactionID string, contentUUID string, err error) {
	event := map[string]interface{}{
		"event": errorEvent,
	}

	if err != nil {
		event["error"] = err
	}
	if transactionID != "" {
		event["transaction_id"] = transactionID
	}
	if contentUUID != "" {
		event["uuid"] = contentUUID
	}

	appLogger.Log.WithFields(event).Warn(message)
}

func (appLogger *AppLogger) Error(message string, transactionID string, contentUUID string, err error) {
	event := map[string]interface{}{
		"event": errorEvent,
		"error": err,
	}

	if transactionID != "" {
		event["transaction_id"] = transactionID
	}
	if contentUUID != "" {
		event["uuid"] = contentUUID
	}

	appLogger.Log.WithFields(event).Error(message)
}

func (appLogger *AppLogger) Fatal(message string, err error) {
	event := map[string]interface{}{
		"event": errorEvent,
		"error": err,
	}

	appLogger.Log.WithFields(event).Fatal(message)
}
