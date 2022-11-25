package fasthttp

import (
	"github.com/valyala/fasthttp"
)

type app interface {
	Route(ctx *fasthttp.RequestCtx) bool

	Name() string
}

type config interface {
	GetPortHTTP() int
	GetPortHTTPS() int

	GetShutdownTimeout() int
	SetShutdownTimeout(int)

	GetPemPath() string
	GetKeyPath() string
	GetCertPath() string
}

type logger interface {
	Debugf(string, interface{})
	Infof(string, interface{})
	Warnf(string, interface{})
	Errorf(string, interface{})
	Printf(string, ...interface{})
}
