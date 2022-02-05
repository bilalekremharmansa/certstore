package windns

import (
	"fmt"
	"strings"

	"bilalekrem.com/certstore/internal/logging"
	"github.com/go-acme/lego/v4/challenge/dns01"
)

// Windows Dns Provider
type winDNSProvider struct {
}

func NewWinDnsProvider() *winDNSProvider {
	return &winDNSProvider{}
}

func (d *winDNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	name, zone, err := getDnsNameAndZoneByFqdn(fqdn)
	if err != nil {
		return err
	}

	// ----

	err = addTxtRecord(zone, name, value)
	if err != nil {
		return err
	}

	return nil
}

func (d *winDNSProvider) CleanUp(domain, token, keyAuth string) error {
	fqdn, _ := dns01.GetRecord(domain, keyAuth)
	name, zone, err := getDnsNameAndZoneByFqdn(fqdn)
	if err != nil {
		return err
	}

	// ----

	err = removeTxtRecord(zone, name)
	if err != nil {
		return err
	}

	return nil
}

// ----

// return dns name and zone for windows dns server
// for a fqdn: live.certstore.com.
// zone: cerstore.com
// name: live
func getDnsNameAndZoneByFqdn(fqdn string) (string, string, error) {
	zone, err := dns01.FindZoneByFqdn(fqdn)
	if err != nil {
		return "", "", err
	}

	dnsNameIndex := strings.LastIndex(fqdn, zone)
	dnsName := fqdn[:dnsNameIndex]
	logging.GetLogger().Infof("zone: %s, name: %s", zone, dnsName)

	// subtract one because, win dns do net expect . at end
	dnsName = dnsName[:len(dnsName)-1]
	zone = zone[:len(zone)-1]
	return dnsName, zone, nil
}

func addTxtRecord(zone string, name string, txt string) error {
	logging.GetLogger().Infof("Adding txt record name %s, zone %s, value %s", name, zone, txt)

	cmd := fmt.Sprintf("Add-DnsServerResourceRecord -AllowUpdateAny -Txt -ZoneName %s -Name %s -DescriptiveText %s", zone, name, txt)
	_, err := runPowershellCmd(cmd)
	return err
}

func removeTxtRecord(zone string, name string) error {
	logging.GetLogger().Infof("Removing txt record name %s, zone %s", name, zone)

	cmd := fmt.Sprintf("Remove-DnsServerResourceRecord -Force -RRType Txt -ZoneName %s -Name %s", zone, name)
	_, err := runPowershellCmd(cmd)
	return err
}
