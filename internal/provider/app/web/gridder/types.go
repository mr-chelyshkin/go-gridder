package gridder

import (
	"github.com/valyala/fasthttp"
)

type writer interface {
	RespOK(ctx *fasthttp.RequestCtx, i interface{})
	RespNotFound(ctx *fasthttp.RequestCtx, i interface{})
	RespBadRequest(ctx *fasthttp.RequestCtx, i interface{})
	RespInternalErr(ctx *fasthttp.RequestCtx, i interface{})
}
