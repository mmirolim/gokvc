package api

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/mmirolim/gokvc/cache"
	"github.com/valyala/fasthttp"
)

func TestGETSET(t *testing.T) {

	key := "keySetGet"
	val := "valSetGet"

	if err := setString(key, val, 0); err != nil {
		t.Error(err)
	}

}

func TestDEL(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	key := "keyDel"
	val := "valDel"

	if err := setString(key, val, 0); err != nil {
		t.Error(err)
		return
	}

	// now del
	req.SetRequestURI(fmt.Sprintf("%s?k=%s", SDEL, key))
	ctx.Init(&req, nil, nil)

	del(&ctx)

	// get and check that it deleted
	req.SetRequestURI(fmt.Sprintf("%s?k=%s", SGET, key))
	ctx.Init(&req, nil, nil)

	get(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	if res.Header.StatusCode() != fasthttp.StatusNotFound {
		t.Errorf("expect status code %d got %d", fasthttp.StatusNotFound, res.Header.StatusCode())
	}
}

func TestSTTL(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	key := "keySTTL"
	val := "valSTTL"
	ttl := 1
	if err := setString(key, val, ttl); err != nil {
		t.Error(err)
		return
	}

	// get and check that it deleted
	req.SetRequestURI(fmt.Sprintf("%s?k=%s", SGET, key))
	ctx.Init(&req, nil, nil)
	// wait till expire
	time.Sleep(time.Duration(ttl) * time.Second)

	get(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	if res.Header.StatusCode() != fasthttp.StatusNotFound {
		t.Errorf("expect status code %d got %d", fasthttp.StatusNotFound, res.Header.StatusCode())
	}
}

func TestSLEN(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	key1 := "key1SLEN"
	key2 := "key2SLEN"
	val := "valSLEN"

	if err := setString(key1, val, 0); err != nil {
		t.Error(err)
		return
	}
	if err := setString(key2, val, 0); err != nil {
		t.Error(err)
		return
	}

	// get keys number
	// get and check that it deleted
	req.SetRequestURI(fmt.Sprintf("%s", SLEN))
	ctx.Init(&req, nil, nil)

	slen(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	body := res.Body()
	// response should be OK
	keysNum, err := strconv.Atoi(string(body))
	if err != nil {
		t.Error(err)
	}

	if keysNum < 2 {
		t.Errorf("expect keys number > %d got %s", 2, body)
	}

}

func TestMain(m *testing.M) {
	cache.Init()
	os.Exit(m.Run())
}

func setString(key, val string, ttl int) error {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	req.SetRequestURI(fmt.Sprintf("%s?k=%s&v=%s&t=%d", SSET, key, val, ttl))

	ctx.Init(&req, nil, nil)

	set(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		return err
	}

	if res.Header.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("expect status code %d got %d", fasthttp.StatusOK, res.Header.StatusCode())
	}

	body := res.Body()
	// response should be OK
	if string(body) != string(OK) {
		return fmt.Errorf("expect resp body %s got %s", OK, body)
	}

	// get value
	req.SetRequestURI(fmt.Sprintf("%s?k=%s", SGET, key))
	ctx.Init(&req, nil, nil)

	get(&ctx)

	s = ctx.Response.String()

	br = bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		return err
	}

	if res.Header.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("expect status code %d got %d", fasthttp.StatusOK, res.Header.StatusCode())
	}

	body = res.Body()
	// response should be value we previously set
	if string(body) != val {
		return fmt.Errorf("expect resp body %s got %s", val, body)
	}

	return nil
}
