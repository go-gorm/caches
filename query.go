package caches

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Query struct {
	Dest         interface{}
	RowsAffected int64
}

func (q *Query) Marshal() ([]byte, error) {
	return json.Marshal(q)
}

func (q *Query) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, q)
}

func (q *Query) replaceOn(db *gorm.DB) {
	SetPointedValue(db.Statement.Dest, q.Dest)
	SetPointedValue(&db.Statement.RowsAffected, &q.RowsAffected)
}
