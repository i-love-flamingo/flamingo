package pugjs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCase_Render(t *testing.T) {
	var buffer = new(bytes.Buffer)
	var node = new(Case)
	var renderState = new(renderState)

	assert.Error(t, node.Render(renderState, buffer))

	node.Expr = JavaScriptExpression("obj")
	node.Block.Nodes = []Node{
		&When{ExpressionNode: ExpressionNode{Expr: "default"}},
		&When{ExpressionNode: ExpressionNode{Expr: "case1"}},
		&When{ExpressionNode: ExpressionNode{Expr: "case2"}},
	}

	assert.NoError(t, node.Render(renderState, buffer))

	assert.Equal(t, "{{- if __op__eql $obj $case1 }}{{- else if __op__eql $obj $case2 }}{{- else }}{{- end }}", buffer.String())
}
