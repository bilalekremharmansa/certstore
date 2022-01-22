package issuecertificate

import (
	go_ctx "context"
	b64 "encoding/base64"
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	grpc "bilalekrem.com/certstore/internal/certstore/grpc/gen"
	"bilalekrem.com/certstore/internal/pipeline/context"
	"github.com/golang/mock/gomock"
)

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc.NewMockCertificateServiceClient(ctrl)
	action := NewIssueCertificateAction(mockClient)

	expectedCertificate := "cert payload"
	expectedPrivateKey := "cert key payload"
	mockClient.
		EXPECT().
		IssueCertificate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ go_ctx.Context, _ interface{}, opts ...interface{}) (*grpc.CertificateResponse, error) {
			b64EncodedCertificate := b64.StdEncoding.EncodeToString([]byte(expectedCertificate))
			b64EncodedPrivateKey := b64.StdEncoding.EncodeToString([]byte(expectedPrivateKey))
			return &grpc.CertificateResponse{
				Certificate: b64EncodedCertificate,
				PrivateKey:  b64EncodedPrivateKey,
			}, nil
		})

	// ----

	ctx := context.New()
	args := getValidArgs()

	err := action.Run(ctx, args)
	assert.NotError(t, err, "running action")

	// -----

	certificate := ctx.GetValue(ISSUED_CERTIFICATE_CTX_KEY).([]byte)
	assert.Equal(t, string(certificate), expectedCertificate)

	privateKey := ctx.GetValue(ISSUED_PRIVATE_KEY_CTX_KEY).([]byte)
	assert.Equal(t, string(privateKey), expectedPrivateKey)
}

func TestRequiredArgumentIssuer(t *testing.T) {
	testRequiredArgument(t, ARGS_ISSUER)
}

func TestRequiredArgumentCommonName(t *testing.T) {
	testRequiredArgument(t, ARGS_COMMON_NAME)
}

func TestRequiredArgumentEmail(t *testing.T) {
	testRequiredArgument(t, ARGS_EMAIL)
}

func TestRequiredArgumentExpirationDays(t *testing.T) {
	testRequiredArgument(t, ARGS_EXPIRATION_DAYS)
}

func TestExpirationDaysConvertableInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc.NewMockCertificateServiceClient(ctrl)
	action := NewIssueCertificateAction(mockClient)

	mockClient.
		EXPECT().
		IssueCertificate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ go_ctx.Context, req *grpc.CertificateRequest, opts ...interface{}) (*grpc.CertificateResponse, error) {
			assert.Equal(t, int32(30), req.ExpirationDays)

			return &grpc.CertificateResponse{
				Certificate: "",
				PrivateKey:  "",
			}, nil
		})

	args := getValidArgs()
	args[ARGS_EXPIRATION_DAYS] = "30"

	err := action.Run(context.New(), args)
	assert.NotError(t, err, "running action")
}

func TestExpirationDaysNotConvertableInt(t *testing.T) {
	action := NewIssueCertificateAction(nil)

	args := getValidArgs()
	args[ARGS_EXPIRATION_DAYS] = "thirty"

	err := action.Run(context.New(), args)
	assert.ErrorContains(t, err, "invalid syntax")
}

func TestMultipleSANs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc.NewMockCertificateServiceClient(ctrl)
	action := NewIssueCertificateAction(mockClient)

	mockClient.
		EXPECT().
		IssueCertificate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ go_ctx.Context, req *grpc.CertificateRequest, opts ...interface{}) (*grpc.CertificateResponse, error) {
			expectedSans := []string{"a.com", "b.com"}
			assert.DeepEqual(t, expectedSans, req.SANs)

			return &grpc.CertificateResponse{
				Certificate: "",
				PrivateKey:  "",
			}, nil
		})

	args := getValidArgs()

	err := action.Run(context.New(), args)
	assert.NotError(t, err, "running action")

}

func TestSingleSAN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc.NewMockCertificateServiceClient(ctrl)
	action := NewIssueCertificateAction(mockClient)

	mockClient.
		EXPECT().
		IssueCertificate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ go_ctx.Context, req *grpc.CertificateRequest, opts ...interface{}) (*grpc.CertificateResponse, error) {
			expectedSans := []string{"a.com"}
			assert.DeepEqual(t, expectedSans, req.SANs)

			return &grpc.CertificateResponse{
				Certificate: "",
				PrivateKey:  "",
			}, nil
		})

	args := getValidArgs()
	args[ARGS_SANS] = "a.com"

	err := action.Run(context.New(), args)
	assert.NotError(t, err, "running action")

}

func TestEmptySAN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc.NewMockCertificateServiceClient(ctrl)
	action := NewIssueCertificateAction(mockClient)

	mockClient.
		EXPECT().
		IssueCertificate(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ go_ctx.Context, req *grpc.CertificateRequest, opts ...interface{}) (*grpc.CertificateResponse, error) {
			assert.Equal(t, 0, len(req.SANs))

			return &grpc.CertificateResponse{
				Certificate: "",
				PrivateKey:  "",
			}, nil
		})

	args := getValidArgs()
	delete(args, ARGS_SANS)

	err := action.Run(context.New(), args)
	assert.NotError(t, err, "running action")

}

// -----

func testRequiredArgument(t *testing.T, arg string) {
	action := NewIssueCertificateAction(nil)

	args := getValidArgs()
	delete(args, arg)

	err := action.Run(context.New(), args)
	assert.ErrorContains(t, err, "required argument")
}

func getValidArgs() map[string]string {
	args := make(map[string]string)
	args[ARGS_ISSUER] = "issuer"
	args[ARGS_COMMON_NAME] = "common"
	args[ARGS_EMAIL] = "common"
	args[ARGS_ORGANIZATION] = "common"
	args[ARGS_EXPIRATION_DAYS] = "30"
	args[ARGS_SANS] = "a.com;b.com"
	return args
}
