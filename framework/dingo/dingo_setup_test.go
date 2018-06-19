package dingo

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

type (
	setupT1 struct {
		member1 string
		member2 string
		member3 string
		Member4 string `inject:"annotation4"`
	}
)

func (s *setupT1) Inject(member1 string, annotated *struct {
	Member2 string `inject:"annotation2"`
	Member3 string `inject:"annotation3"`
}) {
	s.member1 = member1
	s.member2 = annotated.Member2
	s.member3 = annotated.Member3
}

func Test_Dingo_Setup(t *testing.T) {
	injector := NewInjector()
	injector.Bind((*string)(nil)).ToInstance("Member 1")
	injector.Bind((*string)(nil)).AnnotatedWith("annotation2").ToInstance("Member 2")
	injector.Bind((*string)(nil)).AnnotatedWith("annotation3").ToInstance("Member 3")
	injector.Bind((*string)(nil)).AnnotatedWith("annotation4").ToInstance("Member 4")

	test := injector.GetInstance((*setupT1)(nil)).(*setupT1)

	assert.Equal(t, test.member1, "Member 1")
	assert.Equal(t, test.member2, "Member 2")
	assert.Equal(t, test.member3, "Member 3")
	assert.Equal(t, test.Member4, "Member 4")
}
