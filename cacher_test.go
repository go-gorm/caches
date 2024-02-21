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

func (c *cacherMock) Get(_ context.Context, key string, _ *Query[any]) (*Query[any], error) {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil, nil
	}

	return val.(*Query[any]), nil
}

func (c *cacherMock) Store(_ context.Context, key string, val *Query[any]) error {
	c.init()
	c.store.Store(key, val)
	return nil
}

func (c *cacherMock) Invalidate(context.Context) error {
	return nil
}

type cacherStoreErrorMock struct{}

func (c *cacherStoreErrorMock) Get(context.Context, string, *Query[any]) (*Query[any], error) {
	return nil, nil
}

func (c *cacherStoreErrorMock) Store(context.Context, string, *Query[any]) error {
	return errors.New("store-error")
}

func (c *cacherStoreErrorMock) Invalidate(context.Context) error {
	return nil
}

type cacherGetErrorMock struct{}

func (c *cacherGetErrorMock) Get(context.Context, string, *Query[any]) (*Query[any], error) {
	return nil, errors.New("get-error")
}

func (c *cacherGetErrorMock) Store(context.Context, string, *Query[any]) error {
	return nil
}

func (c *cacherGetErrorMock) Invalidate(context.Context) error {
	return nil
}
