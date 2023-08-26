package presence

import (
	"fmt"
	"os"
	"reflect"
)

type configKeyType string

var ConfigKey configKeyType = "config"

type Config struct {

	// Where to take DHCP leases from
	DHCPSource             string `env:"PRESENCE_DHCP_SOURCE"`
	DHCPSourceHost         string `env:"PRESENCE_DHCP_SOURCE_HOST"`
	DHCPSourceUsername     string `env:"PRESENCE_DHCP_SOURCE_USERNAME"`
	DHCPSourcePassword     string `env:"PRESENCE_DHCP_SOURCE_PASSWORD"`
	DHCPSourcePasswordFile string `env:"PRESENCE_DHCP_SOURCE_PASSWORD_FILE"`

	// HTTP server
	HTTPListen string `env:"PRESENCE_HTTP_LISTEN"`

	PageTitle string `env:"PRESENCE_PAGE_TITLE"`
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
		value.SetString(os.Getenv(tag))
	}
	if err := ValidateConfig(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func ValidateConfig(cfg Config) error {
	if cfg.DHCPSource != "openwrt" && cfg.DHCPSource != "mikrotik" {
		return fmt.Errorf("PRESENCE_DHCP_SOURCE must be either 'openwrt' or 'mikrotik'")
	}
	if cfg.DHCPSourceHost == "" {
		return fmt.Errorf("PRESENCE_DHCP_SOURCE_HOST must be set")
	}
	if cfg.PageTitle == "" {
		cfg.PageTitle = "at"
	}
	return nil
}
