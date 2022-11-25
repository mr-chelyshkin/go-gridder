// automate connect components part.
// https://github.com/google/wire

package web

import (
	"github.com/mr-chelyshkin/go-gridder/internal/provider/app/web/gridder"
	"github.com/mr-chelyshkin/go-gridder/internal/provider/service/fasthttp/writer"

	"github.com/google/wire"
	"github.com/valyala/fasthttp"
)

var (
	AppWebApiSet = wire.NewSet(
		ProvideAppWebBlank,
	)
)

// App ...
type App interface {
	Route(ctx *fasthttp.RequestCtx) bool
	Name() string
}

// ProvideAppWebBlank return 'blank' as App object.
func ProvideAppWebBlank(name string, w writer.Writer) (App, error) {
	return gridder.NewApp(name, w)
}
