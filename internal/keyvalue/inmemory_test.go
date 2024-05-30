package keyvalue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryFloatCache(t *testing.T) {
	// Create a new instance of the InMemoryFloatCache
	t.Parallel()

	cache := NewInMemory()

	tests := []struct {
		expectedError error
		setup         func()
		name          string
		key           string
		expectedValue float64
	}{
		{
			name: "key not found",
			setup: func() {
				// No setup required, key should not exist
			},
			key:           "nonexistent",
			expectedValue: 0,
			expectedError: ErrKeyNotFound,
		},
		{
			name: "key expired",
			setup: func() {
				cache.Set("expired", 123.45) //nolint: errcheck
				// Manually expire the key
				cache.data["expired"] = entry{
					value:  123.45,
					expiry: time.Now().Add(-time.Minute),
				}
			},
			key:           "expired",
			expectedValue: 0,
			expectedError: ErrExpired,
		},
		{
			name: "key valid",
			setup: func() {
				cache.Set("valid", 123.45) //nolint: errcheck
			},
			key:           "valid",
			expectedValue: 123.45,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup()

			value, err := cache.Get(tt.key)

			assert.Equal(t, tt.expectedValue, value) //nolint: testifylint
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestInMemoryFloatCacheSet(t *testing.T) {
	t.Parallel()

	cache := NewInMemory()

	key := "testkey"
	value := 123.45

	err := cache.Set(key, value)
	require.NoError(t, err)

	storedValue, err := cache.Get(key)
	require.NoError(t, err)
	assert.InDelta(t, value, storedValue, 1e-5)
}

func TestInMemoryFloatCacheDelete(t *testing.T) {
	t.Parallel()

	cache := NewInMemory()

	key := "testkey"
	value := 123.45

	cache.Set(key, value) //nolint: errcheck

	err := cache.Delete(key)
	require.NoError(t, err)

	_, err = cache.Get(key)
	require.Error(t, err)
	assert.Equal(t, ErrKeyNotFound, err)
}
