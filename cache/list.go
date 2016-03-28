package cache

import (
	"container/list"
	"sync"
)

// List cache item struct
type List struct {
	item
	ll *list.List
}

// NewList initialize List
// creates list and sets first element value
// to val and sets ttl
func NewList(val []byte, ttl int) List {
	var it item
	it.SetTTL(ttl)
	ll := list.New()
	ll.PushFront(val)

	return List{it, ll}
}

// ListCache is bucket holding list caches
// in stripped map
type ListCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]List
		pad [128]byte
	}
}

// LGET returns [][]byte of values stored in List
// by a key, return nil, false if not exists or key is nil
func LGET(key []byte) ([][]byte, bool) {
	return globalListCache.get(key)
}

// LPUSH create List item in cache if not exists with ttl
// and prepends element to a list
func LPUSH(key, val []byte, ttl int) bool {
	return globalListCache.push(key, val, ttl)
}

// LTTL returns ttl of key
// or ttl codes
func LTTL(key []byte) int {
	return getTtl(LIST_CACHE, key)
}

// LPOP returns and removes first element in list
// if it was last element, list item deleted from cache
// returns nil, false if not exist
func LPOP(key []byte) ([]byte, bool) {
	return globalListCache.pop(key)
}

// LDEL removes list element by key
func LDEL(key []byte) bool {
	return globalListCache.del(key)
}

// LLEN returns number of not expired cached lists
func LLEN() int {
	return globalListCache.countKeys()
}

// LKEYS returns []string of all not expired keys
func LKEYS() []string {
	return globalListCache.keys()
}

func (c *ListCache) get(key []byte) ([][]byte, bool) {
	var vals [][]byte
	if key == nil {
		return vals, false
	}
	shard := &c.shards[hash(key)&_MASK]
	shard.RLock()

	v, ok := shard.m[string(key)]
	if !ok || (ok && v.IsExpired()) {
		shard.RUnlock()
		return nil, false
	} else {
		vals = make([][]byte, 0, 10)
		for e := v.ll.Front(); e != nil; e = e.Next() {
			vals = append(vals, e.Value.([]byte))
		}
	}

	shard.RUnlock()

	return vals, ok
}

func (c *ListCache) push(key, val []byte, ttl int) bool {
	if key == nil || val == nil {
		return false
	}

	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()

	v, ok := shard.m[string(key)]
	if ok {
		v.ll.PushFront(val)
		v.SetTTL(ttl)
		shard.m[string(key)] = v
	} else {
		// create list if not exists
		shard.m[string(key)] = NewList(val, ttl)
	}

	shard.Unlock()
	return true
}

func (c *ListCache) pop(key []byte) ([]byte, bool) {
	var val []byte
	if key == nil {
		return val, false
	}

	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()

	v, ok := shard.m[string(key)]
	if !ok || (ok && v.IsExpired()) {
		shard.Unlock()
		return val, false
	} else {
		val = v.ll.Remove(v.ll.Front()).([]byte)
		if v.ll.Len() == 0 {
			// no elements del item
			delete(shard.m, string(key))
		}
	}

	shard.Unlock()

	return val, ok
}

func (c *ListCache) del(key []byte) bool {
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
func (c *ListCache) keys() []string {
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
