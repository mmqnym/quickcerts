package model

// SerialNumber: The new serial number to be uploaded
//
// Reason: The reason for uploading the serial number
type SNInfo struct {
	SerialNumber string `json:"serial_number" binding:"required" example:"779f-4e90-aebd-4295-881a-f8d7"`
	Reason       string `json:"reason" example:"For testing."`
}

// Count: The new serial number to be uploaded
//
// Reason: The reason for uploading the serial number(s)
type SNsInfo struct {
	Count  int    `json:"count" binding:"required" example:"1"`
	Reason string `json:"reason" example:"For testing."`
}

// SerialNumber: Serial number obtained from purchasing software
//
// Note: Additional information
type CertNote struct {
	SerialNumber string `json:"serial_number" binding:"required" example:"779f-4e90-aebd-4295-881a-f8d7"`
	Note         string `json:"note" binding:"required" example:"Additional information"`
}
