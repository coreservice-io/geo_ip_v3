package searcher

import (
	"encoding/csv"
	"io"
	"math/big"
	"os"
	"sort"
	"strconv"
)

type CitySimpleSearcher struct {
	country_ip_list []SORT_COUNTRY_IP
}

func NewCountrySearcher() *CitySimpleSearcher {
	return &CitySimpleSearcher{}
}

func (s *CitySimpleSearcher) LoadFile(file_path string) error {

	fd, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer fd.Close()

	country_ip_list := []SORT_COUNTRY_IP{}

	csvReader := csv.NewReader(fd)

	line_no := 0
	for {
		lines, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line_no = line_no + 1

		record, perr := line_parser_ip(lines, line_no)
		if perr != nil {
			return perr
		}

		if record == nil {
			continue
		}

		country_ip_list = append(country_ip_list, *record)
	}

	//////// sort  start ip desc ///////////////////
	sort.SliceStable(country_ip_list, func(i, j int) bool {
		return country_ip_list[i].Start_ip_score.Cmp(country_ip_list[j].Start_ip_score) == 1
	})

	// if len(country_ip_list) == 0 {
	// 	return errors.New("country_ipv4 len :0 ")
	// }

	s.country_ip_list = country_ip_list
	return nil
}

func (s *CitySimpleSearcher) Search(target_ip_score *big.Int) *SORT_COUNTRY_IP {

	// idx, _ := ExtractBucketIdxIpv6(target_ip_score)
	// fmt.Printf("idx: %d (%x)\n", idx, idx)

	c_len := len(s.country_ip_list)
	country_index := sort.Search(c_len, func(j int) bool {
		return s.country_ip_list[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	if country_index >= 0 && country_index < c_len {
		return &(s.country_ip_list[country_index])
	}

	return nil
}

func line_parser_ip(lines []string, lineno int) (*SORT_COUNTRY_IP, error) {

	// network, extid
	network := lines[0]
	if network == "" {
		return nil, nil
	}

	/////////////
	ipint, err := ParseToIpVal(network)
	if err != nil {
		return nil, err
	}

	meta_id, err := strconv.ParseInt(lines[1], 10, 64)
	if err != nil {
		return nil, err
	}

	record := &SORT_COUNTRY_IP{
		Start_ip:       network,
		Start_ip_score: ipint.num,
		meta_id:        meta_id,
	}

	return record, nil
}
