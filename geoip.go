package ip

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/oschwald/maxminddb-golang"
	"net"
	"regexp"
)

// See https://pkg.go.dev/github.com/oschwald/geoip2-golang#City for a full list of options you can use here to modify
// what data is returned for a specific IP.
type ipLookup struct {
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		IsoCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Location struct {
		TimeZone string `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

type Response struct {
	City     string `json:"city,omitempty"`
	IP       string `json:"ip"`
	Country  string `json:"country"`
	IsoCode  string `json:"iso_code"`
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
		IP:       ip,
		IsoCode:  record.Country.IsoCode,
		Timezone: record.Location.TimeZone,
	}
	if val, ok := record.Country.Names["en"]; ok {
		response.Country = val
	}
	if val, ok := record.City.Names["en"]; ok {
		response.City = val
	}
	return response, nil
}

var fetchIpFromString = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
var possibleHeaderes = []string{
	"X-Original-Forwarded-For",
	"X-Forwarded-For",
	"X-Real-Ip",
	"X-Client-Ip",
	"Forwarded-For",
	"Forwarded",
	"Remote-Addr",
	"Client-Ip",
	"CF-Connecting-IP",
}

// determine user ip
func IP(c *fiber.Ctx) string {
	var headerValue []byte
	if c.App().Config().ProxyHeader == "*" {
		for _, headerName := range possibleHeaderes {
			headerValue = c.Request().Header.Peek(headerName)
			if len(headerValue) > 3 {
				return string(fetchIpFromString.Find(headerValue))
			}
		}
	}
	headerValue = []byte(c.IP())
	if len(headerValue) <= 3 {
		headerValue = []byte("0.0.0.0")
	}

	// find ip address in string
	return string(fetchIpFromString.Find(headerValue))
}

func Detect(c *fiber.Ctx) error {
	ip := IP(c)
	c.Locals("ip", ip)
	return c.Next()
}

func DetectLocation(db *GeoIpDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := IP(c)
		response, _ := db.GetLocation(ip)
		c.Locals("ip", ip)
		c.Locals("location", response)
		return c.Next()
	}
}
