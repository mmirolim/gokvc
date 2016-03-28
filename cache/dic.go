package cache

import "sync"

// Dic cache item struct
type Dic struct {
	item
	dic map[string][]byte
}

// NewDic initialize Dic
// creates map and sets first field and val
func NewDic(fld, val []byte, ttl int) Dic {
	var d Dic
	d.SetTTL(ttl)
	d.dic = make(map[string][]byte)
	d.dic[string(fld)] = val
	return d
}

// DicCache is bucket holding list of caches
// in stripped map
type DicCache struct {
	shards [_CHM_SHARD_NUM]struct {
		sync.RWMutex
		m   map[string]Dic
		pad [128]byte
	}
}

// DGET returns map[string][]byte of values in Dic
// by key, nil, false if not exists or expired
func DGET(key []byte) (map[string][]byte, bool) {
	return globalDicCache.get(key)
}

// DFGET returns field value from Dic defined by key
// or nil, false if not exist
func DFGET(key, fld []byte) ([]byte, bool) {
	return globalDicCache.fget(key, fld)
}

// DFSET add field and val to Dic, if not exists
// it initialize and add to cache
func DFSET(key, fld, val []byte, ttl int) bool {
	return globalDicCache.fset(key, fld, val, ttl)
}

// DFDEL removes fiedl from Dic by key
// false if not exists or expired
func DFDEL(key, fld []byte) bool {
	return globalDicCache.fdel(key, fld)
}

// DDEL removes Dic by key
// returns false if not exists or expired
func DDEL(key []byte) bool {
	return globalDicCache.del(key)
}

// DTTL returns ttl in seconds of key
// of ttl codes
func DTTL(key []byte) int {
	return getTtl(DIC_CACHE, key)
}

// DLEN returns number of Dic cached
func DLEN() int {
	return globalDicCache.countKeys()
}

// DKEYS returns []string of all keys of Dics cached
func DKEYS() []string {
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
		res = make(map[string][]byte)
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
		val, ok = v.dic[string(fld)]
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
	v, ok := shard.m[string(key)]
	if !ok || v.IsExpired() {
		shard.Unlock()
		return false
	}

	delete(shard.m[string(key)].dic, string(fld))
	if len(shard.m[string(key)].dic) == 0 {
		// no elements in dic
		delete(shard.m, string(key))
	}

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
func (c *DicCache) keys() []string {
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
