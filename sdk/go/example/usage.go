package main

import (
	"GoQCS/goqcs"
	"fmt"
)

func main() {
	qcsA := goqcs.NewQCSAdmin(
		"127.0.0.1", 
		33333, 
		"/api/v1", 
		false, 
		"0b09b6dc41f61813346ba76322d19e07a0b71ba939a1bf90211dfff40f552ed0", 
		"",
	)

	// Create a serial number.
	csRes, err := qcsA.CreateSN("f58e-a591-txxx-4944-a1fd-aee7", "none")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(csRes.Msg, csRes.SerialNumber)
	}

	// Generate serial number(s).
	gsRes, err := qcsA.GenerateSN(3, "none")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(gsRes.Msg)
		fmt.Printf("%v", gsRes.SerialNumbers)
	}

	// Get all records.
	garRes, err := qcsA.GetAllRecord()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v", garRes.Data)
	}

	// Get avaliable serial numbers.
	gasRes, err := qcsA.GetAvaliableSN()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v", gasRes.Data)
	}

	// Update note of a serial number.
	usnRes, err := qcsA.UpdateSNNote(
		"93bc-f7a5-b5d2-488c-9c02-bd9d",
		"Additional information",
	)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(usnRes.Msg, usnRes.Note)
	}

	qcsClinet := goqcs.NewQCSClient(
		"127.0.0.1", 
		33333, 
		"/api/v1", 
		false, 
		"QcsTestToken********************************", 
	)

	// Apply certificate.
	acRes, err := qcsClinet.ApplyCert("93bc-f7a5-b5d2-488c-9c02-bd9d", "ASUSTeK Computer Inc.", "ROG STRIX Z790-A GAMING WIFI", "XXXXXXXXXXXX")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(acRes.Key, acRes.Signature)
	}

	// Apply temporary permit.
	atpRes, err := qcsClinet.ApplyTempPermit("ASUSTeK Computer Inc.", "ROG STRIX Z790-A GAMING WIFI", "XXXXXXXXXXXX")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(atpRes.RemainingTime, atpRes.Status)
	}
}