// automate connect components part.
// https://github.com/google/wire

package logger

import (
	"github.com/mr-chelyshkin/go-gridder/internal/provider/logger/zap"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var (
	ZapSet = wire.NewSet(
		ProvideLoggerZap,
		ProvideConfig,
	)
)

// Logger ...
type Logger interface {
	Debugf(string, interface{})
	Infof(string, interface{})
	Warnf(string, interface{})
	Errorf(string, interface{})
	Fatalf(string, interface{})
	Panicf(string, interface{})
	Printf(string, ...interface{})
}

// ProvideConfig return pkg config object.
func ProvideConfig() (cfg config, err error) {
	if err = viper.UnmarshalKey("logger", &cfg); err != nil {
		return
	}
	return cfg, cfg.validate()
}

// ProvideLoggerZap return 'zap' as Logger object.
func ProvideLoggerZap(cfg config) (Logger, error) {
	if cfg.GetLogFilePath() != "" {
		return zap.NewLoggerFile(&cfg)
	}
	return zap.NewLoggerStdout(&cfg)
}
