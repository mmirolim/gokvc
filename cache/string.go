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

func GET(key []byte) ([]byte, bool) {
	return globalStringCache.get(key)
}

func SET(key, val []byte, ttl int) bool {
	return globalStringCache.set(key, val, ttl)
}

func DEL(key []byte) bool {
	return globalStringCache.del(key)
}

func LEN() int {
	return globalStringCache.countKeys()
}

func KEYS() [][]byte {
	return globalStringCache.keys()
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

	str.k = key
	str.b = val
	str.SetTTL(ttl)

	shard := &c.shards[hash(key)&_MASK]
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
func (c *StringCache) keys() [][]byte {
	// init cap for slice
	keys := make([][]byte, 0, 1000)
	for i := 0; i < _CHM_SHARD_NUM; i++ {
		shard := c.shards[i]
		shard.RLock()

		for k := range shard.m {
			keys = append(keys, shard.m[k].k)
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
