package api

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/mmirolim/kvc/cache"
	"github.com/valyala/fasthttp"
)

func dget(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("not impl"))
}

func dkget(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("not impl"))
}

func dkset(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("not impl"))
}

func ddel(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	cache.DDEL(key)
	ctx.SetStatusCode(fasthttp.StatusNotFound)

	if glog.V(3) {
		glog.Infof("key %s", key)
	}

	ctx.SetBody(OK)
}

func dttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	r := cache.TTL(cache.DIC_CACHE, key)

	ctx.SetBody([]byte(strconv.Itoa(r)))

	if r == cache.KeyNotExistCode {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
}

func dkeys(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", cache.DKEYS())
}

func dlen(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte(strconv.Itoa(cache.DLEN())))
}
