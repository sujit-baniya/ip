package main

import (
	"fmt"
	"github.com/sujit-baniya/ip"
	"time"
)

func main() {
	start := time.Now()
	db := ip.NewGeoIpDB("./assets/geoip/GeoLite2-City.mmdb")
	fmt.Println(db.GetLocation("110.44.127.177"))
	fmt.Printf("%s", time.Since(start))
}
