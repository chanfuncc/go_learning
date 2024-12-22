package syncmap

type ChannelMap[K comparable, V any] struct {
	ch chan struct{}
	m  map[K]V
}

func NewChannelMap[K comparable, V any]() *ChannelMap[K, V] {
	return &ChannelMap[K, V]{
		ch: make(chan struct{}, 1),
		m:  make(map[K]V),
	}
}

func (m *ChannelMap[K, V]) Get(key K) (V, bool) {
	m.ch <- struct{}{}
	defer func() { <-m.ch }()
	v, ok := m.m[key]
	return v, ok
}

func (m *ChannelMap[K, V]) Set(key K, value V) {
	m.ch <- struct{}{}
	defer func() { <-m.ch }()
	m.m[key] = value
}

func (m *ChannelMap[K, V]) Delete(key K) {
	m.ch <- struct{}{}
	defer func() { <-m.ch }()
	delete(m.m, key)
}

func (m *ChannelMap[K, V]) Range(f func(key K, value V) bool) {
	m.ch <- struct{}{}
	defer func() { <-m.ch }()
	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
}
