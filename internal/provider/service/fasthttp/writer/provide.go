// automate connect components part.
// https://github.com/google/wire

package writer

import (
	"github.com/google/wire"
	"github.com/valyala/fasthttp"
)

// Writer ...
type Writer interface {
	RespOK(ctx *fasthttp.RequestCtx, i interface{})
	RespNotFound(ctx *fasthttp.RequestCtx, i interface{})
	RespBadRequest(ctx *fasthttp.RequestCtx, i interface{})
	RespInternalErr(ctx *fasthttp.RequestCtx, i interface{})
}

var (
	Set = wire.NewSet(
		ProvideWriter,
	)
)

// ProvideWriter return Writer object for 'fasthttp' service.
func ProvideWriter() Writer {
	return NewWriter()
}
