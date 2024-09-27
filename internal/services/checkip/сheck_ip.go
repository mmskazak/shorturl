package checkip

import (
	"fmt"
	"net"
)

func CheckIPByCIDR(ip, cidr string) (bool, error) {
	ipNet, err := parseCIDR(cidr)
	if err != nil {
		return false, fmt.Errorf("parse cidr error: %w", err)
	}

	clientIP := net.ParseIP(ip)
	if clientIP != nil && ipNet.Contains(clientIP) {
		return true, nil
	}
	return false, nil
}

// parseCIDR парсит CIDR строку и возвращает объект net.IPNet.
func parseCIDR(cidr string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("parse cidr %s error: %w", cidr, err)
	}
	return ipNet, nil
}
