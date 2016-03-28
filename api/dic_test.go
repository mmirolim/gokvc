package api

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestDGET(t *testing.T) {
	// lpush should create cache object if not exist
	key := "keyDGET"
	fld1 := "fld1DGET"
	fld2 := "fld2DGET"
	val1 := "val1DGET"
	val2 := "val2DGET"

	// set two field to dic, it should be created
	if err := setDic(key, fld1, val1, 0); err != nil {
		t.Error(err)
	}

	if err := setDic(key, fld2, val2, 0); err != nil {
		t.Error(err)
	}

	// now get and check that all field and values stored
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	req.SetRequestURI(fmt.Sprintf("%s?k=%s", DGET, key))

	ctx.Init(&req, nil, nil)

	dget(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	body := res.Body()
	// response should be OK
	if !strings.Contains(string(body), val1) || !strings.Contains(string(body), val2) {
		t.Errorf("expect resp body contains %s %s got %s", val1, val2, body)
	}

}

func TestDFGET(t *testing.T) {
	// lpush should create cache object if not exist
	key := "keyDFGET"
	fld1 := "fld1DFGET"
	fld2 := "fld2DFGET"
	val1 := "val1DFGET"
	val2 := "val2DFGET"

	// set two field to dic, it should be created
	if err := setDic(key, fld1, val1, 0); err != nil {
		t.Error(err)
	}

	if err := setDic(key, fld2, val2, 0); err != nil {
		t.Error(err)
	}

	// now get and check that all field and values stored
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	req.SetRequestURI(fmt.Sprintf("%s?k=%s&f=%s", DFGET, key, fld2))

	ctx.Init(&req, nil, nil)

	dfget(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	body := res.Body()
	// response should be OK
	if string(body) != val2 {
		t.Errorf("expect resp body %s got %s", val2, body)
	}

}

func setDic(key, fld, val string, ttl int) error {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	req.SetRequestURI(fmt.Sprintf("%s?k=%s&f=%s&v=%s&t=%d", DFSET, key, fld, val, ttl))

	ctx.Init(&req, nil, nil)

	dfset(&ctx)

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
