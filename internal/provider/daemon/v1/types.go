package v1

import (
	"context"
)

type service interface {
	GetName() string

	Start(ctx context.Context) error
	Shutdown() error
}

type config interface {
	GetMaxPoc() int
	GetGCPercent() int

	GetPprof() bool

	GetPidFilePath() string
	SetPidFilePath(string)

	GetShutdownTimeoutSec() int
	SetShutdownTimeoutSec(int)
}

type logger interface {
	Debugf(string, interface{})
	Infof(string, interface{})
	Warnf(string, interface{})
	Errorf(string, interface{})
	Fatalf(string, interface{})
	Panicf(string, interface{})
}
