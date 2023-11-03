package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Convert the time unit string to time.Duration.
func TimeUnitStrToTimeDuration(unit string) time.Duration {
	defer func() {
		if err := recover(); err != nil {
			color.Red(err.(error).Error())
			os.Exit(1)
		}
	}()

	unit = strings.ToLower(unit)
	var timeUnit time.Duration

	switch unit {
		case "day":
			timeUnit = time.Hour * 24
		case "hour":
			timeUnit = time.Hour
		case "minute":
			timeUnit = time.Minute
		case "second":
			timeUnit = time.Second
		case "millisecond":
			timeUnit = time.Millisecond
		default:
			panic(errors.New("time unit is not valid (Require: day, hour, minute, second, millisecond)"))
	}

	return timeUnit
}