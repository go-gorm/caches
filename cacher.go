package caches

import "context"

type Cacher interface {
	Get(key string) *Query
	GetContext(ctx context.Context, key string) *Query
	Store(key string, val *Query) error
	StoreContext(ctx context.Context, key string, val *Query) error
}
