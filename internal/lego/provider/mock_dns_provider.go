package provider

import (
	"bilalekrem.com/certstore/internal/logging"
	"github.com/go-acme/lego/challenge/dns01"
)

type MockDNSProvider struct {
}

func NewMockDNSProvider() *MockDNSProvider {
	return &MockDNSProvider{}
}

func (d *MockDNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	logging.GetLogger().Infof("received fqdn %s, and value %s", fqdn, value)
	return nil
}

func (d *MockDNSProvider) CleanUp(domain, token, keyAuth string) error {
	// clean up any state you created in Present, like removing the TXT record
	return nil
}
