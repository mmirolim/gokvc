package api

import (
	"fmt"
	"strconv"

	"github.com/mmirolim/gokvc/cache"
	"github.com/valyala/fasthttp"
)

func dget(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	if v, ok := cache.DGET(key); ok {
		fmt.Fprintf(ctx, "%s", v)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func dfget(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	key := args.PeekBytes(PKEY)
	fld := args.PeekBytes(PFLD)

	if v, ok := cache.DFGET(key, fld); ok {
		fmt.Fprintf(ctx, "%s", v)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func dfset(ctx *fasthttp.RequestCtx) {
	var ttl int
	args := ctx.QueryArgs()
	key := args.PeekBytes(PKEY)
	fld := args.PeekBytes(PFLD)
	val := args.PeekBytes(PVAL)
	if args.HasBytes(PTTL) {
		ttl = args.GetUintOrZero("t")
	}

	if cache.DFSET(key, fld, val, ttl) {
		ctx.SetBody(OK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

func dfdel(ctx *fasthttp.RequestCtx) {
	args := ctx.QueryArgs()
	key := args.PeekBytes(PKEY)
	fld := args.PeekBytes(PFLD)

	if cache.DFDEL(key, fld) {
		ctx.SetBody(OK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func ddel(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	if cache.DDEL(key) {
		ctx.SetBody(OK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func dttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	r := cache.DTTL(key)
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

func dkeys(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", cache.DKEYS())
}

func dlen(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte(strconv.Itoa(cache.DLEN())))
}
