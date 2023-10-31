package model

//Token: The credential information that allows calling privileged APIs
//
// SerialNumber: The new serial number to be uploaded
//
// Reason: The reason for uploading the serial number
type SNInfo struct {
	Token        string `json:"token" binding:"required"`
	SerialNumber string `json:"serial_number" binding:"required"`
	Reason       string `json:"reason" binding:"required"`
}

// Token: The credential information that allows calling privileged APIs
//
// Count: The new serial number to be uploaded
//
// Reason: The reason for uploading the serial number(s)
type SNsInfo struct {
	Token  string `json:"token" binding:"required"`
	Count  int    `json:"count" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}
