package utils

import (
	"errors"
	"strings"
	"time"
)

// Convert the time unit string to time.Duration.
func TimeUnitStrToTimeDuration(unit string) (time.Duration, error) {
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
            return 0, errors.New("time unit is not valid (Require: day, hour, minute, second, millisecond)")
    }

    return timeUnit, nil
}