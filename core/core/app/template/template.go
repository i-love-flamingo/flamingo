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
	"flamingo/core/core/app/template/pug-ast"
	"flamingo/core/core/app/web"
	"fmt"
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
)

var (
	assetrewrites map[string]string
	templates     map[string]*template.Template
	templatesLock sync.Mutex
	webpackserver bool
)

func init() {
	loadTemplates()
}

func loadTemplates() {
	start := time.Now()

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

	t.Funcs(template.FuncMap{
		"asset": func(a string) template.URL {
			if webpackserver {
				return template.URL("/assets/" + a)
			}

			url := app.Url("_static")
			aa := strings.Split(a, "/")
			aaa := aa[len(aa)-1]
			var result string
			if assetrewrites[aaa] != "" {
				result = url.String() + "/" + assetrewrites[aaa]
			} else {
				result = url.String() + "/" + a
			}
			ctx.Push(result, nil)
			return template.URL(result)
		},
		"__": fmt.Sprintf, // todo translate
		"__get": func(what string) interface{} {
			if what == "user.name" {
				return "testuser"
			}
			return []map[string]string{{"url": "url1", "name": "item1"}, {"url": "url2", "name": "name2"}}
		},
		"get": func(what string) interface{} {
			log.Println("get", what)
			return app.Get(what, ctx)
		},
	})

	err := t.ExecuteTemplate(buf, tpl, map[string]interface{}{
		"isProductionBuild": !webpackserver,
		"classBody":         "default",
		"title":             "Home",
		"site": map[string]interface{}{
			"title": "Auckland Airport",
		},
	})
	if err != nil {
		panic(err)
	}

	return buf
}
