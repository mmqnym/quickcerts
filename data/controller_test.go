package data

import (
	"testing"
	"time"

	cfg "QuickCertS/configs"
	"QuickCertS/utils"

	"github.com/stretchr/testify/assert"
)

func TestConnectAndDisconnect(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()
	
	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332

	err := ConnectDB()
	assert.Nil(t, err)
	err = DisconnectDB()
	assert.Nil(t, err)

	// Test invalid case
	err = DisconnectDB()
	assert.Equal(t, "currently not connecting the database", err.Error())

	cfg.DB_CONFIG.HOST = "unknown"
	cfg.DB_CONFIG.PORT = 12345
	err = ConnectDB()
	assert.Equal(t, err.Error(), "failed to access the database")
}

func TestAddNewSN(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	err := AddNewSN("XXXX-XXXX-XXXX-XXXX-XXXX-XXXX")
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332

	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	sn := "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
	err = AddNewSN(sn)
	assert.Nil(t, err)

	// Test invalid case
	err = AddNewSN(sn)
	assert.Equal(t, err.Error(), "the s/n already exists")

	// Delete the added test data
	_, err = db.Exec("DELETE FROM certs WHERE sn = $1", sn)
	assert.Nil(t, err)
}

func TestAddNewSNs(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	err := AddNewSNs([]string{"XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"})
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332

	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	snList := []string{
		"XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
		"YYYY-YYYY-YYYY-YYYY-YYYY-YYYY",
		"ZZZZ-ZZZZ-ZZZZ-ZZZZ-ZZZZ-ZZZZ",
	}
	err = AddNewSNs(snList)
	assert.Nil(t, err)

	// Test invalid case
	err = AddNewSNs(snList)
	assert.Equal(t, err.Error(), "some s/ns already exist")

	// Delete the added test data
	_, err = db.Exec("DELETE FROM certs WHERE sn IN ($1, $2, $3)", snList[0], snList[1], snList[2])
	assert.Nil(t, err)
}

func TestIsSNExist(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	_, err := IsSNExist("XXXX-XXXX-XXXX-XXXX-XXXX-XXXX")
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332

	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	sn := "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
	err = AddNewSN(sn)
	assert.Nil(t, err)
	
	_, err = IsSNExist(sn)
	assert.Nil(t, err)

	// Test invalid case
	_, err = IsSNExist("YYYY-YYYY-YYYY-YYYY-YYYY-YYYY")
	assert.Equal(t, err.Error(), "the s/n does not exist")

	// Delete the added test data
	_, err = db.Exec("DELETE FROM certs WHERE sn = $1", sn)
	assert.Nil(t, err)
}

func TestBindSNWithKey(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	err := BindSNWithKey("XXXX-XXXX-XXXX-XXXX-XXXX-XXXX", "key")
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332

	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	sn := "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
	key := "valid key"
	err = AddNewSN(sn)
	assert.Nil(t, err)

	err = BindSNWithKey(sn, key)
	assert.Nil(t, err)

	// Assign the same key again should be ok
	err = BindSNWithKey(sn, key)
	assert.Nil(t, err)

	// Test invalid case
	err = BindSNWithKey(sn, "invalid key")
	assert.Equal(t, err.Error(), "the s/n does not exist or has already been used")

	err = BindSNWithKey("invalid sn", "invalid key")
	assert.Equal(t, err.Error(), "the s/n does not exist or has already been used")

	// Delete the added test data
	_, err = db.Exec("DELETE FROM certs WHERE sn = $1", sn)
	assert.Nil(t, err)
}

func TestAddTemporaryPermit(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	_, err := AddTemporaryPermit("key")
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332

	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	remainingTime, err := AddTemporaryPermit("key")

	timeUnit, _ := utils.TimeUnitStrToTimeDuration(cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT)
	expectedRemainingTime := time.Duration(cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME) * timeUnit / time.Second
	assert.Nil(t, err)
	assert.Equal(t, int64(expectedRemainingTime), remainingTime)

	// Delete the added test data
	_, err = db.Exec("DELETE FROM temporary_permits WHERE key = $1", "key")
	assert.Nil(t, err)
}

func TestGetTemporaryPermitExpiredTime(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	_, err := AddTemporaryPermit("key")
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332

	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	_, err = AddTemporaryPermit("key")
	assert.Nil(t, err)
	remainingTime, err := GetTemporaryPermitExpiredTime("key")
	assert.Nil(t, err)

	timeUnit, _ := utils.TimeUnitStrToTimeDuration(cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT)
	NewRemainingTime := time.Duration(cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME) * timeUnit / time.Second

	assert.LessOrEqual(t, remainingTime, NewRemainingTime)

	// Delete the added test data
	_, err = db.Exec("DELETE FROM temporary_permits WHERE key = $1", "key")
	assert.Nil(t, err)
}

func TestUpdateCertNote(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	err := UpdateCertNote("XXXX-XXXX-XXXX-XXXX-XXXX-XXXX", "note")
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332

	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	sn := "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
	note := "note"
	err = AddNewSN(sn)
	assert.Nil(t, err)

	err = UpdateCertNote(sn, note)
	assert.Nil(t, err)

	// Test invalid case
	err = UpdateCertNote("invalid sn", note)
	assert.Equal(t, err.Error(), "the s/n does not exist")

	// Delete the added test data
	_, err = db.Exec("DELETE FROM certs WHERE sn = $1", sn)
	assert.Nil(t, err)
}

func TestGetAvaliableSN(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	_, err := GetAvaliableSN()
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	resList, err := GetAvaliableSN()
	assert.Nil(t, err)

	snList := []string{
		"XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
		"1234-XXXX-XXXX-XXXX-XXXX-XXXX",
		"1234-1234-XXXX-XXXX-XXXX-XXXX",
	}

	assert.NotContains(t, resList, snList[0])
	assert.NotContains(t, resList, snList[1])
	assert.NotContains(t, resList, snList[2])

	err = AddNewSNs(snList)
	assert.Nil(t, err)

	resList, err = GetAvaliableSN()
	assert.Nil(t, err)

	assert.Contains(t, resList, snList[0])
	assert.Contains(t, resList, snList[1])
	assert.Contains(t, resList, snList[2])

	// Delete the added test data
	_, err = db.Exec("DELETE FROM certs WHERE sn IN ($1, $2, $3)", snList[0], snList[1], snList[2])
}

func TestGetAllCerts(t *testing.T) {
	backupHost := cfg.DB_CONFIG.HOST
	backupPort := cfg.DB_CONFIG.PORT
	defer func() {
		cfg.DB_CONFIG.HOST = backupHost
		cfg.DB_CONFIG.PORT = backupPort
	}()

	// Test invalid case
	_, err := GetAllCerts()
	assert.Equal(t, "currently not connecting the database", err.Error())

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err = ConnectDB()
	assert.Nil(t, err)
	defer func() {
		err = DisconnectDB()
		assert.Nil(t, err)
	}()

	resList, err := GetAllCerts()
	assert.Nil(t, err)

	allCertsLength := len(resList)
	sn := "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
	err = AddNewSN(sn)
	assert.Nil(t, err)

	resList, err = GetAllCerts()
	assert.Nil(t, err)
	assert.Equal(t, allCertsLength + 1, len(resList))

	for _, cert := range resList {
		if cert.SerialNumber == sn {
			assert.Equal(t, cert.Note, "")
			assert.Equal(t, cert.Key, "")
			break
		}
	}

	// Delete the added test data
	_, err = db.Exec("DELETE FROM certs WHERE sn = $1", sn)
}