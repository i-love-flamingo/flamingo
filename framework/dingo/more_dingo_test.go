package dingo

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type(
	mapBindInterface interface {}

	mapBindInterfaceProvider func() map[string]mapBindInterface

	mapBindTest1 struct {
		Mbp mapBindInterfaceProvider `inject:""`
	}

	mapBindTest2 struct {
		Mb mapBindInterface `inject:"map:testkey"`
	}
)

func TestMapBinding(t *testing.T) {
	injector := NewInjector()

	injector.BindMap("testkey", (*mapBindInterface)(nil)).ToInstance("testkey instance")
	injector.BindMap("testkey2", (*mapBindInterface)(nil)).ToInstance("testkey2 instance")
	injector.BindMap("testkey3", (*mapBindInterface)(nil)).ToInstance("testkey3 instance")

	test1 := injector.GetInstance(&mapBindTest1{}).(*mapBindTest1)
	test1map := test1.Mbp()

	assert.Len(t, test1map, 3)
	assert.Equal(t, "testkey instance", test1map["testkey"])
	assert.Equal(t, "testkey2 instance", test1map["testkey2"])
	assert.Equal(t, "testkey3 instance", test1map["testkey3"])

	test2 := injector.GetInstance(&mapBindTest2{}).(*mapBindTest2)
	assert.Equal(t, test2.Mb, "testkey instance")
}
