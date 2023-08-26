package dhcpscraper

import "time"

type DHCPLease struct {
	IP        string
	MAC       string
	Hostname  string
	CreatedAt time.Time
	ExpiresAt time.Time
	Present   bool // Whether the device is still present in the space (Uses last-seen on mikrotik)
}

type DHCPScraper interface {
	Setup(host string, username string, password string) error
	Scrape() ([]DHCPLease, error)
}
