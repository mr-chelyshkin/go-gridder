package fasthttp

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/mr-chelyshkin/go-gridder/internal/sys"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type middleware struct {
	l logger
}

func newMiddleware(l logger) middleware {
	return middleware{
		l: l,
	}
}

func (m *middleware) base(next func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
	fn := func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		requestUUID(ctx)

		defer func() {
			accessLog(ctx, m.l, start)
		}()
		defer func() {
			if rvr := recover(); rvr != nil {
				msg, _ := rvr.(string)
				ctx.Error(msg, fasthttp.StatusInternalServerError)
				return
			}
		}()
		if err := validateURL(ctx); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		}
		next(ctx)
	}
	return fn
}

func requestUUID(ctx *fasthttp.RequestCtx) {
	ctx.Request.Header.Add("uuid", uuid.New().String())
}

func accessLog(ctx *fasthttp.RequestCtx, l logger, start time.Time) {
	fields := map[string]string{
		"status":  strconv.Itoa(ctx.Response.StatusCode()),
		"latency": time.Since(start).String(),
		"method":  sys.ToString(ctx.Method()),
		"uri":     sys.ToString(ctx.RequestURI()),
		"host":    sys.ToString(ctx.Request.Host()),
		"uuid":    sys.ToString(ctx.Request.Header.Peek("uuid")),
	}

	n := ctx.Response.StatusCode()
	switch {
	case n >= 500:
		l.Errorf("server error", fields)
	case n >= 400:
		l.Warnf("client error", fields)
	case n >= 300:
		l.Infof("redirection", fields)
	default:
		l.Infof("success", fields)
	}
	return
}

func validateURL(ctx *fasthttp.RequestCtx) error {
	if bytes.Contains(ctx.RequestURI(), sys.ToBytes("..")) {
		return fmt.Errorf("unsupported dots in path")
	}
	return nil
}
