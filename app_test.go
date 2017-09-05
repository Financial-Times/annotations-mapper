package main

import (
	"regexp"
	"testing"

	"github.com/Financial-Times/kafka-client-go/kafka"
	hooks "github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestHandleMessageWithUnsupportedSystemCode(t *testing.T) {
	whitelist = regexp.MustCompile("http://cmdb\\.ft\\.com/systems/methode-web-pub")

	msg := kafka.FTMessage{}
	msg.Headers = make(map[string]string)
	msg.Headers["Origin-System-Id"] = "http://cmdb.ft.com/systems/pac"

	hook := hooks.NewLocal(logger.Log)

	err := handleMessage(msg)
	assert.NoError(t, err)

	logLine := hook.LastEntry().Message
	assert.Equal(t, "Skipping annotations published with Origin-System-Id \"http://cmdb.ft.com/systems/pac\". It does not match the configured whitelist.", logLine)
}

func TestHandleMessageWithSupportedSystemCode(t *testing.T) {
	whitelist = regexp.MustCompile("http://cmdb\\.ft\\.com/systems/methode-web-pub")

	msg := kafka.FTMessage{}
	msg.Headers = make(map[string]string)
	msg.Headers["Origin-System-Id"] = "http://cmdb.ft.com/systems/methode-web-pub"
	msg.Body = ``

	hook := hooks.NewLocal(logger.Log)

	err := handleMessage(msg)
	assert.Error(t, err) // should fail parsing json

	logLine := hook.LastEntry().Message
	assert.Equal(t, "Cannot unmarshal message body", logLine)
}
