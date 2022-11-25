package main

import (
	"fmt"
	"log"
	"os"

	"etna-notification/pkg/security"
)

// main This tool creates and store rsa public and private keys in the config directory
// as rsa.pub and rsa.private file. If these keys already exists, they are overwritten.
func main() {
	// Create the keys
	private, pub := security.GenerateRsaKeyPair()

	// Export the keys to pem string
	privatePem := security.ExportRsaPrivateKeyAsPemStr(private)
	pubPem, _ := security.ExportRsaPublicKeyAsPemStr(pub)

	pubFile, err := os.Create("config/rsa.pub")
	if err != nil {
		log.Fatal(err)
	}
	defer pubFile.Close()

	privateFile, err := os.Create("config/rsa.private")
	if err != nil {
		log.Fatal(err)
	}
	defer privateFile.Close()

	_, err = pubFile.WriteString(pubPem)
	if err != nil {
		log.Fatal(err)
	}

	_, err = privateFile.WriteString(privatePem)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("RSA file generated under config directory")
}
