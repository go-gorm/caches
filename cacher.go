package caches

import (
	"context"
)

type Cacher interface {
	// Get impl should check if a specific key exists in the cache and return its value
	// look at Query.Marshal
	Get(ctx context.Context, key string, q *Query[any]) (*Query[any], error)
	// Store is supposed to store a cached representation of the val param
	// look at Query.Unmarshal
	Store(ctx context.Context, key string, val *Query[any]) error
}
