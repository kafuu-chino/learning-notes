package concurrent_map

import (
	"sync"

	cmap "github.com/orcaman/concurrent-map"
)

type ConcurrentMap interface {
	Set(k string, v interface{})
	Get(k string) (interface{}, bool)
}

// LockMap 加锁map
type LockMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func NewLockMap() *LockMap {
	return &LockMap{
		m: map[string]interface{}{},
	}
}

func (lm *LockMap) Get(k string) (interface{}, bool) {
	lm.RLock()
	v, ok := lm.m[k]
	lm.RUnlock()

	return v, ok
}

func (lm *LockMap) Set(k string, v interface{}) {
	lm.Lock()
	lm.m[k] = v
	lm.Unlock()
}

// SyncMap 读写分离map
type SyncMap struct {
	m sync.Map
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		m: sync.Map{},
	}
}

func (sm *SyncMap) Get(k string) (interface{}, bool) {
	return sm.m.Load(k)
}

func (sm *SyncMap) Set(k string, v interface{}) {
	sm.m.Store(k, v)
}

// SliceMap 分片map
type SliceMap struct {
	m cmap.ConcurrentMap
}

func NewSliceMap() *SliceMap {
	return &SliceMap{
		m: cmap.New(),
	}
}

func (sm *SliceMap) Get(k string) (interface{}, bool) {
	return sm.m.Get(k)
}

func (sm *SliceMap) Set(k string, v interface{}) {
	sm.m.Set(k, v)
}
