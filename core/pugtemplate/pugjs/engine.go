package pugjs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"go.aoe.com/flamingo/framework/template"
	"go.aoe.com/flamingo/framework/web"

	"github.com/pkg/errors"
	"runtime"
)

type (
	// RenderState holds information about the pug abstract syntax tree
	renderState struct {
		path         string
		mixin        map[string]string
		mixinorder   []string
		mixincounter int
		mixinblocks  []string
		mixinblock   string
		funcs        FuncMap
		rawmode      bool
		doctype      string
		debug        bool
	}

	TemplateFunctionRegistryProvider func() *template.FunctionRegistry

	// Engine is the one and only javascript template engine for go ;)
	Engine struct {
		*sync.Mutex
		Basedir                   string `inject:"config:pug_template.basedir"`
		Debug                     bool   `inject:"config:debug.mode"`
		Assetrewrites             map[string]string
		templates                 map[string]*Template
		TemplateCode              map[string]string
		Webpackserver             bool
		TemplateFunctions         *template.FunctionRegistry
		TemplateFunctionsProvider TemplateFunctionRegistryProvider `inject:""`
	}
)

// NewEngine constructor
func NewEngine() *Engine {
	return &Engine{
		Mutex:        new(sync.Mutex),
		TemplateCode: make(map[string]string),
	}
}

func newRenderState(path string, debug bool) *renderState {
	return &renderState{
		path:  path,
		mixin: make(map[string]string),
		debug: debug,
	}
}

// LoadTemplates with an optional filter
func (e *Engine) LoadTemplates(filtername string) error {
	start := time.Now()

	e.Lock()
	defer e.Unlock()

	manifest, err := ioutil.ReadFile(path.Join(e.Basedir, "manifest.json"))
	if err == nil {
		json.Unmarshal(manifest, &e.Assetrewrites)
	}

	e.TemplateFunctions = e.TemplateFunctionsProvider()
	e.templates, err = e.compileDir(path.Join(e.Basedir, "template", "page"), "", filtername)
	if err != nil {
		return err
	}

	if _, err := http.Get("http://localhost:1337/assets/js/vendor.js"); err == nil {
		e.Webpackserver = true
	} else {
		e.Webpackserver = false
	}

	log.Println("Compiled templates in", time.Since(start))
	return nil
}

// compileDir returns a map of defined templates in directory dirname
func (e *Engine) compileDir(root, dirname, filtername string) (map[string]*Template, error) {
	result := make(map[string]*Template)

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
			tpls, err := e.compileDir(root, path.Join(dirname, filename.Name()), filtername)
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

				renderState := newRenderState(path.Join(e.Basedir, "template", "page"), e.Debug)
				renderState.funcs = FuncMap(e.TemplateFunctions.Populate())
				token, err := renderState.Parse(name)
				if err != nil {
					return nil, err
				}
				result[name], e.TemplateCode[name], err = renderState.TokenToTemplate(name, token)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return result, nil
}

// Render via html/pug_template
func (e *Engine) Render(ctx web.Context, templateName string, data interface{}) (io.Reader, error) {
	defer ctx.Profile("render", templateName)()

	p := strings.Split(templateName, "/")
	for i, v := range p {
		p[i] = strings.Title(v)
	}
	//ctx.WithValue("page.template", "page"+strings.Join(p, ""))
	ctx.WithValue("page.template", "page"+p[len(p)-1])

	// recompile
	if e.templates == nil {
		var finish = ctx.Profile("loadTemplates", "-all-")
		if err := e.LoadTemplates(""); err != nil {
			finish()
			return nil, err
		}
		finish()
	} else if e.Debug {
		var finish = ctx.Profile("debugReloadTemplates", templateName)
		if err := e.LoadTemplates(templateName); err != nil {
			finish()
			return nil, err
		}
		finish()
	}

	result := new(bytes.Buffer)

	tpl, ok := e.templates[templateName]
	if !ok {
		return nil, errors.Errorf(`Template %s not found!`, templateName)
	}

	templateInstance, err := tpl.Clone()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	funcs := make(FuncMap)
	funcs["__"] = fmt.Sprintf // todo translate
	for k, f := range e.TemplateFunctions.ContextAware {
		funcs[k] = f(ctx)
	}
	templateInstance.Funcs(funcs)

	// force GC to lower risk of runtime bugs in reflect.Value
	// should be fixed in go1.9.2
	runtime.GC()
	err = templateInstance.ExecuteTemplate(result, templateName, convert(data))
	if err != nil {
		errstr := err.Error() + "\n"
		for i, l := range strings.Split(e.TemplateCode[templateName], "\n") {
			errstr += fmt.Sprintf("%03d: %s\n", i+1, strings.TrimSpace(strings.TrimSuffix(l, `{{- "" -}}`)))
		}
		return nil, errors.New(errstr)
	}

	return result, nil
}
