package utils

import (
	"net/netip"
	"testing"
)

func TestINET_ATON(t *testing.T) {
	ipv4, _ := netip.ParseAddr("192.168.1.1")
	aton, err := INET_ATON(ipv4)
	if err != nil {
		t.Errorf("err: %s", err.Error())
		return
	}

	t.Logf("number: %d", aton)
}

func TestINET6_ATON(t *testing.T) {
	ipv4, _ := netip.ParseAddr("11:22:33:44:55:66:77:88")
	aton, err := INET6_ATON(ipv4)
	if err != nil {
		t.Errorf("err: %s", err.Error())
		return
	}

	t.Logf("hex: 0x%s", aton)
}

func TestINET_NTOA(t *testing.T) {
	ntoa := INET_NTOA(3232235777)
	t.Logf("ipv4: %s", ntoa.String())
}

func TestINET6_NTOA(t *testing.T) {
	ntoa, err := INET6_NTOA("0x110022003300440055006600770088")
	if err != nil {
		t.Errorf("err: %s", err.Error())
		return
	}

	t.Logf("ipv6: %s", ntoa.String())
}
