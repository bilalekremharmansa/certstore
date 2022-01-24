package savecertificate

import (
	"io/ioutil"

	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/pipeline/action/issuecertificate"
	"bilalekrem.com/certstore/internal/pipeline/context"
)

const (
	ARGS_CERTIFICATE_TARGET_PATH     string = "certificate-target-path"
	ARGS_CERTIFICATE_KEY_TARGET_PATH string = "certificate-key-target-path"
)

type SaveCertificateAction struct {
}

func NewSaveCertificateAction() SaveCertificateAction {
	return SaveCertificateAction{}
}

func (a SaveCertificateAction) Run(ctx *context.Context, args map[string]string) error {
	err := validate(ctx, args)
	if err != nil {
		logging.GetLogger().Errorf("validation args failed, %v", err)
		return err
	}

	// --

	certificate := ctx.GetValue(issuecertificate.ISSUED_CERTIFICATE_CTX_KEY).([]byte)
	privateKey := ctx.GetValue(issuecertificate.ISSUED_PRIVATE_KEY_CTX_KEY).([]byte)

	// --

	targetCertificatePath := args[ARGS_CERTIFICATE_TARGET_PATH]
	logging.GetLogger().Debugf("saving certificate to target path: [%s]", targetCertificatePath)
	err = ioutil.WriteFile(targetCertificatePath, certificate, 0666)
	if err != nil {
		logging.GetLogger().Errorf("writing certificate to file failed, %v", err)
		return err
	}

	targetPrivateKeyPath := args[ARGS_CERTIFICATE_KEY_TARGET_PATH]
	logging.GetLogger().Debugf("saving certificate key to target path: [%s]", targetPrivateKeyPath)
	err = ioutil.WriteFile(targetPrivateKeyPath, privateKey, 0666)
	if err != nil {
		logging.GetLogger().Errorf("writing certificate key to file failed, %v", err)
		return err
	}

	// --

	return nil
}

func validate(ctx *context.Context, args map[string]string) error {
	err := action.ValidateRequiredArgs(args, ARGS_CERTIFICATE_TARGET_PATH, ARGS_CERTIFICATE_KEY_TARGET_PATH)
	if err != nil {
		return err
	}

	err = action.ValidateContextObjectExists(ctx, issuecertificate.ISSUED_CERTIFICATE_CTX_KEY,
		issuecertificate.ISSUED_PRIVATE_KEY_CTX_KEY)
	if err != nil {
		return err
	}

	return nil
}
