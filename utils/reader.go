package utils

import (
	"os"
)

/*
	Get the private key bytes from the local file.
*/
func GetPrivateKeyBytes() ([]byte, error) {
	fileName := "./local/private_key.pem"

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return data, err
}

/*
	Get the public key bytes from the local file.
*/
func GetPublicKeyBytes() ([]byte, error) {
	fileName := "./local/public_key.pem"

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return data, err
}