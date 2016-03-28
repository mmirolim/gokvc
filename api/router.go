package api

import (
	"github.com/golang/glog"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

const (
	PING = "/ping" // ping api
	// system
	EXPVAR = "/expvar"

	// kv commands
	SGET  = "/get"  // GET key
	SDEL  = "/del"  // DEL key
	SSET  = "/set"  // SET key val
	STTL  = "/ttl"  // TTL key gets seconds left o expire
	SKEYS = "/keys" // SKEYS gets all keys of strings
	SLEN  = "/slen" // SLEN gets number of string elements

	// list commands
	LGET  = "/lget"  // LGET key get all list
	LDEL  = "/ldel"  // LDEL key deletes list
	LPUSH = "/lpush" // LPUSH key val prepends list with element
	LPOP  = "/lpop"  // LPOP key removes and returns the first element of the list stored at key
	LTTL  = "/lttl"  // TTL key gets seconds left o expire
	LKEYS = "/lkeys" // LKEYS gets all keys of lists
	LLEN  = "/llen"  // LLEN gets number of list elements

	// dictionary commands
	DGET  = "/dget"  // DGET key get all field from dic
	DDEL  = "/ddel"  // DDEL key delete dic
	DFGET = "/dfget" // DKGET key field get field from dic
	DFDEL = "/dfdel" // DKDEL key field delete field in dic
	DFSET = "/dfset" // DKSET key field val sets field in dic to val
	DTTL  = "/dttl"  // TTL key gets seconds left o expire
	DKEYS = "/dkeys" // DKEYS gets all keys of dics
	DLEN  = "/dlen"  // DLEN return number of dic elements
)

var (
	EQUAL_SIGN = []byte("=")

	PONG = []byte("pong") // ping response

	// TTL passed with http headers
	KEYTTL = []byte("KEYTTL") // KEYTTL seconds

	HTTP_CACHE_CONTROL = []byte("Cache-control")
	HTTP_NO_CACHE      = []byte("private, max-age=0, no-cache")
	OK                 = []byte("OK")

	HandlersMap = map[string]fasthttp.RequestHandler{
		PING:   ping,
		SGET:   get,
		SSET:   set,
		SDEL:   del,
		STTL:   sttl,
		SKEYS:  skeys,
		SLEN:   slen,
		LGET:   lget,
		LPUSH:  lpush,
		LPOP:   lpop,
		LDEL:   ldel,
		LKEYS:  lkeys,
		LLEN:   llen,
		DGET:   dget,
		DFGET:  dfget,
		DFSET:  dfset,
		DFDEL:  dfdel,
		DDEL:   ddel,
		DTTL:   dttl,
		DKEYS:  dkeys,
		DLEN:   dlen,
		EXPVAR: expvarhandler.ExpvarHandler,
	}
)

func New() fasthttp.RequestHandler {
	m := func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetCanonical(HTTP_CACHE_CONTROL, HTTP_NO_CACHE)
		if glog.V(2) {
			glog.Infof("url %s", ctx.Path())
		}

		handler, ok := HandlersMap[string(ctx.Path())]

		if ok {
			handler(ctx)
		} else {
			ctx.Error("cmd not found", fasthttp.StatusNotFound)
		}

	}

	return m
}

func ping(ctx *fasthttp.RequestCtx) {
	ctx.SetBody(PONG)
}
