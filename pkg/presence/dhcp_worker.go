package presence

import (
	"context"
	"fmt"
	"log"
	"time"

	dhcpscraper "github.com/HackerspaceKRK/presence/pkg/presence/dhcp_scraper"
)

type dhcpWorkerKey string

var DHCPWorkerKey dhcpWorkerKey = "dhcpWorker"

type DHCPWorker struct {
	NumFailedScrapes int
	LastError        error
	CurrentLeases    []dhcpscraper.DHCPLease
	Scraper          dhcpscraper.DHCPScraper
}

// WithDHCPWorker creates a new context with DHCPWorker attached to it. And starts the worker.
func WithDHCPWorker(ctx context.Context) context.Context {
	var scraper dhcpscraper.DHCPScraper
	cfg := ctx.Value(ConfigKey).(Config)
	switch cfg.DHCPSource {
	case "mikrotik":
		scraper = &dhcpscraper.MikrotikDHCPScraper{}
	case "openwrt":
		panic("OpenWRT DHCP scraper not implemented")

	}
	scraper.Setup(cfg.DHCPSourceAddr, cfg.DHCPSourceUsername, cfg.DHCPSourcePassword)
	w := &DHCPWorker{
		Scraper: scraper,
	}
	w.NumFailedScrapes = 999
	w.LastError = fmt.Errorf("DHCP scraper is starting up")
	go w.work(ctx)
	return context.WithValue(ctx, DHCPWorkerKey, w)
}

func (w *DHCPWorker) work(ctx context.Context) {
	var firstRun = true
	for {
		cfg := ctx.Value(ConfigKey).(Config)
		select {
		case <-ctx.Done():
			return
		default:
			leases, err := w.Scraper.Scrape()
			if err != nil {
				log.Printf("failed to scrape DHCP leases: %v", err)
				w.NumFailedScrapes++
				w.LastError = err
				time.Sleep(cfg.DHCPSourceScrapeInterval)
				continue
			}
			if firstRun {
				log.Printf("DHCP scraper OK, scraped %d DHCP leases", len(leases))
				firstRun = false
			}
			w.NumFailedScrapes = 0
			w.LastError = nil
			filteredLeases := make([]dhcpscraper.DHCPLease, 0, len(leases))
			for _, lease := range leases {
				if lease.Present {
					filteredLeases = append(filteredLeases, lease)
				}
			}
			w.CurrentLeases = filteredLeases
			time.Sleep(cfg.DHCPSourceScrapeInterval)
		}
	}
}

func (w *DHCPWorker) CountByVendor() map[string]int {
	counts := map[string]int{}
	for _, lease := range w.CurrentLeases {
		name := LookupOuiByMAC(lease.MAC)
		counts[name]++
	}
	return counts
}

func (w *DHCPWorker) ConnectionOK() bool {
	return w.NumFailedScrapes <= 1
}
