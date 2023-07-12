package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

// / the second int is 32 for ipv4 or 128 for ipv6
func IpToBigInt(ip net.IP) (*big.Int, error) {
	val := big.NewInt(0)
	val.SetBytes(ip)
	return val, nil
}

func Ipv4strToBigInt(ip string) (*big.Int, error) {
	ipv := net.ParseIP(ip)
	if ipv.To4() == nil {
		return nil, errors.New(ip + " is not ipv4")
	}
	return IpToBigInt(ipv)
}

func Ipv6strToBigInt(ip string) (*big.Int, error) {
	ipv := net.ParseIP(ip)
	if ipv.To16() == nil {
		return nil, errors.New(ip + " is not ipv6")
	}
	return IpToBigInt(ipv)
}

/****************** Below Method Need To Test ******************/

var shiftIndex = []int{24, 16, 8, 0}

func ConvertIpv4(ip string) (uint32, error) {
	var ps = strings.Split(strings.TrimSpace(ip), ".")
	if len(ps) != 4 {
		return 0, fmt.Errorf("invalid ip address `%s`", ip)
	}

	var val = uint32(0)
	for i, s := range ps {
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("the %dth part `%s` is not an integer", i, s)
		}

		if d < 0 || d > 255 {
			return 0, fmt.Errorf("the %dth part `%s` should be an integer bettween 0 and 255", i, s)
		}

		val |= uint32(d) << shiftIndex[i]
	}

	// convert the ip to integer
	return val, nil
}

// INET_ATON converts an IPv4 netip.Addr object to a 64 bit integer.
func INET_ATON(ip netip.Addr) (int64, error) {
	if !ip.Is4() {
		return 0, fmt.Errorf("ip is not ipv4")
	}
	ipv4Int := big.NewInt(0)
	ipv4 := ip.As4()
	ipv4Int.SetBytes(ipv4[:])
	return ipv4Int.Int64(), nil
}

// INET6_ATON converts an IP Address (IPv4 or IPv6) netip.Addr object to a hexadecimal
// representaiton. This function is the equivalent of
// inet6_aton({{ ip address }}) in MySQL.
func INET6_ATON(ip netip.Addr) (string, error) {
	ipInt := big.NewInt(0)
	switch {
	case ip.Is4():
		ipv4 := ip.As4()
		ipInt.SetBytes(ipv4[:])
		return hex.EncodeToString(ipInt.Bytes()), nil
	case ip.Is6():
		ipv6 := ip.As16()
		ipInt.SetBytes(ipv6[:])
		return hex.EncodeToString(ipInt.Bytes()), nil
	default:
		return "", fmt.Errorf("invalid ip address")
	}
}

// INET_NTOA came from https://go.dev/play/p/JlYJXZnUxl
func INET_NTOA(ipInt64 uint32) (ip netip.Addr) {
	ipArray := [4]byte{byte(ipInt64 >> 24), byte(ipInt64 >> 16), byte(ipInt64 >> 8), byte(ipInt64)}
	ip = netip.AddrFrom4(ipArray)
	return
}

func INET6_NTOA(ipHex string) (ip netip.Addr, err error) {
	ipHex = strings.TrimPrefix(ipHex, "0x")

	HEX, err := hex.DecodeString(ipHex)
	if err != nil {
		err = fmt.Errorf("error decoding hex %w", err)
		return
	}
	if len(HEX) > 16 {
		//adding all caraters that mising of the begening
		HEX = append(make([]byte, 16-len(HEX)), HEX...)
	}

	ip, ok := netip.AddrFromSlice(HEX)
	if !ok {
		err = fmt.Errorf("invalid ip hax")
	}
	return
}
