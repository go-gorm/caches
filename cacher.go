package caches

import "context"

type Cacher interface {
	Get(ctx context.Context, key string) *Query
	Store(ctx context.Context, key string, val *Query) error
}
