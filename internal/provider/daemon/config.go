package daemon

import (
	"fmt"
	"runtime"

	"github.com/mr-chelyshkin/go-gridder/internal/sys"
)

type config struct {
	// alternative path fo pid file.
	PidPath string `json:"pid_path" yaml:"pid_path" toml:"pid_path" mapstructure:"pid_path"`

	// gracefully service shutdown timeout.
	ShutdownTimeoutSec int `json:"shutdown_timeout_sec" yaml:"shutdown_timeout_sec" toml:"shutdown_timeout_sec" mapstructure:"shutdown_timeout_sec"`

	// set GOMAXPROCS value.
	MaxProc int `json:"max_proc" yaml:"max_proc" toml:"max_proc" mapstructure:"max_proc"`

	// set GCPERCENT value.
	GcPercent int `json:"gc_percent" yaml:"gc_percent" toml:"gc_percent" mapstructure:"gc_percent"`

	// bool value for pprof
	Pprof bool `json:"pprof" yaml:"pprof" toml:"pprof" mapstructure:"pprof"`
}

// GetPidFilePath ...
func (c *config) GetPidFilePath() string {
	return c.PidPath
}

// SetPidFilePath ...
func (c *config) SetPidFilePath(v string) {
	c.PidPath = v
}

// GetShutdownTimeoutSec ...
func (c *config) GetShutdownTimeoutSec() int {
	return c.ShutdownTimeoutSec
}

// SetShutdownTimeoutSec ...
func (c *config) SetShutdownTimeoutSec(v int) {
	c.ShutdownTimeoutSec = v
}

// GetMaxPoc ...
func (c *config) GetMaxPoc() int {
	return c.MaxProc
}

// GetGCPercent ...
func (c *config) GetGCPercent() int {
	return c.GcPercent
}

// GetPprof ...
func (c *config) GetPprof() bool {
	return c.Pprof
}

func (c *config) validate() error {
	if c.MaxProc <= 0 {
		c.MaxProc = runtime.NumCPU()
	}
	if c.GcPercent <= 50 {
		c.GcPercent = 100
	}
	if c.PidPath == "" {
		return fmt.Errorf("[daemon] cfgValidationErr: 'pid_path' required attr")
	}

	if c.PidPath != "" && sys.PathIsFile(c.PidPath) {
		return fmt.Errorf("[daemon] cfgValidationErr: 'pid_path', should be a file")
	}
	if c.ShutdownTimeoutSec < 0 {
		return fmt.Errorf("[daemon] cfgValidationErr: 'shutdown_timeout_sec' should be positive")
	}
	return nil
}
