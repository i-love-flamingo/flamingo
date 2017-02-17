package node

import (
	"fmt"
	"strings"

	"github.com/robertkrimen/otto/ast"
	ottoparser "github.com/robertkrimen/otto/parser"
	"github.com/robertkrimen/otto/token"
)

var (
	known map[string]bool

	ops = map[token.Token]string{
		token.PLUS:      "__op__add",   // +
		token.MINUS:     "__op__sub",   // -
		token.MULTIPLY:  "__op__mul",   // *
		token.SLASH:     "__op__slash", // /
		token.REMAINDER: "__op__mod",   // %

		token.AND:                  "__op__b_and",     // &
		token.OR:                   "__op__b_or",      // |
		token.EXCLUSIVE_OR:         "__op__b_xor",     // ^
		token.SHIFT_LEFT:           "__op__b_sleft",   // <<
		token.SHIFT_RIGHT:          "__op__b_sright",  // >>
		token.UNSIGNED_SHIFT_RIGHT: "__op__b_usright", // >>>
		token.AND_NOT:              "__op__b_andnot",  // &^

		token.LOGICAL_AND: "and",       // &&
		token.LOGICAL_OR:  "or",        // ||
		token.INCREMENT:   "__op__inc", // ++
		token.DECREMENT:   "__op__dec", // --

		token.EQUAL:        "eq",  // ==
		token.STRICT_EQUAL: "eq",  // ===
		token.LESS:         "lt",  // <
		token.GREATER:      "gt",  // >
		token.ASSIGN:       "=",   // =
		token.NOT:          "not", // !

		token.BITWISE_NOT: "__op__bitnot", // ~

		token.NOT_EQUAL:        "neq", // !=
		token.STRICT_NOT_EQUAL: "neq", // !==
		token.LESS_OR_EQUAL:    "lte", // <=
		token.GREATER_OR_EQUAL: "gte", // >=

		token.DELETE: "delete",
	}
)

func init() {
	known = make(map[string]bool)
	known["attributes"] = true
}

// JsExpr transforms a javascript expression to go code
func JsExpr(expr string, wrap, rawcode bool) string {
	var finalexpr string
	var stmtlist []ast.Statement

	if rawcode {
		p, err := ottoparser.ParseFile(nil, "", expr, 0)
		if err != nil {
			fmt.Println(expr)
			panic(err)
		}
		stmtlist = p.Body
	} else {
		p, err := ottoparser.ParseFunction("", "return "+expr)
		if err != nil {
			fmt.Println(expr)
			panic(err)
		}
		stmtlist = p.Body.(*ast.BlockStatement).List
	}

	for _, stmt := range stmtlist {
		switch expr := stmt.(type) {
		case *ast.ExpressionStatement:
			finalexpr += renderExpression(expr.Expression, wrap, true)

		case *ast.VariableStatement:
			for _, v := range expr.List {
				finalexpr += renderExpression(v, wrap, true)
			}

		case *ast.ReturnStatement:
			finalexpr += renderExpression(expr.Argument, wrap, true)

		default:
			fmt.Printf("%#v\n", stmt)
			panic("unknown expression")
		}
	}

	return finalexpr
}

func interpolate(input string) string {
	index := 1
	start := 0

	for index < len(input) {
		switch {
		case input[index] == '\\':

		case input[index] == '{' && input[index-1] == '$':
			start = index + 1

		case input[index] == '}' && start != 0:
			substring := JsExpr(input[start:index], false, false)
			input = input[:start-2] + `" ` + substring + ` "` + input[index+1:]
			index = start + len(substring)
			start = 0
		}
		index++
	}
	return input
}

func renderExpression(expr ast.Expression, wrap bool, dot bool) string {
	if expr == nil {
		return ""
	}

	var result string

	switch expr := expr.(type) {
	// Identifier: usually a variable name
	case *ast.Identifier:
		if known[expr.Name] {
			result += `$`
		} else if dot {
			result += `.`
		}
		result += expr.Name
		if wrap {
			result = `{{` + result + ` | raw}}`
		}

	// StringLiteral: "test" or 'test' or `test`
	case *ast.StringLiteral:
		if strings.Index(expr.Value, "${") >= 0 {
			result = `(s "` + interpolate(expr.Value) + `")`
			if wrap {
				result = `{{` + result + `}}`
			}
		} else {
			if wrap {
				result = expr.Value
			} else {
				result = `"` + expr.Value + `"`
			}
		}

	// NumberLiteral: 1 or 1.5
	case *ast.NumberLiteral:
		result = fmt.Sprintf("%v", expr.Value)

	// ArrayLiteral: [1, 2, 3]
	case *ast.ArrayLiteral:
		result += `(__op__array`
		for _, e := range expr.Value {
			result += ` ` + renderExpression(e, false, true)
		}
		result += `)`
		if wrap {
			result = `{{` + result + `}}`
		}

	// BooleanLiteral: true or false
	case *ast.BooleanLiteral:
		result = expr.Literal

	// ObjectLiteral: {"key": "value", "key2": something}
	case *ast.ObjectLiteral:
		result = `(__op__map `
		for _, o := range expr.Value {
			result += ` "` + o.Key + `" ` + renderExpression(o.Value, false, true)
		}
		result += `)`

	// NullLiteral: null
	case *ast.NullLiteral:
		result = ``
		if wrap {
			return `{{null}}`
		}

	// DotExpression: left.right
	case *ast.DotExpression:
		result += renderExpression(expr.Left, false, true) + "." + renderExpression(expr.Identifier, false, true)[1:]
		if wrap {
			result = `{{` + result + ` | raw}}`
		}

	// ConditionalExpression: if (something) { ... }
	case *ast.ConditionalExpression:
		result = `{{if ` + renderExpression(expr.Test, false, true) + ` }}`
		result += renderExpression(expr.Consequent, true, true)
		if renderExpression(expr.Alternate, true, true) != "" {
			result += `{{else}}`
			result += renderExpression(expr.Alternate, true, true)
		}
		result += `{{end}}`

	// BinaryExpression:  left binary-operator right, 1 & 2, 0xff ^ 0x01
	case *ast.BinaryExpression:
		result = fmt.Sprintf(
			`(%s %s %s)`,
			ops[expr.Operator],
			renderExpression(expr.Left, false, true),
			renderExpression(expr.Right, false, true))
		if wrap {
			result = `{{` + result + `}}`
		}

	// CallExpression: calls a function (Callee) with arguments, e.g. url("target", "arg1", 1)
	case *ast.CallExpression:
		result = `(` + renderExpression(expr.Callee, false, false)
		for _, c := range expr.ArgumentList {
			result += ` ` + renderExpression(c, false, true)
		}
		result += `)`
		if wrap {
			result = `{{` + result + `}}`
		}

	// AssignExpression: assigns something to a variable: foo = ...
	case *ast.AssignExpression:
		n := renderExpression(expr.Left, false, false)
		n = strings.TrimLeft(n, "$")
		result = fmt.Sprintf(`$%s :%s %s`,
			n,
			ops[expr.Operator],
			renderExpression(expr.Right, false, true))
		known[n] = true
		if wrap {
			result = `{{` + result + `}}`
		}

	// VariableExpression: creates a new variable, var foo = 1
	case *ast.VariableExpression:
		n := expr.Name
		n = strings.TrimLeft(n, "$")
		result = `$` + n + ` := ` + renderExpression(expr.Initializer, false, true)
		known[n] = true
		if wrap {
			result = `{{` + result + `}}`
		}

	// SequenceExpression, just like ArrayLiteral
	case *ast.SequenceExpression:
		result = `(__op__array `
		for _, s := range expr.Sequence {
			result += ` ` + renderExpression(s, false, true)
		}
		result += `)`

	// BracketExpression: access of array/object members, such ass something[1] or foo[bar]
	case *ast.BracketExpression:
		result += `(index ` + renderExpression(expr.Left, false, true) + ` ` + renderExpression(expr.Member, false, true) + `)`
		if wrap {
			result = `{{` + result + `}}`
		}

	// UnaryExpression: an operation on an operand, such as delete foo[bar]
	case *ast.UnaryExpression:
		result += `(` + ops[expr.Operator] + ` ` + renderExpression(expr.Operand, false, true) + `)`
		if wrap {
			result = `{{` + result + `}}`
		}

	default:
		fmt.Printf("%#v\n", expr)
		panic("unknown expression")
	}

	return result
}
