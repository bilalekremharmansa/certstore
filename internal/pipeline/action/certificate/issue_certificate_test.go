package certificate

import (
	go_ctx "context"
	"reflect"
	"strings"
	"testing"

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
			return &grpc.CertificateResponse{
				Certificate: expectedCertificate,
				PrivateKey:  expectedPrivateKey,
			}, nil
		})

	// ----

	ctx := context.New()
	args := getValidArgs()

	err := action.Run(ctx, args)
	if err != nil {
		t.Fatalf("action failed with error, %v", err)
	}

	// -----

	certificate := ctx.GetValue(ISSUED_CERTIFICATE_CTX_KEY)
	if certificate != expectedCertificate {
		t.Fatalf("certificate is not put into context succesfully, found: %v", certificate)
	}

	privateKey := ctx.GetValue(ISSUED_PRIVATE_KEY_CTX_KEY)
	if expectedPrivateKey != expectedPrivateKey {
		t.Fatalf("private key is not put into context succesfully, found: %v", privateKey)
	}

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
			if req.ExpirationDays != 30 {
				t.Fatalf("expiration days is not correct, expected: 30, found:%d", req.ExpirationDays)
			}

			return &grpc.CertificateResponse{
				Certificate: "",
				PrivateKey:  "",
			}, nil
		})

	args := getValidArgs()
	args[ARGS_EXPIRATION_DAYS] = "30"

	err := action.Run(context.New(), args)
	if err != nil {
		t.Fatalf("action failed with error, %v", err)
	}

}

func TestExpirationDaysNotConvertableInt(t *testing.T) {
	action := NewIssueCertificateAction(nil)

	args := getValidArgs()
	args[ARGS_EXPIRATION_DAYS] = "thirty"

	err := action.Run(context.New(), args)
	if err == nil {
		t.Fatalf("error expected, but could not found")
	} else if !strings.Contains(err.Error(), "invalid syntax") {
		t.Fatalf("invalid syntax statement is expected since str to int should be failed but not found, err: %v", err)
	}

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
			if !reflect.DeepEqual(expectedSans, req.SANs) {
				t.Fatalf("SANs is not correct, expected: %s, found:%s", expectedSans, req.SANs)
			}

			return &grpc.CertificateResponse{
				Certificate: "",
				PrivateKey:  "",
			}, nil
		})

	args := getValidArgs()

	err := action.Run(context.New(), args)
	if err != nil {
		t.Fatalf("action failed with error, %v", err)
	}

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
			if !reflect.DeepEqual(expectedSans, req.SANs) {
				t.Fatalf("SANs is not correct, expected: %s, found:%s", expectedSans, req.SANs)
			}

			return &grpc.CertificateResponse{
				Certificate: "",
				PrivateKey:  "",
			}, nil
		})

	args := getValidArgs()
	args[ARGS_SANS] = "a.com"

	err := action.Run(context.New(), args)
	if err != nil {
		t.Fatalf("action failed with error, %v", err)
	}

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
			if len(req.SANs) != 0 {
				t.Fatalf("SANs is not correct, expected: empty, found:%s", req.SANs)
			}

			return &grpc.CertificateResponse{
				Certificate: "",
				PrivateKey:  "",
			}, nil
		})

	args := getValidArgs()
	delete(args, ARGS_SANS)

	err := action.Run(context.New(), args)
	if err != nil {
		t.Fatalf("action failed with error, %v", err)
	}

}

// -----

func testRequiredArgument(t *testing.T, arg string) {
	action := NewIssueCertificateAction(nil)

	args := getValidArgs()
	delete(args, arg)

	err := action.Run(context.New(), args)
	if err == nil || !strings.Contains(err.Error(), "required argument") {
		t.Fatalf("required arg error is expected but not found")
	}
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
