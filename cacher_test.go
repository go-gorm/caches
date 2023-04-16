package caches

import (
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

func (c *cacherMock) Get(key string) interface{} {
	c.init()
	val, _ := c.store.Load(key)
	return val
}

func (c *cacherMock) Store(key string, val interface{}) error {
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

func (c *cacherStoreErrorMock) Get(key string) interface{} {
	c.init()
	val, _ := c.store.Load(key)
	return val
}

func (c *cacherStoreErrorMock) Store(key string, val interface{}) error {
	return errors.New("store-error")
}
