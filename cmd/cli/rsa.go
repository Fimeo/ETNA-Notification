package main

import (
	"fmt"
	"log"
	"os"

	"etna-notification/internal/infrastructure/security"
)

// Create and store rsa keys
func main() {
	// Create the keys
	priv, pub := security.GenerateRsaKeyPair()

	// Export the keys to pem string
	priv_pem := security.ExportRsaPrivateKeyAsPemStr(priv)
	pub_pem, _ := security.ExportRsaPublicKeyAsPemStr(pub)

	pubf, err := os.Create("config/rsa.pub")
	if err != nil {
		log.Fatal(err)
	}
	defer pubf.Close()

	privatef, err := os.Create("config/rsa.private")
	if err != nil {
		log.Fatal(err)
	}
	defer privatef.Close()

	_, err = pubf.WriteString(pub_pem)
	if err != nil {
		log.Fatal(err)
	}

	_, err = privatef.WriteString(priv_pem)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done")
}
