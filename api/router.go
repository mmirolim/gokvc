package api

import (
	"strings"

	"github.com/golang/glog"
	"github.com/mmirolim/kvc/cache"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

var (
	EQUAL_SIGN = []byte("=")
	SGET       = strings.ToLower("/" + cache.SGET)
	SSET       = strings.ToLower("/" + cache.SSET)
	SDEL       = strings.ToLower("/" + cache.SDEL)

	// TTL passed with http headers
	KEYTTL = []byte("KEYTTL") // KEYTTL seconds

	// system
	EXPVAR = "/expvar"

	HTTP_CACHE_CONTROL = []byte("Cache-control")
	HTTP_NO_CACHE      = []byte("private, max-age=0, no-cache")
	OK                 = []byte("OK")
)

func New() fasthttp.RequestHandler {
	m := func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetCanonical(HTTP_CACHE_CONTROL, HTTP_NO_CACHE)
		if glog.V(2) {
			glog.Infof("url %s", ctx.Path())
		}
		switch string(ctx.Path()) {
		case SGET:
			get(ctx)
		case SSET:
			set(ctx)
		case SDEL:
			del(ctx)
		case EXPVAR:
			expvarhandler.ExpvarHandler(ctx)
		default:
			ctx.Error("cmd not found", fasthttp.StatusNotFound)
		}
	}

	return m
}
