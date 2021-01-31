package ip

import (
	"errors"
	"fmt"
	"github.com/oschwald/maxminddb-golang"
	"net"
)

// See https://pkg.go.dev/github.com/oschwald/geoip2-golang#City for a full list of options you can use here to modify
// what data is returned for a specific IP.
type ipLookup struct {
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	Location struct {
		TimeZone string `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

type Response struct {
	City     string `json:"city,omitempty"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
}

type GeoIpDB struct {
	*maxminddb.Reader
}

func NewGeoIpDB(fileName string) *GeoIpDB {
	db, err := maxminddb.Open(fileName)
	if err != nil {
		fmt.Println("Unable to load 'GeoLite2-City.mmdb'.")
		panic(err)
	}
	return &GeoIpDB{
		Reader: db,
	}
}

func (g *GeoIpDB) GetLocation(ip string) (*Response, error) {
	// Check IP address format
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil, errors.New("Invalid IP address")
	}

	// Perform lookup
	record := new(ipLookup)
	err := g.Lookup(ipAddr, &record)
	if err != nil {
		return nil, err
	}
	response := &Response{
		Country:  record.Country.IsoCode,
		Timezone: record.Location.TimeZone,
	}
	if val, ok := record.City.Names["en"]; ok {
		response.City = val
	}
	return response, nil
}
