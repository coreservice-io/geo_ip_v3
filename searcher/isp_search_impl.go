package searcher

import (
	"encoding/csv"
	"io"
	"math/big"
	"os"
	"sort"
	"strings"
)

type IspSearcherImpl struct {
	isp_ip_list []SORT_ISP_IP
}

func NewIspSearcherImpl() *IspSearcherImpl {
	return &IspSearcherImpl{}
}

func (s *IspSearcherImpl) LoadFile(file_path string) error {

	fd, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer fd.Close()

	isp_ip_list := []SORT_ISP_IP{}

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

		record, perr := line_parser_isp(lines, line_no)
		if perr != nil {
			return perr
		}
		if record == nil {
			continue
		}

		isp_ip_list = append(isp_ip_list, *record)
	}

	//////// sort  start ip desc ///////////////////
	sort.SliceStable(isp_ip_list, func(i, j int) bool {
		return isp_ip_list[i].Start_ip_score.Cmp(isp_ip_list[j].Start_ip_score) == 1
	})

	// if len(isp_ip_list) == 0 {
	// 	return errors.New("isp_ipv6 len :0 ")
	// }

	s.isp_ip_list = isp_ip_list
	return nil
}

func (s *IspSearcherImpl) Search(target_ip_score *big.Int) *SORT_ISP_IP {

	country_index := sort.Search(len(s.isp_ip_list), func(j int) bool {
		return s.isp_ip_list[j].Start_ip_score.Cmp(target_ip_score) <= 0
	})

	if country_index >= 0 && country_index < len(s.isp_ip_list) {
		return &(s.isp_ip_list[country_index])
	}

	return nil
}

func line_parser_isp(lines []string, lineno int) (*SORT_ISP_IP, error) {

	network := lines[0]
	ipint, err := ParseToIpVal(network)
	if err != nil {
		return nil, err
	}

	record := &SORT_ISP_IP{
		Start_ip:       network,
		Start_ip_score: ipint.num,
		Asn:            lines[1],
		Isp:            lines[3],
	}

	if strings.Trim(lines[2], " ") == "1" {
		record.Is_datacenter = true
	} else {
		record.Is_datacenter = false
	}

	return record, nil
}
