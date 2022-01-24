package shouldrenewcertificate

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"bilalekrem.com/certstore/internal/certificate/x509utils"
	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline/context"
)

const (
	ARGS_CERTIFICATE_PATH string = "certificate-path"

	// will be used to decide renew certificate, if certificate will expire in {ACCEPTABLE_NUM_OF_DAYS_TO_EXPIRE} days
	ACCEPTABLE_NUM_OF_DAYS_TO_EXPIRE = 25
)

type shouldRenewCertificateAction struct {
}

func NewShouldRenewCertificateAction() shouldRenewCertificateAction {
	return shouldRenewCertificateAction{}
}

func (a shouldRenewCertificateAction) Run(ctx *context.Context, args map[string]string) error {
	err := validate(args)
	if err != nil {
		logging.GetLogger().Errorf("validation args failed, %v", err)
		return err
	}

	// ----

	certificatePath := args[ARGS_CERTIFICATE_PATH]
	certificate, err := loadCertificate(certificatePath)
	if err != nil {
		logging.GetLogger().Infof("loading certificate failed, certificate should be renewed, %v", err)
		return nil
	}

	// ----

	shouldRenew := shouldRenewCertificate(certificate)
	if shouldRenew {
		logging.GetLogger().Error("certificate should be renewed")
		return nil
	}

	return errors.New("No need to renew certificate...")
}

func validate(args map[string]string) error {
	_, exists := args[ARGS_CERTIFICATE_PATH]
	if !exists {
		return errors.New(fmt.Sprintf("required argument: %s", ARGS_CERTIFICATE_PATH))
	}

	return nil
}

func loadCertificate(certicatePath string) (*x509.Certificate, error) {
	certificateBytes, err := ioutil.ReadFile(certicatePath)
	if err != nil {
		return nil, err
	}

	certificate, err := x509utils.ParsePemCertificate(certificateBytes)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func shouldRenewCertificate(cert *x509.Certificate) bool {
	now := time.Now()
	certificateExpireDate := cert.NotAfter

	daysBetween := certificateExpireDate.Sub(now).Hours() / 24
	logging.GetLogger().Infof("days between expiration date and now: [%f]", daysBetween)
	return (daysBetween <= ACCEPTABLE_NUM_OF_DAYS_TO_EXPIRE)
}
