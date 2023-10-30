package model

// SerialNumber: Serial number obtained from purchasing software
//
// BoardProducer: Motherboard manufacturer
//
// BoardName: Motherboard model
//
// MACAddress: Ethernet MAC address of the motherboard
//
// Note: Additional information
type ApplyCertInfo struct {
	SerialNumber	   string	`json:"serial_number" binding:"required"`
	BoardProducer      string	`json:"board_producer" binding:"required"`
	BoardName 		   string	`json:"board_name" binding:"required"`
	MACAddress 		   string	`json:"mac_address" binding:"required"`
	Note 			   string	`json:"note"`
}

type ApplyTempPermitInfo struct {
	BoardProducer      string	`json:"board_producer" binding:"required"`
	BoardName 		   string	`json:"board_name" binding:"required"`
	MACAddress 		   string	`json:"mac_address" binding:"required"`
	Note 			   string	`json:"note"`
}