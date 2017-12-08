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
}

func NewHealthCheck(c kafka.Consumer) *HealthCheck {
	return &HealthCheck{
		consumer: c,
	}
}

func (h *HealthCheck) Health() func(w http.ResponseWriter, r *http.Request) {
	hc := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{
			SystemCode:  "annotations-mapper",
			Name:        "annotations-mapper",
			Description: "Checks if all the dependent services are reachable and healthy.",
			Checks:      []fthealth.Check{h.readQueueCheck()},
		},
		Timeout: 10 * time.Second,
	}
	return fthealth.Handler(hc)
}

func (h *HealthCheck) readQueueCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "read-message-queue-reachable",
		Name:             "Read Message Queue Reachable",
		Severity:         1,
		BusinessImpact:   "Content V1 Metadata can't be read from queue. This will negatively impact V1 metadata availability.",
		TechnicalSummary: "Read message queue is not reachable/healthy",
		PanicGuide:       "https://dewey.ft.com/annotations-mapper.html",
		Checker:          h.checkKafkaConnectivity,
	}
}

func (h *HealthCheck) GTG() gtg.Status {
	consumerCheck := func() gtg.Status {
		return gtgCheck(h.checkKafkaConnectivity)
	}
	producerCheck := func() gtg.Status {
		return gtgCheck(h.checkKafkaConnectivity)
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

func (h *HealthCheck) checkKafkaConnectivity() (string, error) {
	if err := h.consumer.ConnectivityCheck(); err != nil {
		return "Error connecting with Kafka", err
	}
	return "Successfully connected to Kafka", nil
}
