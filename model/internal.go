package model

// For database table `certs`.
type Cert struct {
	SerialNumber string `json:"serial_number" example:"779f-4e90-aebd-4295-881a-f8d7"`
	Key          string `json:"key" example:"3266cd6a16ca77f9c0f0ff9934eb0e29c4b6bb0729cde98811f9f0caf76d603c"`
	Note         string `json:"note" example:"Updated note."`
}
