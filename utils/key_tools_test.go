package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
)

func TestGenerateSN(t *testing.T) {
	sn, _ := GenerateSN()
	assert.Equal(t, len(sn), 29)
}

func TestGenerateKey(t *testing.T) {
	// Using SHA3-256
	testMsg := "test"
	key, _ := GenerateKey(testMsg)
	assert.Equal(t, key, "8652072d7ffe1e52b9aea293d73b7479e9591d8e05c71acec3f4626cb574e723")

	testMsg = "test2"
	key, _ = GenerateKey(testMsg)
	assert.Equal(t, key, "4c1ebfd5c087adf6ef0a13c79651bea9095404b5465c0f5259b368bbb974e07c")
}

func TestSignMessage(t *testing.T) {
	// Using SHA3-512 with PSS(salt length = hash length)
	testMsg := "test"
	key, _ := GenerateKey(testMsg)
	signature, _ := SignMessage([]byte(key))

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

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		t.Fatal("failed to decode public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	hasher := sha3.New512()
	hasher.Write([]byte(key))
	hash := hasher.Sum(nil)

	opts := &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA3_512,
	}

	err = rsa.VerifyPSS(publicKey.(*rsa.PublicKey), crypto.SHA3_512, hash, signature, opts)
	assert.Nil(t, err)
}

func TestGetHash(t *testing.T) {
	cryptoType, hash := getHash("sha-256", []byte("test"))
	hexHash := hex.EncodeToString(hash)
	assert.Equal(t, cryptoType, crypto.SHA256)
	assert.Equal(t, hexHash, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")

	cryptoType, hash = getHash("sha-384", []byte("test"))
	hexHash = hex.EncodeToString(hash)
	assert.Equal(t, cryptoType, crypto.SHA384)
	assert.Equal(t, hexHash,
		"768412320f7b0aa5812fce428dc4706b3cae50e02a64caa16a782249bfe8efc4b7ef1ccb126255d196047dfedf17a0a9",
	)

	cryptoType, hash = getHash("sha-512", []byte("test"))
	hexHash = hex.EncodeToString(hash)
	assert.Equal(t, cryptoType, crypto.SHA512)
	assert.Equal(t, hexHash,
		"ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff",
	)

	cryptoType, hash = getHash("sha3-256", []byte("test"))
	hexHash = hex.EncodeToString(hash)
	assert.Equal(t, cryptoType, crypto.SHA3_256)
	assert.Equal(t, hexHash, "36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80")

	cryptoType, hash = getHash("sha3-384", []byte("test"))
	hexHash = hex.EncodeToString(hash)
	assert.Equal(t, cryptoType, crypto.SHA3_384)
	assert.Equal(t, hexHash,
		"e516dabb23b6e30026863543282780a3ae0dccf05551cf0295178d7ff0f1b41eecb9db3ff219007c4e097260d58621bd",
	)

	cryptoType, hash = getHash("sha3-512", []byte("test"))
	hexHash = hex.EncodeToString(hash)
	assert.Equal(t, cryptoType, crypto.SHA3_512)
	assert.Equal(t, hexHash,
		"9ece086e9bac491fac5c1d1046ca11d737b92a2b2ebd93f005d7b710110c0a678288166e7fbe796883a4f2e9b3ca9f484f521d0ce464345cc1aec96779149c14",
	)

	cryptoType, hash = getHash("unknown", []byte("test"))
	hexHash = hex.EncodeToString(hash)
	assert.Equal(t, cryptoType, crypto.SHA256)
	assert.Equal(t, hexHash, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
}
