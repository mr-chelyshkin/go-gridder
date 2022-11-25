package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type validator interface {
	Check() error
}

type fileConfig struct {
	v *viper.Viper

	name string
	dir  string
	ext  string

	validator validator
}

// NewFileConfig ...
func NewFileConfig(v *viper.Viper, path string, opts ...FileConfigOpt) (*fileConfig, error) {
	if v == nil {
		return nil, fmt.Errorf("incoming viper object is nil")
	}

	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	dir := filepath.Dir(path)
	file := ""

	// define file type by extension.
	// if extension not in [json, toml, yaml] use toml type as default.
	switch filepath.Ext(path) {
	case ".toml":
		file = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		ext = "toml"
	case ".yaml":
		file = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		ext = "yaml"
	case ".json":
		file = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		ext = "json"
	default:
		file = filepath.Base(path)
		ext = "toml"
	}

	fc := &fileConfig{
		v: v,

		name: file,
		dir:  dir,
		ext:  ext,
	}
	for _, opt := range opts {
		opt(fc)
	}

	// viper set.
	fc.v.SetConfigName(fc.name)
	fc.v.AddConfigPath(fc.dir)
	fc.v.SetConfigType(fc.ext)
	return fc, nil
}

// Read file to viper data.
func (fc *fileConfig) Read() error {
	return fc.v.ReadInConfig()
}

// Parse file to struct.
func (fc *fileConfig) Parse(cfg interface{}) error {
	if err := fc.Read(); err != nil {
		return err
	}
	if err := fc.v.Unmarshal(&cfg); err != nil {
		return err
	}

	if fc.validator != nil {
		return fc.validator.Check()
	}
	return nil
}

// Watch and parse file if changes.
func (fc *fileConfig) Watch(cfg interface{}, ch chan<- error) {
	fc.v.WatchConfig()
	fc.v.OnConfigChange(func(in fsnotify.Event) {
		if err := fc.Parse(cfg); err != nil {
			ch <- err
		}
	})
}
