package dingo

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testScope(t *testing.T, scope Scope) {
	var requestedUnscoped int64

	test := reflect.TypeOf("string")

	unscoped := func(t reflect.Type, annotation string, optional bool) reflect.Value {
		atomic.AddInt64(&requestedUnscoped, 1)

		if optional {
			return reflect.Value{}
		}
		return reflect.New(t).Elem()
	}

	runs := 1000

	wg := new(sync.WaitGroup)
	wg.Add(runs)
	for i := 0; i < runs; i++ {
		go func() {
			t1 := scope.ResolveType(test, "", unscoped)
			t2 := scope.ResolveType(test, "", unscoped)
			assert.Equal(t, t1, t2)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(1), requestedUnscoped)
}

func TestSingleton_ResolveType(t *testing.T) {
	// reset instance
	Singleton = new(SingletonScope)

	testScope(t, Singleton)
}

func TestChildSingleton_ResolveType(t *testing.T) {
	// reset instance
	ChildSingleton = new(ChildSingletonScope)

	testScope(t, ChildSingleton)
}
