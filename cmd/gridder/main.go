package main

import (
	"flag"
	"path/filepath"

	"github.com/mr-chelyshkin/go-gridder/internal/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	applicationName = "gridder"
	defaultCfgPath  = filepath.Join(".local_configs", applicationName, "config.toml")
)

func main() {
	flag.String("config", defaultCfgPath, "daemon config file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// getting configs to viper data.
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	cfg, err := config.NewFileConfig(viper.GetViper(), viper.GetString("config"))
	if err != nil {
		panic(err)
	}
	if err := cfg.Read(); err != nil {
		panic(err)
	}

	// init daemon.
	d, cleanup, err := Init(applicationName)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// run.
	d.Run()
}
