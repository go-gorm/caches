package caches

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Query[T any] struct {
	Dest         T
	RowsAffected int64
}

func (q *Query[T]) Marshal() ([]byte, error) {
	return json.Marshal(q)
}

func (q *Query[T]) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, q)
}

func (q *Query[T]) replaceOn(db *gorm.DB) {
	SetPointedValue(db.Statement.Dest, q.Dest)
	SetPointedValue(&db.Statement.RowsAffected, &q.RowsAffected)
}
