package utils

import (
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func TimeUnitStrToTimeDuration(unit string) time.Duration {
	defer func() {
		if err := recover(); err != nil {
			color.Red(err.(string))
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
			panic("Time unit is not valid (Require: day, hour, minute, second, millisecond).")
	}

	return timeUnit
}