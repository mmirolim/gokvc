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

func (s *DicCache) len() int {
	var counter int
	for i := 0; i < _CHM_SHARD_NUM; i++ {
		s.shards[i].RLock()
		counter += len(s.shards[i].m)
		s.shards[i].RUnlock()
	}
	return counter
}
