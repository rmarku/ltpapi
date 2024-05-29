package keyvalue

import (
	"time"
)

type entry struct {
	expiry time.Time
	value  float64
}
type InMemoryFloatCache struct {
	data     map[string]entry
	duration time.Duration
}

var _ FloatCache = new(InMemoryFloatCache)

func NewInMemory() *InMemoryFloatCache {
	return &InMemoryFloatCache{
		data:     make(map[string]entry),
		duration: time.Minute,
	}
}

func (kv *InMemoryFloatCache) Get(key string) (float64, error) {
	value, ok := kv.data[key]

	if !ok {
		return 0, ErrKeyNotFound
	}

	if time.Now().After(value.expiry) {
		return 0, ErrExpired
	}

	return value.value, nil
}

func (kv *InMemoryFloatCache) Set(key string, value float64) error {
	now := time.Now()
	kv.data[key] = entry{
		value:  value,
		expiry: now.Add(kv.duration),
	}

	return nil
}

func (kv *InMemoryFloatCache) Delete(key string) error {
	delete(kv.data, key)

	return nil
}
