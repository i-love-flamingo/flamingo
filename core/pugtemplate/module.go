package pugtemplate

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"flamingo.me/flamingo/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/core/pugtemplate/templatefunctions"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/template"
	"flamingo.me/flamingo/framework/web"
)

type (
	// Module for framework/pug_template
	Module struct {
		RouterRegistry *router.Registry `inject:""`
		Basedir        string           `inject:"config:pug_template.basedir"`
		DefaultMux     *http.ServeMux   `inject:",optional"`
	}

	// TemplateFunctionInterceptor to use fixtype
	TemplateFunctionInterceptor struct {
		template.ContextFunction
	}
	assetFileSystem struct {
		fs http.FileSystem
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	//m.RouterRegistry.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir(m.Basedir))))
	m.RouterRegistry.Handle("_static", http.StripPrefix("/static/", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		if origin != "" {
			//TODO - configure whitelist
			rw.Header().Add("Access-Control-Allow-Origin", origin)
		}
		http.FileServer(assetFileSystem{http.Dir("frontend/dist/")}).ServeHTTP(rw, req)
	})))
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

	if m.DefaultMux != nil {
		m.DefaultMux.HandleFunc("/assets/", func(rw http.ResponseWriter, req *http.Request) {
			origin := req.Header.Get("Origin")
			if origin != "" {
				//TODO - configure whitelist
				rw.Header().Add("Access-Control-Allow-Origin", origin)
			}
			if r, e := http.Get("http://localhost:1337" + req.RequestURI); e == nil {
				io.Copy(rw, r.Body)
			} else {
				fileServer := http.FileServer(assetFileSystem{http.Dir("frontend/dist/")})
				fileServer.ServeHTTP(rw, req)
			}
		})
	}

	m.RouterRegistry.Route("/assets/*f", "_pugtemplate.assets")
	m.RouterRegistry.Handle("_pugtemplate.assets", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		if origin != "" {
			//TODO - configure whitelist
			rw.Header().Add("Access-Control-Allow-Origin", origin)
		}
		if r, e := http.Get("http://localhost:1337" + req.RequestURI); e == nil {
			io.Copy(rw, r.Body)
		} else {
			http.FileServer(assetFileSystem{http.Dir("frontend/dist/")}).ServeHTTP(rw, req)
		}
	}))

	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.AssetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.JsMath{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.JsObject{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.DebugFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.JsJSON{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.StartsWithFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.TruncateFunc{})

	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.URLFunc{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.TryURLFunc{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.GetFunc{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.DataFunc{})
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
		"pug_template.basedir":  "frontend/dist",
		"pug_template.debug":    true,
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

func copyHeaders(r *http.Response, w http.ResponseWriter) {
	for key, values := range r.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}

func (nfs assetFileSystem) Open(path string) (http.File, error) {
	path = strings.Replace(path, "/assets/", "", 1)
	log.Println(path)
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		return nil, errors.New("not allowed")
	}

	return f, nil
}
