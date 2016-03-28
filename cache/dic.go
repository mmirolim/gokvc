package cache

import "sync"

type Dic struct {
	item
	dic map[string][]byte
}

func NewDic(fld, val []byte, ttl int) Dic {
	var d Dic
	d.SetTTL(ttl)
	d.dic = make(map[string][]byte)
	d.dic[string(fld)] = val
	return d
}

type DicCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]Dic
		pad [128]byte
	}
}

func DGET(key []byte) (map[string][]byte, bool) {
	return globalDicCache.get(key)
}

func DFGET(key, fld []byte) ([]byte, bool) {
	return globalDicCache.fget(key, fld)
}

func DFSET(key, fld, val []byte, ttl int) bool {
	return globalDicCache.fset(key, fld, val, ttl)
}

func DFDEL(key, fld []byte) bool {
	return globalDicCache.fdel(key, fld)
}

func DDEL(key []byte) bool {
	return globalDicCache.del(key)
}

func DLEN() int {
	return globalDicCache.countKeys()
}

func DKEYS() [][]byte {
	return globalDicCache.keys()
}

func (c *DicCache) get(key []byte) (map[string][]byte, bool) {
	var res map[string][]byte
	if key == nil {
		return res, false
	}

	shard := &c.shards[hash(key)&_MASK]
	shard.RLock()

	v, ok := shard.m[string(key)]
	if !ok || (ok && v.IsExpired()) {
		shard.RUnlock()
		return nil, false
	} else {
		res := make(map[string][]byte)
		for k, v := range v.dic {
			res[k] = v
		}
	}

	shard.RUnlock()

	return res, ok
}

func (c *DicCache) fset(key, fld, val []byte, ttl int) bool {
	if key == nil || fld == nil || val == nil {
		return false
	}
	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()

	v, ok := shard.m[string(key)]
	if ok {
		v.dic[string(fld)] = val
		v.SetTTL(ttl)
		shard.m[string(key)] = v
	} else {
		shard.m[string(key)] = NewDic(fld, val, ttl)
	}

	shard.Unlock()

	return true
}

func (c *DicCache) fget(key, fld []byte) ([]byte, bool) {
	var val []byte
	if key == nil || fld == nil {
		return val, false
	}

	shard := &c.shards[hash(key)&_MASK]
	shard.RLock()

	v, ok := shard.m[string(key)]
	if !ok || (ok && v.IsExpired()) {
		shard.RUnlock()
		return nil, false
	} else {
		val = v.dic[string(fld)]
	}

	shard.RUnlock()

	return val, ok
}

func (c *DicCache) fdel(key, fld []byte) bool {
	if key == nil || fld == nil {
		return false
	}

	shard := &c.shards[hash(key)&_MASK]
	shard.Lock()

	delete(shard.m[string(key)].dic, string(fld))

	shard.Unlock()

	return true
}

func (c *DicCache) del(key []byte) bool {
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
