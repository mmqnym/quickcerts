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
	SerialNumber  string `json:"serial_number" binding:"required" example:"779f-4e90-aebd-4295-881a-f8d7"`
	BoardProducer string `json:"board_producer" binding:"required" example:"ASUSTEK COMPUTER INCORPORATION"`
	BoardName     string `json:"board_name" binding:"required" example:"ROG CROSSHAIR X670E HERO"`
	MACAddress    string `json:"mac_address" binding:"required" example:"B42499FE0000"`
}

// BoardProducer: Motherboard manufacturer
//
// BoardName: Motherboard model
//
// MACAddress: Ethernet MAC address of the motherboard
//
// Note: Additional information
type ApplyTempPermitInfo struct {
	BoardProducer string `json:"board_producer" binding:"required" example:"ASUSTEK COMPUTER INCORPORATION"`
	BoardName     string `json:"board_name" binding:"required" example:"ROG CROSSHAIR X670E HERO"`
	MACAddress    string `json:"mac_address" binding:"required" example:"B42499FE0000"`
}
