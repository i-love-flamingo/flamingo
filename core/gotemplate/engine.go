package gotemplate

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

const pathSeparatorString = string(os.PathSeparator)

type (
	templateFuncProvider func() map[string]flamingo.TemplateFunc

	engine struct {
		templatesBasePath  string
		layoutTemplatesDir string
		debug              bool
		tplFuncs           templateFuncProvider
		templates          map[string]*template.Template
		logger             flamingo.Logger
	}

	urlRouter interface {
		URL(name string, params map[string]string) (*url.URL, error)
		Data(ctx context.Context, handler string, params map[interface{}]interface{}) interface{}
	}

	// urlFunc allows templates to access the routers `URL` helper method
	urlFunc struct {
		router urlRouter
	}

	// getFunc allows templates to access the router's `get` method
	dataFunc struct {
		router urlRouter
	}

	getFunc struct {
		router urlRouter
	}
)

var (
	_    flamingo.TemplateFunc = new(urlFunc)
	_    flamingo.TemplateFunc = new(getFunc)
	_    flamingo.TemplateFunc = new(dataFunc)
	lock                       = &sync.Mutex{}
)

// Inject engine dependencies
func (e *engine) Inject(
	tplFuncs templateFuncProvider,
	logger flamingo.Logger,
	config *struct {
		TemplatesBasePath  string `inject:"config:gotemplates.engine.templates.basepath"`
		LayoutTemplatesDir string `inject:"config:gotemplates.engine.layout.dir"`
		Debug              bool   `inject:"config:debug.mode"`
	},
) {
	e.tplFuncs = tplFuncs
	e.templatesBasePath = config.TemplatesBasePath
	e.layoutTemplatesDir = config.LayoutTemplatesDir
	e.debug = config.Debug
	e.logger = logger
}

func (e *engine) Render(ctx context.Context, name string, data interface{}) (io.Reader, error) {
	ctx, span := trace.StartSpan(ctx, "gotemplate/Render")
	defer span.End()

	lock.Lock()
	if e.debug || e.templates == nil {
		err := e.loadTemplates(ctx)
		if err != nil {
			return nil, err
		}
	}
	lock.Unlock()

	_, span = trace.StartSpan(ctx, "gotemplate/Execute")
	buf := &bytes.Buffer{}

	if _, ok := e.templates[name+".html"]; !ok {
		return nil, errors.New("Could not find the template " + name + ".html")
	}
	tpl, err := e.templates[name+".html"].Clone()
	if err != nil {
		return nil, err
	}

	tplFuncs := template.FuncMap{}
	for k, f := range e.tplFuncs() {
		tplFuncs[k] = f.Func(ctx)
	}
	tpl.Funcs(tplFuncs)

	err = tpl.Execute(buf, data)

	defer span.End()

	return buf, err
}

func (e *engine) loadTemplates(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "gotemplate/loadTemplates")
	defer span.End()

	e.templates = make(map[string]*template.Template, 0)

	functionsMap := template.FuncMap{
		"Upper": strings.ToUpper,
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"map": func(p ...interface{}) map[string]interface{} {
			res := make(map[string]interface{})
			for i := 0; i < len(p); i += 2 {
				res[fmt.Sprint(p[i])] = p[i+1]
			}
			return res
		},
	}

	tplFuncs := template.FuncMap{}
	for k, f := range e.tplFuncs() {
		tplFuncs[k] = f.Func(ctx)
	}

	layoutTemplate, err := e.parseLayoutTemplates(functionsMap, tplFuncs)
	if err != nil {
		return err
	}

	err = e.parseSiteTemplateDirectory(layoutTemplate, e.templatesBasePath)
	if err != nil {
		return err
	}

	return nil
}

// parses all layout templates in a template instance which is the base instance for all other templates
func (e *engine) parseLayoutTemplates(functionsMap template.FuncMap, funcs template.FuncMap) (*template.Template, error) {
	tpl := template.New("").Funcs(functionsMap).Funcs(funcs)

	if e.layoutTemplatesDir == "" {
		return tpl, nil
	}

	dir := e.templatesBasePath + pathSeparatorString + e.layoutTemplatesDir

	layoutFilesNames := make([]string, 0)
	err := filepath.Walk(
		dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			layoutFilesNames = append(layoutFilesNames, path)

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	for _, file := range layoutFilesNames {
		tContent, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		templateName, err := filepath.Rel(dir, file)
		t := tpl.New(templateName)

		_, err = t.Parse(string(tContent))
		if err != nil {
			return nil, err
		}
	}

	return tpl, nil
}

// parses all templates from a given directory into a clone of the given layout template, so that all layouts are available
func (e *engine) parseSiteTemplateDirectory(layoutTemplate *template.Template, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		t := template.Must(layoutTemplate.Clone())
		fullName := dir + pathSeparatorString + f.Name()
		if f.IsDir() {
			err = e.parseSiteTemplateDirectory(layoutTemplate, fullName)
			if err != nil {
				return err
			}
			continue
		}
		tContent, err := ioutil.ReadFile(fullName)
		if err != nil {
			return err
		}

		templateName := strings.TrimPrefix(fullName, e.templatesBasePath+pathSeparatorString)
		parsedTemplate, err := t.Parse(string(tContent))
		if err != nil {
			e.logger.WithField("category", "gotemplate").Error(err)
			continue
		}
		e.templates[templateName] = template.Must(parsedTemplate, err)
	}

	return nil
}

func (g *getFunc) Inject(router urlRouter) *getFunc {
	g.router = router
	return g
}

// TemplateFunc as implementation of get method
func (g *getFunc) Func(ctx context.Context) interface{} {
	return func(what string, params ...map[string]interface{}) interface{} {
		var p = make(map[interface{}]interface{})
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return g.router.Data(ctx, what, p)
	}
}

func (d *dataFunc) Inject(router urlRouter) *dataFunc {
	d.router = router
	return d
}

// TemplateFunc as implementation of get method
func (d *dataFunc) Func(ctx context.Context) interface{} {
	return func(what string, params ...map[string]interface{}) interface{} {
		var p = make(map[interface{}]interface{})
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return d.router.Data(ctx, what, p)
	}
}

func (u *urlFunc) Inject(router urlRouter) *urlFunc {
	u.router = router
	return u
}

// TemplateFunc as implementation of url method
func (u *urlFunc) Func(context.Context) interface{} {
	return func(where string, params ...map[string]interface{}) template.URL {
		var p = make(map[string]string)
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		url, _ := u.router.URL(where, p)
		return template.URL(url.String())
	}
}
