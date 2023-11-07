package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Paste the public key here and don't reserve any spaces.
var publicKeyBytes = []byte(
`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt8sWoEyEdJso7GP5jVqY
K+pLu3bUFAsWh3NClHM5CzTH34JKPTInFMQbDTaZ2Q23hmC2uLvYKriX8hFa9UOJ
BXz2uwQhSzCu6RwN0Evrbj1DkWo0p6ifOa4BkYt4+mGDtVrGBGeLQCRtU1CoAVal
AtzOKFfHhrE2xinSZk2uDwUq4lClegfU99hqPKAmAXg3s90mZ+D43cdmn0HkjJ/9
qe4aZwP+u2fdgXow0Z+dRnc8NDVsWfMdfduReuwHiuCOFFjWhh83/Wta9i0JrT0Y
uSYgJswTRHa9bI6uatvIwmHV1mADc0/RUl9uJzc0x/pC/RiMlE/4OYU/exL88Xo4
GwIDAQAB
-----END PUBLIC KEY-----
`)

type VerifyInfo struct {
    hashType         crypto.Hash
    hashMessage      []byte
    signature        []byte
}

func main() {
    defer func() {
        if err := recover(); err != nil {
            errMsg := err.(error).Error()

            if strings.Contains(errMsg, "verification error") {
                println("FAIL")
            } else {
                println(errMsg)
            }
        }
    }()

    v := getVerfiyInfo()
    Verify(v)
}

func getVerfiyInfo() *VerifyInfo {
    // Write your own hash method and message here.
    hashType, hash := getHash("sha3-512", []byte("f39476262640eebefde1bb5ede9a0fc721ab7d9d269002ce95fa89dcbc201b69"))

    // Write your own signature here.
    signatureBase64 := "qJrXQVeoGmoObRj4cqAPuhGRanj1yebFAwP6lxRCCUNqN4pgEv8qiRJXZGJP2ky8dtI67aOx48ij8vbUomxl4a3wEvyxXym1KHAd4vVObw393VQYG5nbKvPAVENQlqfJo3MnkYtTR/B4h3zVj1BQjBKE+kGx2J/4i4W9dnuIOAbtcs05dEWr8woE/JFa4LcFfHv+jJp0Exok5oPxIZ8paFq7/CkNlO91b+W62th35gh4e2bqgCEXdwUifA4I2H0LyuEPscuc2yrqYC0Ve+yQQ58c6g7HLW2SXyCJnXbpcDebMtWeXfp8468iQHj2UE4ykzmrnprQ2jOrnIMv62rF4A=="
    signature, err := base64.StdEncoding.DecodeString(signatureBase64)

    if err != nil {
        panic(err)
    }

    return &VerifyInfo{
        hashType:         hashType,
        hashMessage:      hash,
        signature:        signature,
    }
}

// Verify the given signature with PSS & the admin specified hashing method.
func Verify(v *VerifyInfo) {
    println("Verifing...")

    block, _ := pem.Decode(publicKeyBytes)
    if block == nil {
        panic(errors.New("failed to decode public key"))
    }
    
    publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        panic(err)
    }
    
    opts := &rsa.PSSOptions{
        SaltLength: rsa.PSSSaltLengthAuto,
        Hash:       v.hashType,
    }
    
    err = rsa.VerifyPSS(publicKey.(*rsa.PublicKey), v.hashType, v.hashMessage, v.signature, opts)
    if err != nil {
        panic(err)
    }

    println("PASS")
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