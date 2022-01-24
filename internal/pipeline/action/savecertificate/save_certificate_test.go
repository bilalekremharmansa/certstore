package savecertificate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/pipeline/action/issuecertificate"
	"bilalekrem.com/certstore/internal/pipeline/context"
)

func TestRun(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "test_save_certificate_action")
	assert.NotError(t, err, "creating temp dir")
	defer os.RemoveAll(dir)

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
	assert.NotError(t, err, "running action")

	// ----

	actualCertificateContent, err := ioutil.ReadFile(certificateTargetPath)
	assert.NotError(t, err, "reading file")

	assert.TrueM(t, bytes.Compare(certificateContent, actualCertificateContent) == 0, "certificate content is not correct")

	// ----

	actualCertificateKeyContent, err := ioutil.ReadFile(certificateKeyTargetPath)
	assert.NotError(t, err, "reading file")

	assert.TrueM(t, bytes.Compare(certificateKeyContent, actualCertificateKeyContent) == 0, "certificate key content is not correct")

}

func TestRequiredArgumentCertificate(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_KEY_TARGET_PATH] = "test"

	err := NewSaveCertificateAction().Run(nil, args)
	assert.ErrorContains(t, err, "required argument")
}

func TestRequiredArgumentPrivateKey(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_TARGET_PATH] = "test"

	err := NewSaveCertificateAction().Run(nil, args)
	assert.ErrorContains(t, err, "required argument")
}

func TestIssuedCertificateIsNotInContext(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_TARGET_PATH] = "test"
	args[ARGS_CERTIFICATE_KEY_TARGET_PATH] = "test"

	ctx := context.New()
	ctx.StoreValue(issuecertificate.ISSUED_PRIVATE_KEY_CTX_KEY, []byte("test"))

	err := NewSaveCertificateAction().Run(ctx, args)
	assert.ErrorContains(t, err, "required context object")
}

func TestIssuedCertificateKeyIsNotInContext(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_TARGET_PATH] = "test"
	args[ARGS_CERTIFICATE_KEY_TARGET_PATH] = "test"

	ctx := context.New()
	ctx.StoreValue(issuecertificate.ISSUED_CERTIFICATE_CTX_KEY, []byte("test"))

	err := NewSaveCertificateAction().Run(ctx, args)
	assert.ErrorContains(t, err, "required context object")
}
