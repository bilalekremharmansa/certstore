package issuecertificate

import (
	go_ctx "context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"bilalekrem.com/certstore/internal/certstore/grpc/gen"
	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline/context"
)

const (
	ISSUED_CERTIFICATE_CTX_KEY context.Key = "issued-certificated"
	ISSUED_PRIVATE_KEY_CTX_KEY context.Key = "issued-certificated-private-key"

	ARGS_ISSUER          string = "issuer"
	ARGS_COMMON_NAME     string = "common-name"
	ARGS_EMAIL           string = "email"
	ARGS_ORGANIZATION    string = "organization"
	ARGS_EXPIRATION_DAYS string = "expiration-days"
	ARGS_SANS            string = "sans"
)

type IssueCertificateAction struct {
	client gen.CertificateServiceClient
}

func NewIssueCertificateAction(client gen.CertificateServiceClient) IssueCertificateAction {
	return IssueCertificateAction{client: client}
}

func (a IssueCertificateAction) Run(ctx *context.Context, args map[string]string) error {
	err := validate(args)
	if err != nil {
		logging.GetLogger().Errorf("validation args failed, %v", err)
		return err
	}

	// --

	issuer := args[ARGS_ISSUER]

	expirationDays, err := strconv.Atoi(args[ARGS_EXPIRATION_DAYS])
	if err != nil {
		logging.GetLogger().Errorf("str to int conversion failed for action arg: expiration-days, %v", err)
		return err
	}

	var sans []string
	if args[ARGS_SANS] != "" {
		sans = strings.Split(args[ARGS_SANS], ";")
	}

	// ----

	request := &gen.CertificateRequest{
		Issuer:         issuer,
		CommonName:     args[ARGS_COMMON_NAME],
		Email:          args[ARGS_EMAIL],
		Organization:   args[ARGS_ORGANIZATION],
		ExpirationDays: int32(expirationDays),
		SANs:           sans,
	}

	logging.GetLogger().Debugf("Issuing certificate for issuer: [%s]", issuer)
	response, err := a.client.IssueCertificate(go_ctx.TODO(), request)
	if err != nil {
		logging.GetLogger().Errorf("issuing certificate for issuer: [%s], failed, %v", issuer, err)
		return err
	}

	// ----

	certificate, err := b64.StdEncoding.DecodeString(response.Certificate)
	if err != nil {
		logging.GetLogger().Errorf("decoding issued certificate, failed, %v", err)
		return err
	}
	privateKey, err := b64.StdEncoding.DecodeString(response.PrivateKey)
	if err != nil {
		logging.GetLogger().Errorf("decoding issued certificate private key, failed, %v", err)
		return err
	}

	// ----

	logging.GetLogger().Debugf("Storing issued certificate into context - [%s]", issuer)
	ctx.StoreValue(ISSUED_CERTIFICATE_CTX_KEY, certificate)
	ctx.StoreValue(ISSUED_PRIVATE_KEY_CTX_KEY, privateKey)

	return nil
}

func validate(args map[string]string) error {
	_, exists := args[ARGS_ISSUER]
	if !exists {
		return errors.New(fmt.Sprintf("required argument: %s", ARGS_ISSUER))
	}

	_, exists = args[ARGS_COMMON_NAME]
	if !exists {
		return errors.New(fmt.Sprintf("required argument: %s", ARGS_COMMON_NAME))
	}

	_, exists = args[ARGS_EMAIL]
	if !exists {
		return errors.New(fmt.Sprintf("required argument: %s", ARGS_EMAIL))
	}

	_, exists = args[ARGS_EXPIRATION_DAYS]
	if !exists {
		return errors.New(fmt.Sprintf("required arguments: %s", ARGS_EXPIRATION_DAYS))
	}

	return nil
}
