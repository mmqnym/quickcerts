package main

// This is an example for verifying the signature by the given public key.
// Used arguments:
// hash method: SHA3-512
// PSS salt length: rsa.PSSSaltLengthEqualsHash
// Message: 95e156395687128711f29b68fbc44573667bdfc5f0d65010cb0555b62138d830
// Signature: upagNzGSL3ZqCsxApgG8yiG/x1c+ZZBJgNtzvZR2KYVLP60+hAr5WcnZ129PG486rl6r2kLMwq8jIu4CUSvwpIblqCILWk7kxQzlei+//7JweQxLbkXfWgdmwA1mUflBXyqQ4vAFyL4w3g44GilInp0nT/iswdAFiCgb5RaK8xkmq+HDeghQWHsNxkPjf7ffDU8wnaLxAK0w4vwYm8BdhzKvEyRFbiTFohLwa4F9byVGrTIAEj53CQ0VvbKwQT6SH+LUVAp5Wr5vMPAREebx/0X5Yy63EuXWvCdZwG64n/TAm4qFhMThrtX+8h+zyf+CViDSZ1xAwkPNtfaQ3scN7g==

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
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzTuY9ePxSX533aa54/aY
Qobqzz0/alc40C31/fYgYXLQVeMJ4vXBHKFhWOaf+ZBf2bQBLx2aIa2ODZcH4ZNF
UIbSZu9jmWN6kcSCw5IMPuDW2YF0b0MlxCemPgCPdIioBa/qsgmy4/s6LpZ2JtUG
7+KBOJIBxuzt8k2XtfRK7k8HBL5v3pQI6IqgooN6cq/M9IOWges1RwLTsMcUbISm
pSOGIC57XmreGiOQik3IlWLYaDbo5nOhzhGtnz6FlAOscW3guYuMBiPjYnTERXNz
1rwx1dHM+t+K2/7poB477RoBEHeLYkEF2JkxVZAXdAg+5PKkMj+Cd/U867t83mDG
OQIDAQAB
-----END PUBLIC KEY-----
`)

type verifyInfo struct {
    hashType         crypto.Hash
    hashMessage      []byte
    signature        []byte
}

func VerifyExample() {
    defer func() {
        if err := recover(); err != nil {
            errMsg := err.(error).Error()

            if strings.Contains(errMsg, "verification error") {
                println("FAIL")
            } else {
                println("PASS")
            }
        }
    }()

    v := getVerfiyInfo()
    verify(v)
}

func getVerfiyInfo() *verifyInfo {
    // Write your own hash method and message(key) here.
    hashType, hash := getHash("sha3-512", []byte("95e156395687128711f29b68fbc44573667bdfc5f0d65010cb0555b62138d830"))

    // Write your own signature here.
    signatureBase64 := "upagNzGSL3ZqCsxApgG8yiG/x1c+ZZBJgNtzvZR2KYVLP60+hAr5WcnZ129PG486rl6r2kLMwq8jIu4CUSvwpIblqCILWk7kxQzlei+//7JweQxLbkXfWgdmwA1mUflBXyqQ4vAFyL4w3g44GilInp0nT/iswdAFiCgb5RaK8xkmq+HDeghQWHsNxkPjf7ffDU8wnaLxAK0w4vwYm8BdhzKvEyRFbiTFohLwa4F9byVGrTIAEj53CQ0VvbKwQT6SH+LUVAp5Wr5vMPAREebx/0X5Yy63EuXWvCdZwG64n/TAm4qFhMThrtX+8h+zyf+CViDSZ1xAwkPNtfaQ3scN7g=="
    signature, err := base64.StdEncoding.DecodeString(signatureBase64)

    if err != nil {
        panic(err)
    }

    return &verifyInfo{
        hashType:         hashType,
        hashMessage:      hash,
        signature:        signature,
    }
}

// Verify the given signature with PSS & the admin specified hashing method.
func verify(v *verifyInfo) {
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
        SaltLength: rsa.PSSSaltLengthEqualsHash,
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