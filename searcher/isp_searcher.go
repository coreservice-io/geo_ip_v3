package searcher

import (
	"math/big"
)

type SORT_ISP_IP struct {
	Start_ip       string
	Start_ip_score *big.Int
	Asn            string
	Is_datacenter  bool
	Isp            string
}

type IspSearcher struct {
	isp_ipv4_searcher *IspSearcherImpl
	isp_ipv6_searcher *IspSearcherImpl
}

func NewIspSearcher() *IspSearcher {
	return &IspSearcher{}
}

func (s *IspSearcher) Init(ipv4_path string, ipv6_path string,
	logger func(log_str string), err_logger func(err_log_str string)) error {

	///
	isp_ipv4_searcher := NewIspSearcherImpl()

	if err := isp_ipv4_searcher.LoadFile(ipv4_path); err != nil {
		return err
	} else {
		s.isp_ipv4_searcher = isp_ipv4_searcher
	}
	///
	isp_ipv6_searcher := NewIspSearcherImpl()

	if err := isp_ipv6_searcher.LoadFile(ipv6_path); err != nil {
		return err
	} else {
		s.isp_ipv6_searcher = isp_ipv6_searcher
	}

	return nil
}

func (s *IspSearcher) SearchVal(ipVal *IpVal) (*SORT_ISP_IP, error) {

	//////////////
	search_isp := s.isp_ipv4_searcher

	if ipVal.typ == "ipv6" {
		search_isp = s.isp_ipv6_searcher
	}

	isp_info := search_isp.Search(ipVal.num)

	return isp_info, nil
}

func (s *IspSearcher) Search(target_ip string) (*SORT_ISP_IP, error) {

	ipVal, err := ParseToIpVal(target_ip)
	if err != nil {
		return nil, err
	}

	return s.SearchVal(ipVal)
}
