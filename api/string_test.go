package api

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/mmirolim/gokvc/cache"
	"github.com/valyala/fasthttp"
)

func TestGet(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	key := "keyNotExist"
	req.SetRequestURI(SGET + "?" + key)

	ctx.Init(&req, nil, nil)

	get(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
		return
	}

	if res.Header.StatusCode() != fasthttp.StatusNotFound {
		t.Errorf("expect status code %d got %d", fasthttp.StatusNotFound, res.Header.StatusCode())
	}
}

func TestSet(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	key := "keyApiTestGet"
	val := "valApiTestGet"
	req.SetRequestURI(SSET + "?" + key + "=" + val)

	ctx.Init(&req, nil, nil)

	set(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
		return
	}

	if res.Header.StatusCode() != fasthttp.StatusOK {
		t.Errorf("expect status code %d got %d", fasthttp.StatusOK, res.Header.StatusCode())
	}

	body := res.Body()
	// response should be OK
	if string(body) != string(OK) {
		t.Errorf("expect resp body %s got %s", OK, body)
	}

	// get value
	req.SetRequestURI(SGET + "?" + key)
	ctx.Init(&req, nil, nil)

	get(&ctx)

	s = ctx.Response.String()

	br = bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
		return
	}

	if res.Header.StatusCode() != fasthttp.StatusOK {
		t.Errorf("expect status code %d got %d", fasthttp.StatusOK, res.Header.StatusCode())
	}

	body = res.Body()
	// response should be value we previously set
	if string(body) != val {
		t.Errorf("expect resp body %s got %s", val, body)
	}

}

func TestMain(m *testing.M) {
	cache.Init()
	os.Exit(m.Run())
}
