package searcher

import (
	"fmt"
	"math/big"
)

const NEW_SEARCHER_IPV4 = true

type COUNTRY_IP_INFO struct {
	Start_ip     string
	Country_code string
	Region       string
	City         string
	Latitude     float64
	Longitude    float64
}

type SORT_COUNTRY_IP struct {
	Start_ip       string
	Start_ip_score *big.Int
	meta_id        int64
}

var EmptyCountryIP = &SORT_COUNTRY_IP{
	"0.0.0.0",
	big.NewInt(0),
	0,
}

type CitySearcher struct {
	country_ipv4_searcher *CitySearcherV1
	country_ipv6_searcher *CitySearcherV1
	ipv4_searcher         *CitySearcherV2

	meta_searcher *CityMetaSearcher
}

func NewCitySearcher() *CitySearcher {
	return &CitySearcher{}
}

func (s *CitySearcher) Init(country_ipv4_path string, country_ipv6_path string, country_meta_path string,
	logger func(log_str string), err_logger func(err_log_str string)) error {

	////
	if !NEW_SEARCHER_IPV4 {

		country_ipv4_searcher := NewCountrySearcher()

		if err := country_ipv4_searcher.LoadFile(country_ipv4_path); err != nil {
			return err
		} else {
			s.country_ipv4_searcher = country_ipv4_searcher
		}
	} else {
		ipv4_searcher := NewCitySearcherV2()

		if err := ipv4_searcher.LoadFile(country_ipv4_path); err != nil {
			return err
		} else {
			s.ipv4_searcher = ipv4_searcher
		}
	}

	///
	country_ipv6_searcher := NewCountrySearcher()

	if err := country_ipv6_searcher.LoadFile(country_ipv6_path); err != nil {
		return err
	} else {
		s.country_ipv6_searcher = country_ipv6_searcher
	}

	meta_searcher := NewCityMetaSearcher()

	if err := meta_searcher.LoadFile(country_meta_path); err != nil {
		return err
	} else {
		s.meta_searcher = meta_searcher
	}

	return nil
}

func (s *CitySearcher) SearchVal(ipVal *IpVal) (*COUNTRY_IP_INFO, error) {
	fmt.Printf("old %s %x\n", ipVal.val, ipVal.num.Bytes())

	// actually do
	search_country := s.country_ipv4_searcher

	if ipVal.typ == "ipv6" {
		search_country = s.country_ipv6_searcher
	}

	var country_info *SORT_COUNTRY_IP
	if NEW_SEARCHER_IPV4 && ipVal.typ == "ipv4" {
		ipv4_searcher := s.ipv4_searcher
		country_info = ipv4_searcher.Search(ipVal.val, ipVal.num)
	} else {
		country_info = search_country.Search(ipVal.num)
	}

	if country_info == nil {
		return nil, nil
	}

	meta_searcher := s.meta_searcher
	meta_info := meta_searcher.Search(country_info.meta_id)

	return &COUNTRY_IP_INFO{
		Start_ip:     country_info.Start_ip,
		Country_code: meta_info.Country_code,
		Region:       meta_info.Region,
		City:         meta_info.City,
		Latitude:     meta_info.Latitude,
		Longitude:    meta_info.Longitude,
	}, nil
}

func (s *CitySearcher) Search(target_ip string) (*COUNTRY_IP_INFO, error) {

	ipVal, err := ParseToIpVal(target_ip)
	if err != nil {
		return nil, err
	}

	return s.SearchVal(ipVal)
}
