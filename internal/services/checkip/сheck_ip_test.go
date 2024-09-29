package checkip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIPByCIDR_ValidIPInRange(t *testing.T) {
	ip := "192.168.1.5"
	cidr := "192.168.1.0/24"
	want := true

	got, err := CheckIPByCIDR(ip, cidr)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestCheckIPByCIDR_ValidIPNotInRange(t *testing.T) {
	ip := "192.168.2.5"
	cidr := "192.168.1.0/24"
	want := false

	got, err := CheckIPByCIDR(ip, cidr)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestCheckIPByCIDR_InvalidCIDR(t *testing.T) {
	ip := "192.168.1.5"
	cidr := "192.168.1.0/33"

	_, err := CheckIPByCIDR(ip, cidr)
	assert.Error(t, err)
}

func TestCheckIPByCIDR_InvalidIP(t *testing.T) {
	ip := "invalid-ip"
	cidr := "192.168.1.0/24"
	want := false

	got, err := CheckIPByCIDR(ip, cidr)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestParseCIDR_ValidCIDR(t *testing.T) {
	cidr := "192.168.1.0/24"
	_, expectedNet, _ := net.ParseCIDR(cidr)

	ipNet, err := parseCIDR(cidr)
	assert.NoError(t, err)
	assert.Equal(t, expectedNet, ipNet)
}

func TestParseCIDR_InvalidCIDR(t *testing.T) {
	cidr := "192.168.1.0/33"

	ipNet, err := parseCIDR(cidr)
	assert.Error(t, err)
	assert.Nil(t, ipNet)
	assert.Contains(t, err.Error(), "parse cidr")
}
