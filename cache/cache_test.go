package cache

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func BenchmarkTimeNow(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		time.Now().UnixNano()
	}

}

func BenchmarkSysTime(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		globalSysTimeNow.Get()
	}

}

func TestSetGet(t *testing.T) {
	cases := []struct {
		Key    []byte
		SetVal []byte
		GetVal []byte
		Expect bool
	}{
		{[]byte("k1"), []byte("v1"), []byte("v1"), true},
		{[]byte("k2"), []byte(""), []byte("v2"), false},
	}

	for _, c := range cases {
		SET(c.Key, c.SetVal, 0)
		if v, ok := GET(c.Key); !ok || bytes.Equal(v, c.GetVal) != c.Expect {
			t.Errorf("expected %v got %v\n", c.GetVal, v)
		}
	}
}

func TestDel(t *testing.T) {
	key := []byte("kdel")
	val := []byte("vdel")
	SET(key, val, 0)
	_, ok := GET(key)
	if !ok {
		t.Errorf("key not found", key)
	}

	DEL(key)
	_, ok = GET(key)
	if ok {
		t.Errorf("expected key %s !ok got ok", key)
	}

}

func BenchmarkSet(b *testing.B) {
	key := []byte("k1")
	val := []byte("v1")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		SET(key, val, 0)
	}
}

func BenchmarkGet(b *testing.B) {
	key := []byte("k1")
	val := []byte("v1")
	SET(key, val, 0)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		GET(key)
	}

}

func BenchmarkDel(b *testing.B) {
	key := []byte("k1")
	val := []byte("v1")
	SET(key, val, 0)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		DEL(key)
	}

}

func BenchmarkParallelSET(b *testing.B) {
	key := []byte("k1")
	val := []byte("v1")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			SET(key, val, 0)
		}
	})
}

func BenchmarkParallelGET(b *testing.B) {
	key := []byte("k1")
	val := []byte("v1")
	SET(key, val, 0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GET(key)
		}
	})
}

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}
