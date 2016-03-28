package cache

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	_CHM_SHARD_NUM = 1 << 6
	_MASK          = _CHM_SHARD_NUM - 1

	// TTL codes
	KeyTTLErrCode   = -3
	KeyNotExistCode = -2
	KeyHasNoTTLCode = -1

	// cache types
	STRING_CACHE CacheType = iota + 1
	LIST_CACHE
	DIC_CACHE
)

var (
	// GC params
	gGCRunDelay       = 5 * time.Second
	gGCShardScanDelay = 10 * time.Millisecond

	initOnce sync.Once

	globalSysTimeNow Atomic

	globalStringCache = &StringCache{}
	globalListCache   = &ListCache{}
	globalDicCache    = &DicCache{}

	hasher = fnv.New32a()
)

func CacheTimeNow() int64 {
	return globalSysTimeNow.Get()
}

type CacheType int

type item struct {
	ttl int64
	la  int64
}

func hash(data []byte) uintptr {
	hasher.Write(data)
	h := uintptr(hasher.Sum32())
	hasher.Reset()
	return h
}

func initialize() {
	ticker := time.NewTicker(time.Microsecond)
	go func() {
		for t := range ticker.C {
			globalSysTimeNow.Set(t.UnixNano())
		}
	}()

	for i := 0; i < _CHM_SHARD_NUM; i++ {
		globalStringCache.shards[i].m = make(map[string]String)
		globalListCache.shards[i].m = make(map[string]List)
		globalDicCache.shards[i].m = make(map[string]Dic)
	}

	// start gc
	gc()

}

// GC, marks expired keys and removes them in batches
// each cache type scaned concurrently
func gc() {
	// start gc for string cache
	go func() {
		c := globalStringCache
		expiredKeys := make(map[int][]string)
		for {
			// scan delay
			time.Sleep(gGCRunDelay)

			// select expired keys
			for i := 0; i < _CHM_SHARD_NUM; i++ {
				c.shards[i].RLock()

				for k := range c.shards[i].m {
					if c.shards[i].m[k].IsExpired() {
						expiredKeys[i] = append(expiredKeys[i], k)
					}
				}

				c.shards[i].RUnlock()
				// delay between shard scans
				time.Sleep(gGCShardScanDelay)

			}

			// remove selected keys in batch
			for i := 0; i < _CHM_SHARD_NUM; i++ {
				shard := &c.shards[i]
				shard.Lock()
				for _, key := range expiredKeys[i] {
					// recheck if expired
					v, ok := shard.m[key]
					if ok {
						if v.IsExpired() {
							delete(shard.m, key)
						}
					}

				}
				shard.Unlock()
				// reset slice
				expiredKeys[i] = expiredKeys[i][:0]
			}
		}

	}()
	// start gc for list cache
	go func() {
		c := globalListCache
		expiredKeys := make(map[int][]string)
		for {
			// scan delay
			time.Sleep(gGCRunDelay)
			// select expired keys
			for i := 0; i < _CHM_SHARD_NUM; i++ {
				c.shards[i].RLock()

				for k := range c.shards[i].m {
					if c.shards[i].m[k].IsExpired() {
						expiredKeys[i] = append(expiredKeys[i], k)
					}
				}

				c.shards[i].RUnlock()
				// delay between shard scans
				time.Sleep(gGCShardScanDelay)
			}

			// remove selected keys in batch
			for i := 0; i < _CHM_SHARD_NUM; i++ {
				shard := &c.shards[i]
				shard.Lock()
				for _, key := range expiredKeys[i] {
					// recheck if expired
					v, ok := shard.m[key]
					if ok {
						if v.IsExpired() {
							delete(shard.m, key)
						}
					}

				}
				shard.Unlock()
				// reset slice
				expiredKeys[i] = expiredKeys[i][:0]
			}

		}
	}()
	// start gc for dic cache
	go func() {
		c := globalDicCache
		expiredKeys := make(map[int][]string)
		for {
			// scan delay
			time.Sleep(gGCRunDelay)
			// select expired keys
			for i := 0; i < _CHM_SHARD_NUM; i++ {
				c.shards[i].RLock()

				for k := range c.shards[i].m {
					if c.shards[i].m[k].IsExpired() {
						expiredKeys[i] = append(expiredKeys[i], k)
					}
				}

				c.shards[i].RUnlock()
				// delay between shard scans
				time.Sleep(gGCShardScanDelay)
			}

			// remove selected keys in batch
			for i := 0; i < _CHM_SHARD_NUM; i++ {
				shard := &c.shards[i]
				shard.Lock()
				for _, key := range expiredKeys[i] {
					// recheck if expired
					v, ok := shard.m[key]
					if ok {
						if v.IsExpired() {
							delete(shard.m, key)
						}
					}

				}
				shard.Unlock()
				// reset slice
				expiredKeys[i] = expiredKeys[i][:0]
			}

		}
	}()

}

func Init() {
	// initialize caching system
	initOnce.Do(initialize)
}

type Atomic struct {
	val int64
}

func (a *Atomic) Set(val int64) {
	atomic.StoreInt64(&a.val, val)
}

func (a *Atomic) Get() int64 {
	return atomic.LoadInt64(&a.val)
}

func (a *Atomic) SetSysTS() {
	atomic.StoreInt64(&a.val, globalSysTimeNow.Get())
}

// returns ttl left in seconds if exists
func getTtl(ct CacheType, key []byte) int {
	var it item

	if key == nil {
		return KeyTTLErrCode
	}
	switch ct {
	case STRING_CACHE:
		shard := globalStringCache.shards[hash(key)&_MASK]
		shard.RLock()
		it = shard.m[string(key)].item
		shard.RUnlock()
	case LIST_CACHE:
		shard := globalListCache.shards[hash(key)&_MASK]
		shard.RLock()
		it = shard.m[string(key)].item
		shard.RUnlock()
	case DIC_CACHE:
		shard := globalDicCache.shards[hash(key)&_MASK]
		shard.RLock()
		it = shard.m[string(key)].item
		shard.RUnlock()
	default:
		return KeyTTLErrCode
	}

	if it.IsExpired() {
		return KeyNotExistCode
	}

	if it.ttl == KeyHasNoTTLCode {
		return KeyHasNoTTLCode
	}

	return it.FormatTTL(time.Second)
}

func (it item) IsExpired() bool {
	if it.ttl == KeyHasNoTTLCode || (it.ttl-CacheTimeNow()) > 0 {
		return false
	}
	return true
}

func (it *item) SetTTL(ttl int) {
	if ttl > 0 {
		// current time + ttl time in nanoseconds
		it.ttl = CacheTimeNow() + int64(ttl)*1e9
		return
	}

	it.ttl = KeyHasNoTTLCode
}

func (it *item) FormatTTL(in time.Duration) int {
	switch in {
	case time.Microsecond:
		return int((it.ttl - CacheTimeNow()) / 1e3)
	case time.Millisecond:
		return int((it.ttl - CacheTimeNow()) / 1e6)
	default:
		// in seconds
		return int((it.ttl - CacheTimeNow()) / 1e9)
	}
}
