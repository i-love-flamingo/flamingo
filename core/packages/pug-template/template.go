package template

import (
	"bytes"
	"encoding/json"
	"flamingo/core/flamingo/web"
	"flamingo/core/packages/pug-template/pugast"
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

type PugTemplateEngine struct {
	basedir       string
	assetrewrites map[string]string
	templates     map[string]*template.Template
	templatesLock sync.Mutex
	webpackserver bool
	ast           *pugast.PugAst
}

func NewPugTemplateEngine(basedir string) *PugTemplateEngine {
	return &PugTemplateEngine{
		basedir: basedir,
	}
}

func (t *PugTemplateEngine) loadTemplates() {
	start := time.Now()

	var err error

	t.templatesLock.Lock()
	defer t.templatesLock.Unlock()

	TemplateFunctions.Populate()

	manifest, _ := ioutil.ReadFile(path.Join(t.basedir, "manifest.json"))
	json.Unmarshal(manifest, &t.assetrewrites)

	t.ast = pugast.NewPugAst(path.Join(t.basedir, "templates"))
	t.templates, err = compileDir(t.ast, path.Join(t.basedir, "templates"), "")

	if err != nil {
		panic(err)
	}

	if _, err := http.Get("http://localhost:1337/assets/js/vendor.js"); err == nil {
		t.webpackserver = true
	} else {
		t.webpackserver = false
	}

	log.Println("Compiled templates in", time.Since(start))
}

func compileDir(pugast *pugast.PugAst, root, dirname string) (map[string]*template.Template, error) {
	result := make(map[string]*template.Template)

	dir, _ := os.Open(path.Join(root, dirname))

	filenames, _ := dir.Readdir(-1)

	for _, filename := range filenames {
		if filename.IsDir() {
			tpls, _ := compileDir(pugast, root, path.Join(dirname, filename.Name()))
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

// Render via hmtl/pug-template
func (t *PugTemplateEngine) Render(ctx web.Context, templateName string, data interface{}) io.Reader {
	// recompile
	/*
		if router.Debug {
			loadTemplates()
		}
	*/
	if t.templates == nil {
		t.loadTemplates()
	}

	result := new(bytes.Buffer)

	templateInstance, err := t.templates[templateName].Clone()
	if err != nil {
		panic(err)
	}

	funcs := make(template.FuncMap)
	funcs["__"] = fmt.Sprintf // todo translate
	for k, f := range TemplateFunctions.contextaware {
		funcs[k] = f(ctx)
	}
	templateInstance.Funcs(funcs)

	err = templateInstance.ExecuteTemplate(result, templateName, data)
	if err != nil {
		e := err.Error() + "\n"
		for i, l := range strings.Split(t.ast.TplCode[templateName], "\n") {
			e += fmt.Sprintf("%03d: %s\n", i+1, l)
		}
		panic(e)
	}

	return result
}
