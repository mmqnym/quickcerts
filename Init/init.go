package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
    waitUserConfirm()
    createFolders()
    createKeyFiles()
    
	fmt.Println(decorateColor("\nInitialization completed.", "green"))
}

func waitUserConfirm() {
    if len(os.Args) > 1 && os.Args[1] == "y" {
        return // Skip the confirmation.
    }

    confirmMsg := color.HiCyanString("The initialization process will create the following folders and files:\n")
    willAddFiles := color.HiMagentaString(
        "    ./logs\n" +
        "    ./local\n" +
        "    ./local/private_key.pem\n" +
        "    ./local/public_key.pem\n",
    )
    confirmAns := color.HiCyanString("Press [Y/y] to continue, or [ANY] to cancel: ")

    fmt.Println(confirmMsg)
    fmt.Println(willAddFiles)
    fmt.Print(confirmAns)

    var input string
    fmt.Scanln(&input)

    if strings.ToUpper(input) != "Y" {
        fmt.Println(decorateColor("Initialization canceled.", "red"))
        os.Exit(0)
    } else {
        fmt.Println(decorateColor("Initialization started...", "green"))
    }
}

func exitWithError(err error) {
    fmt.Println(decorateColor("Failed", "red"))
    fmt.Println(decorateColor(err.Error(), "red"))
    os.Exit(1)
}

func createFolders() {
    fmt.Print("Create the logs folder... ")
    err := os.Mkdir("logs", os.FileMode(0700))

    if err != nil {
        if !strings.Contains(err.Error(), "exists") {
            exitWithError(err)
        }
    }

    fmt.Println(decorateColor("OK", "green"))
    fmt.Print("Create the local folder... ")

    err = os.Mkdir("local", os.FileMode(0700))

    if err != nil {
        if !strings.Contains(err.Error(), "exists") {
            exitWithError(err)
        }
    }

    fmt.Println(decorateColor("OK", "green"))
}

func createKeyFiles() {
    fmt.Print("Generating key files... ")
	prvivateKey, publicKey, err := generateKey()

    if err != nil {
        exitWithError(err)
    }

	f0, err := os.Create("./local/private_key.pem")
	if err != nil {
        exitWithError(err)
    }

	defer f0.Close()

	_, err = f0.Write(prvivateKey)

    if err != nil {
        exitWithError(err)
    }

	f1, err := os.Create("./local/public_key.pem")
	if err != nil {
        exitWithError(err)
    }

	defer f1.Close()

	_, err = f1.Write(publicKey)

    if err != nil {
        exitWithError(err)
    }

    fmt.Println(decorateColor("OK", "green"))
}

func generateKey() ([]byte, []byte, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, nil, err
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

func decorateColor(msg string, colorName string) string {
    switch strings.ToLower(colorName) {
		case "green":
            return color.HiGreenString(msg)
        case "red":
			return color.HiRedString(msg)
        case "cyan":
            return color.HiCyanString(msg)
        case "magenta":
            return color.HiMagentaString(msg)
		default:
			return color.HiWhiteString(msg)
    }
}