package utils

import (
	cfg "QuickCertS/configs"
	"crypto/rand"
	"math/big"
	"reflect"
)

// Check if the given data is all not empty except the excluded fields.
func IsValidData(iData interface{}, exclude []string) bool {
    data := reflect.ValueOf(iData)

    for i := 0; i < data.NumField(); i++ {
        fieldName := data.Type().Field(i).Name
        fieldValue := data.Field(i).Interface()

        if fieldValue == "" && !isExcludeField(fieldName, exclude) {
            return false
        }
    }

    return true
}

// Generate a random run time code.
func GenerateRunTimeCode() (string, error) {
    length := cfg.SERVER_CONFIG.RUNTIME_CODE_LENGTH
    code, err := generateRandomString(length)

    if err != nil {
        return "", err
    }

    return code, nil
}

// Generate a random string with the given length.
func generateRandomString(length int) (string, error) {
    const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    b := make([]byte, length)

    for i := range b {
        randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            return "", err
        }
        b[i] = charset[randomIndex.Int64()]
    }

    return string(b), nil
}

// Check if the given field is in the excluded list.
func isExcludeField(field string, exclude []string) bool {
    for _, value := range exclude {
        if field == value {
            return true
        }
    }
    return false
}