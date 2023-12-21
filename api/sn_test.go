package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	cfg "QuickCertS/configs"
	"QuickCertS/data"
	"QuickCertS/model"
	"QuickCertS/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateSN(t *testing.T) {
	backupDBHost := cfg.DB_CONFIG.HOST
	backupDBPort := cfg.DB_CONFIG.PORT

	defer func() {
		cfg.DB_CONFIG.HOST = backupDBHost
		cfg.DB_CONFIG.PORT = backupDBPort
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err := data.ConnectDB()
	assert.Nil(t, err)

	defer func() {
		err = data.DisconnectDB()
		assert.Nil(t, err)
		utils.TestBuffer = ""
	}()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/sn/create", CreateSN)

	// Test valid case
	// Uses docker-compose config
	testSN := "testSN"

	creationInfo := model.SNInfo{
		SerialNumber: testSN,
		Reason:       "testReason",
	}

	jsonValue, _ := json.Marshal(creationInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/sn/create", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res := w.Body.String()
	var createSNResponse model.CreateSNResponse
	err = json.Unmarshal([]byte(res), &createSNResponse)
	
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, createSNResponse.Msg, "Successfully uploaded a new S/N.")
	assert.Equal(t, createSNResponse.SerialNumber, creationInfo.SerialNumber)
	assert.Equal(t, 
		fmt.Sprintf("Successfully uploaded a new S/N [%s] with reason (%s).",
				creationInfo.SerialNumber, creationInfo.Reason), 
		utils.TestBuffer,
	)

	// Test invalid case (Required fields are empty or not exist)
	w = httptest.NewRecorder()
	creationInfo = model.SNInfo{
		SerialNumber: "",
		Reason:       "testReason",
	}

	jsonValue, _ = json.Marshal(creationInfo)
	req, _ = http.NewRequest("POST", "/api/v1/sn/create", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	res = w.Body.String()
	var errorResponse model.ErrorResponse
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid data format.", errorResponse.Error)
	assert.Equal(t, 
		"Key: 'SNInfo.SerialNumber' Error:Field validation for 'SerialNumber' failed on the 'required' tag", 
		utils.TestBuffer,
	)

	// Test invalid case (The S/N already exists)
	creationInfo = model.SNInfo{
		SerialNumber: testSN,
		Reason:       "testReason",
	}

	jsonValue, _ = json.Marshal(creationInfo)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/sn/create", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)
	
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "The S/N already exists.", errorResponse.Error)
	assert.Equal(t, 
		fmt.Sprintf("The S/N [%s] already exists.", creationInfo.SerialNumber), 
		utils.TestBuffer,
	)

	// Delete test data
	err = data.DeleteTestingData("DELETE FROM certs WHERE sn = $1", testSN)
	assert.Nil(t, err)
}

func TestCreateSNs(t *testing.T) {
	backupDBHost := cfg.DB_CONFIG.HOST
	backupDBPort := cfg.DB_CONFIG.PORT

	defer func() {
		cfg.DB_CONFIG.HOST = backupDBHost
		cfg.DB_CONFIG.PORT = backupDBPort
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err := data.ConnectDB()
	assert.Nil(t, err)

	defer func() {
		err = data.DisconnectDB()
		assert.Nil(t, err)
		utils.TestBuffer = ""
	}()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/sn/generate", GenerateSN)

	// Test valid case
	// Uses docker-compose config
	targetGenerationCount := 2

	creationInfo := model.SNsInfo{
		Count: 	targetGenerationCount,
		Reason: "testReason",
	}

	jsonValue, _ := json.Marshal(creationInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/sn/generate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res := w.Body.String()
	var generateSNResponse model.GenerateSNResponse
	err = json.Unmarshal([]byte(res), &generateSNResponse)
	
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, generateSNResponse.Msg, 
		fmt.Sprintf("Successfully uploaded new S/N (%d) with reason (%s).", creationInfo.Count, creationInfo.Reason),
	)
	assert.Equal(t, len(generateSNResponse.SerialNumbers), targetGenerationCount)

	// Test invalid case (Required fields are empty or not exist)
	w = httptest.NewRecorder()
	creationInfo = model.SNsInfo{
		// Count:  0,
		Reason: "testReason",
	}

	jsonValue, _ = json.Marshal(creationInfo)
	req, _ = http.NewRequest("POST", "/api/v1/sn/generate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	res = w.Body.String()
	var errorResponse model.ErrorResponse
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid data format.", errorResponse.Error)
	assert.Equal(t, 
		"Key: 'SNsInfo.Count' Error:Field validation for 'Count' failed on the 'required' tag", 
		utils.TestBuffer,
	)

	// Test invalid case (The S/N already exists)
	creationInfo = model.SNsInfo{
		Count:  -2,
		Reason: "testReason",
	}

	jsonValue, _ = json.Marshal(creationInfo)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/sn/generate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)
	
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "The count must be greater than 0.", errorResponse.Error)
	assert.Equal(t, 
		fmt.Sprintf("Invalid count(<=0) [%d].", -2), 
		utils.TestBuffer,
	)

	// Delete test data
	err = data.DeleteTestingData(
		"DELETE FROM certs WHERE sn IN ($1, $2)", 
		generateSNResponse.SerialNumbers[0], generateSNResponse.SerialNumbers[1],
	)
	assert.Nil(t, err)
}

func TestUpdateCertNote(t *testing.T) {
	backupDBHost := cfg.DB_CONFIG.HOST
	backupDBPort := cfg.DB_CONFIG.PORT

	defer func() {
		cfg.DB_CONFIG.HOST = backupDBHost
		cfg.DB_CONFIG.PORT = backupDBPort
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err := data.ConnectDB()
	assert.Nil(t, err)

	defer func() {
		err = data.DisconnectDB()
		assert.Nil(t, err)
		utils.TestBuffer = ""
	}()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/sn/update", UpdateCertNote)

	testSNList := []string{"testSN1", "testSN2", "testSN3"}
	err = data.AddNewSNs(testSNList)
	assert.Nil(t, err)

	// Test valid case
	// Uses docker-compose config
	updateInfo := model.CertNote{
		SerialNumber: testSNList[1],
		Note:         "testNote",
	}

	jsonValue, _ := json.Marshal(updateInfo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/sn/update", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res := w.Body.String()
	var updateCertNoteResponse model.UpdateCertNoteResponse
	
	err = json.Unmarshal([]byte(res), &updateCertNoteResponse)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	
	assert.Equal(t, "Successfully updated the note of specified S/N.", updateCertNoteResponse.Msg)
	assert.Equal(t, "testNote", updateCertNoteResponse.Note)
	assert.Equal(t,"Successfully updated the note of specified S/N.", utils.TestBuffer,)

	// Test invalid case (Required fields are empty or not exist)
	w = httptest.NewRecorder()
	updateInfo = model.CertNote{
		SerialNumber: "",
		Note:         "testNote",
	}

	jsonValue, _ = json.Marshal(updateInfo)
	req, _ = http.NewRequest("POST", "/api/v1/sn/update", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res = w.Body.String()
	var errorResponse model.ErrorResponse
	err = json.Unmarshal([]byte(res), &errorResponse)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid data format.", errorResponse.Error)
	assert.Equal(t,
		"Key: 'CertNote.SerialNumber' Error:Field validation for 'SerialNumber' failed on the 'required' tag",
		utils.TestBuffer,
	)

	// Test invalid case (The S/N does not exist)
	w = httptest.NewRecorder()
	updateInfo = model.CertNote{
		SerialNumber: "testSN4",
		Note:         "testNote",
	}

	jsonValue, _ = json.Marshal(updateInfo)
	req, _ = http.NewRequest("POST", "/api/v1/sn/update", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	res = w.Body.String()
	err = json.Unmarshal([]byte(res), &errorResponse)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, 
		fmt.Sprintf("The S/N [%s] does not exist.", updateInfo.SerialNumber), 
		errorResponse.Error,
	)
	assert.Equal(t,
		fmt.Sprintf("The S/N [%s] does not exist.", updateInfo.SerialNumber),
		utils.TestBuffer,
	)

	// Delete test data
	err = data.DeleteTestingData(
		"DELETE FROM certs WHERE sn IN ($1, $2, $3)",
		testSNList[0], testSNList[1], testSNList[2],
	)
}

func TestGetAllRecords(t *testing.T) {
	backupDBHost := cfg.DB_CONFIG.HOST
	backupDBPort := cfg.DB_CONFIG.PORT

	defer func() {
		cfg.DB_CONFIG.HOST = backupDBHost
		cfg.DB_CONFIG.PORT = backupDBPort
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err := data.ConnectDB()
	assert.Nil(t, err)

	defer func() {
		err = data.DisconnectDB()
		assert.Nil(t, err)
		utils.TestBuffer = ""
	}()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/sn/get-all", GetAllRecords)

	testSNList := []string{"testSN1", "testSN2", "testSN3"}
	err = data.AddNewSNs(testSNList)
	assert.Nil(t, err)

	// Test valid case
	// Uses docker-compose config
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sn/get-all", nil)

	router.ServeHTTP(w, req)

	res := w.Body.String()
	var getAllRecordResponse model.GetAllRecordsResponse
	
	err = json.Unmarshal([]byte(res), &getAllRecordResponse)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	
	for i, rec := range getAllRecordResponse.Data {
		assert.Equal(t, testSNList[i], rec.SerialNumber)
		assert.Equal(t, "", rec.Key)
		assert.Equal(t, "", rec.Note)
	}

	// Delete test data
	err = data.DeleteTestingData(
		"DELETE FROM certs WHERE sn IN ($1, $2, $3)",
		testSNList[0], testSNList[1], testSNList[2],
	)
}

func TestGetAvaliableSN(t *testing.T) {
	backupDBHost := cfg.DB_CONFIG.HOST
	backupDBPort := cfg.DB_CONFIG.PORT

	defer func() {
		cfg.DB_CONFIG.HOST = backupDBHost
		cfg.DB_CONFIG.PORT = backupDBPort
	}()

	// Test valid case
	// Uses docker-compose config
	cfg.DB_CONFIG.HOST = "localhost"
	cfg.DB_CONFIG.PORT = 33332
	err := data.ConnectDB()
	assert.Nil(t, err)

	defer func() {
		err = data.DisconnectDB()
		assert.Nil(t, err)
		utils.TestBuffer = ""
	}()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/v1/sn/get-available", GetAvaliableSN)

	testSNList := []string{"testSN1", "testSN2", "testSN3"}
	err = data.AddNewSNs(testSNList)
	assert.Nil(t, err)

	// Test valid case
	// Uses docker-compose config
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/sn/get-available", nil)

	router.ServeHTTP(w, req)

	res := w.Body.String()
	var getAvailableSNResponse model.GetAvaliableSNResponse
	
	err = json.Unmarshal([]byte(res), &getAvailableSNResponse)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	
	for _, sn := range getAvailableSNResponse.Data {
		assert.Contains(t, testSNList, sn)
	}

	// Delete test data
	err = data.DeleteTestingData(
		"DELETE FROM certs WHERE sn IN ($1, $2, $3)",
		testSNList[0], testSNList[1], testSNList[2],
	)

	assert.Nil(t, err)
}