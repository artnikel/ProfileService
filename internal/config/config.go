// Package config with environment variables
package config

import "github.com/caarlos0/env"

// Variables is a struct with environment variables
type Variables struct {
	PostgresConnProfile string `env:"POSTGRES_CONN_PROFILE"`
	ProfileAddress      string `env:"PROFILE_ADDRESS"`
}

// New returns parsed object of config
func New() (*Variables, error) {
	cfg := &Variables{}
	err := env.Parse(cfg)
	return cfg, err
}
