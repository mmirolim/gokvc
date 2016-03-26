package api

import (
	"bytes"
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

func New() fasthttp.RequestHandler {
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/get":
			get(ctx)
		case "/set":
			set(ctx)
		case "/del":
			del(ctx)
		case "/expvar":
			expvarhandler.ExpvarHandler(ctx)
		default:
			ctx.Error("cmd not found", fasthttp.StatusNotFound)
		}
	}

	return m
}

func get(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "get key %s\n", ctx.QueryArgs().QueryString())

}

func set(ctx *fasthttp.RequestCtx) {
	qstr := ctx.QueryArgs().QueryString()
	n := bytes.Index(qstr, []byte("="))
	if n > -1 {
		fmt.Fprintf(ctx, "set key %s val %s\n", qstr[:n], ctx.QueryArgs().PeekBytes(qstr[:n]))
	} else {
		fmt.Fprint(ctx, "wrong query format")
	}
}

func del(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	fmt.Fprintf(ctx, "del key %s\n", key)
}
