package router

import (
	"net/http"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type (
	// Registry holds a list of all routes and handlers to be registered in modules.
	//
	// We have:
	// routes: key-params -> path, for reverse routes
	//
	// path: url-pattern -> key+params
	//
	// Handler: key -> controller
	Registry struct {
		handler map[string]Controller
		routes  []*Handler
		alias   map[string]*Handler
	}

	// Handler defines a concrete controller
	Handler struct {
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
		alias:   make(map[string]*Handler),
	}
}

// Handle assigns a controller to a name
func (registry *Registry) Handle(name string, controller Controller) {
	registry.handler[name] = controller
}

// HandleIfNotSet assigns a controller to a name if not already set
func (registry *Registry) HandleIfNotSet(name string, controller Controller) bool {
	if _, ok := registry.handler[name]; ok {
		return false
	}
	registry.handler[name] = controller
	return true
}

// GetControllerForHandle returns Controller for a Handle Name
func (registry *Registry) GetControllerForHandle(name string) (Controller, error) {
	if val, ok := registry.handler[name]; ok {
		return val, nil
	}
	return nil, errors.New("Handle not found")
}

// Route assigns a route to a Handler
func (registry *Registry) Route(path, handler string) {
	var h = parseHandler(handler)
	h.path = NewPath(path)
	if len(h.params) == 0 {
		h.params = parseParams(strings.Join(h.path.params, ", "))
	}
	registry.routes = append(registry.routes, h)
}

// GetRoutes returns registered Routes
func (registry *Registry) GetRoutes() []*Handler {
	return registry.routes
}

// GetHandler returns registered Routes
func (registry *Registry) GetHandler() map[string]Controller {
	return registry.handler
}

// Alias for an existing router definition
func (registry *Registry) Alias(name, to string) {
	registry.alias[name] = parseHandler(to)
}

func parseHandler(h string) *Handler {
	var tmp = strings.SplitN(h, "(", 2)
	h = tmp[0]

	var newHandler = &Handler{
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
	if alias, ok := registry.alias[name]; ok {
		name = alias.handler
		for name, param := range alias.params {
			params[name] = param.value
		}
	}

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

			// set handler default parameters
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

			// add Reverse parameters
			for k, v := range params {
				renderparams[k] = v
			}

			// validate if all parameters have been used
			for key := range params {
				if _, ok := handler.params[key]; !ok {
					continue routeloop
				}
			}

			return handler.path.Render(renderparams)
		}
	}
	return "", errors.Errorf("Reverse for %q not found, parameters: %v", name, params)
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
func (registry *Registry) MatchRequest(req *http.Request) (Controller, map[string]string, *Handler) {
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
			return controller, params, handler
		}
	}
	return nil, nil, nil
}

// GetPath getter
func (handler *Handler) GetPath() string {
	return handler.path.path
}

// GetHandlerName getter
func (handler *Handler) GetHandlerName() string {
	return handler.handler
}
