package caches

type Cacher interface {
	Get(key string) interface{}
	Store(key string, val interface{}) error
}
