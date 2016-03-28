package cache

import (
	"container/list"
	"sync"
)

type List struct {
	item
	ll *list.List
}

func NewList(key, val []byte, ttl int) List {
	it := item{k: key}
	it.SetTTL(ttl)

	ll := list.New()
	ll.PushFront(val)

	return List{it, ll}
}

type ListCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]List
		pad [128]byte
	}
}

func LGET(key []byte) ([][]byte, bool) {
	return globalListCache.get(key)
}

func LPUSH(key, val []byte, ttl int) {
	globalListCache.push(key, val, ttl)
}

func LTTL(key []byte) int {
	return getTtl(LIST_CACHE, key)
}

func LPOP(key []byte) ([]byte, bool) {
	return globalListCache.pop(key)
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

func (c *ListCache) get(key []byte) ([][]byte, bool) {
	var vals [][]byte

	shard := &c.shards[hash(key)&_MASK]
	shard.RLock()

	v, ok := shard.m[string(key)]
	if !ok || (ok && v.IsExpired()) {
		shard.RUnlock()
		return nil, false
	}

	if ok {
		vals := make([][]byte, 0, 10)

		for e := v.ll.Front(); e != nil; e = e.Next() {
			vals = append(vals, e.Value.([]byte))
		}
	}

	shard.RUnlock()

	return vals, ok
}

func (c *ListCache) push(key, val []byte, ttl int) {
	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()

	v, ok := shard.m[string(key)]
	if ok {
		v.ll.PushFront(val)
		v.SetTTL(ttl)
		shard.m[string(key)] = v
	} else {
		// create list if not exists
		shard.m[string(key)] = NewList(key, val, ttl)
	}

	shard.Unlock()
}

func (c *ListCache) pop(key []byte) ([]byte, bool) {
	var val []byte
	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()

	v, ok := shard.m[string(key)]
	if !ok || (ok && v.IsExpired()) {
		shard.Unlock()
		return val, false
	} else {
		val = v.ll.Remove(v.ll.Front()).([]byte)
	}

	shard.Unlock()

	return val, ok
}

func (c *ListCache) del(key []byte) {
	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()

	delete(shard.m, string(key))

	shard.Unlock()
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
