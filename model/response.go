package model

// ErrorMsg: Error message for response
type ErrorResponse struct {
	Error string `json:"error" example:"Error message."`
}

type ApplyCertResponse struct {
	Key       string `json:"key" example:"3266cd6a16ca77f9c0f0ff9934eb0e29c4b6bb0729cde98811f9f0caf76d603c"`
	Signature string `json:"signature" example:"MNj/g7W+X5PmirfgWl5jveV54t50+LZAPmByh5Py880pB2z67Ser0YvZ2G/mTNV4XcIrKmLy1ICFmQ1esjydhvBj1FOuTm3eTIixUIsFLxwlW2co/R6kCIjNRydB3N7L/kWv+ZwSjsSsdHqmMUleXV3OJruxeoXV8TLRCSGE4tHGEwhPULuBLn2aldIehDTgteJx1O1YNJGIcDM3NWVDjJnUA0Bjhq3oRvXWN4M23SnZZG2vT94wJIK0X5q6oNqFTupFjDVBCFcHeWoxQ5xZdPhfXF8rC/VTb4vkZZm5RIiIK1UC9XVaAsXVPEzlxVfYJ0gh+wULx8syE2QyB5GfyQ=="`
}

type ApplyTempPermitResponse struct {
	Status        string `json:"status" example:"activated"`
	RemainingTime int64  `json:"remaining_time" example:"604800"`
}

type CreateSNResponse struct {
	Msg          string `json:"msg" example:"Successfully uploaded a new S/N."`
	SerialNumber string `json:"serial_number" example:"779f-4e90-aebd-4295-881a-f8d7"`
}

type GenerateSNResponse struct {
	Msg           string   `json:"msg" example:"Successfully generated a new S/N."`
	SerialNumbers []string `json:"serial_numbers" example:"[\"779f-4e90-aebd-4295-881a-f8d7\"]"`
}

type UpdateCertNoteResponse struct {
	Msg  string `json:"msg" example:"Successfully updated the note of specified S/N."`
	Note string `json:"note" example:"Updated note."`
}

type GetAllRecordsResponse struct {
	Data []Cert `json:"data"`
}

type GetAvaliableSNResponse struct {
	Data []string `json:"data" example:"779f-4e90-aebd-4295-881a-f8d7"`
}
