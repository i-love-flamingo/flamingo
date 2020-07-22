package config

import (
	"cuelang.org/go/cue/ast"
)

// cueAstTree unifies structs so that
// a: b: 1
// a: c: 1
// becomes
// a: {
//   b: 1
//   c: 1
// }
// currently field comprehensions are not supported
func cueAstTree(in []ast.Decl) []ast.Decl {
	result := make([]ast.Decl, 0, len(in))
	known := make(map[string]*ast.StructLit, len(in))

	for _, d := range in {
		field, ok := d.(*ast.Field)
		if !ok {
			result = append(result, d)
			continue
		}

		ident, ok := field.Label.(*ast.Ident)
		if !ok {
			result = append(result, d)
			continue
		}

		value, ok := field.Value.(*ast.StructLit)
		if !ok {
			result = append(result, d)
			continue
		}

		if known[ident.Name] != nil {
			known[ident.Name].Elts = append(known[ident.Name].Elts, value.Elts...)
		} else {
			known[ident.Name] = value
			result = append(result, field)
			continue
		}
	}
	return result
}

// cueAstMergeFile processes *ast.File declarations
func cueAstMergeFile(base, in *ast.File) *ast.File {
	if base == nil {
		return in
	}
	if in == nil {
		return base
	}
	base.Decls = cueAstTree(base.Decls) // cut base
	in.Decls = cueAstTree(in.Decls)     // cut in

	base.Decls = cueAstMergeDecls(base.Decls, in.Decls) // merge decls
	return base
}

// cueAstMergeDecls merges two ast.Decl lists
func cueAstMergeDecls(base []ast.Decl, in []ast.Decl) []ast.Decl {
	result := make([]ast.Decl, 0, len(base))
	known := make(map[string]*ast.Field, len(in)) // mark and hold reference to known structs

	// process in
	for _, d := range in {
		// ignore non-fields
		field, ok := d.(*ast.Field)
		if !ok {
			continue
		}

		// ignore non-static identifiers
		ident, ok := field.Label.(*ast.Ident)
		if !ok {
			continue
		}

		// mark field identifier as known, and store a reference if available
		known[ident.Name] = field
		if _, ok = field.Value.(*ast.StructLit); !ok {
			result = append(result, field)
		}
	}

	added := make(map[string]struct{}, len(known))

	// walk base declarations
	for _, d := range base {
		// no further processing for non-fields
		field, ok := d.(*ast.Field)
		if !ok {
			result = append(result, d)
			continue
		}

		// no further processing for non-identifier labels
		ident, ok := field.Label.(*ast.Ident)
		if !ok {
			result = append(result, d)
			continue
		}

		// simply append if we didn't mark this field before
		if knownField, ok := known[ident.Name]; !ok {
			result = append(result, field)
			continue
		} else {
			switch value := field.Value.(type) {
			case *ast.StructLit:
				// merge incoming struct into base struct
				// we preserve the original reference
				value.Elts = cueAstMergeDecls(value.Elts, knownField.Value.(*ast.StructLit).Elts)
				result = append(result, field)
				added[ident.Name] = struct{}{}
			default:
				continue
			}
		}
	}

	for k, v := range known {
		// skip added fields
		if _, ok := added[k]; ok {
			continue
		}
		// skip non-structs
		if _, ok := v.Value.(*ast.StructLit); !ok {
			continue
		}
		// add missing structs
		result = append(result, v)
	}

	return result
}
