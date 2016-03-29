package cache

import (
	"bytes"
	"sync"
)

// String cache item struct
type String struct {
	item
	b []byte
}

// StringCache is bucket holding string caches
// in stripped map
type StringCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]String
		pad [128]byte
	}
}

// GET the value of a key,
// returns []byte, true
// or nil, false if expired or not exist
func GET(key []byte) ([]byte, bool) {
	return globalStringCache.get(key)
}

// SET the key with value and ttl
// ttl 0 means no expire
// returns false if key, val is nil
func SET(key, val []byte, ttl int) bool {
	return globalStringCache.set(key, val, ttl)
}

// DEL removes element by key
func DEL(key []byte) bool {
	return globalStringCache.del(key)
}

// LEN returns number of not expired keys
// of stored strings
func LEN() int {
	return globalStringCache.countKeys()
}

// KEYS returns []string of all not expired keys
// of stored strings
func KEYS() []string {
	return globalStringCache.keys()
}

// TTL returns ttl in seconds of key
// of ttl codes
func TTL(key []byte) int {
	return getTtl(STRING_CACHE, key)
}

func (c *StringCache) get(key []byte) ([]byte, bool) {
	if key == nil {
		return nil, false
	}
	shard := &c.shards[hash(key)&_MASK]
	shard.RLock()

	v, ok := shard.m[string(key)]

	shard.RUnlock()

	if !ok || (ok && v.IsExpired()) {
		return nil, false
	}

	return v.b, ok
}

// ttl in seconds
func (c *StringCache) set(key, val []byte, ttl int) bool {
	if key == nil || val == nil {
		return false
	}
	var str String

	str.b = val
	str.SetTTL(ttl)

	shard := &c.shards[hash(key)&_MASK]
	shard.RLock()
	v, ok := shard.m[string(key)]
	if ok {
		if bytes.Equal(str.b, v.b) && v.ttl == str.ttl {
			shard.RUnlock()
			return true
		}
	}
	shard.RUnlock()
	shard.Lock()

	shard.m[string(key)] = str

	shard.Unlock()
	return true
}

func (c *StringCache) del(key []byte) bool {
	if key == nil {
		return false
	}
	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()

	delete(shard.m, string(key))

	shard.Unlock()
	return true
}

// returns slice of keys
func (c *StringCache) keys() []string {
	// init cap for slice
	keys := make([]string, 0, 1000)
	for i := 0; i < _CHM_SHARD_NUM; i++ {
		shard := c.shards[i]
		shard.RLock()

		for k := range shard.m {
			keys = append(keys, k)
		}

		shard.RUnlock()
	}

	return keys
}

func (c *StringCache) countKeys() int {
	var counter int
	for i := 0; i < _CHM_SHARD_NUM; i++ {
		c.shards[i].RLock()

		for k := range c.shards[i].m {
			// count only not expired keys
			if !c.shards[i].m[k].IsExpired() {
				counter++
			}
		}

		c.shards[i].RUnlock()
	}
	return counter
}
