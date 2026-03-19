package store

import "sync"

type Repository[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewRepository[K comparable, V any]() *Repository[K, V] {
	return &Repository[K, V]{data: make(map[K]V)}
}

func (r *Repository[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.data[key]
	return v, ok
}

func (r *Repository[K, V]) Set(key K, value V) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
}

func (r *Repository[K, V]) Delete(key K) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, key)
}

func (r *Repository[K, V]) GetAll() map[K]V {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[K]V, len(r.data))
	for k, v := range r.data {
		out[k] = v
	}
	return out
}

func (r *Repository[K, V]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.data)
}
