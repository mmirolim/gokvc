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
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		fmt.Fprintf(ctx, "key %s not found\n", key)
		return
	}

	if glog.V(3) {
		glog.Infof("key %s ok %b val %s", key, ok, val)
	}

	fmt.Fprintf(ctx, "get key %s val %s\n", key, val)

}

func set(ctx *fasthttp.RequestCtx) {
	var ttl int
	qstr := ctx.QueryArgs().QueryString()
	n := bytes.Index(qstr, EQUAL_SIGN)
	if n == -1 {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprint(ctx, "wrong query format")
		return
	}
	ttlVal := ctx.Request.Header.PeekBytes(cache.KEYTTL)
	if ttlVal != nil {
		ttl, _ = strconv.Atoi(string(ttlVal))
	}
	key := qstr[:n]
	val := ctx.QueryArgs().PeekBytes(key)

	if glog.V(3) {
		glog.Infof("key %s, val %s, ttl %d", key, val, ttl)
	}

	cache.SET(key, val, ttl)
	fmt.Fprintf(ctx, "set key %s val %s\n", key, val)
	if ttl > 0 {
		fmt.Fprintf(ctx, "with ttl %d seconds\n", ttl)
	}

}

func del(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	cache.DEL(key)
	ctx.Response.SetStatusCode(fasthttp.StatusNotFound)

	if glog.V(3) {
		glog.Infof("key %s", key)
	}

	fmt.Fprintf(ctx, "del key %s\n", key)
}
