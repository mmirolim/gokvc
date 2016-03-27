package api

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/golang/glog"
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

	if glog.V(3) {
		glog.Infof("key %s ok %b val %s", key, ok, val)
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

	if glog.V(3) {
		glog.Infof("key %s, val %s, ttl %d", key, val, ttl)
	}

	cache.SET(key, val, ttl)

	ctx.SetBody(OK)
}

func del(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	cache.DEL(key)
	ctx.SetStatusCode(fasthttp.StatusNotFound)

	if glog.V(3) {
		glog.Infof("key %s", key)
	}

	ctx.SetBody(OK)
}

func sttl(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	r := cache.TTL(cache.STRING_CACHE, key)

	ctx.SetBody([]byte(strconv.Itoa(r)))

	if r == cache.KeyNotExistCode {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
}
