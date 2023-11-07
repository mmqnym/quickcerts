package utils

import (
	"os"
)

// Get the private key bytes from the local file.
func GetPrivateKeyBytes() ([]byte, error) {
    fileName := "./local/private_key.pem"

    data, err := os.ReadFile(fileName)
    if err != nil {
        return nil, err
    }

    return data, err
}