package goqcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type QCSAdmin struct {
	accessPrefix string
	accessToken string
	runtimeCode string
}

func NewQCSAdmin(host string, port int, apiPath string, tls bool, accessToken string, runtimeCode string) *QCSAdmin {
	accessPrefix := fmt.Sprintf("%s:%d%s", host, port, apiPath)

	if tls {
		accessPrefix = "https://" + accessPrefix
	} else {
		accessPrefix = "http://" + accessPrefix
	}

	return &QCSAdmin {
		accessPrefix: accessPrefix,
		accessToken: accessToken,
		runtimeCode: runtimeCode,
	}
}

func (qcsA *QCSAdmin) CreateSN(sn string, reason string) (*QCSCreateSNResponse, error) {
	url := qcsA.accessPrefix + "/sn/create"

	body := map[string]string {
		"serial_number": sn,
		"reason": reason,
	}

	jsonfiedBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(jsonfiedBody)))
	
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", qcsA.accessToken)
	req.Header.Add("X-Runtime-Code", qcsA.runtimeCode)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		errorMsg := data["error"].(string)
		return nil, fmt.Errorf("QCS::Error:%s", errorMsg)
	}

	var response QCSCreateSNResponse
	response.Msg, _ = data["msg"].(string)
	response.SerialNumber, _ = data["serial_number"].(string)
	
	return &response, nil
}

func (qcsA *QCSAdmin) GenerateSN(count uint, reason string) (*QCSGnerateSNResponse, error) {
	url := qcsA.accessPrefix + "/sn/generate"

	body := map[string]interface{} {
		"count": count,
		"reason": reason,
	}

	jsonfiedBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(jsonfiedBody)))
	
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", qcsA.accessToken)
	req.Header.Add("X-Runtime-Code", qcsA.runtimeCode)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		errorMsg := data["error"].(string)
		return nil, fmt.Errorf("QCS::Error:%s", errorMsg)
	}

	var response QCSGnerateSNResponse
	response.Msg, _ = data["msg"].(string)
	
	for _, sn := range data["serial_numbers"].([]interface{}) {
		response.SerialNumbers = append(response.SerialNumbers, sn.(string))
	}
	
	return &response, nil
}

func (qcsA *QCSAdmin) GetAllRecord() (*QCSAllRecordResponse, error) {
	url := qcsA.accessPrefix + "/sn/get-all"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", qcsA.accessToken)
	req.Header.Add("X-Runtime-Code", qcsA.runtimeCode)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		errorMsg := data["error"].(string)
		return nil, fmt.Errorf("QCS::Error:%s", errorMsg)
	}

	var response QCSAllRecordResponse
	var records []QCSRecord
	
	for _, irecord := range data["data"].([]interface{}) {
		var record QCSRecord
		record.SerialNumber = irecord.(map[string]interface{})["sn"].(string)
		record.Key = irecord.(map[string]interface{})["key"].(string)
		record.Note = irecord.(map[string]interface{})["note"].(string)
		records = append(records, record)
	}

	response.Data = records
	
	return &response, nil
}

func (qcsA *QCSAdmin) GetAvaliableSN() (*QCSAvaliableSNResponse, error) {
	url := qcsA.accessPrefix + "/sn/get-available"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", qcsA.accessToken)
	req.Header.Add("X-Runtime-Code", qcsA.runtimeCode)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		errorMsg := data["error"].(string)
		return nil, fmt.Errorf("QCS::Error:%s", errorMsg)
	}

	var response QCSAvaliableSNResponse
	
	for _, sn := range data["data"].([]interface{}) {
		response.Data = append(response.Data, sn.(string))
	}
	
	return &response, nil
}

func (qcsA *QCSAdmin) UpdateSNNote(sn string, note string) (*QCSUpdateSNNoteResponse, error) {
	url := qcsA.accessPrefix + "/sn/update"

	body := map[string]string{
		"serial_number": sn,
		"note": note,
	}

	jsonfiedBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(jsonfiedBody)))
	
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", qcsA.accessToken)
	req.Header.Add("X-Runtime-Code", qcsA.runtimeCode)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		errorMsg := data["error"].(string)
		return nil, fmt.Errorf("QCS::Error:%s", errorMsg)
	}

	var response QCSUpdateSNNoteResponse
	response.Msg, _ = data["msg"].(string)
	response.Note, _ = data["note"].(string)

	return &response, nil
}

type QCSClient struct {	
	accessPrefix string
	accessToken string
}

func NewQCSClient(host string, port int, apiPath string, tls bool, accessToken string) *QCSClient {
	accessPrefix := fmt.Sprintf("%s:%d%s", host, port, apiPath)

	if tls {
		accessPrefix = "https://" + accessPrefix
	} else {
		accessPrefix = "http://" + accessPrefix
	}

	return &QCSClient {
		accessPrefix: accessPrefix,
		accessToken: accessToken,
	}
}

func (qcsC *QCSClient) ApplyCert(
	sn string, 
	board_producer string, 
	board_name string, 
	mac_address string,
	) (*QCSApplyCertResponse, error) {

	url := qcsC.accessPrefix + "/apply/cert"

	body := map[string]string{
		"serial_number": sn,
		"board_producer": board_producer,
		"board_name": board_name,
		"mac_address": mac_address,
	}

	jsonfiedBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(jsonfiedBody)))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", qcsC.accessToken)
	
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var data map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		errorMsg := data["error"].(string)
		return nil, fmt.Errorf("QCS::Error:%s", errorMsg)
	}

	var response QCSApplyCertResponse
	response.Key, _ = data["key"].(string)
	response.Signature, _ = data["signature"].(string)

	return &response, nil
}

func (qcsC *QCSClient) ApplyTempPermit(
	board_producer string, 
	board_name string, 
	mac_address string,
	) (*QCSApplyTempPermitResponse, error) {
	
	url := qcsC.accessPrefix + "/apply/temp-permit"

	body := map[string]string{
		"board_producer": board_producer,
		"board_name": board_name,
		"mac_address": mac_address,
	}

	jsonfiedBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(jsonfiedBody)))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Access-Token", qcsC.accessToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var data map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		errorMsg := data["error"].(string)
		return nil, fmt.Errorf("QCS::Error:%s", errorMsg)
	}

	var response QCSApplyTempPermitResponse
	response.RemainingTime, _ = data["remaining_time"].(float64)
	response.Status, _ = data["status"].(string)

	return &response, nil
}