package service

import (
	"errors"
	"net/mail"
)

func validateCertificateRequest(req *NewCertificateRequest) error {
	if req.CommonName == "" {
		return errors.New("Validation error: common name can not be empty")
	}

	if req.Email != "" {
		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			return errors.New("Validation error: email is not valid")
		}
	}

	if req.ExpirationDays < 1 {
		return errors.New("Validation error: expiration days must be bigger than 1")
	}

	return nil
}
