// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/mr-chelyshkin/go-gridder/internal/provider/app/web"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/daemon"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/logger"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/service"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/service/fasthttp/writer"
)

// Injectors from wire.go:

// Init ...
// generate wire_gen.go file.
// https://github.com/google/wire
func Init(name string) (daemon.Daemon, func(), error) {
	context, cleanup := daemon.ProvideContext()
	config, err := daemon.ProvideConfig()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	serviceConfig, err := service.ProvideConfig()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	writerWriter := writer.ProvideWriter()
	app, err := web.ProvideAppWebBlank(name, writerWriter)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	loggerConfig, err := logger.ProvideConfig()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	loggerLogger, err := logger.ProvideLoggerZap(loggerConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	serviceService, err := service.ProvideServiceFasthttp(serviceConfig, app, loggerLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	daemonDaemon, err := daemon.ProvideDaemonV1(context, config, serviceService, loggerLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return daemonDaemon, func() {
		cleanup()
	}, nil
}
