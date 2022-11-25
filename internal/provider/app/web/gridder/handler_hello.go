package gridder

import (
	"github.com/valyala/fasthttp"
)

func hello(ctx *fasthttp.RequestCtx, w writer) {
	w.RespOK(ctx, map[string]string{
		"other": "hello",
	})
}
