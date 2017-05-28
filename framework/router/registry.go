package router

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type (
	/*
	  Registry holds a list of all routes and handlers to be registered in modules.

	  We have:
	  routes: key-params -> path, for reverse routes

	  path: url-pattern -> key+params

	  handler: key -> controller
	*/
	Registry struct {
		handler map[string]Controller
		routes  []*handler
	}

	handler struct {
		path    *Path
		handler string
		params  map[string]*param
	}

	param struct {
		value    string
		optional bool
	}
)

// NewRegistry creates a new Registry
func NewRegistry() *Registry {
	return &Registry{
		handler: make(map[string]Controller),
	}
}

// Handle assigns a controller to a name
func (registry *Registry) Handle(name string, controller Controller) {
	registry.handler[name] = controller
}

// Route assigns a route to a handler
func (registry *Registry) Route(path, handler string) {
	var h = parseHandler(handler)
	h.path = NewPath(path)
	registry.routes = append(registry.routes, h)
}

// Mount auto-generates a controller name from the path
func (registry *Registry) Mount(path string, controller Controller) {
	var p = NewPath(path)
	var name = strings.Replace(strings.Trim(path, "/.: "), "/", ".", -1)

	registry.Handle(name, controller)
	registry.Route(path, fmt.Sprintf("%s(%s)", name, strings.Join(p.params, ", ")))
}

func parseHandler(h string) *handler {
	var tmp = strings.SplitN(h, "(", 2)
	h = tmp[0]

	var newHandler = &handler{
		handler: h,
		params:  make(map[string]*param),
	}

	if len(tmp) == 2 {
		newHandler.params = parseParams(tmp[1][:len(tmp[1])-1])
	}

	return newHandler
}

// list: foo, bar, x ?= "y", z = "a"
func parseParams(list string) map[string]*param {
	var params = make(map[string]*param)

	var name, val string
	var optional bool
	var quote byte
	var readto = &name

	for i := 0; i < len(list); i++ {
		if list[i] != quote && quote != 0 {
			if list[i] != '\\' {
				*readto += string(list[i])
			} else {
				i++
				*readto += string(list[i])
			}
			continue
		}

		switch list[i] {
		case '\\':
			i++
			*readto += string(list[i])

		case quote:
			quote = 0

		case '"', '\'':
			quote = list[i]
			val = ""
			readto = &val

		case ',':
			name = strings.TrimSpace(name)
			params[name] = &param{
				optional: optional,
				value:    val,
			}
			optional = false
			name = ""
			val = ""
			readto = &name

		case '?':
			optional = true

		case '=':
			continue

		default:
			*readto += string(list[i])
		}
	}

	name = strings.TrimSpace(name)
	if name != "" {
		params[name] = &param{
			optional: optional,
			value:    val,
		}
	}

	return params
}

// Reverse builds the path from a named route with params
func (registry *Registry) Reverse(name string, params map[string]string) (string, error) {
	var keys = make([]string, len(params))
	var i = 0
	for k := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

routeloop:
	for _, handler := range registry.routes {
		if handler.handler == name {
			var renderparams = make(map[string]string)

			for key, param := range handler.params {
				v, ok := params[key]

				// unset not-optional param
				if !param.optional && !ok {
					continue routeloop
				}

				// not-optional param set with wrong value
				if !param.optional && ok && param.value != "" && param.value != v {
					continue routeloop
				}
				renderparams[key] = param.value
			}

			for k, v := range params {
				renderparams[k] = v
			}

			return handler.path.Render(renderparams)
		}
	}
	return "", errors.New("Reverse for " + name + " not found")
}

// Match a request path
func (registry *Registry) Match(path string) (controller Controller, params map[string]string) {
	for _, handler := range registry.routes {
		if match := handler.path.Match(path); match != nil {
			controller = registry.handler[handler.handler]
			params = make(map[string]string)
			for k, param := range handler.params {
				params[k] = param.value
			}
			for k, v := range match.Values {
				params[k] = v
			}
			return
		}
	}
	return
}

// MatchRequest matches a http Request (with GET and path parameters)
func (registry *Registry) MatchRequest(req *http.Request) (Controller, map[string]string) {
	var path = req.URL.Path

matchloop:
	for _, handler := range registry.routes {
		if match := handler.path.Match(path); match != nil {
			controller := registry.handler[handler.handler]
			params := make(map[string]string)
			if len(handler.params) > 0 {
				for k, param := range handler.params {
					if !param.optional && param.value != "" {
						params[k] = param.value
					} else if v, ok := match.Values[k]; ok {
						params[k] = v
					} else if val := req.URL.Query().Get(k); val != "" {
						params[k] = val
					} else if !param.optional && param.value == "" {
						continue matchloop
					} else {
						params[k] = param.value
					}
				}
			} else {
				params = match.Values
			}
			return controller, params
		}
	}
	return nil, nil
}
