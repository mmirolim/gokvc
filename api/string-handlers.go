package api

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/mmirolim/kvc/cache"
	"github.com/valyala/fasthttp"
)

func get(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()

	val, ok := cache.GET(key)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		fmt.Fprintf(ctx, "key %s not found\n", key)
		return
	}

	ctx.SetBody(val)
}

func set(ctx *fasthttp.RequestCtx) {
	var ttl int
	qstr := ctx.QueryArgs().QueryString()
	n := bytes.Index(qstr, EQUAL_SIGN)
	if n == -1 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprint(ctx, "wrong query format")
		return
	}

	ttlVal := ctx.Request.Header.PeekBytes(KEYTTL)
	if ttlVal != nil {
		ttl, _ = strconv.Atoi(string(ttlVal))
	}

	key := qstr[:n]
	val := ctx.QueryArgs().PeekBytes(key)

	cache.SET(key, val, ttl)

	ctx.SetBody(OK)
}

func del(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	cache.DEL(key)

	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBody(OK)
}

func sttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	r := cache.TTL(cache.STRING_CACHE, key)

	if r == cache.KeyNotExistCode {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	ctx.SetBody([]byte(strconv.Itoa(r)))
}

func skeys(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "%s", cache.KEYS())
}

func slen(ctx *fasthttp.RequestCtx) {
	ctx.SetBody([]byte(strconv.Itoa(cache.LEN())))
}
