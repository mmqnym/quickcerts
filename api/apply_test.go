package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mmq88/quickcerts/model"
	"github.com/mmq88/quickcerts/utils"

	"github.com/mmq88/quickcerts/data"

	cfg "github.com/mmq88/quickcerts/configs"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestApplyCertificate(t *testing.T) {
	backupDBHost := cfg.DB_CONFIG.HOST
	backupDBPort := cfg.DB_CONFIG.PORT
	backupRDBHost := cfg.CACHE_CONFIG.HOST
	backupRDBPort := cfg.CACHE_CONFIG.PORT

	defer func() {
		cfg.DB_CONFIG.HOST = backupDBHost
		cfg.DB_CONFIG.PORT = backupDBPort
		cfg.CACHE_CONFIG.HOST = backupRDBHost
		cfg.CACHE_CONFIG.PORT = backupRDBPort
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err := data.ConnectDB()
	assert.Nil(t, err)
	cfg.CACHE_CONFIG.HOST = "localhost"
	cfg.CACHE_CONFIG.PORT = 33334
	err = data.ConnectRDB()
	assert.Nil(t, err)

	defer func() {
		err = data.DisconnectDB()
		assert.Nil(t, err)
		err = data.DisconnectRDB()
		assert.Nil(t, err)
		utils.TestBuffer = ""
	}()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/apply/cert", ApplyCertificate)

	testSN := "testSN"
	err = data.AddNewSN(testSN)
	assert.Nil(t, err)

	applyInfo := model.ApplyCertInfo{
		SerialNumber:  testSN,
		BoardProducer: "testBP",
		BoardName:     "testBN",
		MACAddress:    "testMAC",
	}
	jsonValue, _ := json.Marshal(applyInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/apply/cert", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res := w.Body.String()
	var applyCertResponse model.ApplyCertResponse
	err = json.Unmarshal([]byte(res), &applyCertResponse)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedKey := "5578c9d3cd718345af4319f3021157999b993f2e991481524234746f38b84c03"
	assert.Equal(t, expectedKey, applyCertResponse.Key)
	assert.Equal(t, fmt.Sprintf("Successfully updated and sent the key [%s].", expectedKey), utils.TestBuffer)

	// Test invalid case (Required fields are empty or not exist)
	w = httptest.NewRecorder()
	applyInfo = model.ApplyCertInfo{
		SerialNumber:  "",
		BoardProducer: "testBP",
		BoardName:     "testBN",
		MACAddress:    "testMAC",
	}
	jsonValue, _ = json.Marshal(applyInfo)
	req, _ = http.NewRequest("POST", "/api/v1/apply/cert", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	res = w.Body.String()
	var errorResponse model.ErrorResponse
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid data format.", errorResponse.Error)
	assert.Equal(t,
		"Key: 'ApplyCertInfo.SerialNumber' Error:Field validation for 'SerialNumber' failed on the 'required' tag",
		utils.TestBuffer,
	)

	// Test invalid case (The S/N does not exist)
	w = httptest.NewRecorder()
	applyInfo = model.ApplyCertInfo{
		SerialNumber:  "none",
		BoardProducer: "testBP",
		BoardName:     "testBN",
		MACAddress:    "testMAC",
	}
	jsonValue, _ = json.Marshal(applyInfo)
	req, _ = http.NewRequest("POST", "/api/v1/apply/cert", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "The S/N does not exist.", errorResponse.Error)
	assert.Equal(t, "The S/N [none] does not exist.", utils.TestBuffer)

	// Test invalid case (Use the same S/N with different device)
	w = httptest.NewRecorder()
	applyInfo = model.ApplyCertInfo{
		SerialNumber:  testSN,
		BoardProducer: "testInvalidBP",
		BoardName:     "testInvalidBN",
		MACAddress:    "testInvalidMAC",
	}
	jsonValue, _ = json.Marshal(applyInfo)
	req, _ = http.NewRequest("POST", "/api/v1/apply/cert", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "The S/N does not exist or has already been used.", errorResponse.Error)
	assert.Equal(t, "The S/N [testSN] does not exist or has already been used.", utils.TestBuffer)

	// Test invalid case (Disconnect the redis database)
	w = httptest.NewRecorder()
	err = data.DisconnectRDB()
	assert.Nil(t, err)

	applyInfo = model.ApplyCertInfo{
		SerialNumber:  testSN,
		BoardProducer: "testBP",
		BoardName:     "testBN",
		MACAddress:    "testMAC",
	}
	jsonValue, _ = json.Marshal(applyInfo)
	req, _ = http.NewRequest("POST", "/api/v1/apply/cert", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Currently not connecting the redis database.", errorResponse.Error)
	assert.Equal(t, "Currently not connecting the redis database.", utils.TestBuffer)
	err = data.ConnectRDB()
	assert.Nil(t, err)

	// Test invalid case (Disconnect the redis database)
	w = httptest.NewRecorder()
	err = data.DisconnectRDB()
	assert.Nil(t, err)

	applyInfo = model.ApplyCertInfo{
		SerialNumber:  testSN,
		BoardProducer: "testBP",
		BoardName:     "testBN",
		MACAddress:    "testMAC",
	}
	jsonValue, _ = json.Marshal(applyInfo)
	req, _ = http.NewRequest("POST", "/api/v1/apply/cert", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Currently not connecting the redis database.", errorResponse.Error)
	assert.Equal(t, "Currently not connecting the redis database.", utils.TestBuffer)
	err = data.ConnectRDB()
	assert.Nil(t, err)

	// Delete the testing data
	err = data.DeleteTestingData("DELETE FROM certs WHERE sn = $1", testSN)
	assert.Nil(t, err)
	err = data.DeleteTestingCache(testSN)
	assert.Nil(t, err)
}

func TestApplyTemporaryPermit(t *testing.T) {
	backupDBHost := cfg.DB_CONFIG.HOST
	backupDBPort := cfg.DB_CONFIG.PORT
	backupRDBHost := cfg.CACHE_CONFIG.HOST
	backupRDBPort := cfg.CACHE_CONFIG.PORT
	backupTemporaryPermitTime := cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME
	backupTemporaryPermitTimeUnit := cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT

	defer func() {
		cfg.DB_CONFIG.HOST = backupDBHost
		cfg.DB_CONFIG.PORT = backupDBPort
		cfg.CACHE_CONFIG.HOST = backupRDBHost
		cfg.CACHE_CONFIG.PORT = backupRDBPort
		cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME = backupTemporaryPermitTime
		cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT = backupTemporaryPermitTimeUnit
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err := data.ConnectDB()
	assert.Nil(t, err)
	cfg.CACHE_CONFIG.HOST = "localhost"
	cfg.CACHE_CONFIG.PORT = 33334
	err = data.ConnectRDB()
	assert.Nil(t, err)

	defer func() {
		err = data.DisconnectDB()
		assert.Nil(t, err)
		err = data.DisconnectRDB()
		assert.Nil(t, err)
		utils.TestBuffer = ""
	}()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/apply/temp-permit", ApplyTemporaryPermit)

	// Test valid case
	cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME = 1
	cfg.SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT = "second"

	applyInfo := model.ApplyTempPermitInfo{
		BoardProducer: "testBP",
		BoardName:     "testBN",
		MACAddress:    "testMAC",
	}

	jsonValue, _ := json.Marshal(applyInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/apply/temp-permit", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res := w.Body.String()
	var applyTempPermitResponse model.ApplyTempPermitResponse
	err = json.Unmarshal([]byte(res), &applyTempPermitResponse)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedKey := "94acb9791b49e5e9d92673fa4c909377973dc65f172463fd1750107450615530"
	assert.Equal(t, int64(1), applyTempPermitResponse.RemainingTime)
	assert.Equal(t, "activated", applyTempPermitResponse.Status)
	assert.Equal(t,
		fmt.Sprintf("Authorized [%s] temporary use of the product remaining [%d s].", expectedKey, 1),
		utils.TestBuffer,
	)

	// Test case (Same device use the same key after 1 second)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/apply/temp-permit", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &applyTempPermitResponse)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.GreaterOrEqual(t, int64(604800), applyTempPermitResponse.RemainingTime)
	assert.Equal(t, "activated", applyTempPermitResponse.Status)

	// Test invalid case (timeout)
	time.Sleep(1 * time.Second)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/apply/temp-permit", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res = w.Body.String()
	var errorResponse model.ErrorResponse
	err = json.Unmarshal([]byte(res), &errorResponse)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t,
		fmt.Sprintf("The authorization for [%s] to use the product has expired.", expectedKey),
		utils.TestBuffer,
	)
	assert.Equal(t, "The authorization has expired.", errorResponse.Error)

	// Test invalid case (Required fields are empty or not exist)
	w = httptest.NewRecorder()
	applyInfo = model.ApplyTempPermitInfo{
		BoardProducer: "",
		BoardName:     "testBN",
		MACAddress:    "testMAC",
	}

	jsonValue, _ = json.Marshal(applyInfo)

	req, _ = http.NewRequest("POST", "/api/v1/apply/temp-permit", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid data format.", errorResponse.Error)
	assert.Equal(t,
		"Key: 'ApplyTempPermitInfo.BoardProducer' Error:Field validation for 'BoardProducer' failed on the 'required' tag",
		utils.TestBuffer,
	)

	// Test invalid case (Disconnect the redis database)
	w = httptest.NewRecorder()
	err = data.DisconnectRDB()
	assert.Nil(t, err)

	applyInfo = model.ApplyTempPermitInfo{
		BoardProducer: "testBP",
		BoardName:     "testBN",
		MACAddress:    "testMAC",
	}

	jsonValue, _ = json.Marshal(applyInfo)
	req, _ = http.NewRequest("POST", "/api/v1/apply/temp-permit", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Currently not connecting the redis database.", errorResponse.Error)
	assert.Equal(t, "Currently not connecting the redis database.", utils.TestBuffer)
	err = data.ConnectRDB()
	assert.Nil(t, err)

	// Delete the testing data
	err = data.DeleteTestingData("DELETE FROM temporary_permits WHERE key = $1", expectedKey)
	assert.Nil(t, err)
	err = data.DeleteTestingCache(expectedKey)
	assert.Nil(t, err)
}
