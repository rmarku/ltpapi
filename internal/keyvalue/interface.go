package keyvalue

import "errors"

type FloatCache interface {
	Get(key string) (float64, error)
	Set(key string, value float64) error
	Delete(key string) error
}

var ErrKeyNotFound = errors.New("key not found")
var ErrExpired = errors.New("value expired")
