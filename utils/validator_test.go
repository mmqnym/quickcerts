package utils

import (
	"testing"

	cfg "QuickCertS/configs"

	"github.com/stretchr/testify/assert"
)

func TestIsExcludeField(t *testing.T) {
	// Test valid case
	testList := []string{
		"id0",
		"id1",
		"id2",
	}
	res := isExcludeField("id0", testList)
	assert.Equal(t, res, true)
	// Test invalid case
	res = isExcludeField("id3", testList)
	assert.Equal(t, res, false)
}

func TestIsValidData(t *testing.T) {
	// Test valid case
	type TestData struct {
		Name string
		ID   string
		Age  string
	}

	test0 := TestData{
		Name: "test",
		ID:   "1",
		Age:  "20",
	}

	res := IsValidData(test0, nil)
	assert.Equal(t, res, true)

	test1 := TestData{
		Name: "test",
		ID:   "",
		Age:  "",
	}

	res = IsValidData(test1, []string{"ID", "Age"})
	assert.Equal(t, res, true)

	// Test invalid case
	test2 := TestData{
		Name: "test",
		ID:   "",
		Age:  "",
	}

	res = IsValidData(test2, nil)
	assert.Equal(t, res, false)
}

func TestGenerateRandomString(t *testing.T) {
	res, _ := generateRandomString(6)
	assert.Equal(t, len(res), 6)
}

func TestGenerateRunTimeCode(t *testing.T) {
	code, _ := GenerateRunTimeCode()
	assert.Equal(t, len(code), cfg.SERVER_CONFIG.RUNTIME_CODE_LENGTH)
}
