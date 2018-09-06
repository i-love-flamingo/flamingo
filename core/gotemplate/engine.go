package gotemplate

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	flamingotemplate "flamingo.me/flamingo/framework/template"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

const pathSeparatorString = string(os.PathSeparator)

type (
	engine struct {
		templatesBasePath  string
		layoutTemplatesDir string
		debug              bool
		tplFuncs           flamingotemplate.FuncProvider
		tplCtxFuncs        flamingotemplate.CtxFuncProvider
		templates          map[string]*template.Template
		logger             flamingo.Logger
	}

	// urlFunc allows templates to access the routers `URL` helper method
	urlFunc struct {
		Router *router.Router `inject:""`
	}

	// getFunc allows templates to access the router's `get` method
	dataFunc struct {
		Router *router.Router `inject:""`
	}

	getFunc struct {
		Router *router.Router `inject:""`
	}
)

var (
	_    flamingotemplate.Func    = new(urlFunc)
	_    flamingotemplate.CtxFunc = new(getFunc)
	_    flamingotemplate.CtxFunc = new(dataFunc)
	lock                          = &sync.Mutex{}
)

// Inject engine dependencies
func (e *engine) Inject(
	tplFuncs flamingotemplate.FuncProvider,
	tplCtxFuncs flamingotemplate.CtxFuncProvider,
	logger flamingo.Logger,
	config *struct {
		TemplatesBasePath  string `inject:"config:gotemplates.engine.templates.basepath"`
		LayoutTemplatesDir string `inject:"config:gotemplates.engine.layout.dir"`
		Debug              bool   `inject:"config:debug.mode"`
	},
) {
	e.tplFuncs = tplFuncs
	e.tplCtxFuncs = tplCtxFuncs
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
		e.loadTemplates(ctx)
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
	for k, f := range e.tplCtxFuncs() {
		tplFuncs[k] = f.Func(ctx)
	}
	tpl.Funcs(tplFuncs)

	err = tpl.Execute(buf, data)

	defer span.End()

	return buf, err
}

func (e *engine) loadTemplates(ctx context.Context) {
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
		tplFuncs[k] = f.Func()
	}
	for k, f := range e.tplCtxFuncs() {
		tplFuncs[k] = f.Func(ctx)
	}

	layoutTemplate := template.Must(e.parseLayoutTemplates(functionsMap, tplFuncs))

	err := e.parseSiteTemplateDirectory(layoutTemplate, e.templatesBasePath)
	if err != nil {
		panic(err)
	}
}

// parses all layout templates in a template instance which is the base instance for all other templates
func (e *engine) parseLayoutTemplates(functionsMap template.FuncMap, funcs template.FuncMap) (*template.Template, error) {
	tpl := template.New("").Funcs(functionsMap).Funcs(funcs)

	if e.layoutTemplatesDir == "" {
		return tpl, nil
	}

	dir := e.templatesBasePath + pathSeparatorString + e.layoutTemplatesDir
	layoutFilesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	layoutFilesNames := make([]string, 0)
	for _, file := range layoutFilesInfo {
		if file.IsDir() {
			continue
		}
		layoutFilesNames = append(layoutFilesNames, dir+pathSeparatorString+file.Name())
	}

	return tpl.ParseFiles(layoutFilesNames...)
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

// Func as implementation of get method
func (g *getFunc) Func(ctx context.Context) interface{} {
	return func(what string, params ...map[string]interface{}) interface{} {
		var p = make(map[interface{}]interface{})
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return g.Router.Data(ctx, what, p)
	}
}

// Func as implementation of get method
func (d *dataFunc) Func(ctx context.Context) interface{} {
	return func(what string, params ...map[string]interface{}) interface{} {
		var p = make(map[interface{}]interface{})
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return d.Router.Data(ctx, what, p)
	}
}

// Func as implementation of url method
func (u *urlFunc) Func() interface{} {
	return func(where string, params ...map[string]interface{}) template.URL {
		var p = make(map[string]string)
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return template.URL(u.Router.URL(where, p).String())
	}
}
