package main

import (
	"fmt"
	"github.com/sujit-baniya/ip"
)

func main() {
	db := ip.NewGeoIpDB("./assets/geoip/GeoLite2-City.mmdb")
	fmt.Println(db.GetLocation("110.44.127.177"))
}
