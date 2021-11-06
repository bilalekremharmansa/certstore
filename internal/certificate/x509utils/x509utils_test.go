package x509utils

import (
	"testing"

	"bilalekrem.com/certstore/internal/testutils"
)

func TestParsePrivateKey(t *testing.T) {
	pemPrivateKey := testutils.GetCAPrivateKey()
	_, err := ParsePemPrivateKey([]byte(pemPrivateKey))
	if err != nil {
		t.Fatalf("parsing certificate failed: [%v]", err)
	}
}

func TestParseCertificate(t *testing.T) {
	pemCert := testutils.GetCAPem()

	cert, err := ParsePemCertificate([]byte(pemCert))
	if err != nil {
		t.Fatalf("parsing certificate failed: [%v]", err)
	}

	expected := "test"
	if cert.Subject.CommonName != expected {
		t.Fatalf("cert common name is not expected: [%s], actual: [%s]", expected, cert.Subject.CommonName)
	}
}

func TestRandomCertSerialNumber(t *testing.T) {
	first, err := GetRandomCertificateSerialNumber()
	if err != nil {
		t.Fatal("Creating random serial number failed")
	}
	second, err := GetRandomCertificateSerialNumber()
	if err != nil {
		t.Fatal("Creating random serial number failed")
	}

	if first == second {
		t.Fatal("Created *random* serial numbers are same, should've been different!")
	}
}
