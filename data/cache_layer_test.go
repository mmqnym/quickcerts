package data

import (
	"testing"

	cfg "QuickCertS/configs"

	"github.com/stretchr/testify/assert"
)

func TestConnectAndDisconnectRDB(t *testing.T) {
	backupHost := cfg.CACHE_CONFIG.HOST
	backupPort := cfg.CACHE_CONFIG.PORT

	defer func() {
		cfg.CACHE_CONFIG.HOST = backupHost
		cfg.CACHE_CONFIG.PORT = backupPort
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.CACHE_CONFIG.HOST = "localhost"
	cfg.CACHE_CONFIG.PORT = 33334

	err := ConnectRDB()
	assert.Nil(t, err)
	err = DisconnectRDB()
	assert.Nil(t, err)

	// Test invalid case
	err = DisconnectRDB()
	assert.Equal(t, "currently not connecting the redis database", err.Error())

	cfg.CACHE_CONFIG.HOST = "unknown"
	cfg.CACHE_CONFIG.PORT = 12345
	err = ConnectRDB()
	assert.Equal(t, err.Error(), "failed to access the redis database")
}

func TestSetAndGetKeyCache(t *testing.T) {
	backupHost := cfg.CACHE_CONFIG.HOST
	backupPort := cfg.CACHE_CONFIG.PORT

	defer func() {
		cfg.CACHE_CONFIG.HOST = backupHost
		cfg.CACHE_CONFIG.PORT = backupPort
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.CACHE_CONFIG.HOST = "localhost"
	cfg.CACHE_CONFIG.PORT = 33334

	err := ConnectRDB()
	assert.Nil(t, err)

	base := "testSN&testBP&testBN&testMAC"
	expectedKey := "5578c9d3cd718345af4319f3021157999b993f2e991481524234746f38b84c03"

	err = SetDeviceKeyCache(base, expectedKey)
	assert.Nil(t, err)

	actualKey, err := GetDeviceKeyCache(base)
	assert.Nil(t, err)
	assert.Equal(t, actualKey, expectedKey)

	err = DeleteTestingCache(base)
	assert.Nil(t, err)

	err = DisconnectRDB()
	assert.Nil(t, err)

	// Test invalid case
	err = SetDeviceKeyCache("test", "test")
	assert.Equal(t, "currently not connecting the redis database", err.Error())

	_, err = GetDeviceKeyCache("test")
	assert.Equal(t, "currently not connecting the redis database", err.Error())

	err = ConnectRDB()
	assert.Nil(t, err)
	actualKey, err = GetDeviceKeyCache("test")
	assert.Equal(t, "the key not exist in the cache", err.Error())
	assert.Equal(t, "", actualKey)

	err = DisconnectRDB()
	assert.Nil(t, err)
}

func TestDeleteTestingCache(t *testing.T) {
	err := DeleteTestingCache("test")
	assert.Equal(t, "currently not connecting the redis database", err.Error())
}
