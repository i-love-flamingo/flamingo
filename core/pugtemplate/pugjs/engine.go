package pugjs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"flamingo.me/flamingo/v3/framework/event"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/template"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

// BUG: the template loading is far from optimal, if debug is false and the loading fails we might end up in a broken situation

type (
	// RenderState holds information about the pug abstract syntax tree
	renderState struct {
		path         string
		mixin        map[string]string
		mixincalls   map[string]struct{}
		mixinorder   []string
		mixincounter int
		mixinblocks  []string
		mixinblock   string
		funcs        FuncMap
		rawmode      bool
		doctype      string
		debug        bool
		eventRouter  event.Router
		logger       flamingo.Logger
	}

	// Engine is the one and only javascript template engine for go ;)
	Engine struct {
		*sync.RWMutex
		Basedir         string `inject:"config:pug_template.basedir"`
		Debug           bool   `inject:"config:debug.mode"`
		Assetrewrites   map[string]string
		templatesLoaded int32
		templates       map[string]*Template
		TemplateCode    map[string]string
		Webpackserver   bool
		EventRouter     event.Router             `inject:""`
		FuncProvider    template.FuncProvider    `inject:""`
		CtxFuncProvider template.CtxFuncProvider `inject:""`
		Logger          flamingo.Logger          `inject:""`
	}

	// EventSubscriber is the event subscriber for Engine
	EventSubscriber struct {
		engine *Engine
		logger flamingo.Logger
	}
)

var (
	rt             = stats.Int64("flamingo/pugtemplate/render", "pugtemplate render times", stats.UnitMilliseconds)
	templateKey, _ = tag.NewKey("template")
)

func init() {
	opencensus.View("flamingo/pugtemplate/render", rt, view.Distribution(50, 100, 250, 500, 1000, 2000), templateKey)
}

// NewEngine constructor
func NewEngine() *Engine {
	return &Engine{
		RWMutex:      new(sync.RWMutex),
		TemplateCode: make(map[string]string),
	}
}

func newRenderState(path string, debug bool, eventRouter event.Router, logger flamingo.Logger) *renderState {
	return &renderState{
		path:        path,
		mixin:       make(map[string]string),
		mixincalls:  make(map[string]struct{}),
		debug:       debug,
		eventRouter: eventRouter,
		logger:      logger,
	}
}

// Inject injects the EventSubscibers dependencies
func (e *EventSubscriber) Inject(engine *Engine, logger flamingo.Logger) {
	e.engine = engine
	e.logger = logger
}

// Notify the event subscriper
func (e *EventSubscriber) Notify(event event.Event) {
	if _, ok := event.(*flamingo.AppStartupEvent); ok {
		e.logger.Info("preloading templates on flamingo.AppStartupEvent")
		go e.engine.LoadTemplates("")
	}
}

// LoadTemplates with an optional filter
func (e *Engine) LoadTemplates(filtername string) error {
	e.Lock()
	defer e.Unlock()

	if !atomic.CompareAndSwapInt32(&e.templatesLoaded, 0, 1) && filtername == "" {
		return errors.New("Can not preload all templates again")
	}

	start := time.Now()

	manifest, err := ioutil.ReadFile(path.Join(e.Basedir, "manifest.json"))
	if err == nil {
		json.Unmarshal(manifest, &e.Assetrewrites)
	}

	e.templates, err = e.compileDir(path.Join(e.Basedir, "template", "page"), "", filtername)
	if err != nil {
		atomic.StoreInt32(&e.templatesLoaded, 0) // bail out :(
		return err
	}

	if _, err := http.Get("http://localhost:1337/assets/js/vendor.js"); err == nil {
		e.Webpackserver = true
	} else {
		e.Webpackserver = false
	}

	e.Logger.Info("Compiled templates in", time.Since(start))
	return nil
}

// compileDir returns a map of defined templates in directory dirname
func (e *Engine) compileDir(root, dirname, filtername string) (map[string]*Template, error) {
	result := make(map[string]*Template)

	dir, err := os.Open(path.Join(root, dirname))
	if err != nil {
		return nil, err
	}

	defer dir.Close()

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

				renderState := newRenderState(path.Join(e.Basedir, "template", "page"), e.Debug, e.EventRouter, e.Logger)
				renderState.funcs = FuncMap{}

				for k, f := range e.FuncProvider() {
					renderState.funcs[k] = f.Func()
				}
				for k, f := range e.CtxFuncProvider() {
					renderState.funcs[k] = f.Func
				}

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

var renderChan = make(chan struct{}, 8)

// RenderPartials is used for progressive enhancements / rendering of partial template areas
// usually this is requested via the appropriate javascript headers and taken care of in the framework renderer
func (e *Engine) RenderPartials(ctx context.Context, templateName string, data interface{}, partials []string) (map[string]string, error) {
	res := make(map[string]string, len(partials))

	for _, partial := range partials {
		buf, err := e.Render(ctx, templateName+".partial/"+partial, data)
		if err != nil {
			return nil, err
		}
		str, err := ioutil.ReadAll(buf)
		if err != nil {
			return nil, err
		}
		res[partial] = string(str)
	}

	return res, nil
}

// Render via html/pug_template
func (e *Engine) Render(ctx context.Context, templateName string, data interface{}) (io.Reader, error) {
	ctx, span := trace.StartSpan(ctx, "pug/render")
	defer span.End()

	span.Annotate(nil, templateName)

	//block if buffered channel size is reached
	renderChan <- struct{}{}
	defer func() {
		//release one entry from channel (will release one block)
		<-renderChan
	}()

	p := strings.Split(templateName, "/")
	for i, v := range p {
		p[i] = strings.Title(v)
	}
	page := p[len(p)-1]
	if len(p) >= 2 && p[len(p)-2] != page {
		page = p[len(p)-2] + p[len(p)-1]
	}
	ctx = context.WithValue(ctx, "page.template", "page"+page)

	// recompile, make sure to fully load only once!
	if atomic.LoadInt32(&e.templatesLoaded) == 0 && !e.Debug {
		_, spanLoad := trace.StartSpan(ctx, "pug/loadAllTemplates")
		if err := e.LoadTemplates(""); err != nil {
			spanLoad.End()
			return nil, err
		}
		spanLoad.End()
	} else if e.Debug {
		_, spanLoad := trace.StartSpan(ctx, "pug/loadTemplate")
		spanLoad.Annotate(nil, templateName)
		if err := e.LoadTemplates(templateName); err != nil {
			spanLoad.End()
			return nil, err
		}
		spanLoad.End()
	}

	// make sure template loading has finished by now!
	e.RLock()

	result := new(bytes.Buffer)

	templateInstance, ok := e.templates[templateName]
	e.RUnlock()
	if !ok {
		return nil, errors.Errorf(`Template %s not found!`, templateName)
	}

	ctx, execSpan := trace.StartSpan(ctx, "pug/execute")
	execSpan.Annotate(nil, templateName)
	start := time.Now()
	err := templateInstance.ExecuteTemplate(ctx, result, templateName, convert(data))
	execSpan.End()
	ctx, _ = tag.New(ctx, tag.Upsert(templateKey, templateName))
	stats.Record(ctx, rt.M(time.Since(start).Nanoseconds()/1000000))

	if err != nil {
		errstr := err.Error() + "\n"
		for i, l := range strings.Split(e.TemplateCode[templateName], "\n") {
			errstr += fmt.Sprintf("%03d: %s\n", i+1, strings.TrimSpace(strings.TrimSuffix(l, `{{- "" -}}`)))
		}
		return nil, errors.New(errstr)
	}

	return result, nil
}
