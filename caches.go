package caches

import (
	"sync"

	"gorm.io/gorm"
)

type Caches struct {
	Conf *Config

	queue   *sync.Map
	queryCb func(*gorm.DB)
}

type Config struct {
	Easer  bool
	Cacher Cacher
}

func (c *Caches) Name() string {
	return "gorm:caches"
}

func (c *Caches) Initialize(db *gorm.DB) error {
	if c.Conf == nil {
		c.Conf = &Config{
			Easer:  false,
			Cacher: nil,
		}
	}

	if c.Conf.Easer {
		c.queue = &sync.Map{}
	}

	c.queryCb = db.Callback().Query().Get("gorm:query")

	err := db.Callback().Query().Replace("gorm:query", c.Query)
	if err != nil {
		return err
	}

	return nil
}

func (c *Caches) Query(db *gorm.DB) {
	if c.Conf.Easer == false && c.Conf.Cacher == nil {
		c.queryCb(db)
		return
	}

	identifier := buildIdentifier(db)

	if c.checkCache(db, identifier) {
		return
	}

	c.ease(db, identifier)
	if db.Error != nil {
		return
	}

	c.storeInCache(db, identifier)
	if db.Error != nil {
		return
	}
}

func (c *Caches) ease(db *gorm.DB, identifier string) {
	if c.Conf.Easer == false {
		c.queryCb(db)
		return
	}

	res := ease(&queryTask{
		id:      identifier,
		db:      db,
		queryCb: c.queryCb,
	}, c.queue).(*queryTask)

	if db.Error != nil {
		return
	}

	if res.db.Statement.Dest == db.Statement.Dest {
		return
	}

	SetPointedValue(db.Statement.Dest, res.db.Statement.Dest)
	//TODO: when dealing with timestamps the reflection fails, investigate further and find a durable solution for it
	//if err := deepCopy(res.db.Statement.Dest, db.Statement.Dest); err != nil {
	//	_ = db.AddError(err)
	//	return
	//}
}

func (c *Caches) checkCache(db *gorm.DB, identifier string) bool {
	if c.Conf.Cacher != nil {
		if res := c.Conf.Cacher.Get(identifier); res != nil {
			SetPointedValue(db.Statement.Dest, res)
			//TODO: when dealing with timestamps the reflection fails, investigate further and find a durable solution for it
			//if err := deepCopy(res, db.Statement.Dest); err != nil {
			//	_ = db.AddError(err)
			//}
			return true
		}
	}
	return false
}

func (c *Caches) storeInCache(db *gorm.DB, identifier string) {
	if c.Conf.Cacher != nil {
		err := c.Conf.Cacher.Store(identifier, db.Statement.Dest)
		if err != nil {
			_ = db.AddError(err)
		}
	}
}
