package cache

import "sync"

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

func (s *StringCache) GET(key []byte) ([]byte, bool) {
	shard := &s.shards[hash(key)&_MASK]
	shard.RLock()
	v, ok := shard.m[string(key)]
	shard.RUnlock()
	// check ttl
	if ok && v.ttl != 0 && CacheTimeNow() > v.ttl {
		return nil, false
	}
	return v.b, ok
}

// ttl in seconds
func (s *StringCache) SET(key, val []byte, ttl int) {
	var str String
	copy(str.b, val)
	if ttl > 0 {
		// item ttl in nano seconds from epoch
		str.ttl = CacheTimeNow() + int64(ttl)*1e9
	}
	shard := &s.shards[hash(key)&_MASK]
	shard.Lock()
	shard.m[string(key)] = str
	shard.Unlock()
}

func (s *StringCache) DEL(key []byte) {
	shard := &s.shards[hash(key)&_MASK]
	shard.Lock()
	delete(shard.m, string(key))
	shard.Unlock()
}

// returns ttl left in seconds if exists
func (s *StringCache) TTL(key []byte) int {
	shard := &s.shards[hash(key)&_MASK]
	shard.RLock()
	v, ok := shard.m[string(key)]
	shard.RUnlock()
	switch {
	case !ok:
		return KeyNotExistCode
	case v.ttl == 0:
		return KeyHasNoTTLCode
	case v.ttl > 0:
		return int(v.ttl - CacheTimeNow()/1e9)
	default:
		return KeyTTLErrCode
	}
}
