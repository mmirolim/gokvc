package cache

import "sync"

type Dic struct {
	item
	dic map[string]element
}

type element struct {
	k []byte // key
	v []byte // value
}

type DicCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]Dic
		pad [128]byte
	}
}

func DDEL(key []byte) {
	globalDicCache.del(key)
}

func DLEN() int {
	return globalDicCache.countKeys()
}

func DKEYS() [][]byte {
	return globalDicCache.keys()
}

func (c *DicCache) del(key []byte) {
	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()
	delete(shard.m, string(key))
	shard.Unlock()
}

// returns slice of keys
func (c *DicCache) keys() [][]byte {
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

func (c *DicCache) countKeys() int {
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
