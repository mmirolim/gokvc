package cache

import (
	"os"
	"testing"
	"time"
)

func BenchmarkTimeNow(b *testing.B) {
	b.ReportAllocs()
	var val int64
	for i := 0; i < b.N; i++ {
		val = time.Now().UnixNano()
	}

	b.Logf("val %d\n", val)
}

func BenchmarkSysTime(b *testing.B) {
	b.ReportAllocs()
	var val int64
	for i := 0; i < b.N; i++ {
		val = globalSysTimeNow.Get()
	}

	b.Logf("val %d\n", val)
}

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}
