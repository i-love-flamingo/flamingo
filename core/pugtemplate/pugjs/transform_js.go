package pugjs

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/pkg/errors"
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

		token.LOGICAL_AND: "__op__and", // &&
		token.LOGICAL_OR:  "__op__or",  // ||
		token.INCREMENT:   "__op__inc", // ++
		token.DECREMENT:   "__op__dec", // --

		token.EQUAL:        "__op__eql", // ==
		token.STRICT_EQUAL: "__op__eql", // ===
		token.LESS:         "__op__lt",  // <
		token.GREATER:      "__op__gt",  // >
		token.ASSIGN:       "=",         // =
		token.NOT:          "__op__not", // !

		token.BITWISE_NOT: "__op__bitnot", // ~

		token.NOT_EQUAL:        "__op__neq", // !=
		token.STRICT_NOT_EQUAL: "__op__neq", // !==
		token.LESS_OR_EQUAL:    "__op__lte", // <=
		token.GREATER_OR_EQUAL: "__op__gte", // >=

		token.DELETE: "__op__delete",
	}
)

// StrToStatements reads Javascript Statements and returns an AST representation
func StrToStatements(expr string) []ast.Statement {
	p, err := ottoparser.ParseFile(nil, "", expr, 0)
	if err != nil {
		panic(errors.Wrap(err, expr))
	}
	return p.Body
}

// FuncToStatements reads Javascript Statements and evaluates them as the return of a function
func FuncToStatements(expr string) []ast.Statement {
	p, err := ottoparser.ParseFunction("", "return "+expr)
	if err != nil {
		panic(errors.Wrap(err, expr))
	}
	return p.Body.(*ast.BlockStatement).List
}

// JsExpr transforms a javascript expression to go code
func (p *renderState) JsExpr(expr string, wrap, rawcode bool) string {
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
		finalexpr += p.renderStatement(stmt, wrap, true)
		if p.debug && wrap && len(stmtlist) > 1 {
			finalexpr += "     {{- \"\" -}}\n"
		}
	}

	return finalexpr
}

// interpolate a string, in the format of `something something ${arbitrary js code resuting in a string} blah`
// we use a helper function called `s` to merge them later
func (p *renderState) interpolate(input string) string {
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

func (p *renderState) renderStatement(stmt ast.Statement, wrap bool, dot bool) string {
	var finalexpr string

	if stmt == nil {
		return ""
	}

	switch expr := stmt.(type) {
	// an expression is just any javascript expression
	case *ast.ExpressionStatement:
		finalexpr += p.renderExpression(expr.Expression, wrap, dot)

		// a variable statement is a list of expressions, usually variable assignments (var foo = 1, bar = 2)
	case *ast.VariableStatement:
		for _, v := range expr.List {
			finalexpr += p.renderExpression(v, wrap, dot)
		}

		// the return statement is created by ParseFunction
	case *ast.ReturnStatement:
		finalexpr += p.renderExpression(expr.Argument, wrap, dot)

	case *ast.IfStatement:
		finalexpr = `{{if ` + p.renderExpression(expr.Test, false, true) + `}}`
		finalexpr += p.renderStatement(expr.Consequent, true, true)
		elsebranch := p.renderStatement(expr.Alternate, true, true)
		if elsebranch != "" && elsebranch != "{{null}}" {
			finalexpr += `{{else}}`
			finalexpr += elsebranch
		}
		finalexpr += `{{end}}`

	case *ast.ThrowStatement:
		finalexpr += p.renderExpression(expr.Argument, wrap, true)

	case *ast.BlockStatement:
		for _, s := range expr.List {
			finalexpr += p.renderStatement(s, wrap, true)
		}

	case *ast.ForInStatement:
		finalexpr = `{{ range ` + p.renderExpression(expr.Into, false, true) + ` := (__range_helper_keys__ ` + p.renderExpression(expr.Source, false, true) + `) }}`
		finalexpr += p.renderStatement(expr.Body, true, true)
		finalexpr += `{{ end }}`

	//case *ast.BranchStatement:
	//	finalexpr = `{{ if `

	// we cannot deal with other expressions at the moment, and we don'e expect them ayway
	default:
		fmt.Printf("%#v\n", stmt)
		panic("unknown expression")
	}

	return finalexpr
}

// renderExpression renders the javascript expression into go template
func (p *renderState) renderExpression(expr ast.Expression, wrap bool, dot bool) string {
	if expr == nil {
		return ""
	}

	var result string

	switch expr := expr.(type) {
	// Identifier: usually a variable name
	case *ast.Identifier:
		if _, known := p.funcs[expr.Name]; dot && !known {
			result += `$`
		} else if dot && !known {
			panic("ain'e no dot allowed")
		}
		if expr.Name == "range" {
			expr.Name = "__Range"
		}
		result += expr.Name
		if wrap {
			if !p.rawmode {
				result += ` | __pug__html`
			}
			result = `{{` + result + `}}`
		}

	// StringLiteral: "test" or 'test' or `test`
	case *ast.StringLiteral:
		if strings.Index(expr.Value, "${") >= 0 {
			result = `(__str "` + p.interpolate(expr.Value) + `")`
			result = strings.Replace(result, `""`, ``, -1)
			if wrap {
				result = `{{` + result + `}}`
			}
		} else {
			if wrap {
				result = template.HTMLEscapeString(expr.Value)
			} else {
				result = fmt.Sprintf(`%q`, expr.Value)
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
		result += p.renderExpression(expr.Left, false, true) + "."
		identifier := p.renderExpression(expr.Identifier, false, true)
		if identifier[0] == '.' || identifier[0] == '$' {
			identifier = identifier[1:]
		}
		result += identifier
		if wrap {
			if !p.rawmode {
				result += ` | __pug__html`
			}
			result = `{{` + result + `}}`
		}

	// ConditionalExpression: if (something) { ... } or foo ? a : b
	case *ast.ConditionalExpression:
		cons := p.renderExpression(expr.Consequent, false, true)
		if cons == "" {
			cons = "null"
		}
		alternate := p.renderExpression(expr.Alternate, false, true)
		if alternate == "" {
			alternate = "null"
		}
		result = `(__if (` + p.renderExpression(expr.Test, false, true) + `) (` + cons + `) (` + alternate + `) )`
		if wrap {
			result = `{{` + result + `}}`
		}

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
			if !p.rawmode {
				result += ` | __pug__html`
			}
			result = `{{` + result + `}}`
		}

	// AssignExpression: assigns something to a variable: foo = ...
	case *ast.AssignExpression:
		n := p.renderExpression(expr.Left, false, false)
		n = strings.TrimLeft(n, "$")
		right := p.renderExpression(expr.Right, false, true)
		if len(right) == 0 {
			right = "null"
		}

		// special case: assign into object
		if strings.Index(n, ".") > 0 {
			ns := strings.Split(n, ".")
			n = strings.Join(ns[:len(ns)-1], ".")
			r := ns[len(ns)-1]
			result = fmt.Sprintf(`($%s.__assign "%s" %s)`,
				n,
				r,
				right)
		} else {
			result = fmt.Sprintf(`$%s :%s %s`,
				n,
				ops[expr.Operator],
				right)
		}
		if wrap {
			result = `{{ ` + result + ` -}}`
		}

	// VariableExpression: creates a new variable, var foo = 1
	case *ast.VariableExpression:
		n := expr.Name
		n = strings.TrimLeft(n, "$")
		init := p.renderExpression(expr.Initializer, false, true)
		if len(init) == 0 {
			init = "null"
		}
		result = `$` + n + ` := ` + init
		if wrap {
			result = `{{ ` + result + ` -}}`
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
		result += `(__pug__index ` + p.renderExpression(expr.Left, false, true) + ` ` + p.renderExpression(expr.Member, false, true) + `)`
		if wrap {
			if !p.rawmode {
				result += ` | __pug__html`
			}
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
			result = `{{ ` + result + ` -}}`
		} else {
			result = `(` + result + `)`
		}

	case *ast.NewExpression:
		result = `(__op__array`
		for _, o := range expr.ArgumentList {
			result += ` ` + p.renderExpression(o, false, true)
		}
		result += `)`
		if wrap {
			result = `{{` + result + `}}`
		}

	default:
		fmt.Printf("%#v\n", expr)
		panic("unknown expression")
	}

	return result
}
