// Package config with environment variables
package config

// Variables is a struct with environment variables
type Variables struct {
	PostgresConnProfile string `env:"POSTGRES_CONN_PROFILE"`
	TokenSignature      string `env:"TOKEN_SIGNATURE"`
}
