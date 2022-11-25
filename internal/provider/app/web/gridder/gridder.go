package gridder

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type app struct {
	name string

	w writer
}

// NewApp return blank fasthttp application object.
func NewApp(name string, w writer) (*app, error) {
	if w == nil {
		return nil, fmt.Errorf("incoming writer object is nil")
	}
	return &app{
		w: w,

		name: name,
	}, nil
}

// Name of service.
func (a *app) Name() string {
	return a.name
}

// Route by application handlers.
func (a *app) Route(ctx *fasthttp.RequestCtx) bool {
	next := true

	switch string(ctx.Path()) {
	case "/hello":
		hello(ctx, a.w)
	default:
		next = false
	}
	return next
}
