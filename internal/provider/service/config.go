package service

import (
	"fmt"

	"github.com/mr-chelyshkin/go-gridder/internal/sys"
)

type config struct {
	// set port for HTTP serving service.
	PortHTTP int `json:"port_http" yaml:"port_http" toml:"port_http" mapstructure:"port_http"`

	// set port for HTTPS serving service.
	PortHTTPS int `json:"port_https" yaml:"port_https" toml:"port_http" mapstructure:"port_https"`

	// gracefully service shutdown timeout.
	ShutdownTimeoutSec int `json:"shutdown_timeout_sec" yaml:"shutdown_timeout_sec" toml:"shutdown_timeout_sec" mapstructure:"shutdown_timeout_sec"`

	// set path for 'pem' file for https serving.
	PemPath string `json:"pem_path" yaml:"pem_path" toml:"pem_path" mapstructure:"pem_path"`

	// set path for 'cert' file for https serving.
	CertPath string `json:"cert_path" yaml:"cert_path" toml:"cert_path" mapstructure:"cert_path"`

	// set path for 'key' file for https serving.
	KeyPath string `json:"key_path" yaml:"key_path" toml:"key_path" mapstructure:"key_path"`
}

// GetPortHTTP ...
func (c *config) GetPortHTTP() int {
	return c.PortHTTP
}

// GetPortHTTPS ...
func (c *config) GetPortHTTPS() int {
	return c.PortHTTPS
}

// GetShutdownTimeout ...
func (c *config) GetShutdownTimeout() int {
	return c.ShutdownTimeoutSec
}

// SetShutdownTimeout ...
func (c *config) SetShutdownTimeout(v int) {
	c.ShutdownTimeoutSec = v
}

// GetPemPath ...
func (c *config) GetPemPath() string {
	return c.PemPath
}

// GetCertPath ...
func (c *config) GetCertPath() string {
	return c.CertPath
}

// GetKeyPath ...
func (c *config) GetKeyPath() string {
	return c.KeyPath
}

func (c *config) validate() error {
	if c.ShutdownTimeoutSec < 0 {
		return fmt.Errorf("[service] cfgValidationErr: 'shutdown_timeout_sec' should be positive")
	}

	if c.PemPath != "" && !sys.PathIsFile(c.PemPath) {
		return fmt.Errorf("[service] cfgValidationErr: 'service.pem_path' in config, should be a file")
	}
	if c.CertPath != "" && !sys.PathIsFile(c.CertPath) {
		return fmt.Errorf("[service] cfgValidationErr: 'service.cert_path' in config, should be a file")
	}
	if c.KeyPath != "" && !sys.PathIsFile(c.KeyPath) {
		return fmt.Errorf("[service] cfgValidationErr: 'service.key_path' in config, should be a file")
	}
	return nil
}
