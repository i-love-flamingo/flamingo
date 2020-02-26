package gotemplate

import (
	"bytes"
	"context"
	"errors"
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
		Relative(name string, params map[string]string) (*url.URL, error)
		Data(ctx context.Context, handler string, params map[interface{}]interface{}) interface{}
	}
)

var (
	lock = &sync.Mutex{}
)

// Inject engine dependencies
func (e *engine) Inject(
	tplFuncs templateFuncProvider,
	logger flamingo.Logger,
	config *struct {
		TemplatesBasePath  string `inject:"config:core.gotemplate.engine.templates.basepath"`
		LayoutTemplatesDir string `inject:"config:core.gotemplate.engine.layout.dir"`
		Debug              bool   `inject:"config:flamingo.debug.mode"`
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
			lock.Unlock()
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
			for i := 0; i < len(p) && len(p)%2 == 0; i += 2 {
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
		if err != nil {
			return nil, err
		}
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
		// operating systems like windows use \ instead of /, so `c.responder.Render("foo/bar", ...)` will not
		// resolve the template properly, as it is known as `foo\bar`. this makes sure we register the template
		// as `foo/bar` as well as `foo\bar`.
		if pathSeparatorString != "/" {
			e.templates[strings.Replace(templateName, pathSeparatorString, "/", -1)] = template.Must(parsedTemplate, err)
		}
	}

	return nil
}
