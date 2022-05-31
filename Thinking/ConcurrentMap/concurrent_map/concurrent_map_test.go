package concurrent_map

import (
	"fmt"
	"sync/atomic"
	"testing"
)

// keyPool 生成不同的key
type keyPool struct {
	i    int32
	Keys []string
}

func newKeyPool(len int) *keyPool {
	Keys := make([]string, len, len)
	for i := 0; i < len; i++ {
		Keys[i] = fmt.Sprint("k", i)
	}

	return &keyPool{
		i:    0,
		Keys: Keys,
	}
}

func (kp *keyPool) Get() string {
	i := atomic.AddInt32(&kp.i, 1)
	return kp.Keys[int(i)%len(kp.Keys)]
}

func BenchmarkKeyPool(b *testing.B) {
	keyPool := newKeyPool(100000)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			keyPool.Get()
		}
	})
}

// 裸map性能
func BenchmarkMap(b *testing.B) {
	m := map[string]interface{}{}

	keyPool := newKeyPool(1000000)
	//
	for i := 0; i < 1000000; i++ {
		m[keyPool.Get()] = "v"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m[keyPool.Get()]
		// m[keyPool.Get()] = "v"
	}
}

func BenchmarkConcurrenceMap(b *testing.B) {
	type test struct {
		name string
		m    ConcurrentMap
	}

	tests := []struct {
		name string
		m    ConcurrentMap
	}{
		{
			name: "LockMap",
			m:    NewLockMap(),
		},
		{
			name: "SyncMap",
			m:    NewSyncMap(),
		},
		{
			name: "SliceMap",
			m:    NewSliceMap(),
		},
	}

	keyPool := newKeyPool(100000)

	multiRun := func(name string, fn func(m ConcurrentMap, b *testing.B)) {
		for _, tt := range tests {
			b.Run(tt.name+name, func(b *testing.B) { fn(tt.m, b) })
		}
	}

	// 命中写
	// multiRun("_HitSet", func(m ConcurrentMap, b *testing.B) {
	// 	for _, v := range keyPool.Keys {
	// 		m.Set(v, "v")
	// 	}
	// 	b.ReportAllocs()
	// 	b.ResetTimer()
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			m.Set(keyPool.Get(), "v")
	// 		}
	// 	})
	// })
	//
	// 未命中写
	multiRun("_MissSet", func(m ConcurrentMap, b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m.Get(keyPool.Get())
			}
		})
	})

	// 命中读
	multiRun("_HitGet", func(m ConcurrentMap, b *testing.B) {
		for _, v := range keyPool.Keys {
			m.Set(v, "v")
		}

		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m.Get(keyPool.Get())
			}
		})
	})
	//
	// // 未命中读
	// multiRun("_MissGet", func(m ConcurrentMap, b *testing.B) {
	// 	b.ReportAllocs()
	// 	b.ResetTimer()
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			m.Get(keyPool.Get())
	// 		}
	// 	})
	// })

	// 异步读多于写
	multiRun("_GetMoreThanSet", func(m ConcurrentMap, b *testing.B) {
		var i int32

		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				key := keyPool.Get()
				ni := atomic.AddInt32(&i, 1)
				if ni%2000 != 0 {
					m.Get(key)
				} else {
					m.Set(key, "k")
				}
			}
		})
	})

	// 异步写多于读
	// multiRun("_SetMoreThanGet", func(m ConcurrentMap, b *testing.B) {
	// 	var i int32
	//
	// 	b.ReportAllocs()
	// 	b.ResetTimer()
	//
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			ni := atomic.AddInt32(&i, 1)
	// 			if ni%10000 != 0 {
	// 				m.Set(keyPool.Get(), "k")
	// 			} else {
	// 				m.Get(keyPool.Get())
	// 			}
	// 		}
	// 	})
	// })
}
