//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/app/web"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/daemon"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/logger"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/service"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/service/fasthttp/writer"
)

// Init ...
// generate wire_gen.go file.
// https://github.com/google/wire
func Init(name string) (daemon.Daemon, func(), error) {
	wire.Build(
		wire.NewSet(
			service.FasthttpSet,
			web.AppWebApiSet,
			logger.ZapSet,
			daemon.V1Set,
			writer.Set,
		),
	)
	return nil, nil, nil
}
