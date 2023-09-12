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
func (c *cacherMock) Get(key string) *Query {
	return c.GetContext(context.Background(), key)
}

func (c *cacherMock) GetContext(ctx context.Context, key string) *Query {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	return val.(*Query)
}

func (c *cacherMock) Store(key string, val *Query) error {
	return c.StoreContext(context.Background(), key, val)
}

func (c *cacherMock) StoreContext(ctx context.Context, key string, val *Query) error {
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

func (c *cacherStoreErrorMock) Get(key string) *Query {
	return c.GetContext(context.Background(), key)
}

func (c *cacherStoreErrorMock) GetContext(ctx context.Context, key string) *Query {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	return val.(*Query)
}

func (c *cacherStoreErrorMock) Store(key string, query *Query) error {
	return c.StoreContext(context.Background(), key, query)
}

func (c *cacherStoreErrorMock) StoreContext(context.Context, string, *Query) error {
	return errors.New("store-error")
}
