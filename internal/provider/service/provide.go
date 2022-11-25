// automate connect components part.
// https://github.com/google/wire

package service

import (
	"context"

	"github.com/mr-chelyshkin/go-gridder/internal/provider/app/web"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/logger"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/service/fasthttp"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var (
	FasthttpSet = wire.NewSet(
		ProvideServiceFasthttp,
		ProvideConfig,
	)
)

// Service ...
type Service interface {
	GetName() string

	Start(ctx context.Context) error
	Shutdown() error
}

// ProvideConfig return pkg config object.
func ProvideConfig() (cfg config, err error) {
	if err = viper.UnmarshalKey("service", &cfg); err != nil {
		return
	}
	return cfg, cfg.validate()
}

// ProvideServiceFasthttp return 'fasthttp' as Service object.
func ProvideServiceFasthttp(cfg config, a web.App, l logger.Logger) (Service, error) {
	return fasthttp.NewService(&cfg, a, l)
}
