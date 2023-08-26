package main

import (
	"flag"
	"fmt"

	dhcpscraper "github.com/HackerspaceKRK/presence/pkg/presence/dhcp_scraper"
)

var (
	address  = flag.String("address", "127.0.0.1:8728", "RouterOS address and port")
	username = flag.String("username", "admin", "User name")
	password = flag.String("password", "admin", "Password")
)

func main() {

	flag.Parse()
	scrap := &dhcpscraper.MikrotikDHCPScraper{}
	scrap.Setup(*address, *username, *password)
	res, err := scrap.Scrape()
	if err != nil {
		panic(err)
	}
	for _, r := range res {
		fmt.Printf("%#v\n", r)
	}

}
