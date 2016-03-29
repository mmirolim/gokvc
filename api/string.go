package api

import (
	"fmt"
	"strconv"

	"github.com/mmirolim/gokvc/cache"
	"github.com/valyala/fasthttp"
)

func get(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	if v, ok := cache.GET(key); ok {
		ctx.SetBody(v)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func set(ctx *fasthttp.RequestCtx) {
	var ttl int
	args := ctx.QueryArgs()
	key := args.PeekBytes(PKEY)
	val := args.PeekBytes(PVAL)
	if args.HasBytes(PTTL) {
		ttl = args.GetUintOrZero("t")
	}

	if cache.SET(key, val, ttl) {
		ctx.SetBody(OK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

func del(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	if cache.DEL(key) {
		ctx.SetBody(OK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func sttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	r := cache.TTL(key)
	if r == cache.KeyNotExistCode {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
	// do not make alloc for common ttl codes
	v, ok := cache.TtlKeyCodes[r]
	if ok {
		ctx.SetBody(v)
		return
	}

	ctx.SetBody([]byte(strconv.Itoa(r)))
}

func skeys(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", cache.KEYS())
}

func slen(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte(strconv.Itoa(cache.LEN())))
}
