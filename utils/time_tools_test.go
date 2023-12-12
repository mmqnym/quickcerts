package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeUnitStrToTimeDuration(t *testing.T) {
	// Test valid case
	duration, _ := TimeUnitStrToTimeDuration("day")
	assert.Equal(t, time.Hour * 24, duration)
	duration, _ = TimeUnitStrToTimeDuration("hour")
	assert.Equal(t, time.Hour, duration)
	duration, _ = TimeUnitStrToTimeDuration("minute")
	assert.Equal(t, time.Minute, duration)
	duration, _ = TimeUnitStrToTimeDuration("second")
	assert.Equal(t, time.Second, duration)
	duration, _ = TimeUnitStrToTimeDuration("millisecond")
	assert.Equal(t, time.Millisecond, duration)

	// Test invalid case
	_, err := TimeUnitStrToTimeDuration("invalid")
	assert.Equal(t, err.Error(), "time unit is not valid (Require: day, hour, minute, second, millisecond)")
}