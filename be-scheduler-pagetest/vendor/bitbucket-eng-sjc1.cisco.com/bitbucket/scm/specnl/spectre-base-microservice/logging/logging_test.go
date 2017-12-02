package logging

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.Error("Hello error")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "Hello error", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestLoggingInitWrongFormat(t *testing.T) {
	err := Log.LoggingInit("xyz", ioutil.Discard, "Warn")
	assert.Equal(t, "Incorrect log format provided", err.Error())
}

func TestLoggingInitWrongLevel(t *testing.T) {
	err := Log.LoggingInit("Text", ioutil.Discard, "XYZ")
	assert.Equal(t, "Incorrect log level provided", err.Error())
}

func TestLoggingInitFormatTextLevelInfo(t *testing.T) {
	Log.LoggingInit("Text", os.Stderr, "Info")
	assert.Equal(t, os.Stderr, Log.Logger.Out)
	assert.Equal(t, &logrus.TextFormatter{}, Log.Logger.Formatter)
	assert.Equal(t, logrus.InfoLevel, Log.Logger.Level)
}

func TestLoggingInitFormatTextLevelDebug(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Debug")
	assert.Equal(t, logrus.DebugLevel, Log.Logger.Level)
}

func TestLoggingInitFormatTextLevelError(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Error")
	assert.Equal(t, logrus.ErrorLevel, Log.Logger.Level)
}

func TestLoggingInitFormatTextLevelWarn(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Warn")
	assert.Equal(t, logrus.WarnLevel, Log.Logger.Level)
}

func TestLoggingInitFormatTextLevelPanic(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Panic")
	assert.Equal(t, logrus.PanicLevel, Log.Logger.Level)
}

func TestLoggingInitFormatTextLevelFatal(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Fatal")
	assert.Equal(t, logrus.FatalLevel, Log.Logger.Level)
}

func TestLoggingInitFormatJSON(t *testing.T) {
	Log.LoggingInit("Json", os.Stderr, "Info")
	assert.Equal(t, os.Stderr, Log.Logger.Out)
	assert.Equal(t, &logrus.JSONFormatter{}, Log.Logger.Formatter)
	assert.Equal(t, logrus.InfoLevel, Log.Logger.Level)
}

func TestLoggingInitFormatLogStash(t *testing.T) {
	Log.LoggingInit("LogStash", os.Stderr, "Info")
	assert.Equal(t, os.Stderr, Log.Logger.Out)
	assert.Equal(t, &logstash.LogstashFormatter{}, Log.Logger.Formatter)
	assert.Equal(t, logrus.InfoLevel, Log.Logger.Level)
}

func TestGenerateLogLevelDebug(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Debug")
	entry, _ := Log.GenerateLog("Debug", "test_microservice", "1234 2345 Test Msg")
	assert.Equal(t, logrus.DebugLevel, entry.Level)
	assert.Equal(t, "1234 2345 Test Msg", entry.Message)
}

func TestGenerateLogLevelError(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Error")
	entry, _ := Log.GenerateLog("Error", "test_microservice", "1234 2345 Test Msg")
	assert.Equal(t, logrus.ErrorLevel, entry.Level)
	assert.Equal(t, "1234 2345 Test Msg", entry.Message)
}

func TestGenerateLogLevelInfo(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Info")
	entry, _ := Log.GenerateLog("Info", "test_microservice", "1234 2345 Test Msg")
	assert.Equal(t, logrus.InfoLevel, entry.Level)
	assert.Equal(t, "1234 2345 Test Msg", entry.Message)
}

func TestGenerateLogLevelWarning(t *testing.T) {
	Log.LoggingInit("Text", ioutil.Discard, "Warning")
	entry, _ := Log.GenerateLog("Warning", "test_microservice", "1234 2345 Test Msg")
	assert.Equal(t, logrus.WarnLevel, entry.Level)
	assert.Equal(t, "1234 2345 Test Msg", entry.Message)
}

func TestGenerateLogWithEmptyLevel(t *testing.T) {
	_, err := Log.GenerateLog("", "test_microservice", "1234 2345 Test Msg")
	assert.Equal(t, "Not enough input provided", err.Error())
}

func TestGenerateLogWithEmptyLog(t *testing.T) {
	_, err := Log.GenerateLog("Error", "test_microservice", nil)
	assert.Equal(t, "Not enough input provided", err.Error())
}

func TestGenerateLogWithEmptyMicroserviceName(t *testing.T) {
	_, err := Log.GenerateLog("Error", "", nil)
	assert.Equal(t, "Not enough input provided", err.Error())
}

func TestGenerateLogWithWrongLevel(t *testing.T) {
	_, err := Log.GenerateLog("Xyz", "test_microservice", "1234 2345 Test Msg")
	assert.Equal(t, "Incorrect Log Level provided", err.Error())
}
