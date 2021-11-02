package main

import (
	"bilalekrem.com/certstore/internal/certificate/service"
)

func main() {
	caCertService := service.CACertificateService{}

	request := &service.NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response, _ := caCertService.CreateCertificate(request)
	println(string(response.Certificate))
	println(string(response.PrivateKey))
}
