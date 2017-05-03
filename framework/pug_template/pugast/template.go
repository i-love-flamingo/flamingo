package pugast

import (
	"bytes"
	"encoding/json"
	coretemplate "flamingo/framework/template"
	"flamingo/framework/web"
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

// PugTemplateEngine is the one and only javascript template engine for go ;)
type PugTemplateEngine struct {
	Basedir                   string `inject:"config:pug_template.basedir"`
	Debug                     bool   `inject:"config:debug.mode"`
	Assetrewrites             map[string]string
	templates                 map[string]*template.Template
	templatesLock             sync.Mutex
	Webpackserver             bool
	Ast                       *PugAst
	TemplateFunctions         *coretemplate.FunctionRegistry
	TemplateFunctionsProvider func() *coretemplate.FunctionRegistry `inject:""`
}

// loadTemplate gathers configuration and templates for the Engine
func (t *PugTemplateEngine) loadTemplates(filtername string) {
	start := time.Now()

	var err error

	t.templatesLock.Lock()
	defer t.templatesLock.Unlock()

	manifest, _ := ioutil.ReadFile(path.Join(t.Basedir, "manifest.json"))
	json.Unmarshal(manifest, &t.Assetrewrites)

	t.Ast = NewPugAst(path.Join(t.Basedir, "templates"))

	t.TemplateFunctions = t.TemplateFunctionsProvider()
	t.Ast.FuncMap = t.TemplateFunctions.Populate()

	t.templates, err = compileDir(t.Ast, path.Join(t.Basedir, "templates"), "", filtername)

	if err != nil {
		panic(err)
	}

	if _, err := http.Get("http://localhost:1337/assets/js/vendor.js"); err == nil {
		t.Webpackserver = true
	} else {
		t.Webpackserver = false
	}

	log.Println("Compiled templates in", time.Since(start))
}

// compileDir returns a map of defined templates in directory dirname
func compileDir(pugast *PugAst, root, dirname, filtername string) (map[string]*template.Template, error) {
	result := make(map[string]*template.Template)

	dir, err := os.Open(path.Join(root, dirname))
	if err != nil {
		return nil, err
	}

	filenames, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, filename := range filenames {
		if filename.IsDir() {
			tpls, err := compileDir(pugast, root, path.Join(dirname, filename.Name()), filtername)
			if err != nil {
				return nil, err
			}
			for k, v := range tpls {
				if result[k] == nil {
					result[k] = v
				}
			}
		} else {
			if strings.HasSuffix(filename.Name(), ".ast.json") {
				name := path.Join(dirname, filename.Name())
				name = name[:len(name)-len(".ast.json")]

				if filtername != "" && !strings.HasPrefix(name, filtername) {
					continue
				}

				result[name] = pugast.TokenToTemplate(name, pugast.Parse(name))
			}
		}
	}

	return result, nil
}

// Render via html/pug_template
func (t *PugTemplateEngine) Render(ctx web.Context, templateName string, data interface{}) io.Reader {
	defer ctx.Profile("render", templateName)()

	// recompile
	if t.templates == nil {
		var finish = ctx.Profile("loadTemplates", "-all-")
		t.loadTemplates("")
		finish()
	} else if t.Debug {
		var finish = ctx.Profile("debugReloadTemplates", templateName)
		t.loadTemplates(templateName)
		finish()
	}

	result := new(bytes.Buffer)

	templateInstance, err := t.templates[templateName].Clone()
	if err != nil {
		panic(err)
	}

	funcs := make(template.FuncMap)
	funcs["__"] = fmt.Sprintf // todo translate
	for k, f := range t.TemplateFunctions.ContextAware {
		funcs[k] = f(ctx)
	}
	templateInstance.Funcs(funcs)

	err = templateInstance.ExecuteTemplate(result, templateName, data)
	if err != nil {
		e := err.Error() + "\n"
		for i, l := range strings.Split(t.Ast.TplCode[templateName], "\n") {
			e += fmt.Sprintf("%03d: %s\n", i+1, l)
		}
		panic(e)
	}

	return result
}
