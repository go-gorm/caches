package caches

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

type mockDest struct {
	Result string
}

func TestCaches_Name(t *testing.T) {
	caches := &Caches{
		Conf: &Config{
			Easer:  true,
			Cacher: nil,
		},
	}
	expectedName := "gorm:caches"
	if act := caches.Name(); act != expectedName {
		t.Errorf("Name on caches did not return the expected value, expected: %s, actual: %s",
			expectedName, act)
	}
}

func TestCaches_Initialize(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		caches := &Caches{}
		db, err := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		if err != nil {
			t.Fatalf("gorm initialization resulted into an unexpected error, %s", err.Error())
		}

		originalQueryCb := db.Callback().Query().Get("gorm:query")

		err = db.Use(caches)
		if err != nil {
			t.Fatalf("gorm:caches loading resulted into an unexpected error, %s", err.Error())
		}

		newQueryCallback := db.Callback().Query().Get("gorm:query")

		if db.Callback().Create().Get("gorm:query") == nil {
			t.Errorf("loading of gorm:caches, expected to replace the `gorm:query` callback for Create")
		}
		if db.Callback().Update().Get("gorm:query") == nil {
			t.Errorf("loading of gorm:caches, expected to replace the `gorm:query` callback for Update")
		}
		if db.Callback().Delete().Get("gorm:query") == nil {
			t.Errorf("loading of gorm:caches, expected to replace the `gorm:query` callback for Delete")
		}
		if _, found := caches.callbacks[uponQuery]; !found {
			t.Errorf("loading of gorm:caches, expected to store the default Query `gorm:query` callback in the callbacks map")
		}
		if _, found := caches.callbacks[uponCreate]; !found {
			t.Errorf("loading of gorm:caches, expected to store the default Create `gorm:query` callback in the callbacks map")
		}
		if _, found := caches.callbacks[uponUpdate]; !found {
			t.Errorf("loading of gorm:caches, expected to store the default Update `gorm:query` callback in the callbacks map")
		}
		if _, found := caches.callbacks[uponDelete]; !found {
			t.Errorf("loading of gorm:caches, expected to store the default Delete `gorm:query` callback in the callbacks map")
		}
		if reflect.ValueOf(originalQueryCb).Pointer() == reflect.ValueOf(newQueryCallback).Pointer() {
			t.Errorf("loading of gorm:caches, expected to replace the `gorm:query` callback for Query")
		}
		if reflect.ValueOf(newQueryCallback).Pointer() != reflect.ValueOf(caches.query).Pointer() {
			t.Errorf("loading of gorm:caches, expected to replace the `gorm:query` callback, with caches.query")
		}
	})
	t.Run("config - easer", func(t *testing.T) {
		caches := &Caches{
			Conf: &Config{
				Easer:  true,
				Cacher: nil,
			},
		}
		db, err := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		if err != nil {
			t.Fatalf("gorm initialization resulted into an unexpected error, %s", err.Error())
		}

		originalQueryCb := db.Callback().Query().Get("gorm:query")

		err = db.Use(caches)
		if err != nil {
			t.Fatalf("gorm:caches loading resulted into an unexpected error, %s", err.Error())
		}

		newQueryCallback := db.Callback().Query().Get("gorm:query")

		if reflect.ValueOf(originalQueryCb).Pointer() == reflect.ValueOf(newQueryCallback).Pointer() {
			t.Errorf("loading of gorm:caches, expected to replace the `gorm:query` callback")
		}

		if reflect.ValueOf(newQueryCallback).Pointer() != reflect.ValueOf(caches.query).Pointer() {
			t.Errorf("loading of gorm:caches, expected to replace the `gorm:query` callback, with caches.query")
		}

		if reflect.ValueOf(originalQueryCb).Pointer() != reflect.ValueOf(caches.callbacks[uponQuery]).Pointer() {
			t.Errorf("loading of gorm:caches, expected to load original `gorm:query` callback, to caches.queryCb")
		}
	})
}

func TestCaches_query(t *testing.T) {
	t.Run("nothing enabled", func(t *testing.T) {
		conf := &Config{
			Easer:  false,
			Cacher: nil,
		}
		db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		db.Statement.Dest = &mockDest{}
		caches := &Caches{
			Conf: conf,
			callbacks: map[queryType]func(db *gorm.DB){
				uponQuery: func(db *gorm.DB) {
					db.Statement.Dest.(*mockDest).Result = db.Statement.SQL.String()
				},
			},
		}

		// Set the query SQL into something specific
		exampleQuery := "demo-query"
		db.Statement.SQL.WriteString(exampleQuery)

		caches.query(db) // Execute the query

		if db.Error != nil {
			t.Fatalf("an unexpected error has occurred, %v", db.Error)
		}

		if db.Statement.Dest == nil {
			t.Fatal("no query result was set after caches Query was executed")
		}

		if res := db.Statement.Dest.(*mockDest); res.Result != exampleQuery {
			t.Errorf("the execution of the Query expected a result of `%s`, got `%s`", exampleQuery, res)
		}
	})

	t.Run("easer only", func(t *testing.T) {
		conf := &Config{
			Easer:  true,
			Cacher: nil,
		}

		t.Run("one query", func(t *testing.T) {
			db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
			db.Statement.Dest = &mockDest{}
			caches := &Caches{
				Conf: conf,

				queue: &sync.Map{},
				callbacks: map[queryType]func(db *gorm.DB){
					uponQuery: func(db *gorm.DB) {
						db.Statement.Dest.(*mockDest).Result = db.Statement.SQL.String()
					},
				},
			}

			// Set the query SQL into something specific
			exampleQuery := "demo-query"
			db.Statement.SQL.WriteString(exampleQuery)

			caches.query(db) // Execute the query

			if db.Error != nil {
				t.Fatalf("an unexpected error has occurred, %v", db.Error)
			}

			if db.Statement.Dest == nil {
				t.Fatal("no query result was set after caches Query was executed")
			}

			if res := db.Statement.Dest.(*mockDest); res.Result != exampleQuery {
				t.Errorf("the execution of the Query expected a result of `%s`, got `%s`", exampleQuery, res)
			}
		})

		t.Run("two identical queries", func(t *testing.T) {
			t.Run("without error", func(t *testing.T) {
				var incr int32
				db1, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
				db1.Statement.Dest = &mockDest{}
				db2, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
				db2.Statement.Dest = &mockDest{}

				caches := &Caches{
					Conf: conf,

					queue: &sync.Map{},
					callbacks: map[queryType]func(db *gorm.DB){
						uponQuery: func(db *gorm.DB) {
							time.Sleep(1 * time.Second)
							atomic.AddInt32(&incr, 1)

							db.Statement.Dest.(*mockDest).Result = fmt.Sprintf("%d", atomic.LoadInt32(&incr))
						},
					},
				}

				// Set the queries' SQL into something specific
				exampleQuery := "demo-query"
				db1.Statement.SQL.WriteString(exampleQuery)
				db2.Statement.SQL.WriteString(exampleQuery)

				wg := &sync.WaitGroup{}
				wg.Add(2)
				go func() {
					caches.query(db1) // Execute the query
					wg.Done()
				}()
				go func() {
					time.Sleep(500 * time.Millisecond) // Execute the second query half a second later
					caches.query(db2)                  // Execute the query
					wg.Done()
				}()
				wg.Wait()

				if db1.Error != nil {
					t.Fatalf("an unexpected error has occurred, %v", db1.Error)
				}

				if db2.Error != nil {
					t.Fatalf("an unexpected error has occurred, %v", db2.Error)
				}

				if act := atomic.LoadInt32(&incr); act != 1 {
					t.Errorf("when executing two identical queries, expected to run %d time, but %d", 1, act)
				}
			})
		})

		t.Run("two different queries", func(t *testing.T) {
			var incr int32
			db1, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
			db1.Statement.Dest = &mockDest{}
			db2, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
			db2.Statement.Dest = &mockDest{}

			caches := &Caches{
				Conf: conf,

				queue: &sync.Map{},
				callbacks: map[queryType]func(db *gorm.DB){
					uponQuery: func(db *gorm.DB) {
						time.Sleep(1 * time.Second)
						atomic.AddInt32(&incr, 1)

						db.Statement.Dest.(*mockDest).Result = fmt.Sprintf("%d", atomic.LoadInt32(&incr))
					},
				},
			}

			// Set the queries' SQL into something specific
			exampleQuery1 := "demo-query-1"
			db1.Statement.SQL.WriteString(exampleQuery1)
			exampleQuery2 := "demo-query-2"
			db2.Statement.SQL.WriteString(exampleQuery2)

			wg := &sync.WaitGroup{}
			wg.Add(2)
			go func() {
				caches.query(db1) // Execute the query
				wg.Done()
			}()
			go func() {
				time.Sleep(500 * time.Millisecond) // Execute the second query half a second later
				caches.query(db2)                  // Execute the query
				wg.Done()
			}()
			wg.Wait()

			if db1.Error != nil {
				t.Fatalf("an unexpected error has occurred, %v", db1.Error)
			}

			if db2.Error != nil {
				t.Fatalf("an unexpected error has occurred, %v", db2.Error)
			}

			if act := atomic.LoadInt32(&incr); act != 2 {
				t.Errorf("when executing two identical queries, expected to run %d times, but %d", 2, act)
			}
		})
	})

	t.Run("cacher only", func(t *testing.T) {
		t.Run("one query", func(t *testing.T) {
			t.Run("with error", func(t *testing.T) {
				t.Run("store", func(t *testing.T) {
					db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
					db.Statement.Dest = &mockDest{}

					caches := &Caches{
						Conf: &Config{
							Easer:  false,
							Cacher: &cacherStoreErrorMock{},
						},

						queue: &sync.Map{},
						callbacks: map[queryType]func(db *gorm.DB){
							uponQuery: func(db *gorm.DB) {
								db.Statement.Dest.(*mockDest).Result = db.Statement.SQL.String()
							},
						},
					}

					// Set the query SQL into something specific
					exampleQuery := "demo-query"
					db.Statement.SQL.WriteString(exampleQuery)

					caches.query(db) // Execute the query

					if db.Error == nil {
						t.Error("an error was expected, got none")
					}
				})
				t.Run("get", func(t *testing.T) {
					db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
					db.Statement.Dest = &mockDest{}

					caches := &Caches{
						Conf: &Config{
							Easer:  false,
							Cacher: &cacherGetErrorMock{},
						},

						queue: &sync.Map{},
						callbacks: map[queryType]func(db *gorm.DB){
							uponQuery: func(db *gorm.DB) {
								db.Statement.Dest.(*mockDest).Result = db.Statement.SQL.String()
							},
						},
					}

					// Set the query SQL into something specific
					exampleQuery := "demo-query"
					db.Statement.SQL.WriteString(exampleQuery)

					caches.query(db) // Execute the query

					if db.Error == nil {
						t.Error("an error was expected, got none")
					}
				})
			})
			t.Run("without error", func(t *testing.T) {
				db, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
				db.Statement.Dest = &mockDest{}

				caches := &Caches{
					Conf: &Config{
						Easer:  false,
						Cacher: &cacherMock{},
					},

					queue: &sync.Map{},
					callbacks: map[queryType]func(db *gorm.DB){
						uponQuery: func(db *gorm.DB) {
							db.Statement.Dest.(*mockDest).Result = db.Statement.SQL.String()
						},
					},
				}

				// Set the query SQL into something specific
				exampleQuery := "demo-query"
				db.Statement.SQL.WriteString(exampleQuery)

				caches.query(db) // Execute the query

				if db.Error != nil {
					t.Fatalf("an unexpected error has occurred, %v", db.Error)
				}

				if db.Statement.Dest == nil {
					t.Fatal("no query result was set after caches Query was executed")
				}

				if res := db.Statement.Dest.(*mockDest); res.Result != exampleQuery {
					t.Errorf("the execution of the Query expected a result of `%s`, got `%s`", exampleQuery, res)
				}
			})
		})

		t.Run("two identical queries", func(t *testing.T) {
			var incr int32
			db1, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
			db1.Statement.Dest = &mockDest{}
			db2, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
			db2.Statement.Dest = &mockDest{}

			caches := &Caches{
				Conf: &Config{
					Easer:  false,
					Cacher: &cacherMock{},
				},

				queue: &sync.Map{},
				callbacks: map[queryType]func(db *gorm.DB){
					uponQuery: func(db *gorm.DB) {
						time.Sleep(1 * time.Second)
						atomic.AddInt32(&incr, 1)

						db.Statement.Dest.(*mockDest).Result = fmt.Sprintf("%d", atomic.LoadInt32(&incr))
					},
				},
			}

			// Set the queries' SQL into something specific
			exampleQuery := "demo-query"
			db1.Statement.SQL.WriteString(exampleQuery)
			db2.Statement.SQL.WriteString(exampleQuery)

			caches.query(db1)
			caches.query(db2)

			if db1.Error != nil {
				t.Fatalf("an unexpected error has occurred, %v", db1.Error)
			}

			if db2.Error != nil {
				t.Fatalf("an unexpected error has occurred, %v", db2.Error)
			}

			if act := atomic.LoadInt32(&incr); act != 1 {
				t.Errorf("when executing two identical queries, expected to run %d time, but %d", 1, act)
			}
		})

		t.Run("two different queries", func(t *testing.T) {
			var incr int32
			db1, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
			db1.Statement.Dest = &mockDest{}
			db2, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
			db2.Statement.Dest = &mockDest{}

			caches := &Caches{
				Conf: &Config{
					Easer:  false,
					Cacher: &cacherMock{},
				},

				queue: &sync.Map{},
				callbacks: map[queryType]func(db *gorm.DB){
					uponQuery: func(db *gorm.DB) {
						time.Sleep(1 * time.Second)
						atomic.AddInt32(&incr, 1)

						db.Statement.Dest.(*mockDest).Result = fmt.Sprintf("%d", atomic.LoadInt32(&incr))
					},
				},
			}

			// Set the queries' SQL into something specific
			exampleQuery1 := "demo-query-1"
			db1.Statement.SQL.WriteString(exampleQuery1)
			exampleQuery2 := "demo-query-2"
			db2.Statement.SQL.WriteString(exampleQuery2)

			caches.query(db1)
			if db1.Error != nil {
				t.Fatalf("an unexpected error has occurred, %v", db1.Error)
			}

			caches.query(db2)
			if db2.Error != nil {
				t.Fatalf("an unexpected error has occurred, %v", db2.Error)
			}

			if act := atomic.LoadInt32(&incr); act != 2 {
				t.Errorf("when executing two identical queries, expected to run %d times, but %d", 2, act)
			}
		})
	})
}

func TestCaches_getMutatorCb(t *testing.T) {
	testCases := map[string]queryType{
		"upon create": uponCreate,
		"upon update": uponUpdate,
		"upon delete": uponDelete,
	}

	for testName, qt := range testCases {
		t.Run(testName, func(t *testing.T) {
			expectedDb, err := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
			if err != nil {
				t.Fatalf("gorm initialization resulted into an unexpected error, %s", err.Error())
			}
			caches := &Caches{
				Conf: &Config{
					Cacher: &cacherMock{},
				},
				callbacks: map[queryType]func(db *gorm.DB){
					qt: func(db *gorm.DB) {
						if act, exp := reflect.ValueOf(db).Pointer(), reflect.ValueOf(expectedDb).Pointer(); exp != act {
							t.Errorf("the mutator did not get called with the same db instance as expected: expected %d, actual %d", exp, act)
						}
					},
				},
			}
			mutator := caches.getMutatorCb(qt)
			if mutator == nil {
				t.Errorf("loading of gorm:caches, expected generate mutator but it did not")
			}
			mutator(expectedDb)
		})
	}
}

func TestCaches_canCacheTables(t *testing.T) {
	t.Run("two identical queries cached", func(t *testing.T) {
		var incr int32
		db1, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		db1.Statement.Dest = &mockDest{}
		db2, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		db2.Statement.Dest = &mockDest{}

		caches := &Caches{
			Conf: &Config{
				Easer:           false,
				Cacher:          &cacherMock{},
				CanCachedTables: []any{&mockDest{}},
			},

			queue:          &sync.Map{},
			cacheDecisions: &sync.Map{},
			callbacks: map[queryType]func(db *gorm.DB){
				uponQuery: func(db *gorm.DB) {
					time.Sleep(1 * time.Second)
					atomic.AddInt32(&incr, 1)

					db.Statement.Dest.(*mockDest).Result = fmt.Sprintf("%d", atomic.LoadInt32(&incr))
				},
			},
		}

		// Set the queries' SQL into something specific
		exampleQuery := "demo-query"
		db1.Statement.SQL.WriteString(exampleQuery)
		db2.Statement.SQL.WriteString(exampleQuery)

		caches.query(db1)
		caches.query(db2)

		if db1.Error != nil {
			t.Fatalf("an unexpected error has occurred, %v", db1.Error)
		}

		if db2.Error != nil {
			t.Fatalf("an unexpected error has occurred, %v", db2.Error)
		}

		if act := atomic.LoadInt32(&incr); act != 1 {
			t.Errorf("when executing two identical queries, expected to run %d time, but %d", 1, act)
		}
	})

	t.Run("two identical queries cache failed", func(t *testing.T) {
		var incr int32
		db1, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		db1.Statement.Dest = &mockDest{}
		db2, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		db2.Statement.Dest = &mockDest{}

		caches := &Caches{
			Conf: &Config{
				Easer:           false,
				Cacher:          &cacherMock{},
				CanCachedTables: []any{"faker_table"},
			},

			queue:          &sync.Map{},
			cacheDecisions: &sync.Map{},
			callbacks: map[queryType]func(db *gorm.DB){
				uponQuery: func(db *gorm.DB) {
					time.Sleep(1 * time.Second)
					atomic.AddInt32(&incr, 1)

					db.Statement.Dest.(*mockDest).Result = fmt.Sprintf("%d", atomic.LoadInt32(&incr))
				},
			},
		}

		// Set the queries' SQL into something specific
		exampleQuery := "demo-query"
		db1.Statement.SQL.WriteString(exampleQuery)
		db2.Statement.SQL.WriteString(exampleQuery)

		caches.query(db1)
		caches.query(db2)

		if db1.Error != nil {
			t.Fatalf("an unexpected error has occurred, %v", db1.Error)
		}

		if db2.Error != nil {
			t.Fatalf("an unexpected error has occurred, %v", db2.Error)
		}

		if act := atomic.LoadInt32(&incr); act != 2 {
			t.Errorf("when executing two identical queries, expected to run %d time, but %d", 2, act)
		}
	})

	t.Run("two different queries can not cache", func(t *testing.T) {
		var incr int32
		db1, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		db1.Statement.Dest = &mockDest{}
		db2, _ := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
		db2.Statement.Dest = &mockDest{}

		caches := &Caches{
			Conf: &Config{
				Easer:           false,
				Cacher:          &cacherMock{},
				CanCachedTables: []any{&mockDest{}},
			},

			queue:          &sync.Map{},
			cacheDecisions: &sync.Map{},
			callbacks: map[queryType]func(db *gorm.DB){
				uponQuery: func(db *gorm.DB) {
					time.Sleep(1 * time.Second)
					atomic.AddInt32(&incr, 1)

					db.Statement.Dest.(*mockDest).Result = fmt.Sprintf("%d", atomic.LoadInt32(&incr))
				},
			},
		}

		// Set the queries' SQL into something specific
		exampleQuery1 := "demo-query-1"
		db1.Statement.SQL.WriteString(exampleQuery1)
		exampleQuery2 := "demo-query-2"
		db2.Statement.SQL.WriteString(exampleQuery2)

		caches.query(db1)
		if db1.Error != nil {
			t.Fatalf("an unexpected error has occurred, %v", db1.Error)
		}

		caches.query(db2)
		if db2.Error != nil {
			t.Fatalf("an unexpected error has occurred, %v", db2.Error)
		}

		if act := atomic.LoadInt32(&incr); act != 2 {
			t.Errorf("when executing two identical queries, expected to run %d times, but %d", 2, act)
		}
	})
}
