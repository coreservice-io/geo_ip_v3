package searcher

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
)

type CitySearcherV2 struct {
	country_ip_map map[uint32]([]SORT_COUNTRY_IP)
}

func NewCitySearcherV2() *CitySearcherV2 {
	return &CitySearcherV2{}
}

func (s *CitySearcherV2) LoadFile(file_path string) error {

	fd, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer fd.Close()

	csvReader := csv.NewReader(fd)

	country_ip_map := make(map[uint32]([]SORT_COUNTRY_IP))

	var last_record *SORT_COUNTRY_IP
	var last_bucket_idx uint32
	last_record = EmptyCountryIP
	last_bucket_idx = 0
	country_ip_map[0] = append(country_ip_map[0], *last_record)

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

		bucket_idx, _ := ExtractBucketIdx(record.Start_ip)

		if last_bucket_idx != bucket_idx {
			for idx := last_bucket_idx + 1; idx <= bucket_idx; idx++ {
				country_ip_map[idx] = append(country_ip_map[idx], *last_record)
			}
		}

		country_ip_map[bucket_idx] = append(country_ip_map[bucket_idx], *record)

		last_bucket_idx = bucket_idx
		last_record = record
	}

	//////// sort  start ip desc ///////////////////
	for _, country_ip_list := range country_ip_map {
		sort.SliceStable(country_ip_list, func(i, j int) bool {
			return country_ip_list[i].Start_ip_score.Cmp(country_ip_list[j].Start_ip_score) == 1
		})
	}

	s.country_ip_map = country_ip_map
	return nil
}

func (s *CitySearcherV2) Search(target_ip string, target_ip_score *big.Int) *SORT_COUNTRY_IP {

	idx, _ := ExtractBucketIdx(target_ip)

	if group, ok := s.country_ip_map[idx]; !ok {
		return nil
	} else {
		country_index := sort.Search(len(group), func(j int) bool {
			return group[j].Start_ip_score.Cmp(target_ip_score) <= 0
		})

		if country_index >= 0 && country_index < len(group) {
			return &(group[country_index])
		}

		return nil
	}
}

func ExtractBucketIdx(ip string) (uint32, error) {
	var ps = strings.Split(strings.TrimSpace(ip), ".")
	if len(ps) != 4 {
		return 0, fmt.Errorf("invalid ip address `%s`", ip)
	}

	d0, _ := strconv.Atoi(ps[0])
	d1, _ := strconv.Atoi(ps[1])

	var val = uint32(0)
	val |= uint32(d0) << 8
	val |= uint32(d1)

	// convert the ip to integer
	return val, nil
}
