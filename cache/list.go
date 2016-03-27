package cache

import (
	"container/list"
	"sync"
)

type List struct {
	item
	ll *list.List
}

type ListCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]List
		pad [128]byte
	}
}

func LDEL(key []byte) {
	globalListCache.del(key)
}

func LLEN() int {
	return globalListCache.countKeys()
}

func LKEYS() [][]byte {
	return globalListCache.keys()
}

func (c *ListCache) del(key []byte) {
	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()
	delete(shard.m, string(key))
	shard.Unlock()
}

func (c *ListCache) countKeys() int {
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

// returns slice of keys
func (c *ListCache) keys() [][]byte {
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
