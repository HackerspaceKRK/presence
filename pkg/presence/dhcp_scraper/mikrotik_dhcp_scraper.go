package dhcpscraper

import (
	"fmt"
	"time"

	"github.com/go-routeros/routeros"
)

type MikrotikDHCPScraper struct {
	host     string
	username string
	password string
}

func (s *MikrotikDHCPScraper) Setup(host string, username string, password string) error {
	s.host = host
	s.username = username
	s.password = password
	return nil
}

func (s *MikrotikDHCPScraper) Scrape() ([]DHCPLease, error) {
	conn, err := routeros.Dial(s.host, s.username, s.password)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Mikrotik at %v: %w", s.host, err)
	}
	defer conn.Close()

	reply, err := conn.Run("/ip/dhcp-server/lease/print")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DHCP leases from Mikrotik at %v: %w", s.host, err)
	}
	leases := make([]DHCPLease, 0, len(reply.Re))
	for _, re := range reply.Re {

		age, _ := time.ParseDuration(re.Map["age"])
		expiresAfter, _ := time.ParseDuration(re.Map["expires-after"])
		lastSeen, _ := time.ParseDuration(re.Map["last-seen"])
		createdAt := time.Now().Add(-age)
		expiresAt := time.Now().Add(expiresAfter)

		lease := DHCPLease{
			IP:       re.Map["address"],
			MAC:      re.Map["mac-address"],
			Hostname: re.Map["host-name"],

			CreatedAt: createdAt,
			ExpiresAt: expiresAt,
			Present:   lastSeen < time.Minute*5,
		}
		leases = append(leases, lease)

	}

	return leases, nil
}
