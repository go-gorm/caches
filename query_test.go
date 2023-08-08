package caches

import (
	"encoding/json"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestQuery(t *testing.T) {
	t.Run("replaceOn", func(t *testing.T) {
		db := &gorm.DB{
			Statement: &gorm.Statement{
				DB: &gorm.DB{},
				Dest: &struct {
					Name string
					gorm.Model
				}{},
			},
		}

		expectedDestValue := &struct {
			Name string
			gorm.Model
		}{
			Name: "ktsivkov",
		}
		expectedAffectedRows := int64(2)

		query := Query{
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
		query := Query{
			Dest: &struct {
				Name string
				gorm.Model
			}{
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
		marshalled := "{\"Dest\":{\"Name\":\"ktsivkov\",\"ID\":0,\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null},\"RowsAffected\":2}"
		expected := Query{
			Dest: &struct {
				Name string
				gorm.Model
			}{
				Name: "ktsivkov",
			},
			RowsAffected: 2,
		}
		query := Query{
			Dest: &struct {
				Name string
				gorm.Model
			}{},
			RowsAffected: 0,
		}
		err := query.Unmarshal([]byte(marshalled))
		if err != nil {
			t.Fatalf("Unmarshal resulted to an unexpected error. %v", err)
		}

		if !reflect.DeepEqual(expected, query) {
			t.Fatalf("Unmarshal was expected to shape the query into the expected, but failed.")
		}
	})
}
