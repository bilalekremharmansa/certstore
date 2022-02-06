package service

import (
	"context"
	b64 "encoding/base64"

	certificate_service "bilalekrem.com/certstore/internal/certificate/service"
	certstore_pac "bilalekrem.com/certstore/internal/certstore"
	grpc "bilalekrem.com/certstore/internal/certstore/grpc/gen"
	"bilalekrem.com/certstore/internal/logging"
)

type certificateService struct {
	grpc.UnimplementedCertificateServiceServer

	certstore certstore_pac.CertStore
}

func NewCertificateService(certstore certstore_pac.CertStore) *certificateService {
	return &certificateService{
		certstore: certstore,
	}
}

func (s *certificateService) IssueCertificate(_ context.Context, req *grpc.CertificateRequest) (*grpc.CertificateResponse, error) {
	certificateRequest := convertServiceRequestInternalRequest(req)

	certificateResponse, err := s.certstore.IssueCertificate(req.Issuer, certificateRequest)
	if err != nil {
		logging.GetLogger().Debugf("Error occurred while issuing certificate in grpc service, %v", err)
		return nil, err
	}

	// ---

	resp := convertInternalResponseToServiceResponse(certificateResponse)
	return resp, nil
}

// ----

func convertServiceRequestInternalRequest(req *grpc.CertificateRequest) *certificate_service.NewCertificateRequest {
	return &certificate_service.NewCertificateRequest{
		CommonName:              req.CommonName,
		Email:                   []string{req.Email},
		Organization:            []string{req.Organization},
		ExpirationDays:          int(req.ExpirationDays),
		SubjectAlternativeNames: req.SANs,
	}
}

func convertInternalResponseToServiceResponse(res *certificate_service.NewCertificateResponse) *grpc.CertificateResponse {
	certificate := b64.StdEncoding.EncodeToString(res.Certificate)
	privateKey := b64.StdEncoding.EncodeToString(res.PrivateKey)

	return &grpc.CertificateResponse{
		Certificate: certificate,
		PrivateKey:  privateKey,
	}
}
