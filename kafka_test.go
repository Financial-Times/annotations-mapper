package main

import (
	"regexp"
	"testing"

	"github.com/Financial-Times/go-logger"
	"github.com/pborman/uuid"

	"github.com/Financial-Times/kafka-client-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleMessage__UnsupportedSystemCode(t *testing.T) {
	whitelist = regexp.MustCompile("http://cmdb\\.ft\\.com/systems/methode-web-pub")

	msg := kafka.FTMessage{}
	msg.Headers = make(map[string]string)
	msg.Headers["Origin-System-Id"] = "http://cmdb.ft.com/systems/pac"

	hook := logger.NewTestHook("")
	err := handleMessage(msg)
	assert.NoError(t, err)

	logLine := hook.LastEntry().Message
	assert.Equal(t, "Skipping annotations published with Origin-System-Id \"http://cmdb.ft.com/systems/pac\". It does not match the configured whitelist.", logLine)
}

func TestHandleMessage__InvalidJSONBody(t *testing.T) {
	whitelist = regexp.MustCompile("http://cmdb\\.ft\\.com/systems/methode-web-pub")

	msg := kafka.FTMessage{}
	msg.Headers = make(map[string]string)
	msg.Headers["Origin-System-Id"] = "http://cmdb.ft.com/systems/methode-web-pub"
	msg.Body = ``

	hook := logger.NewTestHook("")

	err := handleMessage(msg)
	assert.Error(t, err) // should fail parsing json

	logLine := hook.LastEntry()
	assert.Equal(t, "Cannot unmarshal message body", logLine.Message)

	isValid, ok := logLine.Data["isValid"]
	require.True(t, ok)
	require.Equal(t, "false", isValid.(string))
}

func TestHandleMessage__InvalidBase64Value(t *testing.T) {
	msg := kafka.FTMessage{}
	msg.Headers = make(map[string]string)
	msg.Headers["Origin-System-Id"] = "http://cmdb.ft.com/systems/methode-web-pub"

	testUUID := uuid.New()
	msg.Body = `{"uuid":"` + testUUID + `","value":"I AM NOT BASE64!"}`

	hook := logger.NewTestHook("")

	err := handleMessage(msg)
	assert.Error(t, err) // should fail parsing json

	logLine := hook.LastEntry()
	assert.Equal(t, "Error decoding body", logLine.Message)

	require.Equal(t, "false", logLine.Data["isValid"].(string))
	require.Equal(t, testUUID, logLine.Data["uuid"].(string))
}

func TestHandleMessage__InvalidXML(t *testing.T) {
	msg := kafka.FTMessage{}
	msg.Headers = make(map[string]string)
	msg.Headers["Origin-System-Id"] = "http://cmdb.ft.com/systems/methode-web-pub"

	testUUID := uuid.New()
	msg.Body = `{"uuid":"` + testUUID + `","value":"eyJtc2ciOiJOb3QgWE1MIn0="}`

	hook := logger.NewTestHook("")

	err := handleMessage(msg)
	assert.Error(t, err) // should fail parsing json

	logLine := hook.LastEntry()
	assert.Equal(t, "Message is not valid due to parsing issues.", logLine.Message)

	require.Equal(t, "false", logLine.Data["isValid"].(string))
	require.Equal(t, testUUID, logLine.Data["uuid"].(string))
}
