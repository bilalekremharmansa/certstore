package x509utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"errors"
)

func GetRandomCertificateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	return serialNumber, nil
}

func EncodePEMCertAndKey(privateKey *rsa.PrivateKey, cert []byte) (*bytes.Buffer, *bytes.Buffer) {
	privateKeyPem := new(bytes.Buffer)
	pem.Encode(privateKeyPem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// ---

	certPem := new(bytes.Buffer)
	pem.Encode(certPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	return privateKeyPem, certPem
}

func ParsePemPrivateKey(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("decoding pem failed for private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func ParsePemCertificate(certPem []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certPem)
	if block == nil {
		return nil, errors.New("decoding pem failed for certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}