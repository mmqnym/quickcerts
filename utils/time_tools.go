package utils

import (
	"strings"
	"time"
)

func TimeUnitStrToTimeDuration(unit string) time.Duration {
	unit = strings.ToUpper(unit)
	var timeUnit time.Duration

	switch unit {
		case "DAY":
			timeUnit = time.Hour * 24
		case "HOUR":
			timeUnit = time.Hour
		case "MINUTE":
			timeUnit = time.Minute
		case "SECOND":
			timeUnit = time.Second
		case "MILLISECOND":
			timeUnit = time.Millisecond
		default:
			panic("Time unit is not valid (Require: day, hour, minute, second, millisecond).")
	}

	return timeUnit
}