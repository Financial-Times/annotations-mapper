package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"time"
	"unicode/utf8"

	logger "github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/kafka-client-go/kafka"
	"github.com/twinj/uuid"
)

func startKafkaConsumer(messageConsumer kafka.Consumer) {
	messageConsumer.Consume(handleMessage)
}

func handleMessage(msg kafka.FTMessage) error {
	tid := msg.Headers["X-Request-Id"]
	log := logger.NewEntry(tid)

	systemCode := msg.Headers["Origin-System-Id"]
	if !whitelist.MatchString(systemCode) {
		log.Infof("Skipping annotations published with Origin-System-Id \"%v\". It does not match the configured whitelist.", systemCode)
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

	log.WithUUID(metadataPublishEvent.UUID).Info("Processing metadata publish event")

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
		entry.WithError(err).Error("Message is not valid due to parsing issues.")
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
