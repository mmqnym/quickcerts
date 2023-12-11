package configs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTimeCodeLength(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()
	
	// Test valid case
	backup_use_runtime_code := SERVER_CONFIG.USE_RUNTIME_CODE
	backup_runtime_code_length := SERVER_CONFIG.RUNTIME_CODE_LENGTH
	SERVER_CONFIG.USE_RUNTIME_CODE = true
	SERVER_CONFIG.RUNTIME_CODE_LENGTH = 6
	checkRunTimeCodeLength()
	assert.GreaterOrEqual(t, SERVER_CONFIG.RUNTIME_CODE_LENGTH, 6, "RUNTIME_CODE_LENGTH should be bigger or equal to 6")
	
	// Test invalid case
	SERVER_CONFIG.RUNTIME_CODE_LENGTH = 5
	checkRunTimeCodeLength()
	assert.GreaterOrEqual(t, SERVER_CONFIG.RUNTIME_CODE_LENGTH, 6, "RUNTIME_CODE_LENGTH should be bigger or equal to 6")

	SERVER_CONFIG.RUNTIME_CODE_LENGTH = backup_runtime_code_length
	SERVER_CONFIG.USE_RUNTIME_CODE = backup_use_runtime_code
}

func TestKeepAliveTimeout(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()
	
	// Test valid case
	checkKeepAliveTimeout()
	assert.GreaterOrEqual(t, int(SERVER_CONFIG.KEEP_ALIVE_TIMEOUT), 0, "KEEP_ALIVE_TIMEOUT should be bigger or equal to 0")
	
	// Test invalid case
	backup_keep_alive_timeout := SERVER_CONFIG.KEEP_ALIVE_TIMEOUT
	SERVER_CONFIG.KEEP_ALIVE_TIMEOUT = -1
	checkKeepAliveTimeout()
	assert.GreaterOrEqual(t, int(SERVER_CONFIG.KEEP_ALIVE_TIMEOUT), 0, "KEEP_ALIVE_TIMEOUT should be bigger or equal to 0")

	SERVER_CONFIG.KEEP_ALIVE_TIMEOUT = backup_keep_alive_timeout
}

func TestCheckKeepAliveTimeoutUnit(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()

	// Test valid case
	checkKeepAliveTimeoutUnit()
	timeUnits := []string{"hour", "minute", "second", "millisecond"}
	assert.Contains(t, timeUnits, SERVER_CONFIG.KEEP_ALIVE_TIMEOUT_UNIT, 
		"KEEP_ALIVE_TIMEOUT_UNIT should be one of hour, minute, second, millisecond",
	)

	// Test invalid case
	backup_keep_alive_timeout_unit := SERVER_CONFIG.KEEP_ALIVE_TIMEOUT_UNIT
	SERVER_CONFIG.KEEP_ALIVE_TIMEOUT_UNIT = "invalid"
	checkKeepAliveTimeoutUnit()
	assert.Contains(t, timeUnits, SERVER_CONFIG.KEEP_ALIVE_TIMEOUT_UNIT,
		"KEEP_ALIVE_TIMEOUT_UNIT should be one of hour, minute, second, millisecond",
	)

	SERVER_CONFIG.KEEP_ALIVE_TIMEOUT_UNIT = backup_keep_alive_timeout_unit
}

func TestCheckTemporaryPermitTime(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()

	// Test valid case
	checkTemporaryPermitTime()
	assert.Greater(t, SERVER_CONFIG.TEMPORARY_PERMIT_TIME, 0, "TEMPORARY_PERMIT_TIME should be bigger than 0")

	// Test invalid case
	backup_temporary_permit_time := SERVER_CONFIG.TEMPORARY_PERMIT_TIME
	SERVER_CONFIG.TEMPORARY_PERMIT_TIME = 0
	checkTemporaryPermitTime()
	assert.Greater(t, SERVER_CONFIG.TEMPORARY_PERMIT_TIME, 0, "TEMPORARY_PERMIT_TIME should be bigger than 0")

	SERVER_CONFIG.TEMPORARY_PERMIT_TIME = backup_temporary_permit_time
}

func TestCheckTemporaryPermitTimeUnit(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()

	// Test valid case
	checkTemporaryPermitTimeUnit()
	timeUnits := []string{"day", "hour", "minute"}
	assert.Contains(t, timeUnits, SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT,
		"TEMPORARY_PERMIT_TIME_UNIT should be one of day, hour, minute",
	)

	// Test invalid case
	backup_temporary_permit_time_unit := SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT
	SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT = "invalid"
	checkTemporaryPermitTimeUnit()
	assert.Contains(t, timeUnits, SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT,
		"TEMPORARY_PERMIT_TIME_UNIT should be one of day, hour, minute",
	)

	SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT = backup_temporary_permit_time_unit
}

func TestCheckLogMaxAge(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()

	// Test valid case
	checkLogMaxAge()
	assert.Greater(t, SERVER_CONFIG.LOG_MAX_AGE, 0, "LOG_MAX_AGE should be bigger than 0")

	// Test invalid case
	backup_log_max_age := SERVER_CONFIG.LOG_MAX_AGE
	SERVER_CONFIG.LOG_MAX_AGE = 0
	checkLogMaxAge()
	assert.Greater(t, SERVER_CONFIG.LOG_MAX_AGE, 0, "LOG_MAX_AGE should be bigger than 0")

	SERVER_CONFIG.LOG_MAX_AGE = backup_log_max_age
}

func TestCheckLogRotationTime(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()

	// Test valid case
	checkLogRotationTime()
	assert.Greater(t, SERVER_CONFIG.LOG_ROTATION_TIME, 0, "LOG_ROTATION_TIME should be bigger than 0")

	// Test invalid case
	backup_log_rotation_time := SERVER_CONFIG.LOG_ROTATION_TIME
	SERVER_CONFIG.LOG_ROTATION_TIME = 0
	checkLogRotationTime()
	assert.Greater(t, SERVER_CONFIG.LOG_ROTATION_TIME, 0, "LOG_ROTATION_TIME should be bigger than 0")

	SERVER_CONFIG.LOG_ROTATION_TIME = backup_log_rotation_time
}

func TestCheckLogTimeUnit(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()

	// Test valid case
	checkLogTimeUnit()
	timeUnits := []string{"day", "hour", "minute", "second"}
	assert.Contains(t, timeUnits, SERVER_CONFIG.LOG_TIME_UNIT,
		"LOG_TIME_UNIT should be one of day, hour, minute, second",
	)

	// Test invalid case
	backup_log_time_unit := SERVER_CONFIG.LOG_TIME_UNIT
	SERVER_CONFIG.LOG_TIME_UNIT = "invalid"
	checkLogTimeUnit()
	assert.Contains(t, timeUnits, SERVER_CONFIG.LOG_TIME_UNIT,
		"LOG_TIME_UNIT should be one of day, hour, minute, second",
	)

	SERVER_CONFIG.LOG_TIME_UNIT = backup_log_time_unit
}

func TestCheckCacheExpiration(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()

	// Test valid case
	checkCacheExpiration()
	assert.Greater(t, CACHE_CONFIG.EXPIRATION, 0, "EXPIRATION should be bigger than 0")

	// Test invalid case
	backup_expiration := CACHE_CONFIG.EXPIRATION
	CACHE_CONFIG.EXPIRATION = 0
	checkCacheExpiration()
	assert.Greater(t, CACHE_CONFIG.EXPIRATION, 0, "EXPIRATION should be bigger than 0")

	CACHE_CONFIG.EXPIRATION = backup_expiration
}

func TestCheckCacheExpirationUnit(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in function:", r)
		}
	}()

	// Test valid case
	checkCacheExpirationUnit()
	timeUnits := []string{"day", "hour", "minute", "second"}
	assert.Contains(t, timeUnits, CACHE_CONFIG.EXPIRATION_UNIT,
		"EXPIRATION_UNIT should be one of day, hour, minute, second",
	)

	// Test invalid case
	backup_expiration_unit := CACHE_CONFIG.EXPIRATION_UNIT
	CACHE_CONFIG.EXPIRATION_UNIT = "invalid"
	checkCacheExpirationUnit()
	assert.Contains(t, timeUnits, CACHE_CONFIG.EXPIRATION_UNIT,
		"EXPIRATION_UNIT should be one of day, hour, minute, second",
	)

	CACHE_CONFIG.EXPIRATION_UNIT = backup_expiration_unit
}