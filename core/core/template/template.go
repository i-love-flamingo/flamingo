package template

/*

WARNING!!!

This is a work in progress!

Please do not judge this file! Please :)

*/

import (
	"bytes"
	"encoding/json"
	"flamingo/core/core/app"
	"flamingo/core/core/app/web"
	"flamingo/core/core/template/pug-ast"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"fmt"

	"github.com/fatih/structs"
)

var (
	assetrewrites map[string]string
	templates     map[string]*template.Template
	templatesLock sync.Mutex
	webpackserver bool
)

func init() {
	//loadTemplates()
}

func loadTemplates() {
	start := time.Now()

	TemplateFunctions.Populate()

	var err error

	templatesLock.Lock()
	defer templatesLock.Unlock()

	manifest, _ := ioutil.ReadFile("frontend/dist/manifest.json")
	json.Unmarshal(manifest, &assetrewrites)

	pugast := node.PugAst{
		Path: "frontend/dist/templates",
	}
	templates, err = compile(&pugast, "frontend/dist/templates", "")

	if err != nil {
		panic(err)
	}

	if _, err := http.Get("http://localhost:1337/assets/js/vendor.js"); err == nil {
		webpackserver = true
	} else {
		webpackserver = false
	}

	log.Println("Compiled templates in", time.Since(start))
}

func compile(pugast *node.PugAst, root, dirname string) (map[string]*template.Template, error) {
	result := make(map[string]*template.Template)

	dir, _ := os.Open(path.Join(root, dirname))

	filenames, _ := dir.Readdir(-1)

	for _, filename := range filenames {
		if filename.IsDir() {
			tpls, _ := compile(pugast, root, path.Join(dirname, filename.Name()))
			for k, v := range tpls {
				if result[k] == nil {
					result[k] = v
				}
			}
		} else {
			if strings.HasSuffix(filename.Name(), ".ast.json") {
				name := path.Join(dirname, filename.Name())
				name = name[:len(name)-len(".ast.json")]
				result[name] = pugast.TokenToTemplate(name, pugast.Parse(name))
			}
		}
	}

	return result, nil
}

// Render via hmtl/template
func Render(app *app.App, ctx web.Context, tpl string, data interface{}) io.Reader {
	buf := new(bytes.Buffer)

	// recompile
	if app.Debug {
		loadTemplates()
	}

	t, _ := templates[tpl].Clone()

	funcs := make(template.FuncMap)
	funcs["__"] = fmt.Sprintf // todo translate
	for k, f := range TemplateFunctions.contextaware {
		funcs[k] = f(ctx)
	}
	t.Funcs(funcs)

	var d interface{}
	if data != nil {
		d = structs.Map(data)
	} else {
		d = data
	}
	err := t.ExecuteTemplate(buf, tpl, d)
	if err != nil {
		panic(err)
	}

	return buf
}
