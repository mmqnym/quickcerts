package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	cfg "QuickCertS/configs"

	"github.com/fatih/color"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger // The logger for the server.
var fileWriter *rotatelogs.RotateLogs // The file writer for the server logger & access logger.

func init() {
	path := "./logs/qcs-%Y-%m-%d@%H_%M_%S"

	timeUnit := TimeUnitStrToTimeDuration(cfg.SERVER_CONFIG.LOG_TIME_UNIT)
	var err error

	fileWriter, err = rotatelogs.New(
		path,
		rotatelogs.WithLinkName("./logs/qcs-latest"),
		rotatelogs.WithMaxAge(time.Duration(cfg.SERVER_CONFIG.LOG_MAX_AGE) * timeUnit),
		rotatelogs.WithRotationTime(time.Duration(cfg.SERVER_CONFIG.LOG_ROTATION_TIME) * timeUnit),
	)

	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	var formatter logrus.Formatter

	switch strings.ToUpper(cfg.SERVER_CONFIG.LOG_FORMATTER) {
		case "JSON":
			formatter = &logrus.JSONFormatter{}
		case "TEXT":
			fallthrough
		default:
			formatter = &QCSTextFormatter{
				TextFormatter: &logrus.TextFormatter{
					DisableColors:   false,
					TimestampFormat: "2006/01/02 - 15:04:05",
					FullTimestamp:   true,
				},
			}
	}

	Logger = &logrus.Logger{
		Out:       io.MultiWriter(fileWriter, os.Stdout),
		Formatter: formatter,
		Level:     logrus.InfoLevel,
	}
}

type QCSTextFormatter struct {
    *logrus.TextFormatter
}

func (f *QCSTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    level := getServerLogConfig(entry.Level)

	time := entry.Time.Format(f.TimestampFormat)
    message := entry.Message

	return []byte(
		fmt.Sprintf("[QCS] %s | %s |\t%s\n", 
			time, 
			level, 
			message,
	)), nil
}

func getServerLogConfig(level logrus.Level) string {
    switch level {
		case logrus.InfoLevel:
			return color.New(color.FgGreen).Sprint("INFO ")
		case logrus.WarnLevel:
			return color.New(color.FgYellow).Sprint("WARN ")
		case logrus.ErrorLevel:
			return color.New(color.FgRed).Sprint("ERROR")
		case logrus.FatalLevel:
			return color.New(color.FgMagenta).Sprint("FATAL")
		default:
			return color.New(color.FgWhite).Sprint("NONE ")
    }
}

type QCSExtractGINCtx struct {
    StatusCode int
	Latency time.Duration
	ClientIP string
	Method string
	FullPath string
}

func OverwriteGinLog(ctx *QCSExtractGINCtx) {
	level, statusCode := getAccessLogConfig(ctx.StatusCode)

	msg := fmt.Sprintf("[QCS] %s | %s |\t%s | %12s | %15s | %6s | %s\n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		level,
		statusCode,
		ctx.Latency,
		ctx.ClientIP,
		color.New(color.FgHiCyan).Sprint(ctx.Method),
		ctx.FullPath,
	)

	logOutput := io.MultiWriter(fileWriter, os.Stdout)
	logOutput.Write([]byte(msg))
}

func getAccessLogConfig(statusCode int) (string, string) {
    switch {
		case statusCode >= 200 && statusCode < 300:
			return color.New(color.FgGreen).Sprint("INFO "), color.New(color.FgHiGreen).Sprint(statusCode)
		case statusCode >= 400 && statusCode < 500:
			return color.New(color.FgYellow).Sprint("WARN "), color.New(color.FgHiYellow).Sprint(statusCode)
		case statusCode >= 500:
			return color.New(color.FgRed).Sprint("ERROR"), color.New(color.FgHiRed).Sprint(statusCode)
		default:
			return color.New(color.FgWhite).Sprint("NONE "), color.New(color.FgHiWhite).Sprint(statusCode)
    }
}