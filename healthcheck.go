package main

import (
	"net/http"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/kafka-client-go/kafka"
	"github.com/Financial-Times/service-status-go/gtg"
)

type HealthCheck struct {
	consumer kafka.Consumer
	producer kafka.Producer
}

func NewHealthCheck(c kafka.Consumer, p kafka.Producer) *HealthCheck {
	return &HealthCheck{
		consumer: c,
		producer: p,
	}
}

func (h *HealthCheck) Health() func(w http.ResponseWriter, r *http.Request) {
	hc := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{
			SystemCode:  "annotations-mapper",
			Name:        "annotations-mapper",
			Description: "Checks if all the dependent services are reachable and healthy.",
			Checks:      []fthealth.Check{h.kafkaConsumerCheck(), h.kafkaProducerCheck()},
		},
		Timeout: 10 * time.Second,
	}
	return fthealth.Handler(hc)
}

func (h *HealthCheck) kafkaConsumerCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "kafka-consumer-connected",
		Name:             "Kafka Consumer is Connected",
		Severity:         2,
		BusinessImpact:   "Content V1 Metadata can't be read from queue. This will negatively impact V1 metadata availability.",
		TechnicalSummary: "Read message queue is not reachable/healthy",
		PanicGuide:       "https://dewey.ft.com/annotations-mapper.html",
		Checker:          h.checkKafkaConsumerConnectivity,
	}
}

func (h *HealthCheck) kafkaProducerCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "kafka-producer-connected",
		Name:             "Kafka Producer is Connected",
		Severity:         2,
		BusinessImpact:   "Content V1 Metadata can't be written to the queue. This will negatively impact V1 metadata availability.",
		TechnicalSummary: "This will cause map failures as messages will not be produced to the next queue after being mapped.",
		PanicGuide:       "https://dewey.ft.com/annotations-mapper.html",
		Checker:          h.checkKafkaProducerConnectivity,
	}
}

func (h *HealthCheck) GTG() gtg.Status {
	consumerCheck := func() gtg.Status {
		return gtgCheck(h.checkKafkaConsumerConnectivity)
	}
	producerCheck := func() gtg.Status {
		return gtgCheck(h.checkKafkaProducerConnectivity)
	}

	return gtg.FailFastParallelCheck([]gtg.StatusChecker{
		consumerCheck,
		producerCheck,
	})()
}

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	return gtg.Status{GoodToGo: true}
}

func (h *HealthCheck) checkKafkaConsumerConnectivity() (string, error) {
	if err := h.consumer.ConnectivityCheck(); err != nil {
		return "Error connecting with Kafka", err
	}
	return "Successfully connected to Kafka", nil
}

func (h *HealthCheck) checkKafkaProducerConnectivity() (string, error) {
	if err := h.producer.ConnectivityCheck(); err != nil {
		return "Error checking kafka producer connectivity", err
	}
	return "Successfully connected to Kafka", nil
}
