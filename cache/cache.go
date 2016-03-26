package cache

import (
	"container/list"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	_CHM_SHARD_NUM = 1 << 6
	_MASK          = _CHM_SHARD_NUM - 1

	// kv commands
	SGET = "GET" // GET key
	SDEL = "DEL" // DEL key
	SSET = "SET" // SET key val
	STTL = "TTL" // TTL key gets seconds left o expire

	// list commands
	LGET  = "LGET"  // LGET key get all list
	LDEL  = "LDEL"  // LDEL key deletes list
	LPUSH = "LPUSH" // LPUSH key val prepends list with val
	LLEN  = "LLEN"  // LLEN key gets length of list
	LTTL  = "LTTL"  // TTL key gets seconds left o expire

	// dictionary commands
	DGET   = "DGET"   // DGET key get all field from dic
	DDEL   = "DDEL"   // DDEL key delete dic
	DKGET  = "DKGET"  // DKGET key field get field from dic
	DKDEL  = "DKDEL"  // DKDEL key field delete field in dic
	DKSET  = "DKSET"  // DKSET key field val sets field in dic to val
	DKSGET = "DKSGET" // DKSGET key gets all fields in dic
	DTTL   = "DTTL"   // TTL key gets seconds left o expire

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
	// TTL passed with http headers
	KEYTTL = []byte("KEYTTL") // KEYTTL seconds

	initOnce sync.Once

	globalSysTimeNow Atomic

	globalStringCache = &StringCache{}
	globalListCache   = &ListCache{}
	globalDicCache    = &DicCache{}

	hasher = fnv.New32()
)

func CacheTimeNow() int64 {
	return globalSysTimeNow.Get()
}

type CacheType int

type item struct {
	ttl int64 // 0 is no ttl
	la  int64
}

type List struct {
	item
	ll *list.List
}

type Dic struct {
	item
	dic map[string]element
}

type element struct {
	k []byte // key
	v []byte // value
}

type Cacher interface {
	Len() int
	DEL(key []byte)
	TTL(key []byte) int
}

type ListCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]List
		pad [128]byte
	}
}

type DicCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]Dic
		pad [128]byte
	}
}

// returns ttl left in seconds if exists
func getTtl(ct CacheType, key []byte) int {
	var ttl int64
	switch ct {
	case STRING_CACHE:
		shard := globalStringCache.shards[hash(key)&_MASK]
		shard.RLock()
		ttl = shard.m[string(key)].ttl
		shard.RUnlock()
	case LIST_CACHE:
		shard := globalListCache.shards[hash(key)&_MASK]
		shard.RLock()
		ttl = shard.m[string(key)].ttl
		shard.RUnlock()
	case DIC_CACHE:
		shard := globalDicCache.shards[hash(key)&_MASK]
		shard.RLock()
		ttl = shard.m[string(key)].ttl
		shard.RUnlock()
	}

	if ttl > 0 {
		return int(ttl - CacheTimeNow()/1e9)
	}

	return KeyNotExistCode
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

func (s *ListCache) Len() int {
	var counter int
	for i := 0; i < _CHM_SHARD_NUM; i++ {
		s.shards[i].RLock()
		counter += len(s.shards[i].m)
		s.shards[i].RUnlock()
	}
	return counter
}

func (s *DicCache) Len() int {
	var counter int
	for i := 0; i < _CHM_SHARD_NUM; i++ {
		s.shards[i].RLock()
		counter += len(s.shards[i].m)
		s.shards[i].RUnlock()
	}
	return counter
}
