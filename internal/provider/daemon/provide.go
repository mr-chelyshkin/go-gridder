// automate connect components part.
// https://github.com/google/wire

package daemon

import (
	"context"
	"os/signal"
	"syscall"

	v1 "github.com/mr-chelyshkin/go-gridder/internal/provider/daemon/v1"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/logger"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/service"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var (
	V1Set = wire.NewSet(
		ProvideContext,
		ProvideConfig,
		ProvideDaemonV1,
	)
)

// Daemon ...
type Daemon interface {
	Run()
}

// ProvideContext return Daemon context.
func ProvideContext() (ctx context.Context, cancel func()) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
}

// ProvideConfig return pkg config object.
func ProvideConfig() (cfg config, err error) {
	if err = viper.UnmarshalKey("daemon", &cfg); err != nil {
		return
	}
	return cfg, cfg.validate()
}

// ProvideDaemonV1 return 'daemon_v1' as Daemon object.
func ProvideDaemonV1(
	ctx context.Context,
	cfg config,

	s service.Service,
	l logger.Logger,
) (Daemon, error) {
	return v1.NewDaemon(ctx, &cfg, s, l)
}
