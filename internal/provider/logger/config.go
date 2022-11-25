package logger

import (
	"fmt"
)

type config struct {
	// set logger level: [DEBUG, INFO, WARN, ERROR].
	Level string `json:"level" yaml:"level" toml:"level" mapstructure:"level"`

	// set file to sync logs.
	Path string `json:"path" yaml:"path" toml:"path" mapstructure:"path"`

	// rotation settings, worked only if logs sync to file.
	MaxSize    int `json:"max_size"    yaml:"max_size"    toml:"max_size"    mapstructure:"max_size"`
	MaxAge     int `json:"max_age"     yaml:"max_age"     toml:"max_age"     mapstructure:"max_age"`
	MaxBackups int `json:"max_backups" yaml:"max_backups" toml:"max_backups" mapstructure:"max_backups"`
}

// GetLogLevel ...
func (c *config) GetLogLevel() string {
	return c.Level
}

// SetLogLevel ...
func (c *config) SetLogLevel(v string) {
	c.Level = v
}

// GetLogFilePath ...
func (c *config) GetLogFilePath() string {
	return c.Path
}

// GetLogFileMaxSize ...
func (c *config) GetLogFileMaxSize() int {
	return c.MaxSize
}

// SetLogFileMaxSize ...
func (c *config) SetLogFileMaxSize(v int) {
	c.MaxSize = v
}

// GetLogFileMaxAge ...
func (c *config) GetLogFileMaxAge() int {
	return c.MaxAge
}

// SetLogFileMaxAge ...
func (c *config) SetLogFileMaxAge(v int) {
	c.MaxAge = v
}

// GetLogFileMaxBackups ...
func (c *config) GetLogFileMaxBackups() int {
	return c.MaxBackups
}

// SetLogFileMaxBackups ...
func (c *config) SetLogFileMaxBackups(v int) {
	c.MaxBackups = v
}

func (c *config) validate() error {
	if c.MaxAge < 0 {
		return fmt.Errorf("[logger] cfgValidationErr: 'logger.max_size' should be positive")
	}
	if c.MaxAge < 0 {
		return fmt.Errorf("[logger] cfgValidationErr: 'logger.max_age' should be positive")
	}
	if c.MaxBackups < 0 {
		return fmt.Errorf("[logger] cfgValidationErr: 'logger.max_backups' should be positive")
	}
	return nil
}
