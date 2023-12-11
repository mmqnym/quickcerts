package goqcs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getQCSAdmin() *QCSAdmin {
	return NewQCSAdmin(
		"127.0.0.1", 
		33333, 
		"/api/v1", 
		false, 
		"0b09b6dc41f61813346ba76322d19e07a0b71ba939a1bf90211dfff40f552ed0", 
		"",
	)
}

func getQCSClient() *QCSClient {
	return NewQCSClient(
		"127.0.0.1", 
		33333, 
		"/api/v1", 
		false, 
		"QcsTestToken********************************", 
	)
}

func TestCreateSN(t *testing.T) {
	qcsA := getQCSAdmin()

	snList := []string {
		"XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
		"XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
		"YYYY-YYYY-YYYY-YYYY-YYYY-YYYY",
	}

	for _, sn := range snList {
		res, err := qcsA.CreateSN(sn, "none")

		if err != nil {
			if err.Error() != "QCS::Error:the S/N already exists" {
				t.Fatal(err)
			}
			continue
		}

		assert.Equal(t, "Successfully uploaded a new S/N.", res.Msg)
		assert.Equal(t, sn, res.SerialNumber)
	}
}

func TestGenerateSN(t *testing.T) {
	qcsA := getQCSAdmin()

	reason := "test"
	generateCount := 3

	res, err := qcsA.GenerateSN(3, reason)
	if err != nil {
		t.Fatal(err)
	}
	
	targetMsg := fmt.Sprintf("Successfully uploaded new S/N (%d) with reason (%s).", generateCount, reason)

	assert.Equal(t, targetMsg, res.Msg)
	assert.Equal(t, generateCount, len(res.SerialNumbers))
}

func TestGetAllRecords(t *testing.T) {
	qcsA := getQCSAdmin()

	_, err := qcsA.GetAllRecords()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAvailableSN(t *testing.T) {
	qcsA := getQCSAdmin()

	_, err := qcsA.GetAvailableSN()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateSNNote(t *testing.T) {
	qcsA := getQCSAdmin()
	
	snList := []string {
		"XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
		"XXXX-XXXX-XXXX-XXXX-XXXX-1234", // Not Exist
	}

	for _, sn := range snList {
		res, err := qcsA.UpdateSNNote(
			sn,
			"Additional information",
		)

		if err != nil {
			assert.Equal(t, "QCS::Error:The S/N [XXXX-XXXX-XXXX-XXXX-XXXX-1234] does not exist.", err.Error())
			continue
		}

		assert.Equal(t, "Successfully updated the note of specified S/N.", res.Msg)
		assert.Equal(t, "Additional information", res.Note)
	}
}

func TestApplyCert(t *testing.T) {
	qcsC := getQCSClient()

	snList := []string {
		"XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
		"XXXX-XXXX-XXXX-XXXX-XXXX-1234", // Not Exist
	}

	for _, sn := range snList {
		res, err := qcsC.ApplyCert(
			sn,
			"ASUSTeK Computer Inc.",
			"ROG STRIX Z790-A GAMING WIFI",
			"XXXXXXXXXXXX",
		)

		if err != nil {
			if err.Error() != "QCS::Error:The S/N does not exist." {
				t.Fatal(err)
			}
			continue
		}

		assert.Equal(t, "74f996b5670352cab3e8749e7074a158dc716deb3bbc681dd0b79f763d2396f6", res.Key)
		assert.NotEmpty(t, res.Signature, "Signature should not be empty.")
	}
}

func TestApplyTempPermit(t *testing.T) {
	qcsC := getQCSClient()
	res, err := qcsC.ApplyTempPermit(
		"ASUSTeK Computer Inc.", 
		"ROG STRIX Z790-A GAMING WIFI", 
		"XXXXXXXXXXXX",
	)

	if err != nil {
		t.Fatal(err)
	} 

	// Default value of temporary permit is 7 days(604800 seconds).
	assert.Equal(t, 604800.0, res.RemainingTime)
	assert.NotEmpty(t, res.Status, "activated")
}