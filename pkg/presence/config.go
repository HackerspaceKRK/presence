package presence

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

type configKeyType string

var ConfigKey configKeyType = "config"

type Config struct {

	// Where to take DHCP leases from
	DHCPSource               string        `env:"PRESENCE_DHCP_SOURCE"`
	DHCPSourceAddr           string        `env:"PRESENCE_DHCP_SOURCE_ADDR"`
	DHCPSourceUsername       string        `env:"PRESENCE_DHCP_SOURCE_USERNAME"`
	DHCPSourcePassword       string        `env:"PRESENCE_DHCP_SOURCE_PASSWORD"`
	DHCPSourcePasswordFile   string        `env:"PRESENCE_DHCP_SOURCE_PASSWORD_FILE"`
	DHCPSourceScrapeInterval time.Duration `env:"PRESENCE_DHCP_SOURCE_SCRAPE_INTERVAL" default:"1m"`

	// HTTP server
	HTTPListen string `env:"PRESENCE_HTTP_LISTEN" default:":8080"`

	PageTitle string `env:"PRESENCE_PAGE_TITLE" default:"presence"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(&cfg)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")
		if tag == "" {
			continue
		}
		value := v.Elem().FieldByName(field.Name)
		if !value.IsValid() {
			continue
		}
		strval := os.Getenv(tag)
		if strval == "" {
			strval = field.Tag.Get("default")
		}

		switch value.Kind() {
		case reflect.String:
			value.SetString(strval)
		case reflect.Int64:
			if value.Type().Name() == "Duration" {
				var durationval time.Duration
				durationval, err := time.ParseDuration(strval)
				if err != nil {
					return Config{}, fmt.Errorf("failed to parse duration from %v: %w", strval, err)
				}
				value.Set(reflect.ValueOf(durationval))
				break
			}
			var intval int
			fmt.Sscanf(strval, "%d", &intval)
			value.SetInt(int64(intval))
		case reflect.Bool:
			var boolval bool
			if strval == "true" {
				boolval = true
			}
			value.SetBool(boolval)
		}
	}
	if err := ValidateConfig(cfg); err != nil {
		return Config{}, err
	}
	log.Printf("Loaded config: %#v", cfg)
	return cfg, nil
}

func ValidateConfig(cfg Config) error {
	if cfg.DHCPSource != "openwrt" && cfg.DHCPSource != "mikrotik" {
		return fmt.Errorf("PRESENCE_DHCP_SOURCE must be either 'openwrt' or 'mikrotik'")
	}
	if cfg.DHCPSourceAddr == "" {
		return fmt.Errorf("PRESENCE_DHCP_SOURCE_ADDR must be set")
	}

	return nil
}
