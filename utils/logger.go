package utils

import (
	"encoding/json"
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

var (
    Logger *logrus.Logger // The logger for the server.
    fileWriter *rotatelogs.RotateLogs // The file writer for the server logger & access logger.
)

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

    switch strings.ToLower(cfg.SERVER_CONFIG.LOG_FORMATTER) {
        case "json":
            formatter = &QCSJSONFormatter{
                TextFormatter: &logrus.TextFormatter{
                    DisableColors:   true,
                    TimestampFormat: "2006/01/02 - 15:04:05",
                    FullTimestamp:   true,
                },
            }
        case "text":
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
        Out:       io.MultiWriter(os.Stderr, fileWriter),
        Formatter: formatter,
        Level:     logrus.InfoLevel,
    }
}

type QCSJSONFormatter struct {
    *logrus.TextFormatter
}

func (f *QCSJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    logType := "SERVER"
    time := entry.Time.Format(f.TimestampFormat)
    level := getServerLogConfig(entry.Level)
    message := entry.Message

    return []byte(
        fmt.Sprintf("{\"type\":\"%s\",\"time\":\"%s\",\"level\":\"%s\",\"message\":\"%s\"}\n",
            logType,
            time, 
            level, 
            message,
    )), nil
}

type QCSTextFormatter struct {
    *logrus.TextFormatter
}

func (f *QCSTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    level := getServerLogConfigWithColor(entry.Level)

    time := entry.Time.Format(f.TimestampFormat)
    message := entry.Message

    return []byte(
        fmt.Sprintf("[QCS] %s | %s |\t%s\n", 
            time, 
            level, 
            message,
    )), nil
}

// Get the corresponding log level string from the given log level.
func getServerLogConfig(level logrus.Level) string {
    switch level {
        case logrus.InfoLevel:
            return "INFO"
        case logrus.WarnLevel:
            return "WARN"
        case logrus.ErrorLevel:
            return "ERROR"
        case logrus.FatalLevel:
            return "FATAL"
        default:
            return "NONE"
    }
}

// Get the corresponding log level string with color from the given log level.
func getServerLogConfigWithColor(level logrus.Level) string {
    switch level {
        case logrus.InfoLevel:
            return color.New(color.FgGreen).Sprint(" INFO ")
        case logrus.WarnLevel:
            return color.New(color.FgYellow).Sprint(" WARN ")
        case logrus.ErrorLevel:
            return color.New(color.FgRed).Sprint(" ERROR")
        case logrus.FatalLevel:
            return color.New(color.FgMagenta).Sprint(" FATAL")
        default:
            return color.New(color.FgWhite).Sprint(" NONE ")
    }
}

type QCSExtractGINCtx struct {
    StatusCode int
    Latency    time.Duration
    ClientIP   string
    Method     string
    FullPath   string
}

type AccessLog struct {
    Type     string    `json:"type"`
    Time     string    `json:"time"`
    Level    string    `json:"level"`  
    Status   int       `json:"status"`
    Latency  string    `json:"latency"`
    ClientIP string    `json:"client_ip"`
    Method   string    `json:"method"`
    Path     string    `json:"path"`
}

// Overwrite the default logger of Gin Framework.
func OverwriteGinLog(ctx *QCSExtractGINCtx) {
    if strings.ToLower(cfg.SERVER_CONFIG.LOG_FORMATTER) == "text" {
        level, statusCode := getAccessLogConfig(ctx.StatusCode)

        message := fmt.Sprintf("[QCS] %s | %s |\t%s | %12s | %15s | %6s | %s\n",
            time.Now().Format("2006/01/02 - 15:04:05"),
            level,
            statusCode,
            ctx.Latency,
            ctx.ClientIP,
            color.New(color.FgHiCyan).Sprintf(" %-4s", ctx.Method),
            ctx.FullPath,
        )

        logOutput := io.MultiWriter(fileWriter, os.Stdout)
        logOutput.Write([]byte(message))
    } else {
        message := AccessLog{
            Type: "ACCESS",
            Time: time.Now().Format("2006/01/02 - 15:04:05"),
            Level: getAccessLogLevel(ctx.StatusCode),
            Status: ctx.StatusCode,
            Latency: fmt.Sprint(ctx.Latency),
            ClientIP: ctx.ClientIP,
            Method: ctx.Method,
            Path: ctx.FullPath,
        }

        jsonBytes, _ := json.Marshal(message)

        logOutput := io.MultiWriter(fileWriter, os.Stderr)
        logOutput.Write(jsonBytes)
        logOutput.Write([]byte("\n"))
    }
}

// Get the corresponding log level string from the given status code.
func getAccessLogLevel(statusCode int) string {
    switch {
        case statusCode >= 200 && statusCode < 300:
            return "INFO"
        case statusCode >= 400 && statusCode < 500:
            return "WARN"
        case statusCode >= 500:
            return"ERROR"
        default:
            return "NONE"
    }
}

// Get the corresponding log level string with color and status code with color from the given status code.
func getAccessLogConfig(statusCode int) (string, string) {
    switch {
        case statusCode >= 200 && statusCode < 300:
            return color.New(color.FgGreen).Sprint(" INFO "), color.New(color.FgHiGreen).Sprint(statusCode)
        case statusCode >= 400 && statusCode < 500:
            return color.New(color.FgYellow).Sprint(" WARN "), color.New(color.FgHiYellow).Sprint(statusCode)
        case statusCode >= 500:
            return color.New(color.FgRed).Sprint(" ERROR"), color.New(color.FgHiRed).Sprint(statusCode)
        default:
            return color.New(color.FgWhite).Sprint(" NONE "), color.New(color.FgHiWhite).Sprint(statusCode)
    }
}