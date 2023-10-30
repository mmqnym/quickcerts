package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
    fmt.Print("Create the logs folder... ")
    os.Mkdir("logs", os.FileMode(0700))
    fmt.Println("Done.")
    fmt.Print("Create the local folder... ")
    os.Mkdir("local", os.FileMode(0700))
    fmt.Println("Done.")

    fmt.Print("Generating key files... ")
	prvivateKey, publicKey, err := generateKey()

    if err != nil {
        log.Fatal(err)
    }

	f_0, err := os.Create("./local/private_key.pem")
	if err != nil {
        log.Fatal(err)
    }

	defer f_0.Close()

	_, err = f_0.Write(prvivateKey)

    if err != nil {
        log.Fatal(err)
    }

	f_1, err := os.Create("./local/public_key.pem")
	if err != nil {
        log.Fatal(err)
    }

	defer f_1.Close()

	_, err = f_1.Write(publicKey)

    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Done.")
	fmt.Println("Initialization completed.")
}

func generateKey() ([]byte, []byte, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        panic(err)
    }

    privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
    if err != nil {
        return nil, nil, err
    }

    privateKeyBlock := &pem.Block{
        Type:  "PRIVATE KEY",
        Bytes: privateKeyBytes,
    }

    privateKeyPEM := pem.EncodeToMemory(privateKeyBlock)

    publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
    if err != nil {
        return nil, nil, err
    }

    publicKeyBlock := &pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: publicKeyBytes,
    }
    
    publicKeyPEM := pem.EncodeToMemory(publicKeyBlock)

    return privateKeyPEM, publicKeyPEM, nil
}