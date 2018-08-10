package dingo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinding_To(t *testing.T) {
	b := &Binding{typeof: reflect.TypeOf(new(string)).Elem()}

	b.To(new(string))
	assert.Equal(t, b.to, b.typeof)

	assert.Panics(t, func() {
		b.To(new(int))
	})
}

func TestBinding_ToInstance(t *testing.T) {
	b := &Binding{typeof: reflect.TypeOf(new(string)).Elem()}

	b.ToInstance("test")
	assert.Equal(t, b.instance.itype, b.typeof)

	assert.Panics(t, func() {
		b.ToInstance(123)
	})
}

func TestBinding_ToProvider(t *testing.T) {
	b := &Binding{typeof: reflect.TypeOf(new(string)).Elem()}

	b.ToProvider(func() string { return "test" })
	assert.Equal(t, b.provider.fnctype, b.typeof)

	assert.Panics(t, func() {
		b.ToProvider(b.ToProvider(func() int { return 123 }))
	})
}

func TestBinding_equal(t *testing.T) {
	b := &Binding{typeof: reflect.TypeOf(new(string)).Elem()}
	b2 := &Binding{typeof: reflect.TypeOf(new(string)).Elem()}
	b3 := &Binding{typeof: reflect.TypeOf(new(int)).Elem()}

	assert.True(t, b.equal(b2))
	assert.False(t, b.equal(b3))
}
