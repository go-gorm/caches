package caches

import (
	"encoding/json"
	"sync"

	"gorm.io/gorm"
)

type Query[T any] struct {
	Dest         T
	RowsAffected int64
	mu           sync.RWMutex
}

func (q *Query[T]) Marshal() ([]byte, error) {
	q.mu.RUnlock()
	defer q.mu.RUnlock()
	return json.Marshal(q)
}

func (q *Query[T]) Unmarshal(bytes []byte) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	return json.Unmarshal(bytes, q)
}

func (q *Query[T]) replaceOn(db *gorm.DB) {
	q.mu.Lock()
	defer q.mu.Unlock()
	SetPointedValue(db.Statement.Dest, q.Dest)
	SetPointedValue(&db.Statement.RowsAffected, &q.RowsAffected)
}
