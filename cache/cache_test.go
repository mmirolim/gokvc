package cache

import (
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

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}
