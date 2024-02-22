package caches

import (
	"testing"

	"gorm.io/gorm"
)

func Test_buildIdentifier(t *testing.T) {
	db := &gorm.DB{}
	db.Statement = &gorm.Statement{}
	db.Statement.SQL.WriteString("TEST-SQL")
	db.Statement.Vars = append(db.Statement.Vars, "test", 123, 12.3, true, false, []string{"test", "me"})

	actual := buildIdentifier(db)
	expected := "gorm-caches::TEST-SQL-[test 123 12.3 true false [test me]]"
	if actual != expected {
		t.Errorf("buildIdentifier expected to return `%s` but got `%s`", expected, actual)
	}
}

func Test_sliceToString(t *testing.T) {
	expected := "[test-val test-val 1 1 true true [test-val] [1] [true] [test-val] [1] [true] [test-val] [1] [true] [test-val] [1] [true] {test-val: test-val} {1: 1} {true: true} {test-val: test-val} {1: 1} {true: true} {test-val: test-val} {1: 1} {true: true} {test-val: test-val} {1: 1} {true: true}]"

	strVal := "test-val"
	intVal := 1
	boolVal := true
	sliceOfStr := []string{strVal}
	sliceOfInt := []int{intVal}
	sliceOfBool := []bool{boolVal}
	sliceOfPointerStr := []*string{&strVal}
	sliceOfPointerInt := []*int{&intVal}
	sliceOfPointerBool := []*bool{&boolVal}
	mapOfStr := map[string]string{strVal: strVal}
	mapOfInt := map[int]int{intVal: intVal}
	mapOfBool := map[bool]bool{boolVal: boolVal}
	mapOfPointerStr := map[*string]*string{&strVal: &strVal}
	mapOfPointerInt := map[*int]*int{&intVal: &intVal}
	mapOfPointerBool := map[*bool]*bool{&boolVal: &boolVal}
	actual := valueToString([]interface{}{
		strVal, &strVal, intVal, // Primitives
		&intVal, boolVal, &boolVal, // Pointer passed primitives
		sliceOfStr, sliceOfInt, sliceOfBool, // Slices of primitives
		&sliceOfStr, &sliceOfInt, &sliceOfBool, // Pointer passed slices of primitives
		sliceOfPointerStr, sliceOfPointerInt, sliceOfPointerBool, // Slices of pointer primitives
		&sliceOfPointerStr, &sliceOfPointerInt, &sliceOfPointerBool, // Pointer passed slices of pointer primitives
		mapOfStr, mapOfInt, mapOfBool, // Map of primitives
		mapOfPointerStr, mapOfPointerInt, mapOfPointerBool, // Map of pointer primitives
		&mapOfStr, &mapOfInt, &mapOfBool, // Map of primitives
		&mapOfPointerStr, &mapOfPointerInt, &mapOfPointerBool, // Map of pointer primitives
	})
	if expected != actual {
		t.Errorf("sliceToString expected to return `%s` but got `%s`", expected, actual)
	}
}
