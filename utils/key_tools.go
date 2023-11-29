package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"

	cfg "QuickCertS/configs"
)

var privateKeyBytes []byte

func init() {
    var err error
    privateKeyBytes, err = GetPrivateKeyBytes()
    if err != nil {
        Logger.Fatal("Failed to load the private key.")
    }
}

// Generate a serial number by uuid v4 and custom rule(24 bits).
func GenerateSN() (string, error) {
    uuid := make([]byte, 16)
    _, err := rand.Read(uuid)
    if err != nil {
        return "", err
    }
    uuid[6] = (uuid[6] & 0x0f) | 0x40
    uuid[8] = (uuid[8] & 0x3f) | 0x80
    plainUUID := fmt.Sprintf("%x", uuid)[:24]
    return fmt.Sprintf("%s-%s-%s-%s-%s-%s",
        plainUUID[0:4], plainUUID[4:8], plainUUID[8:12], plainUUID[12:16], plainUUID[16:20], plainUUID[20:24]), nil
}


// Generate an APP key by SHA3-256 for the device.
func GenerateKey(base string) (string, error) {
    hash := sha3.New256()
    _, err := hash.Write([]byte(base + "SALT"))

    if err != nil {
        return "", err
    }

    sum := hash.Sum(nil)
    key := fmt.Sprintf("%x", sum)
    return key, err
}

// Sign the given message with specified hashing method.
func SignMessage(message []byte) ([]byte, error) {
    privateKey, err := keyBytesToPrivateKey(privateKeyBytes)

    if err != nil {
        return []byte{}, err
    }

    sinature, err := signMessage(cfg.SERVER_CONFIG.HASHING_METHOD, message, privateKey)

    if err != nil {
        return []byte{}, err
    }

    return sinature, err
}

// Convert the private key bytes to a *rsa.PrivateKey.
func keyBytesToPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode(keyBytes)
    if block == nil {
        return nil, errors.New("private key error: unable to decode PEM block")
    }

    key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    privateKey, ok := key.(*rsa.PrivateKey)
    if !ok {
        return nil, errors.New("private key error: not an RSA private key")
    }

    return privateKey, nil
}

// Get the hash type and hash value by the given method name.
func getHash(methodName string, message []byte) (cryptoType crypto.Hash, hash []byte) {
    switch strings.ToLower(methodName) {
        case "sha-256":
            hash := sha256.Sum256(message)
            return crypto.SHA256, hash[:]
        case "sha-384":
            hash := sha512.Sum384(message)
            return crypto.SHA384, hash[:]
        case "sha-512":
            hash := sha512.Sum512(message)
            return crypto.SHA512, hash[:]
        case "sha3-256":
            hasher := sha3.New256()
            hasher.Write(message)
            return crypto.SHA3_256, hasher.Sum(nil)
        case "sha3-384":
            hasher := sha3.New384()
            hasher.Write(message)
            return crypto.SHA3_384, hasher.Sum(nil)
        case "sha3-512":
            hasher := sha3.New512()
            hasher.Write(message)
            return crypto.SHA3_512, hasher.Sum(nil)
        default:
            // Default to SHA-256
            hash := sha256.Sum256(message)
            return crypto.SHA256, hash[:]
    }
}

// Sign the given message with PSS & the admin specified hashing method.
func signMessage(methodName string, data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    cryptoType, hash := getHash(methodName, data)
    
    opts := &rsa.PSSOptions{
        SaltLength: rsa.PSSSaltLengthEqualsHash,
        Hash:       cryptoType,
    }

    signature, err := rsa.SignPSS(rand.Reader, privateKey, cryptoType, hash[:], opts)

    if err != nil {
        return []byte{}, err
    }

    return signature, err
}