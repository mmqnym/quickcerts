package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestQCSJSONFormatter(t *testing.T) {
	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Time:    time.Now(),
		Level:   logrus.InfoLevel,
		Message: "test message",
	}

	formatter := QCSJSONFormatter{&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05",
	}}

	logBytes, err := formatter.Format(entry)
	assert.NoError(t, err)
	fmt.Println(string(logBytes))
}

func TestQCSTextFormatter(t *testing.T) {
	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Time:    time.Now(),
		Level:   logrus.WarnLevel,
		Message: "test warning",
	}

	formatter := QCSTextFormatter{&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05",
	}}

	logBytes, err := formatter.Format(entry)
	assert.NoError(t, err)
	fmt.Println(string(logBytes))
}

func TestGetServerLogConfig(t *testing.T) {
	// Test valid case
	assert.Equal(t, "INFO", getServerLogConfig(logrus.InfoLevel))
	assert.Equal(t, "WARN", getServerLogConfig(logrus.WarnLevel))
	assert.Equal(t, "ERROR", getServerLogConfig(logrus.ErrorLevel))
	assert.Equal(t, "FATAL", getServerLogConfig(logrus.FatalLevel))
	assert.Equal(t, "NONE", getServerLogConfig(logrus.DebugLevel))
}

func TestGetServerLogConfigWithColor(t *testing.T) {
	color.NoColor = false
	tests := []struct {
		level     logrus.Level
		colorCode string
		text      string
	}{
		{logrus.InfoLevel, "\x1b[32m", " INFO "},  // ANSI Green
		{logrus.WarnLevel, "\x1b[33m", " WARN "},  // ANSI Yellow
		{logrus.ErrorLevel, "\x1b[31m", " ERROR"}, // ANSI Red
		{logrus.FatalLevel, "\x1b[35m", " FATAL"}, // ANSI Magenta
		{logrus.PanicLevel, "\x1b[37m", " NONE "}, // ANSI White
	}

	for _, test := range tests {
		output := getServerLogConfigWithColor(test.level)
		expected := test.colorCode + test.text + "\x1b[0m" // ANSI Reset

		if output != expected {
			t.Errorf("For level %v: expected %v, got %v", test.level, expected, output)
		}
	}
}

func TestGetAccessLogLevel(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   string
	}{
		{200, "INFO"},
		{204, "INFO"},
		{301, "NONE"},
		{400, "WARN"},
		{500, "ERROR"},
	}

	for _, test := range tests {
		output := getAccessLogLevel(test.statusCode)
		expected := test.expected

		assert.Equal(t, expected, output)
	}
}

func TestGetAccessLogConfig(t *testing.T) {
	color.NoColor = false

	tests := []struct {
		statusCode         int
		expectedLevel      string
		expectedStatusCode string
	}{
		{200, "\x1b[32m INFO \x1b[0m", "\x1b[92m200\x1b[0m"}, // ANSI Green & High Green
		{400, "\x1b[33m WARN \x1b[0m", "\x1b[93m400\x1b[0m"}, // ANSI Yellow & High Yellow
		{500, "\x1b[31m ERROR\x1b[0m", "\x1b[91m500\x1b[0m"}, // ANSI Red & High Red
		{100, "\x1b[37m NONE \x1b[0m", "\x1b[97m100\x1b[0m"}, // ANSI White & High White
	}

	for _, test := range tests {
		gotLevel, gotStatusCode := getAccessLogConfig(test.statusCode)
		assert.Equal(t, test.expectedLevel, gotLevel)
		assert.Equal(t, test.expectedStatusCode, gotStatusCode)
	}
}
