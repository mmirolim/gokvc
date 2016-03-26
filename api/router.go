package api

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/mmirolim/kvc/cache"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

var (
	SGET = strings.ToLower("/" + cache.SGET)
	SSET = strings.ToLower("/" + cache.SSET)
	SDEL = strings.ToLower("/" + cache.SDEL)

	// system
	EXPVAR = "/expvar"
)

func New() fasthttp.RequestHandler {
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case SGET:
			get(ctx)
		case SSET:
			set(ctx)
		case SDEL:
			del(ctx)
		case EXPVAR:
			expvarhandler.ExpvarHandler(ctx)
		default:
			ctx.Error("cmd not found", fasthttp.StatusNotFound)
		}
	}

	return m
}

func get(ctx *fasthttp.RequestCtx) {
	qstr := ctx.QueryArgs().QueryString()

	val, ok := cache.GET(qstr)
	if !ok {
		fmt.Fprintf(ctx, "key %s not found\n", qstr)
		return
	}

	fmt.Fprintf(ctx, "get key %s val %s\n", qstr, val)

}

func set(ctx *fasthttp.RequestCtx) {
	var ttl int
	qstr := ctx.QueryArgs().QueryString()
	n := bytes.Index(qstr, []byte("="))
	if n == -1 {
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
		glog.Infof("url %s, key %s, val %s, ttl %d", SSET, key, val, ttl)
	}
	cache.SET(key, val, ttl)
	fmt.Fprintf(ctx, "set key %s val %s\n", key, val)
	if ttl > 0 {
		fmt.Fprintf(ctx, "with ttl %d seconds\n", ttl)
	}

}

func del(ctx *fasthttp.RequestCtx) {
	key := ctx.QueryArgs().QueryString()
	fmt.Fprintf(ctx, "del key %s\n", key)
}
