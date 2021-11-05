package service

import (
	"errors"
	"fmt"
	"net/mail"
)

func validateCertificateRequest(req *NewCertificateRequest) error {
	if req.CommonName == "" {
		return errors.New("Validation error: common name can not be empty")
	}

	for _, email := range req.Email {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return errors.New(fmt.Sprintf("Validation error: email is not valid: [%s]", email))
		}
	}

	if req.ExpirationDays < 1 {
		return errors.New("Validation error: expiration days must be bigger than 1")
	}

	return nil
}
