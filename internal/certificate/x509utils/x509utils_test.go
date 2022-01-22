package x509utils

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/testutils"
)

func TestParsePrivateKey(t *testing.T) {
	pemPrivateKey := testutils.GetCAPrivateKey()
	_, err := ParsePemPrivateKey([]byte(pemPrivateKey))
	assert.NotError(t, err, "parse private key")
}

func TestParseCertificate(t *testing.T) {
	pemCert := testutils.GetCAPem()

	cert, err := ParsePemCertificate([]byte(pemCert))
	assert.NotError(t, err, "parsing pem certificate")

	assert.Equal(t, "test", cert.Subject.CommonName)
}

func TestRandomCertSerialNumber(t *testing.T) {
	first, err := GetRandomCertificateSerialNumber()
	assert.NotError(t, err, "random seraial number creation")
	second, err := GetRandomCertificateSerialNumber()
	assert.NotError(t, err, "random seraial number creation 2")

	assert.Falsef(t, first == second, "Created *random* serial numbers are same, should've been different!")
}
