package searcher

import (
	"errors"
	"math/big"
	"net"

	"github.com/coreservice-io/geo_ip_v3/utils"
)

type IpVal struct {
	val string
	typ string
	num *big.Int
}

func ParseToIpVal(ip string) (*IpVal, error) {
	ip_type := ""
	target_net_ip := net.ParseIP(ip)

	if target_net_ip.To4() != nil {
		ip_type = "ipv4"
	} else if target_net_ip.To16() != nil {
		ip_type = "ipv6"
	} else {
		return nil, errors.New("ip format error:" + ip)
	}

	target_ip_score, err := utils.IpToBigInt(target_net_ip)
	if err != nil {
		return nil, err
	}

	return &IpVal{
		ip,
		ip_type,
		target_ip_score,
	}, nil
}

func ExtractBucketIdxIpv4(num *big.Int) (uint32, error) {

	ps := num.Bytes()

	d0 := uint(ps[2])
	d1 := uint(ps[3])

	var val = uint32(0)
	val |= uint32(d0) << 8
	val |= uint32(d1)

	// convert the ip to integer
	return val, nil
}

// 002::/3 . from 512 (0x0200) to 1023 (0x03FF)
func ExtractBucketIdxIpv6(num *big.Int) (uint32, error) {

	ps := num.Bytes()

	d0 := uint(ps[0])
	d1 := uint(ps[1])

	var val = uint32(0)
	val |= uint32(d0) << 8
	val |= uint32(d1)

	// convert the ip to integer
	return val, nil
}
