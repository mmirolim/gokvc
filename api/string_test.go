package api

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mmirolim/gokvc/cache"
	"github.com/valyala/fasthttp"
)

func TestGET(t *testing.T) {

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
	time.Sleep(time.Duration(ttl)*time.Second + 10*time.Millisecond)

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

func TestKEYS(t *testing.T) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	var res fasthttp.Response

	key1 := "key1KEYS"
	key2 := "key2KEYS"
	val := "valKEYS"

	if err := setString(key1, val, 0); err != nil {
		t.Error(err)
		return
	}
	if err := setString(key2, val, 0); err != nil {
		t.Error(err)
		return
	}

	// get keys
	req.SetRequestURI(fmt.Sprintf("%s", SKEYS))
	ctx.Init(&req, nil, nil)

	skeys(&ctx)

	s := ctx.Response.String()

	br := bufio.NewReader(bytes.NewBufferString(s))
	if err := res.Read(br); err != nil {
		t.Error(err)
	}

	body := res.Body()
	// response should be OK
	if !(strings.Count(string(body), "key") > 2) {
		t.Errorf("expected number of keys > %d, got %s", 2, string(body))
	}
}

func BenchmarkSET(b *testing.B) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request

	req.SetRequestURI(fmt.Sprintf("%s?k=%s&v=%s&t=%d", SSET, "key1", "val1", 0))
	ctx.Init(&req, nil, nil)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		set(&ctx)
	}

}

func BenchmarkGET(b *testing.B) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request

	req.SetRequestURI(fmt.Sprintf("%s?k=%s&v=%s&t=%d", SSET, "key1", "val1", 0))
	ctx.Init(&req, nil, nil)
	// set data
	set(&ctx)
	// prepare get request
	req.SetRequestURI(fmt.Sprintf("%s?k=%s", SGET, "key1"))
	ctx.Init(&req, nil, nil)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		get(&ctx)
	}
}

func BenchmarkTTL(b *testing.B) {
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request

	req.SetRequestURI(fmt.Sprintf("%s?k=%s&v=%s&t=%d", SSET, "key1", "val1", 0))
	ctx.Init(&req, nil, nil)
	// set data
	set(&ctx)
	// prepare get request
	req.SetRequestURI(fmt.Sprintf("%s?k=%s", STTL, "key1"))
	ctx.Init(&req, nil, nil)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sttl(&ctx)
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

	return nil
}
