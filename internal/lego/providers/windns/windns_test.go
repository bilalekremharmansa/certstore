package windns

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
)

func TestGetNameAndZone(t *testing.T) {
	fqdn := "_acme-challenge.live.certstore.com."

	name, zone, err := getDnsNameAndZoneByFqdn(fqdn)
	assert.NotError(t, err, "failed finding name and zone")

	assert.Equal(t, "_acme-challenge.live", name)
	assert.Equal(t, "certstore.com", zone)
}

func TestGetNameAndZoneNotSubdomain(t *testing.T) {
	fqdn := "_acme-challenge.certstore.com."

	name, zone, err := getDnsNameAndZoneByFqdn(fqdn)
	assert.NotError(t, err, "failed finding name and zone")

	assert.Equal(t, "_acme-challenge", name)
	assert.Equal(t, "certstore.com", zone)
}

func TestGetNameAndZoneWildcard(t *testing.T) {
	fqdn := "_acme-challenge.*.certstore.com."

	name, zone, err := getDnsNameAndZoneByFqdn(fqdn)
	assert.NotError(t, err, "failed finding name and zone")

	assert.Equal(t, "_acme-challenge.*", name)
	assert.Equal(t, "certstore.com", zone)
}
