package cache

import (
	"sync"

	"github.com/golang/glog"
)

type String struct {
	item
	b []byte
}

type StringCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]String
		pad [128]byte
	}
}

func (s *StringCache) Len() int {
	var counter int
	for i := 0; i < _CHM_SHARD_NUM; i++ {
		s.shards[i].RLock()
		counter += len(s.shards[i].m)
		s.shards[i].RUnlock()
	}
	return counter
}

func GET(key []byte) ([]byte, bool) {
	return globalStringCache.get(key)
}

func SET(key, val []byte, ttl int) {
	globalStringCache.set(key, val, ttl)
}

func DEL(key []byte) {
	globalStringCache.del(key)
}

func (s *StringCache) get(key []byte) ([]byte, bool) {
	shard := &s.shards[hash(key)&_MASK]
	shard.RLock()
	v, ok := shard.m[string(key)]
	shard.RUnlock()
	// check ttl
	if ok && v.ttl != 0 && CacheTimeNow() > v.ttl {
		return nil, false
	}

	if glog.V(4) {
		glog.Infof("sget key %s found %#v val %s", key, v, v.b)
	}
	return v.b, ok
}

// ttl in seconds
func (s *StringCache) set(key, val []byte, ttl int) {
	var str String
	str.b = append(str.b, val...)
	if ttl > 0 {
		// item ttl in nano seconds from epoch
		str.ttl = CacheTimeNow() + int64(ttl)*1e9
	}
	if glog.V(4) {
		glog.Infof("sset key %s, val %s, ttl %d, element %#v", key, val, ttl, str)
	}
	shard := &s.shards[hash(key)&_MASK]
	shard.Lock()
	shard.m[string(key)] = str
	shard.Unlock()
}

func (s *StringCache) del(key []byte) {
	shard := &s.shards[hash(key)&_MASK]
	shard.Lock()
	delete(shard.m, string(key))
	shard.Unlock()
}
