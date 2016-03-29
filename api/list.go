package api

import (
	"fmt"
	"strconv"

	"github.com/mmirolim/gokvc/cache"
	"github.com/valyala/fasthttp"
)

func lget(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	if v, ok := cache.LGET(key); ok {
		fmt.Fprintf(ctx, "%s", v)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func lpush(ctx *fasthttp.RequestCtx) {
	var ttl int
	args := ctx.QueryArgs()
	key := args.PeekBytes(PKEY)
	val := args.PeekBytes(PVAL)
	if args.HasBytes(PTTL) {
		ttl = args.GetUintOrZero("t")
	}

	if cache.LPUSH(key, val, ttl) {
		ctx.SetBody(OK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

func lpop(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	if v, ok := cache.LPOP(key); ok {
		fmt.Fprintf(ctx, "%s", v)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func ldel(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	if cache.LDEL(key) {
		ctx.SetBody(OK)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func lttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().PeekBytes(PKEY)

	r := cache.LTTL(key)
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

func lkeys(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%v", cache.LKEYS())
}

func llen(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte(strconv.Itoa(cache.LLEN())))
}
