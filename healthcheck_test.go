package main

import (
	"net/http/httptest"
	"testing"

	"github.com/Financial-Times/kafka-client-go/kafka"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mockConsumer struct {
	err error
}

func (mc mockConsumer) ConnectivityCheck() error {
	return mc.err
}

func (mc mockConsumer) StartListening(messageHandler func(message kafka.FTMessage) error) {
	return
}

func (mc mockConsumer) Shutdown() {
	return
}

func TestHappyHealthCheck(t *testing.T) {
	hc := HealthCheck{mockConsumer{nil}}

	req := httptest.NewRequest("GET", "http://example.com/__health", nil)
	w := httptest.NewRecorder()

	hc.Health()(w, req)

	assert.Equal(t, 200, w.Code, "It should return HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"name":"Read Message Queue Proxy Reachable","ok":true`, "Read message queue proxy healthcheck should be happy")
}

func TestHealthCheckWithUnhappyConsumer(t *testing.T) {
	hc := HealthCheck{mockConsumer{errors.New("Error connecting to the queue")}}

	req := httptest.NewRequest("GET", "http://example.com/__health", nil)
	w := httptest.NewRecorder()

	hc.Health()(w, req)

	assert.Equal(t, 200, w.Code, "It should return HTTP 200 OK")
	assert.Contains(t, w.Body.String(), `"name":"Read Message Queue Proxy Reachable","ok":false`, "Read message queue proxy healthcheck should be unhappy")
}

func TestGTGHappyFlow(t *testing.T) {
	hc := HealthCheck{mockConsumer{nil}}

	status := hc.GTG()
	assert.True(t, status.GoodToGo)
	assert.Empty(t, status.Message)
}

func TestGTGBrokenConsumer(t *testing.T) {
	hc := HealthCheck{mockConsumer{errors.New("Error connecting to the queue")}}

	status := hc.GTG()
	assert.False(t, status.GoodToGo)
	assert.Equal(t, "Error connecting to the queue", status.Message)
}
