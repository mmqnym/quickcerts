package utils

import (
	"testing"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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
        {logrus.InfoLevel, "\x1b[32m", " INFO "}, // ANSI Green
        {logrus.WarnLevel, "\x1b[33m", " WARN "}, // ANSI Yellow
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