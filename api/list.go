package api

import (
	"fmt"
	"strconv"

	"github.com/mmirolim/gokvc/cache"
	"github.com/valyala/fasthttp"
)

func lget(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()

	v, ok := cache.LGET(key)
	if ok {
		fmt.Fprintf(ctx, "%s", v)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func lpush(ctx *fasthttp.RequestCtx) {
	var ttl int
	key := ctx.QueryArgs().QueryString()
	val := ctx.QueryArgs().PeekBytes(key)

	ttlVal := ctx.Request.Header.PeekBytes(KEYTTL)
	if ttlVal != nil {
		ttl, _ = strconv.Atoi(string(ttlVal))
	}

	cache.LPUSH(key, val, ttl)

	ctx.SetBody(OK)
}

func lpop(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()

	v, ok := cache.LPOP(key)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(v)
}

func ldel(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	cache.LDEL(key)

	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBody(OK)
}

func lttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	r := cache.TTL(cache.LIST_CACHE, key)

	if r == cache.KeyNotExistCode {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	ctx.SetBody([]byte(strconv.Itoa(r)))
}

func lkeys(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%v", cache.LKEYS())
}

func llen(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte(strconv.Itoa(cache.LLEN())))
}
