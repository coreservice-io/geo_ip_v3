package searcher

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/coreservice-io/geo_ip_v3/data"
)

type META_CITY_INFO struct {
	meta_id      int64
	Country_code string
	Region       string
	City         string
	Latitude     float64
	Longitude    float64
}

var EmptyCityInfo = &META_CITY_INFO{
	0,
	"ZZ",
	"",
	"",
	0.000000,
	0.000000,
}

type CityMetaSearcher struct {
	data_map map[int64](*META_CITY_INFO)
}

func NewCityMetaSearcher() *CityMetaSearcher {
	return &CityMetaSearcher{}
}

func (s *CityMetaSearcher) LoadFile(file_path string) error {

	fd, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer fd.Close()

	country_ip_map := make(map[int64](*META_CITY_INFO))

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

		record, perr := line_parser_meta(lines, line_no)
		if perr != nil {
			return perr
		}

		if record == nil {
			continue
		}

		country_ip_map[record.meta_id] = record
	}

	s.data_map = country_ip_map
	return nil
}

func (s *CityMetaSearcher) Search(meta_id int64) *META_CITY_INFO {

	val, isOk := s.data_map[meta_id]
	if isOk {
		return val
	}
	return EmptyCityInfo
}

func line_parser_meta(lines []string, lineno int) (*META_CITY_INFO, error) {

	if lines[1] == "" {
		return nil, nil
	}
	if _, exist := data.CountryList[lines[1]]; !exist {
		return nil, fmt.Errorf("parser line err '%#v'", lines)
	}

	meta_id, err := strconv.ParseInt(lines[0], 10, 64)
	if err != nil {
		return nil, err
	}

	record := &META_CITY_INFO{
		meta_id:      meta_id,
		Country_code: lines[1],
		Region:       lines[2],
		City:         lines[3],
	}

	if lines[4] == "" || lines[5] == "" {
		record.Latitude = 0
		record.Longitude = 0
	} else {
		if lati, err := strconv.ParseFloat(lines[4], 64); err != nil {
			return nil, fmt.Errorf("parser line err '%#v':%d. Err: %s", lines, lineno, err.Error())
		} else {
			record.Latitude = lati
		}

		if longti, err := strconv.ParseFloat(lines[5], 64); err != nil {
			return nil, fmt.Errorf("parser line err '%#v':%d. Err: %s", lines, lineno, err.Error())
		} else {
			record.Longitude = longti
		}
	}

	return record, nil
}
