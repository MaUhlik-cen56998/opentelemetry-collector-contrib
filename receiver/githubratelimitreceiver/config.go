package githubratelimitreceiver

import (
	"errors"
)

var (
	errMissingEndpoint = errors.New(`"endpoint" must be specified`)
	errInvalidEndpoint = errors.New(`"endpoint" must be in the form of <scheme>://<hostname>[:<port>]`)
)

type Config struct {
	Endpoint       string `mapstructure:"endpoint"`
	Token          string `mapstructure:"token"`
	Name           string `mapstructure:"name"`
	Target         string `mapstructure:"target"`
	LogLevel       string `mapstructure:"log_level"`
	ScrapeInterval int    `mapstructure:"scrape_interval"`
}

func (cfg *Config) Validate() error {
	if cfg.Token == "" {
		return errors.New("GitHub token is required")
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://api.github.com/rate_limit"
	}
	if cfg.Name == "" {
		return errors.New("Name of token is required")
	}
	if cfg.Target == "" {
		return errors.New("GitHub target is required")
	}
	if cfg.ScrapeInterval <= 0 {
		cfg.ScrapeInterval = 60
	}

	return nil
}
