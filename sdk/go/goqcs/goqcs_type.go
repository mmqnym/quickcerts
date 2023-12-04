package goqcs

type QCSCreateSNResponse struct {
	Msg          string `json:"msg"`
	SerialNumber string `json:"serial_number"`
}

type QCSGnerateSNResponse struct {
	Msg           string   `json:"msg"`
	SerialNumbers []string `json:"serial_numbers"`
}

type QCSRecord struct {
	SerialNumber string `json:"sn"`
	Key          string `json:"key"`
	Note         string `json:"note"`
}

type QCSAllRecordResponse struct {
	Data []QCSRecord `json:"data"`
}

type QCSAvaliableSNResponse struct {
	Data []string `json:"data"`
}

type QCSUpdateSNNoteResponse struct {
	Msg  string `json:"msg"`
	Note string `json:"note"`
}

type QCSApplyCertResponse struct {
	Key       string `json:"key"`
	Signature string `json:"signature"`
}

type QCSApplyTempPermitResponse struct {
	RemainingTime float64 `json:"remaining_time"`
	Status        string  `json:"status"`
}
