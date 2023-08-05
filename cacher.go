package caches

type Cacher interface {
	Get(key string) *Query
	Store(key string, val *Query) error
}
