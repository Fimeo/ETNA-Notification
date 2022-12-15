package security

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestGenerateKey(t *testing.T) {
	// Create the keys
	priv, pub := GenerateRsaKeyPair()

	// Export the keys to pem string
	privPem := ExportRsaPrivateKeyAsPemStr(priv)
	pubPem, _ := ExportRsaPublicKeyAsPemStr(pub)

	// Import the keys from pem string
	privParsed, _ := ParseRsaPrivateKeyFromPemStr(privPem)
	pubParsed, _ := ParseRsaPublicKeyFromPemStr(pubPem)

	// Export the newly imported keys
	privParsedPem := ExportRsaPrivateKeyAsPemStr(privParsed)
	pubParsedPem, _ := ExportRsaPublicKeyAsPemStr(pubParsed)

	fmt.Println(privParsedPem)
	fmt.Println(pubParsedPem)

	// Check that the exported/imported keys match the original keys
	td.CmpFalse(t, privPem != privParsedPem, "Failure: Export and Import did not result in same Keys")
	td.CmpFalse(t, pubPem != pubParsedPem, "Failure: Export and Import did not result in same Keys")
}

func TestCryptAndDecrypt(t *testing.T) {
	// Create the keys
	priv, pub := GenerateRsaKeyPair()

	message := "hello"
	bytes, err := Encrypt([]byte(message), *pub)
	td.CmpNoError(t, err)
	td.CmpNotNil(t, bytes)

	decrypt, err := Decrypt(bytes, *priv)
	td.CmpNoError(t, err)
	td.CmpNotNil(t, decrypt)

	td.Cmp(t, string(decrypt), message)
}
