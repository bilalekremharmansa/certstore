package issuecertificate

import (
	go_ctx "context"
	b64 "encoding/base64"
	"strconv"
	"strings"

	"bilalekrem.com/certstore/internal/certstore/grpc/gen"
	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline/action"
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
	request, err := createCertificateRequest(args)
	if err != nil {
		logging.GetLogger().Errorf("creating certificate request %v", err)
		return err
	}

	// -----

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
	err := action.ValidateRequiredArgs(args, ARGS_ISSUER, ARGS_COMMON_NAME)
	if err != nil {
		return err
	}

	return nil
}

func createCertificateRequest(args map[string]string) (*gen.CertificateRequest, error) {
	issuer := args[ARGS_ISSUER]

	request := &gen.CertificateRequest{
		Issuer:     issuer,
		CommonName: args[ARGS_COMMON_NAME],
	}

	expirationDaysStr, exists := args[ARGS_EXPIRATION_DAYS]
	if exists {
		expirationDays, err := strconv.Atoi(expirationDaysStr)
		if err != nil {
			logging.GetLogger().Errorf("str to int conversion failed for action arg: expiration-days, %v", err)
			return nil, err
		}

		request.ExpirationDays = int32(expirationDays)
	}

	email, exists := args[ARGS_EMAIL]
	if exists {
		request.Email = email
	}

	organization, exists := args[ARGS_ORGANIZATION]
	if exists {
		request.Organization = organization
	}

	sansStr, exists := args[ARGS_SANS]
	if exists && sansStr != "" {
		request.SANs = strings.Split(args[ARGS_SANS], ";")
	}

	return request, nil
}
