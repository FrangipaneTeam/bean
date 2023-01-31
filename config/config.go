// Package config provides a simple way to load configuration files
package config

import "os"

// Provider is the configuration provider.
type Provider struct {
	Path       string
	Version    string
	NewVersion string
	Debug      bool
	Logger     *os.File
}
