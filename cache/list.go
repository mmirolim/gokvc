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

func (s *ListCache) len() int {
	var counter int
	for i := 0; i < _CHM_SHARD_NUM; i++ {
		s.shards[i].RLock()
		counter += len(s.shards[i].m)
		s.shards[i].RUnlock()
	}
	return counter
}
