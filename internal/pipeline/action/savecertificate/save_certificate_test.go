package savecertificate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"bilalekrem.com/certstore/internal/pipeline/action/issuecertificate"
	"bilalekrem.com/certstore/internal/pipeline/context"
)

func TestRun(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "test_save_certificate_action")
	if err != nil {
		t.Fatalf("error occurred while creating temp dir, %v", err)
	}
	// defer os.RemoveAll(dir)

	// ----

	certificateContent := []byte("certificate content\n")
	certificateKeyContent := []byte("private key content\n")

	ctx := context.New()
	ctx.StoreValue(issuecertificate.ISSUED_CERTIFICATE_CTX_KEY, certificateContent)
	ctx.StoreValue(issuecertificate.ISSUED_PRIVATE_KEY_CTX_KEY, certificateKeyContent)

	// ----

	certificateTargetPath := fmt.Sprintf("%s/test.crt", dir)
	certificateKeyTargetPath := fmt.Sprintf("%s/test.key", dir)

	args := make(map[string]string)
	args[ARGS_CERTIFICATE_TARGET_PATH] = certificateTargetPath
	args[ARGS_CERTIFICATE_KEY_TARGET_PATH] = certificateKeyTargetPath

	// ----

	action := NewSaveCertificateAction()
	err = action.Run(ctx, args)
	if err != nil {
		t.Fatalf("error occurred while running action, %v", err)
	}

	// ----

	actualCertificateContent, err := ioutil.ReadFile(certificateTargetPath)
	if err != nil {
		t.Fatalf("error occurred while reading file, %v", err)
	}

	if bytes.Compare(certificateContent, actualCertificateContent) != 0 {
		t.Fatalf("certificate content is not correct, expected: %s, found: %s", certificateContent, actualCertificateContent)
	}

	// ----

	actualCertificateKeyContent, err := ioutil.ReadFile(certificateKeyTargetPath)
	if err != nil {
		t.Fatalf("error occurred while reading file, %v", err)
	}

	if bytes.Compare(certificateKeyContent, actualCertificateKeyContent) != 0 {
		t.Fatalf("certificate key content is not correct, expected: %s, found: %s", certificateContent, actualCertificateKeyContent)
	}

}

func TestRequiredArgumentCertificate(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_KEY_TARGET_PATH] = "test"

	err := NewSaveCertificateAction().Run(nil, args)
	if err == nil || !strings.Contains(err.Error(), "required argument") {
		t.Fatalf("required arg error is expected but not found")
	}
}

func TestRequiredArgumentPrivateKey(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_TARGET_PATH] = "test"

	err := NewSaveCertificateAction().Run(nil, args)
	if err == nil || !strings.Contains(err.Error(), "required argument") {
		t.Fatalf("required arg error is expected but not found")
	}
}

func TestIssuedCertificateIsNotInContext(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_TARGET_PATH] = "test"
	args[ARGS_CERTIFICATE_KEY_TARGET_PATH] = "test"

	ctx := context.New()
	ctx.StoreValue(issuecertificate.ISSUED_PRIVATE_KEY_CTX_KEY, []byte("test"))

	err := NewSaveCertificateAction().Run(ctx, args)
	if err == nil || !strings.Contains(err.Error(), "not found in context") {
		t.Fatalf("not found in context error is expected but not found")
	}
}

func TestIssuedCertificateKeyIsNotInContext(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_TARGET_PATH] = "test"
	args[ARGS_CERTIFICATE_KEY_TARGET_PATH] = "test"

	ctx := context.New()
	ctx.StoreValue(issuecertificate.ISSUED_CERTIFICATE_CTX_KEY, []byte("test"))

	err := NewSaveCertificateAction().Run(ctx, args)
	if err == nil || !strings.Contains(err.Error(), "not found in context") {
		t.Fatalf("not found in context error is expected but not found")
	}
}
