package syncmap

import (
	"sync"
	"testing"
)

type MapInterface[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
	Delete(key K)
	Range(func(key K, value V) bool)
}

func testMap[K comparable, V any](t *testing.T, name string, createMap func() MapInterface[int, int]) {
	t.Run(name, func(t *testing.T) {
		t.Run("basic operations", func(t *testing.T) {
			m := createMap()

			m.Set(1, 1)
			if v, ok := m.Get(1); !ok || v != 1 {
				t.Errorf("expected 1, got %v", v)
			}

			m.Delete(1)
			if _, ok := m.Get(1); ok {
				t.Error("key should be deleted")
			}
		})

		t.Run("concurrent operations", func(t *testing.T) {
			m := createMap()
			const workers = 10
			const ops = 1000

			var wg sync.WaitGroup
			wg.Add(workers * 2)

			for i := 0; i < workers; i++ {
				go func(id int) {
					defer wg.Done()
					for j := 0; j < ops; j++ {
						key := id*ops + j
						m.Set(key, key)
					}
				}(i)
			}

			for i := 0; i < workers; i++ {
				go func(id int) {
					defer wg.Done()
					for j := 0; j < ops; j++ {
						key := id*ops + j
						if val, ok := m.Get(key); ok && val != key {
							t.Errorf("got %v, want %v", val, key)
						}
					}
				}(i)
			}

			wg.Wait()

			count := 0
			m.Range(func(k, v int) bool {
				if k != v {
					t.Errorf("key %v != value %v", k, v)
				}
				count++
				return true
			})

			if count != workers*ops {
				t.Errorf("expected %d entries, got %d", workers*ops, count)
			}
		})
	})
}

func TestMaps(t *testing.T) {
	// 测试 MutexMap
	testMap[int, int](t, "MutexMap", func() MapInterface[int, int] {
		return NewMutexMap[int, int]()
	})

	// 测试 ChannelMap
	testMap[int, int](t, "ChannelMap", func() MapInterface[int, int] {
		return NewChannelMap[int, int]()
	})
}
