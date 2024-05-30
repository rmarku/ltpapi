package keyvalue

import (
	"github.com/stretchr/testify/mock"
)

type MockFloatCache struct {
	mock.Mock
}

func (m *MockFloatCache) Get(key string) (float64, error) {
	args := m.Called(key)

	return args.Get(0).(float64), args.Error(1) //nolint: forcetypeassert, wrapcheck
}

func (m *MockFloatCache) Set(key string, value float64) error {
	args := m.Called(key, value)

	return args.Error(0) //nolint:  wrapcheck
}

func (m *MockFloatCache) Delete(key string) error {
	args := m.Called(key)

	return args.Error(0) //nolint:  wrapcheck
}
