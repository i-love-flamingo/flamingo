package node

import (
	"fmt"
	"strings"

	"github.com/robertkrimen/otto/ast"
	ottoparser "github.com/robertkrimen/otto/parser"
	"github.com/robertkrimen/otto/token"
)

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
		if expr, ok := stmt.(*ast.ExpressionStatement); ok {
			finalexpr += renderExpression(expr.Expression, wrap, true)
		} else if expr, ok := stmt.(*ast.VariableStatement); ok {
			for _, v := range expr.List {
				finalexpr += renderExpression(v, wrap, true)
			}
		} else if expr, ok := stmt.(*ast.ReturnStatement); ok {
			finalexpr += renderExpression(expr.Argument, wrap, true)
		} else {
			fmt.Printf("%#v\n", stmt)
			panic("unknown expression")
		}
	}

	return finalexpr
}

var ops = map[token.Token]string{
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

var known map[string]bool

func init() {
	known = make(map[string]bool)
	known["attributes"] = true
}

// a#{b}c
// start 3
// index 4

func interpolate(s string) string {
	index := 1
	start := 0

	for index < len(s) {
		switch {
		case s[index] == '\\':

		case s[index] == '{' && s[index-1] == '$':
			start = index + 1

		case s[index] == '}' && start != 0:
			ss := JsExpr(s[start:index], false, false)
			s = s[:start-2] + `" ` + ss + ` "` + s[index+1:]
			index = start + len(ss)
			start = 0
		}
		index++
	}
	return s
}

func renderExpression(expr ast.Expression, wrap bool, dot bool) string {
	if expr == nil {
		return ""
	}

	var finalexpr string

	if str, ok := expr.(*ast.StringLiteral); ok {
		if strings.Index(str.Value, "${") >= 0 {
			finalexpr = `(s "` + interpolate(str.Value) + `")`
		} else {
			finalexpr = `"` + str.Value + `"`
		}

		if wrap {
			lf := len(finalexpr)
			if finalexpr[0] == '"' && finalexpr[lf-1] == '"' {
				return finalexpr[1 : lf-1]
			}
			finalexpr = `{{` + finalexpr + `}}`
		}

	} else if identifier, ok := expr.(*ast.Identifier); ok {
		if known[identifier.Name] {
			finalexpr += `$`
		} else if dot {
			finalexpr += `.`
		}
		finalexpr += identifier.Name
		if wrap {
			finalexpr = `{{` + finalexpr + ` | raw}}`
		}
	} else if de, ok := expr.(*ast.DotExpression); ok {
		finalexpr += renderExpression(de.Left, false, true) + "." + renderExpression(de.Identifier, false, true)[1:]
		if wrap {
			finalexpr = `{{` + finalexpr + ` | raw}}`
		}
	} else if conditional, ok := expr.(*ast.ConditionalExpression); ok {
		finalexpr = `{{if ` + renderExpression(conditional.Test, false, true) + ` }}`
		finalexpr += renderExpression(conditional.Consequent, true, true)
		if renderExpression(conditional.Alternate, true, true) != "" {
			finalexpr += `{{else}}`
			finalexpr += renderExpression(conditional.Alternate, true, true)
		}
		finalexpr += `{{end}}`
	} else if be, ok := expr.(*ast.BinaryExpression); ok {
		finalexpr = fmt.Sprintf(
			`(%s %s %s)`,
			ops[be.Operator],
			renderExpression(be.Left, false, true),
			renderExpression(be.Right, false, true))
		if wrap {
			finalexpr = `{{` + finalexpr + `}}`
		}
	} else if ce, ok := expr.(*ast.CallExpression); ok {
		fn := renderExpression(ce.Callee, false, false)
		if fn == "range" {
			fn = "__op__array"
		}
		finalexpr = `(` + fn
		for _, c := range ce.ArgumentList {
			finalexpr += ` ` + renderExpression(c, false, true)
		}
		finalexpr += `)`
		if wrap {
			finalexpr = `{{` + finalexpr + `}}`
		}
	} else if ae, ok := expr.(*ast.AssignExpression); ok {
		n := renderExpression(ae.Left, false, false)
		n = strings.TrimLeft(n, "$")
		finalexpr = fmt.Sprintf(`$%s :%s %s`,
			n,
			ops[ae.Operator],
			renderExpression(ae.Right, false, true))
		known[n] = true
		if wrap {
			finalexpr = `{{` + finalexpr + `}}`
		}
	} else if nl, ok := expr.(*ast.NumberLiteral); ok {
		finalexpr = fmt.Sprintf("%v", nl.Value)
	} else if _, ok := expr.(*ast.FunctionLiteral); ok {
		finalexpr = `FUNC//CLOSURE//`
	} else if al, ok := expr.(*ast.ArrayLiteral); ok {
		finalexpr += `(__op__array`
		for _, e := range al.Value {
			finalexpr += ` ` + renderExpression(e, false, true)
		}
		finalexpr += `)`
		if wrap {
			finalexpr = `{{` + finalexpr + `}}`
		}
	} else if bl, ok := expr.(*ast.BooleanLiteral); ok {
		finalexpr = bl.Literal
	} else if ve, ok := expr.(*ast.VariableExpression); ok {
		n := ve.Name
		n = strings.TrimLeft(n, "$")
		finalexpr = `$` + n + ` := ` + renderExpression(ve.Initializer, false, true)
		known[n] = true
		if wrap {
			finalexpr = `{{` + finalexpr + `}}`
		}
	} else if ol, ok := expr.(*ast.ObjectLiteral); ok {
		finalexpr = `(__op__map `
		for _, o := range ol.Value {
			finalexpr += ` "` + o.Key + `" ` + renderExpression(o.Value, false, true)
		}
		finalexpr += `)`
	} else if _, ok := expr.(*ast.NullLiteral); ok {
		finalexpr = ``
		if wrap {
			return `{{null}}`
		}
	} else if se, ok := expr.(*ast.SequenceExpression); ok {
		finalexpr = `(__op__array `
		for _, s := range se.Sequence {
			finalexpr += ` ` + renderExpression(s, false, true)
		}
		finalexpr += `)`
	} else if be, ok := expr.(*ast.BracketExpression); ok {
		finalexpr += `(index ` + renderExpression(be.Left, false, true) + ` ` + renderExpression(be.Member, false, true) + `)`
		if wrap {
			finalexpr = `{{` + finalexpr + `}}`
		}
	} else if ue, ok := expr.(*ast.UnaryExpression); ok {
		finalexpr += `(` + ops[ue.Operator] + ` ` + renderExpression(ue.Operand, false, true) + `)`
		if wrap {
			finalexpr = `{{` + finalexpr + `}}`
		}
	} else {
		fmt.Printf("%#v\n", expr)
		panic("unknown expression")
	}

	return finalexpr
}
