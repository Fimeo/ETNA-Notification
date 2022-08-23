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
	priv_pem := ExportRsaPrivateKeyAsPemStr(priv)
	pub_pem, _ := ExportRsaPublicKeyAsPemStr(pub)

	// Import the keys from pem string
	priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
	pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)

	// Export the newly imported keys
	priv_parsed_pem := ExportRsaPrivateKeyAsPemStr(priv_parsed)
	pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)

	fmt.Println(priv_parsed_pem)
	fmt.Println(pub_parsed_pem)

	// Check that the exported/imported keys match the original keys
	td.CmpFalse(t, priv_pem != priv_parsed_pem, "Failure: Export and Import did not result in same Keys")
	td.CmpFalse(t, pub_pem != pub_parsed_pem, "Failure: Export and Import did not result in same Keys")
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
