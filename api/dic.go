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
	key := ctx.QueryArgs().PeekBytes(PKEY)
	fld := ctx.QueryArgs().PeekBytes(PFLD)

	if v, ok := cache.DFGET(key, fld); ok {
		fmt.Fprintf(ctx, "%s", v)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func dfset(ctx *fasthttp.RequestCtx) {
	var ttl int
	key := ctx.QueryArgs().PeekBytes(PKEY)
	fld := ctx.QueryArgs().PeekBytes(PFLD)
	val := ctx.QueryArgs().PeekBytes(PVAL)

	ttlVal := ctx.QueryArgs().PeekBytes(PTTL)
	if ttlVal != nil {
		ttl, _ = strconv.Atoi(string(ttlVal))
	}

	if cache.DFSET(key, fld, val, ttl) {
		ctx.SetBody(OK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

func dfdel(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)
	fld := ctx.QueryArgs().PeekBytes(PFLD)

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

	r := cache.TTL(cache.DIC_CACHE, key)
	if r == cache.KeyNotExistCode {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	ctx.SetBody([]byte(strconv.Itoa(r)))
}

func dkeys(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", cache.DKEYS())
}

func dlen(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte(strconv.Itoa(cache.DLEN())))
}
