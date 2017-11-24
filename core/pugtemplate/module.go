package pugtemplate

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
	"go.aoe.com/flamingo/core/pugtemplate/templatefunctions"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/template"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// Module for framework/pug_template
	Module struct {
		RouterRegistry *router.Registry `inject:""`
		Basedir        string           `inject:"config:pug_template.basedir"`
	}

	// TemplateFunctionInterceptor to use fixtype
	TemplateFunctionInterceptor struct {
		template.ContextFunction
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir(m.Basedir))))
	m.RouterRegistry.Route("/static/*n", "_static")

	m.RouterRegistry.Route("/_pugtpl/debug", "pugtpl.debug")
	m.RouterRegistry.Handle("pugtpl.debug", new(DebugController))

	m.RouterRegistry.Handle("page.template", func(ctx web.Context) interface{} {
		return ctx.Value("page.template")
	})

	// We bind the Template Engine to the ChildSingleton level (in case there is different config handling
	// We use the provider to make sure both are always the same injected type
	injector.Bind(pugjs.Engine{}).In(dingo.ChildSingleton).ToProvider(pugjs.NewEngine)
	injector.Bind((*template.Engine)(nil)).In(dingo.ChildSingleton).ToProvider(
		func(t *pugjs.Engine, i *dingo.Injector) template.Engine {
			return (template.Engine)(t)
		},
	)

	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.AssetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.JsMath{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.JsObject{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.DebugFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.JsJSON{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.URLFunc{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.GetFunc{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.DataFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.PriceFormatFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.StriptagsFunc{})

	m.loadmock("../src/layout/*")
	m.loadmock("../src/layout/*/*")
	m.loadmock("../src/layout/*/*/*")
	m.loadmock("../src/atom/*")
	m.loadmock("../src/molecule/*/*")
	m.loadmock("../src/organism/*")
	m.loadmock("../src/page/*/*")
	m.loadmock("../src/mock")
}

// DefaultConfig for setting pug-related config options
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"pug_template.basedir": "frontend/dist",
		"pug_template.debug":   true,
		"accounting": config.Map{
			"decimal":    ",",
			"thousand":   ".",
			"formatZero": "%s -,-",
			"format":     "%s %v",
		},
		"imageservice.base_url": "-",
		"imageservice.secret":   "-",
	}
}

func (m *Module) loadmock(where string) (interface{}, error) {
	matches, err := filepath.Glob(m.Basedir + "/" + where + "/*.mock.json")
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		b, e := ioutil.ReadFile(match)
		if e != nil {
			continue
		}
		var res interface{}
		json.Unmarshal(b, &res)
		name := strings.Replace(filepath.Base(match), ".mock.json", "", 1)
		if m.RouterRegistry.HandleIfNotSet(name, mockcontroller(name, res)) {
			log.Println("mocking because not set:", name)
		}
	}
	return nil, nil
}

func mockcontroller(name string, data interface{}) func(web.Context) interface{} {
	return func(ctx web.Context) interface{} {
		defer ctx.Profile("pugmock", name)()
		return data
	}
}
