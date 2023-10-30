package model

//Token: The credential information that allows calling privileged APIs
//
// SerialNumber: The new serial number to be uploaded
//
// Reason: The reason for uploading the serial number
type SNInfo struct {
	Token        string `json:"token"`
	SerialNumber string `json:"serial_number"`
	Reason       string `json:"reason"`
}

// Token: The credential information that allows calling privileged APIs
//
// Count: The new serial number to be uploaded
//
// Reason: The reason for uploading the serial number(s)
type SNsInfo struct {
	Token  string `json:"token"`
	Count  int    `json:"count"`
	Reason string `json:"reason"`
}
