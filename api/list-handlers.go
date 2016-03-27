package api

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/mmirolim/kvc/cache"
	"github.com/valyala/fasthttp"
)

func lget(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("not impl"))
}

func lpush(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("not impl"))
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
