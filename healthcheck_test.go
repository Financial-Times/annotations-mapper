package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Times/kafka-client-go/kafka"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
)

type mockKafkaConnection struct {
	err error
}

func (mc mockKafkaConnection) ConnectivityCheck() error {
	return mc.err
}

func (mc mockKafkaConnection) StartListening(messageHandler func(message kafka.FTMessage) error) {
	return
}

func (mc mockKafkaConnection) SendMessage(message kafka.FTMessage) error {
	return nil
}

func (mc mockKafkaConnection) Shutdown() {
	return
}

func unmarshalHealthcheck(r io.Reader) (fthealth.HealthResult, error) {
	dec := json.NewDecoder(r)
	res := fthealth.HealthResult{}
	err := dec.Decode(&res)
	return res, err
}

func assertHealthcheckOk(t *testing.T, healthchecks fthealth.HealthResult, id string) {
	checkExists := false
	for _, check := range healthchecks.Checks {
		if check.ID == id {
			assert.True(t, check.Ok)
			checkExists = true
		}
	}

	assert.True(t, checkExists)
}

func assertHealthcheckFailed(t *testing.T, healthchecks fthealth.HealthResult, id string, expectedOutput string) {
	assert.False(t, healthchecks.Ok)

	checkExists := false
	for _, check := range healthchecks.Checks {
		if check.ID == id {
			assert.False(t, check.Ok)
			assert.Equal(t, expectedOutput, check.CheckOutput)
			checkExists = true
		}
	}

	assert.True(t, checkExists)
}

func TestHappyHealthCheck(t *testing.T) {
	hc := NewHealthCheck(mockKafkaConnection{err: nil}, mockKafkaConnection{err: nil})

	req := httptest.NewRequest("GET", "http://example.com/__health", nil)
	w := httptest.NewRecorder()

	hc.Health()(w, req)

	assert.Equal(t, 200, w.Code, "It should return HTTP 200 OK")

	health, err := unmarshalHealthcheck(w.Body)
	require.NoError(t, err)
	assert.True(t, health.Ok)

	assertHealthcheckOk(t, health, "kafka-consumer-connected")
	assertHealthcheckOk(t, health, "kafka-producer-connected")
}

func TestHealthCheck__UnhappyConsumer(t *testing.T) {
	hc := NewHealthCheck(mockKafkaConnection{err: errors.New("uh-oh consumption failed")}, mockKafkaConnection{err: nil})

	req := httptest.NewRequest("GET", "http://example.com/__health", nil)
	w := httptest.NewRecorder()

	hc.Health()(w, req)

	assert.Equal(t, 200, w.Code, "It should return HTTP 200 OK")
	health, err := unmarshalHealthcheck(w.Body)
	require.NoError(t, err)

	assertHealthcheckFailed(t, health, "kafka-consumer-connected", "uh-oh consumption failed")
	assertHealthcheckOk(t, health, "kafka-producer-connected")
}

func TestHealthCheck__UnhappyProducer(t *testing.T) {
	hc := NewHealthCheck(mockKafkaConnection{err: nil}, mockKafkaConnection{err: errors.New("uh-oh production failed")})
	req := httptest.NewRequest("GET", "http://example.com/__health", nil)
	w := httptest.NewRecorder()

	hc.Health()(w, req)

	assert.Equal(t, 200, w.Code, "It should return HTTP 200 OK")
	health, err := unmarshalHealthcheck(w.Body)
	require.NoError(t, err)

	assertHealthcheckOk(t, health, "kafka-consumer-connected")
	assertHealthcheckFailed(t, health, "kafka-producer-connected", "uh-oh production failed")
}

func TestGTGHappyFlow(t *testing.T) {
	hc := NewHealthCheck(mockKafkaConnection{err: nil}, mockKafkaConnection{err: nil})

	status := hc.GTG()
	assert.True(t, status.GoodToGo)
	assert.Empty(t, status.Message)
}

func TestGTG__BrokenConsumer(t *testing.T) {
	hc := NewHealthCheck(mockKafkaConnection{err: errors.New("Error connecting to the queue")}, mockKafkaConnection{err: nil})

	status := hc.GTG()
	assert.False(t, status.GoodToGo)
	assert.Equal(t, "Error connecting to the queue", status.Message)
}

func TestGTG__BrokenProducer(t *testing.T) {
	hc := NewHealthCheck(mockKafkaConnection{err: nil}, mockKafkaConnection{err: errors.New("Error connecting to the queue")})

	status := hc.GTG()
	assert.False(t, status.GoodToGo)
	assert.Equal(t, "Error connecting to the queue", status.Message)
}
