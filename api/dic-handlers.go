package api

import (
	"fmt"
	"strconv"

	"github.com/mmirolim/kvc/cache"
	"github.com/valyala/fasthttp"
)

func dget(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()

	v, ok := cache.DGET(key)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	fmt.Fprintf(ctx, "%s", v)
}

func dfget(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("not impl"))
}

func dfset(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte("not impl"))
}

func dfdel(ctx *fasthttp.RequestCtx) {
	//	key := ctx.QueryArgs().QueryString()
	//	ctx.SetBody(cache.DFDEL(key, fld))
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBody(OK)
}

func ddel(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()

	cache.DDEL(key)

	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBody(OK)
}

func dttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()

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
