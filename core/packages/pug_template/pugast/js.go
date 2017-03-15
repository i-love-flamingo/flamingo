package pugast

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/robertkrimen/otto/ast"
	ottoparser "github.com/robertkrimen/otto/parser"
	"github.com/robertkrimen/otto/token"
)

var (
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

		token.EQUAL:        "__op__eql", // ==
		token.STRICT_EQUAL: "__op__eql", // ===
		token.LESS:         "lt",        // <
		token.GREATER:      "gt",        // >
		token.ASSIGN:       "=",         // =
		token.NOT:          "not",       // !

		token.BITWISE_NOT: "__op__bitnot", // ~

		token.NOT_EQUAL:        "neq", // !=
		token.STRICT_NOT_EQUAL: "neq", // !==
		token.LESS_OR_EQUAL:    "lte", // <=
		token.GREATER_OR_EQUAL: "gte", // >=

		token.DELETE: "delete",
	}
)

// StrToStatements reads Javascript Statements and returns an AST representation
func StrToStatements(expr string) []ast.Statement {
	p, err := ottoparser.ParseFile(nil, "", expr, 0)
	if err != nil {
		panic(err)
	}
	return p.Body
}

// FuncToStatements reads Javascript Statements and evaluates them as the return of a function
func FuncToStatements(expr string) []ast.Statement {
	p, err := ottoparser.ParseFunction("", "return "+expr)
	if err != nil {
		panic(err)
	}
	return p.Body.(*ast.BlockStatement).List
}

// JsExpr transforms a javascript expression to go code
func (p *PugAst) JsExpr(expr string, wrap, rawcode bool) string {
	var finalexpr string
	var stmtlist []ast.Statement

	if rawcode {
		// Expect the input to be raw js code. This makes `{ ... }` being treated as a logical block
		stmtlist = StrToStatements(expr)
	} else {
		// Expect the input to be a value, this makes `{ ... }` being treated as a map.
		// Essentially we create a function with one return-statement and inject our return value
		stmtlist = FuncToStatements(expr)
	}

	for _, stmt := range stmtlist {
		switch expr := stmt.(type) {
		// an expression is just any javascript expression
		case *ast.ExpressionStatement:
			finalexpr += p.renderExpression(expr.Expression, wrap, true)

		// a variable statement is a list of expressions, usually variable assignments (var foo = 1, bar = 2)
		case *ast.VariableStatement:
			for _, v := range expr.List {
				finalexpr += p.renderExpression(v, wrap, true)
			}

		// the return statement is created by ParseFunction
		case *ast.ReturnStatement:
			finalexpr += p.renderExpression(expr.Argument, wrap, true)

		// we cannot deal with other expressions at the moment, and we don't expect them ayway
		default:
			fmt.Printf("%#v\n", stmt)
			panic("unknown expression")
		}
	}

	return finalexpr
}

// interpolate a string, in the format of `something something ${arbitrary js code resuting in a string} blah`
// we use a helper function called `s` to merge them later
func (p *PugAst) interpolate(input string) string {
	index := 1
	start := 0

	for index < len(input) {
		switch {
		case input[index] == '\\':
			break

		case input[index] == '{' && input[index-1] == '$':
			start = index + 1

		case input[index] == '}' && start != 0:
			substring := p.JsExpr(input[start:index], false, false)
			input = input[:start-2] + `" ` + substring + ` "` + input[index+1:]
			index = start + len(substring)
			start = 0
		}
		index++
	}
	return input
}

// renderExpression renders the javascript expression into go template
func (p *PugAst) renderExpression(expr ast.Expression, wrap bool, dot bool) string {
	if expr == nil {
		return ""
	}

	var result string

	switch expr := expr.(type) {
	// Identifier: usually a variable name
	case *ast.Identifier:
		if _, known := p.FuncMap[expr.Name]; !known && p.knownVar[expr.Name] {
			result += `$`
		} else if dot && !known {
			result += `.`
		}
		result += expr.Name
		if wrap {
			if p.rawmode {
				result += ` | raw`
			}
			result = `{{` + result + `}}`
		}

	// StringLiteral: "test" or 'test' or `test`
	case *ast.StringLiteral:
		if strings.Index(expr.Value, "${") >= 0 {
			result = `(s "` + p.interpolate(expr.Value) + `")`
			result = strings.Replace(result, `""`, ``, -1)
			if wrap {
				result = `{{` + result + `}}`
			}
		} else {
			if wrap {
				result = template.HTMLEscapeString(expr.Value)
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
			result += ` ` + p.renderExpression(e, false, true)
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
		result = `(__op__map`
		for _, o := range expr.Value {
			result += ` "` + o.Key + `" ` + p.renderExpression(o.Value, false, true)
		}
		result += `)`
		if wrap {
			result = `{{` + result + `}}`
		}

	// NullLiteral: null
	case *ast.NullLiteral:
		result = ``
		if wrap {
			return `{{null}}`
		}

	// DotExpression: left.right
	case *ast.DotExpression:
		result += p.renderExpression(expr.Left, false, true) + "." + p.renderExpression(expr.Identifier, false, true)[1:]
		if wrap {
			if p.rawmode {
				result += ` | raw`
			}
			result = `{{` + result + `}}`
		}

	// ConditionalExpression: if (something) { ... } or foo ? a : b
	case *ast.ConditionalExpression:
		result = `{{if ` + p.renderExpression(expr.Test, false, true) + `}}`
		result += p.renderExpression(expr.Consequent, true, true)
		elsebranch := p.renderExpression(expr.Alternate, true, true)
		if elsebranch != "" && elsebranch != "{{null}}" {
			result += `{{else}}`
			result += elsebranch
		}
		result += `{{end}}`

	// BinaryExpression:  left binary-operator right, 1 & 2, 0xff ^ 0x01
	case *ast.BinaryExpression:
		result = fmt.Sprintf(
			`(%s %s %s)`,
			ops[expr.Operator],
			p.renderExpression(expr.Left, false, true),
			p.renderExpression(expr.Right, false, true))
		if wrap {
			result = `{{` + result + `}}`
		}

	// CallExpression: calls a function (Callee) with arguments, e.g. url("target", "arg1", 1)
	case *ast.CallExpression:
		result = `(` + p.renderExpression(expr.Callee, false, false)
		for _, c := range expr.ArgumentList {
			result += ` ` + p.renderExpression(c, false, true)
		}
		result += `)`
		if wrap {
			result = `{{` + result + `}}`
		}

	// AssignExpression: assigns something to a variable: foo = ...
	case *ast.AssignExpression:
		n := p.renderExpression(expr.Left, false, false)
		n = strings.TrimLeft(n, "$")
		result = fmt.Sprintf(`$%s :%s %s`,
			n,
			ops[expr.Operator],
			p.renderExpression(expr.Right, false, true))
		p.knownVar[n] = true
		if wrap {
			result = `{{- ` + result + ` -}}`
		}

	// VariableExpression: creates a new variable, var foo = 1
	case *ast.VariableExpression:
		n := expr.Name
		n = strings.TrimLeft(n, "$")
		result = `$` + n + ` := ` + p.renderExpression(expr.Initializer, false, true)
		p.knownVar[n] = true
		if wrap {
			result = `{{- ` + result + ` -}}`
		}

	// SequenceExpression, just like ArrayLiteral
	case *ast.SequenceExpression:
		result = `(__op__array`
		for _, s := range expr.Sequence {
			result += ` ` + p.renderExpression(s, false, true)
		}
		result += `)`

	// BracketExpression: access of array/object members, such ass something[1] or foo[bar]
	case *ast.BracketExpression:
		result += `(index ` + p.renderExpression(expr.Left, false, true) + ` ` + p.renderExpression(expr.Member, false, true) + `)`
		if wrap {
			result = `{{` + result + `}}`
		}

	// UnaryExpression: an operation on an operand, such as delete foo[bar]
	case *ast.UnaryExpression:
		if expr.Operator == token.INCREMENT {
			result += p.renderExpression(expr.Operand, false, true) + ` := ` + ops[expr.Operator] + ` ` + p.renderExpression(expr.Operand, false, true)
		} else {
			result += ops[expr.Operator] + ` ` + p.renderExpression(expr.Operand, false, true)
		}
		if wrap {
			result = `{{- ` + result + ` -}}`
		} else {
			result = `(` + result + `)`
		}

	default:
		fmt.Printf("%#v\n", expr)
		panic("unknown expression")
	}

	return result
}
