package api

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestLPUSH(t *testing.T) {
	// lpush should create cache object if not exist
	key := "keyNotExistLPUSH"
	val := "valLPUSH"

	if err := setList(key, val, 0); err != nil {
		t.Error(err)
	}
}

func TestLGET(t *testing.T) {
	// lpush should create cache object if not exist
	key := "keyLGET"
	val1 := "val1LGET"
	val2 := "val2LGET"

	// push two values to list
	if err := setList(key, val1, 0); err != nil {
		t.Error(err)
	}

	if err := setList(key, val2, 0); err != nil {
		t.Error(err)
	}

	// get all elements of a list
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	req.SetRequestURI(fmt.Sprintf("%s?k=%s", LGET, key))

	ctx.Init(&req, nil, nil)

	lget(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	body := res.Body()
	// response should be OK
	if !strings.Contains(string(body), val1) && !strings.Contains(string(body), val2) {
		t.Errorf("expect resp body conatains %s and %s, got %s", val1, val2, body)
	}

}

func TestLPOP(t *testing.T) {
	// lpush should create cache object if not exist
	key := "keyLPOP"
	val := "valLPOP"

	// push two values to list
	if err := setList(key, val, 0); err != nil {
		t.Error(err)
	}

	// first pop should return val
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	req.SetRequestURI(fmt.Sprintf("%s?k=%s", LPOP, key))

	ctx.Init(&req, nil, nil)

	lpop(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	body := res.Body()
	// response should be OK
	if string(body) != val {
		t.Errorf("expect resp body %s, got %s", val, body)
	}
	// second one should delete list

	req.SetRequestURI(fmt.Sprintf("%s?k=%s", LPOP, key))

	ctx.Init(&req, nil, nil)

	lpop(&ctx)

	s = ctx.Response.String()

	br = bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	if res.Header.StatusCode() != fasthttp.StatusNotFound {
		t.Errorf("expect status code %d got %d", fasthttp.StatusNotFound, res.Header.StatusCode())
	}

}

func setList(key, val string, ttl int) error {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	req.SetRequestURI(fmt.Sprintf("%s?k=%s&v=%s&t=%d", LPUSH, key, val, ttl))

	ctx.Init(&req, nil, nil)

	lpush(&ctx)

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

	return nil
}
