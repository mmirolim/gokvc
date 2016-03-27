package api

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/mmirolim/kvc/cache"
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
	if ok {
		fmt.Fprintf(ctx, "%s", v)
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func ldel(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	cache.LDEL(key)
	ctx.SetStatusCode(fasthttp.StatusNotFound)

	if glog.V(3) {
		glog.Infof("key %s", key)
	}

	ctx.SetBody(OK)
}

func lttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	r := cache.TTL(cache.LIST_CACHE, key)

	ctx.SetBody([]byte(strconv.Itoa(r)))

	if r == cache.KeyNotExistCode {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
}

func lkeys(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", cache.LKEYS())
}

func llen(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte(strconv.Itoa(cache.LLEN())))
}
