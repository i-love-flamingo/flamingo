package pugast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

type (
	// Attr is a simple key-value pair
	Attr struct {
		Name, Val  string
		MustEscape bool
	}

	// Fileref is used by include/extends
	Fileref struct {
		Type, Path string
		Line       int
	}

	// Token defines the basic token read by the tokenizer
	// Tokens form a tree, where the beginning root node starts the document
	Token struct {
		// default
		Type, Name string
		Mode, Val  string
		Line       int

		// subblock
		Block *Token
		// subblock childs
		Nodes []*Token

		// specific information
		AttributeBlocks []string
		Attrs           []*Attr
		MustEscape      bool
		File            *Fileref
		Filename        string
		SelfClosing     bool
		IsInline        *bool
		Obj             string
		Key             string

		// mixin
		Call bool   // mixin call?
		Args string // call args

		// if
		Test                  string // if
		Consequent, Alternate *Token // if result

		// Interpolated
		Expr string
	}
)

// Parse parses a filename into a Token-tree
func (p *PugAst) Parse(file string) *Token {
	bytes, err := ioutil.ReadFile(path.Join(p.Path, file) + ".ast.json")

	if err != nil {
		fmt.Println(file)
		panic(err)
	}

	return p.ParseJson(bytes, file)
}

// ParseJson parses a json into a Token-tree
func (p *PugAst) ParseJson(bytes []byte, file string) *Token {
	token := new(Token)

	err := json.Unmarshal(bytes, token)
	if err != nil {
		fmt.Println(file)
		panic(err)
	}

	return token
}
