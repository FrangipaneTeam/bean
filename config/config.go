// Package config provides a simple way to load configuration files
package config

// Provider is the configuration provider.
type Provider struct {
	Path       string
	Version    string
	NewVersion string
	Debug      bool
}
