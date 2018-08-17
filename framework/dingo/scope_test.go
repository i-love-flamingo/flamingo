package dingo

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testScope(t *testing.T, scope Scope) {
	var requestedUnscoped int64

	test := reflect.TypeOf("string")
	test2 := reflect.TypeOf("int")

	unscoped := func(t reflect.Type, annotation string, optional bool) reflect.Value {
		atomic.AddInt64(&requestedUnscoped, 1)

		if optional {
			return reflect.Value{}
		}
		return reflect.New(t).Elem()
	}

	runs := 100 // change to 10? 100? 1000? to trigger a bug? todo investigate

	wg := new(sync.WaitGroup)
	wg.Add(runs)
	for i := 0; i < runs; i++ {
		go func() {
			t1 := scope.ResolveType(test, "", unscoped)
			t12 := scope.ResolveType(test2, "", unscoped)
			t2 := scope.ResolveType(test, "", unscoped)
			t22 := scope.ResolveType(test2, "", unscoped)
			assert.Equal(t, t1, t2)
			assert.Equal(t, t12, t22)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(1), requestedUnscoped)

}

func TestSingleton_ResolveType(t *testing.T) {
	// reset instance
	Singleton = NewSingletonScope()

	testScope(t, Singleton)
}

func TestChildSingleton_ResolveType(t *testing.T) {
	// reset instance
	ChildSingleton = NewChildSingletonScope()

	testScope(t, ChildSingleton)
}

type (
	singletonA struct {
		B *singletonB `inject:""`
	}

	singletonB struct {
		C singletonC `inject:""`
	}

	singletonC string
)

func TestScopeWithSubDependencies(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			injector := NewInjector()

			injector.Bind(new(singletonA)).In(Singleton)
			injector.Bind(new(singletonB)).In(Singleton)
			injector.Bind(singletonC("")).In(Singleton).ToInstance(singletonC("singleton C"))

			runs := 10

			wg := new(sync.WaitGroup)
			wg.Add(runs)
			for i := 0; i < runs; i++ {
				go func() {
					a := injector.GetInstance(new(singletonA)).(*singletonA)
					assert.Equal(t, a.B.C, singletonC("singleton C"))
					wg.Done()
				}()
			}
			wg.Wait()
		})
	}
}
