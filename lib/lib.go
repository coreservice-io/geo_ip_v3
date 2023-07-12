package lib

import (
	"errors"
	"path/filepath"

	"github.com/coreservice-io/package_client"

	"github.com/coreservice-io/geo_ip_v3/data"
	"github.com/coreservice-io/geo_ip_v3/searcher"
	"github.com/coreservice-io/geo_ip_v3/utils"
)

type GeoIpClient struct {
	city_searcher *searcher.CitySearcher
	isp_searcher  *searcher.IspSearcher
	pc            *package_client.PackageClient
}

func (geoip_c *GeoIpClient) ReloadCsv(datafolder string,
	logger func(log_str string), err_logger func(err_log_str string)) error {

	country_ipv4_file_abs := filepath.Join(datafolder, "ipv4_city_data.csv")
	country_ipv6_file_abs := filepath.Join(datafolder, "ipv6_city_data.csv")
	country_meta_file_abs := filepath.Join(datafolder, "ip_loc_data.csv")

	isp_ipv4_file_abs := filepath.Join(datafolder, "ipv4_isp_data.csv")
	isp_ipv6_file_abs := filepath.Join(datafolder, "ipv6_isp_data.csv")

	////
	city_searcher := searcher.NewCitySearcher()

	if err := city_searcher.Init(country_ipv4_file_abs, country_ipv6_file_abs, country_meta_file_abs,
		logger, err_logger); err != nil {
		return err
	} else {
		geoip_c.city_searcher = city_searcher
	}

	///
	isp_searcher := searcher.NewIspSearcher()

	if err := isp_searcher.Init(isp_ipv4_file_abs, isp_ipv6_file_abs, logger, err_logger); err != nil {
		return err
	} else {
		geoip_c.isp_searcher = isp_searcher
	}

	return nil
}

func NewClient(update_key string, current_version string, datafolder string, ignore_data_exist bool,
	logger func(log_str string), err_logger func(err_log_str string)) (GeoIpInterface, error) {

	client := &GeoIpClient{}
	if !ignore_data_exist {
		load_err := client.ReloadCsv(datafolder, logger, err_logger)
		if load_err != nil {
			logger("load_err:" + load_err.Error())
			return nil, load_err
		}
	}
	///
	// pc, err := StartAutoUpdate(update_key, current_version, false, datafolder, func() {
	// 	client.ReloadCsv(datafolder, logger, err_logger)
	// }, logger, err_logger)

	// if err != nil {
	// 	logger("StartAutoUpdate err:" + err.Error())
	// }

	// client.pc = pc
	////////////////////////
	return client, nil
}

func (geoip_c *GeoIpClient) Upgrade(ignore_version bool) error {
	return geoip_c.pc.Update(ignore_version)
}

func (geoip_c *GeoIpClient) GetInfo(target_ip string) (*GeoInfo, error) {

	// pre check ip
	if isLan, err := utils.IsLanIp(target_ip); err != nil {
		return nil, err
	} else if isLan {
		return nil, errors.New("is lan ip")
	}

	ipVal, err := searcher.ParseToIpVal(target_ip)
	if err != nil {
		return nil, err
	}

	//////////////
	////
	result := GeoInfo{
		Ip:             target_ip,
		Latitude:       0,
		Longitude:      0,
		Country_code:   data.NA,
		Country_name:   data.NA,
		Continent_code: data.NA,
		Continent_name: data.NA,
		Region:         data.NA,
		City:           data.NA,
		Asn:            data.NA,
		Isp:            data.NA,
		Is_datacenter:  false,
	}

	city_searcher := geoip_c.city_searcher
	country_info, err := city_searcher.SearchVal(ipVal)
	if err != nil {
		return nil, err
	}

	isp_search := geoip_c.isp_searcher
	isp_info, err := isp_search.SearchVal(ipVal)
	if err != nil {
		return nil, err
	}

	if country_info != nil {
		fillGeoInfo(&result, country_info)
	}

	if isp_info != nil {
		fillIspInfo(&result, isp_info)
	}

	return &result, nil
}

func fillGeoInfo(result *GeoInfo, info *searcher.COUNTRY_IP_INFO) {
	result.Latitude = info.Latitude
	result.Longitude = info.Longitude
	result.Country_code = info.Country_code
	result.Region = info.Region
	result.City = info.City

	if val, ok := data.CountryList[result.Country_code]; ok {
		result.Continent_code = val.ContinentCode
		result.Continent_name = val.ContinentName
		result.Country_name = val.CountryName
	}
}

func fillIspInfo(result *GeoInfo, info *searcher.SORT_ISP_IP) {
	result.Asn = info.Asn
	result.Isp = info.Isp
	result.Is_datacenter = info.Is_datacenter
}
