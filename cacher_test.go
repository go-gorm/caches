package caches

import (
	"context"
	"errors"
	"sync"
)

type cacherMock struct {
	store *sync.Map
}

func (c *cacherMock) init() {
	if c.store == nil {
		c.store = &sync.Map{}
	}
}

func (c *cacherMock) Get(ctx context.Context, key string) *Query {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	return val.(*Query)
}

func (c *cacherMock) Store(ctx context.Context, key string, val *Query) error {
	c.init()
	c.store.Store(key, val)
	return nil
}

type cacherStoreErrorMock struct {
	store *sync.Map
}

func (c *cacherStoreErrorMock) init() {
	if c.store == nil {
		c.store = &sync.Map{}
	}
}

func (c *cacherStoreErrorMock) Get(ctx context.Context, key string) *Query {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	return val.(*Query)
}

func (c *cacherStoreErrorMock) Store(context.Context, string, *Query) error {
	return errors.New("store-error")
}
