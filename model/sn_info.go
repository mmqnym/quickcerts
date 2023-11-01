package model

// SerialNumber: The new serial number to be uploaded
//
// Reason: The reason for uploading the serial number
type SNInfo struct {
	SerialNumber string `json:"serial_number" binding:"required"`
	Reason       string `json:"reason" binding:"required"`
}

// Count: The new serial number to be uploaded
//
// Reason: The reason for uploading the serial number(s)
type SNsInfo struct {
	Count  int    `json:"count" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}
