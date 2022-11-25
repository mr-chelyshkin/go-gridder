package writer

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type writer struct{}

// NewWriter return writer object for fasthttp service.
func NewWriter() *writer {
	return &writer{}
}

// RespOK send JSON response with code 200.
func (w *writer) RespOK(ctx *fasthttp.RequestCtx, i interface{}) {
	writeJSON(ctx, i, fasthttp.StatusOK)
	return
}

// RespNotFound send JSON response with code 404.
func (w *writer) RespNotFound(ctx *fasthttp.RequestCtx, i interface{}) {
	writeJSON(ctx, i, fasthttp.StatusNotFound)
	return
}

// RespInternalErr send JSON response with code 500.
func (w *writer) RespInternalErr(ctx *fasthttp.RequestCtx, i interface{}) {
	writeJSON(ctx, i, fasthttp.StatusInternalServerError)
	return
}

// RespBadRequest send JSON response with code 400.
func (w *writer) RespBadRequest(ctx *fasthttp.RequestCtx, i interface{}) {
	writeJSON(ctx, i, fasthttp.StatusBadRequest)
	return
}

func writeJSON(ctx *fasthttp.RequestCtx, i interface{}, s int) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(s)

	if err := json.NewEncoder(ctx).Encode(i); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
