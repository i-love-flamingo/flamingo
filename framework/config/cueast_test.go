package config

import (
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/parser"
	"github.com/stretchr/testify/assert"
)

func Test_cueAstTree(t *testing.T) {
	base, err := parser.ParseFile("base", `
a: b: 1
a: c: 2
a: {
	b: 1
	c: 2
}
a: {
	b: int
	c: int
	c: >0
}
b: int
b: 2
// foo bar
`)
	assert.NoError(t, err)
	base.Decls = cueAstTree(base.Decls)
	assert.Len(t, base.Decls, 3)

	var testData struct {
		A struct {
			B int
			C int
		}
		B int
	}

	buildInstance := build.NewContext().NewInstance("", nil)
	assert.NoError(t, buildInstance.AddSyntax(base))
	instance, err := new(cue.Runtime).Build(buildInstance)
	assert.NoError(t, err)
	assert.NoError(t, instance.Err)
	assert.NoError(t, instance.Value().Decode(&testData))

	assert.Equal(t, 1, testData.A.B)
	assert.Equal(t, 2, testData.A.C)
	assert.Equal(t, 2, testData.B)
}

func Test_cueAstMerge(t *testing.T) {
	predefined, err := parser.ParseFile("predefined", `
predef :: {x: 1}
`)
	assert.NoError(t, err)

	base, err := parser.ParseFile("base", `
a: int
b: 2
a: 1
struct: d: e: {
	f1: 11
	f2: 12
}
list: [1,2,3]
struct: x: 1
defbase :: { c: 3 }
def :: { defbase, a: 1,	b: 2}
defined: def
predefined: predef
`)
	assert.NoError(t, err)

	in, err := parser.ParseFile("in", `
b: 22
struct: d: e: f2: 223
list: [3,2,1]
struct: y: 2
struct: z: {z1: 51, z2: 52}
def :: { b: 3 }
`)
	assert.NoError(t, err)

	res := cueAstMergeFile(base, in)

	buildInstance := build.NewContext().NewInstance("", nil)
	assert.NoError(t, buildInstance.AddSyntax(predefined))
	assert.NoError(t, buildInstance.AddSyntax(res))
	instance, err := new(cue.Runtime).Build(buildInstance)
	assert.NoError(t, err)

	var testData struct {
		A      int
		B      int
		Struct struct {
			D struct {
				E struct {
					F1 int
					F2 int
				}
			}
			X int
			Y int
			Z struct {
				Z1 int
				Z2 int
			}
		}
		List    []int
		Defined struct {
			A int
			B int
			C int
		}
		Predefined struct {
			X int
		}
	}
	assert.NoError(t, instance.Value().Decode(&testData))

	assert.Equal(t, 1, testData.A)
	assert.Equal(t, 22, testData.B)
	assert.Equal(t, 1, testData.Struct.X)
	assert.Equal(t, 2, testData.Struct.Y)
	assert.Equal(t, 11, testData.Struct.D.E.F1)
	assert.Equal(t, 223, testData.Struct.D.E.F2)
	assert.Equal(t, 51, testData.Struct.Z.Z1)
	assert.Equal(t, 52, testData.Struct.Z.Z2)
	assert.Equal(t, []int{3, 2, 1}, testData.List)
	assert.Equal(t, 1, testData.Defined.A)
	assert.Equal(t, 3, testData.Defined.B)
	assert.Equal(t, 3, testData.Defined.C)
	assert.Equal(t, 1, testData.Predefined.X)
}
