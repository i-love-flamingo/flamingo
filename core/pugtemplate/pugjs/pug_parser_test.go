package pugjs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderState_buildNode(t *testing.T) {
	cases := []struct {
		from *Token
		to   Node
	}{
		{&Token{Type: "Text", Val: "testText"}, &Text{ValueNode{Val: "testText"}}},
		{&Token{Type: "Code", Val: "testCode"}, &Code{ValueNode: ValueNode{Val: "testCode"}}},
		{&Token{Type: "Each"}, &Each{}},
		{&Token{Type: "NamedBlock"}, &Block{}},
		{&Token{Type: "Block"}, &Block{}},
		{&Token{Type: "Case"}, &Case{}},
		{&Token{Type: "When"}, &When{}},
		{&Token{Type: "Doctype", Val: "html"}, &Doctype{ValueNode{Val: "html"}}},

		// we ignore comments
		{&Token{Type: "Comment"}, nil},
		{&Token{Type: "BlockComment"}, nil},

		// tags
		{&Token{Type: "Tag", Name: "testTag"}, &Tag{Name: "testTag"}},
		{&Token{Type: "InterpolatedTag"}, &InterpolatedTag{}},

		// mixins
		{&Token{Type: "MixinBlock"}, &MixinBlock{}},
		{&Token{Type: "Mixin", Call: false}, &Mixin{Args: "[]"}},
		{&Token{Type: "Mixin", Call: true}, &Mixin{Args: "[]", Call: true}},

		// conditional
		{
			&Token{Type: "Conditional", Consequent: &Token{Type: "Text"}},
			&Conditional{Consequent: &Text{}},
		},
	}

	var renderState = new(renderState)
	for _, testCase := range cases {
		t.Run("Testing "+testCase.from.Type, func(t *testing.T) {
			assert.Equal(t, testCase.to, renderState.buildNode(testCase.from))
		})
	}
}
