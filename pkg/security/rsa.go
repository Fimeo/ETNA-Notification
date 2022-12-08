package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Security struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewSecurity() Security {
	public, err := RsaPublicKeyFromFile(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		panic(err)
	}
	private, err := RsaPrivateKeyFromFile(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}
	return Security{
		PublicKey:  public,
		PrivateKey: private,
	}
}

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privateKey, &privateKey.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	return string(privateKeyPem)
}

func RsaPrivateKeyFromFile(path string) (*rsa.PrivateKey, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, errors.New(fmt.Sprintf("file not found : %+v", err))
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot read input file : %+v", err))
	}
	return ParseRsaPrivateKeyFromPemStr(string(content))
}

func RsaPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, errors.New(fmt.Sprintf("file not found : %+v", err))
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot read input file : %+v", err))
	}
	return ParseRsaPublicKeyFromPemStr(string(content))
}

func ParseRsaPrivateKeyFromPemStr(privatePEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return private, nil
}

func ExportRsaPublicKeyAsPemStr(pubKey *rsa.PublicKey) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	pubKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyBytes,
		},
	)

	return string(pubKeyPem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

func Encrypt(data []byte, publicKey rsa.PublicKey) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&publicKey,
		data,
		nil)
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}

func Decrypt(data []byte, privateKey rsa.PrivateKey) ([]byte, error) {
	decryptedBytes, err := privateKey.Decrypt(nil, data, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}
