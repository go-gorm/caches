package caches

import (
	"encoding/json"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestQuery(t *testing.T) {
	t.Run("replaceOn", func(t *testing.T) {
		type User struct {
			Name string
			gorm.Model
		}
		db := &gorm.DB{
			Statement: &gorm.Statement{
				DB:   &gorm.DB{},
				Dest: &User{},
			},
		}

		expectedDestValue := &User{
			Name: "ktsivkov",
		}
		expectedAffectedRows := int64(2)

		query := Query[*User]{
			Dest:         expectedDestValue,
			RowsAffected: expectedAffectedRows,
		}
		query.replaceOn(db)

		if !reflect.DeepEqual(db.Statement.Dest, expectedDestValue) {
			t.Fatalf("replaceOn was expected to replace the destination value with the one contained inside the query.")
		}

		if !reflect.DeepEqual(db.Statement.RowsAffected, expectedAffectedRows) {
			t.Fatalf("replaceOn was expected to replace the affected rows value with the one contained inside the query.")
		}
	})

	t.Run("Marshal", func(t *testing.T) {
		type User struct {
			Name string
			gorm.Model
		}
		query := Query[User]{
			Dest: User{
				Name: "ktsivkov",
			},
			RowsAffected: 2,
		}
		res, err := query.Marshal()
		if err != nil {
			t.Fatalf("Marshal resulted to an unexpected error. %v", err)
		}

		if !json.Valid(res) {
			t.Fatalf("Marshal returned an invalid json result. %v", err)
		}
	})
	t.Run("Unmarshal", func(t *testing.T) {
		type User struct {
			Name string
			gorm.Model
		}
		marshalled := "{\"Dest\":{\"Name\":\"ktsivkov\",\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null},\"RowsAffected\":2}"
		expected := Query[User]{
			Dest: User{
				Name: "ktsivkov",
			},
			RowsAffected: 2,
		}
		var q Query[User]
		err := q.Unmarshal([]byte(marshalled))
		if err != nil {
			t.Fatalf("Unmarshal resulted to an unexpected error. %v", err)
		}

		if !reflect.DeepEqual(expected, q) {
			t.Fatalf("Unmarshal was expected to shape the query into the expected, but failed.")
		}
	})
}
