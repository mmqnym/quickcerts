package utils

import (
	cfg "QuickCertS/configs"
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

// Check if the given field is in the excluded list.
func isExcludeField(field string, exclude []string) bool {
	for _, value := range exclude {
		if field == value {
			println("exclude field: ", field)
			return true
		}
	}
	return false
}

func GetValidTokenOwner(permit string) string {
	for _, token := range cfg.ALLOWEDLIST.TOKENS {
		if token.PERMIT == permit {
			return token.NAME
		}
	}

	return ""
}