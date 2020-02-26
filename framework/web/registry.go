package web

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"flamingo.me/dingo"
)

type (
	// RouterRegistry holds a list of all routes and handlers to be registered in modules.
	//
	// We have:
	// routes: key-params -> path, for reverse routes
	//
	// path: url-pattern -> key+params
	//
	// Handler: key -> Controller
	RouterRegistry struct {
		handler map[string]handlerAction
		routes  []*Handler
		alias   map[string]*Handler
	}

	// Handler defines a concrete Controller
	Handler struct {
		path     *Path
		handler  string
		params   map[string]*param
		catchall bool
	}

	handlerAction struct {
		method map[string]Action
		any    Action
		data   DataAction
	}

	matchedHandler struct {
		handlerAction handlerAction
		handler       *Handler
		match         *Match
	}

	matchedHandlers []*matchedHandler

	param struct {
		value    string
		optional bool
	}

	// RoutesModule defines a router RoutesModule, which is able to register routes
	RoutesModule interface {
		Routes(registry *RouterRegistry)
	}
)

// BindRoutes is a convenience helper to multi-bind router modules
func BindRoutes(injector *dingo.Injector, m RoutesModule) {
	injector.BindMulti(new(RoutesModule)).To(m)
}

// NewRegistry creates a new RouterRegistry
func NewRegistry() *RouterRegistry {
	return &RouterRegistry{
		handler: make(map[string]handlerAction),
		alias:   make(map[string]*Handler),
	}
}

func (ha *handlerAction) set(method string, action Action) {
	if ha.method == nil {
		ha.method = make(map[string]Action, 1)
	}
	ha.method[method] = action
}

func (ha *handlerAction) setAny(any Action) {
	ha.any = any
}

func (ha *handlerAction) setData(data DataAction) {
	ha.data = data
}

func (mh matchedHandlers) getHandleAny() *matchedHandler {
	for _, matched := range mh {
		if matched.handlerAction.any != nil {
			return matched
		}
	}

	return nil
}

func (mh matchedHandlers) hasMethod(method string) bool {
	for _, matched := range mh {
		if _, ok := matched.handlerAction.method[method]; ok {
			return true
		}
	}

	return false
}

// HandleAny serves as a fallback to handle HTTP requests which are not taken care of by other handlers
func (registry *RouterRegistry) HandleAny(name string, action Action) {
	ha := registry.handler[name]
	ha.setAny(action)
	registry.handler[name] = ha
}

// HandleData sets the controllers data action
func (registry *RouterRegistry) HandleData(name string, action DataAction) {
	ha := registry.handler[name]
	ha.setData(action)
	registry.handler[name] = ha
}

// HandleMethod handles requests for the specified HTTP Method
func (registry *RouterRegistry) HandleMethod(method, name string, action Action) {
	ha := registry.handler[name]
	ha.set(method, action)
	registry.handler[name] = ha
}

// HandleGet handles a HTTP GET request
func (registry *RouterRegistry) HandleGet(name string, action Action) {
	registry.HandleMethod(http.MethodGet, name, action)
}

// HandlePost handles HTTP POST requests
func (registry *RouterRegistry) HandlePost(name string, action Action) {
	registry.HandleMethod(http.MethodPost, name, action)
}

// HandlePut handles HTTP PUT requests
func (registry *RouterRegistry) HandlePut(name string, action Action) {
	registry.HandleMethod(http.MethodPut, name, action)
}

// HandleDelete handles HTTP DELETE requests
func (registry *RouterRegistry) HandleDelete(name string, action Action) {
	registry.HandleMethod(http.MethodDelete, name, action)
}

// HandleOptions handles HTTP OPTIONS requests
func (registry *RouterRegistry) HandleOptions(name string, action Action) {
	registry.HandleMethod(http.MethodOptions, name, action)
}

// HandleHead handles HTTP HEAD requests
func (registry *RouterRegistry) HandleHead(name string, action Action) {
	registry.HandleMethod(http.MethodHead, name, action)
}

// Has checks if a method is set for a given handler name
func (registry *RouterRegistry) Has(method, name string) bool {
	la, ok := registry.handler[name]
	_, methodSet := la.method[method]
	return ok && methodSet
}

// HasAny checks if an any handler is set for a given name
func (registry *RouterRegistry) HasAny(name string) bool {
	la, ok := registry.handler[name]
	return ok && la.any != nil
}

// HasData checks if a data handler is set for a given name
func (registry *RouterRegistry) HasData(name string) bool {
	la, ok := registry.handler[name]
	return ok && la.data != nil
}

// MustRoute checks the result of a `Route` call
func MustRoute(handler *Handler, err error) *Handler {
	if err != nil {
		panic(err)
	}
	return handler
}

// MustRoute makes a checked Route call
func (registry *RouterRegistry) MustRoute(path, handler string) *Handler {
	return MustRoute(registry.Route(path, handler))
}

// Route assigns a route to a Handler
func (registry *RouterRegistry) Route(path, handler string) (*Handler, error) {
	var h = parseHandler(handler)
	var err error

	h.path, err = NewPath(path)
	if err != nil {
		return nil, err
	}

	if len(h.params) == 0 {
		h.params, h.catchall = parseParams(strings.Join(h.path.params, ", "))
	}

	registry.routes = append(registry.routes, h)
	return h, nil
}

// GetRoutes returns registered Routes
func (registry *RouterRegistry) GetRoutes() []*Handler {
	return registry.routes
}

// getHandler returns registered Routes
func (registry *RouterRegistry) getHandler() map[string]handlerAction {
	return registry.handler
}

// Alias for an existing router definition
func (registry *RouterRegistry) Alias(name, to string) {
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
		newHandler.params, newHandler.catchall = parseParams(tmp[1][:len(tmp[1])-1])
	}

	return newHandler
}

// list: foo, bar, x ?= "y", z = "a"
func parseParams(list string) (params map[string]*param, catchall bool) {
	// try to get enough space for the list
	params = make(map[string]*param, strings.Count(list, ","))

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

		case '*':
			catchall = true

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

	return params, catchall
}

// Reverse builds the path from a named route with params
func (registry *RouterRegistry) Reverse(name string, params map[string]string) (string, error) {
	if alias, ok := registry.alias[name]; ok {
		name = alias.handler
		if params == nil {
			params = make(map[string]string, len(alias.params))
		}
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
		if handler.handler != name {
			continue
		}
		var renderparams = make(map[string]string, len(handler.params)+len(params))
		var usedValues = make(map[string]struct{}, len(handler.params))

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
			usedValues[key] = struct{}{}
		}

		// add Reverse parameters
		for k, v := range params {
			if v != renderparams[k] {
				delete(usedValues, k)
			}

			renderparams[k] = v
		}

		// validate if all parameters have been used
		for key := range params {
			if _, ok := handler.params[key]; !ok {
				continue routeloop
			}
		}

		return handler.path.Render(renderparams, usedValues)

	}

catchallrouteloop:
	for _, handler := range registry.routes {
		if handler.handler != name || !handler.catchall {
			continue
		}
		var renderparams = make(map[string]string, len(handler.params)+len(params))
		var usedValues = make(map[string]struct{}, len(handler.params))

		// set handler default parameters
		for key, param := range handler.params {
			v, ok := params[key]

			// unset not-optional param
			if !param.optional && !ok {
				continue catchallrouteloop
			}

			// not-optional param set with wrong value
			if !param.optional && ok && param.value != "" && param.value != v {
				continue catchallrouteloop
			}
			renderparams[key] = param.value
			usedValues[key] = struct{}{}
		}

		// add Reverse parameters
		for k, v := range params {
			if v != renderparams[k] {
				delete(usedValues, k)
			}

			renderparams[k] = v
		}

		return handler.path.Render(renderparams, usedValues)
	}

	return "", fmt.Errorf("reverse for %q not found, parameters: %v", name, params)
}

// Match a request path
func (registry *RouterRegistry) match(path string) (handler handlerAction, params map[string]string) {
	for _, route := range registry.routes {
		if match := route.path.Match(path); match != nil {
			handler = registry.handler[route.handler]
			params = make(map[string]string)
			for k, param := range route.params {
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

// matchRequest matches a http Request (with query and path parameters)
func (registry *RouterRegistry) matchRequest(req *http.Request) (handlerAction, map[string]string, *Handler) {
	var path = req.URL.Path
	if req.URL.RawPath != "" {
		path = req.URL.RawPath
	}

	path = "/" + strings.TrimLeft(path, "/")

	var matchedHandlers matchedHandlers
	for _, handler := range registry.routes {
		if match := handler.path.Match(path); match != nil {
			controller := registry.handler[handler.handler]
			matchedHandler := &matchedHandler{
				handlerAction: controller,
				handler:       handler,
				match:         match,
			}
			matchedHandlers = append(matchedHandlers, matchedHandler)
		}
	}

	if any := matchedHandlers.getHandleAny(); any != nil && !matchedHandlers.hasMethod(req.Method) {
		return registry.makeHandler(req, *any)
	}

	for _, matched := range matchedHandlers {
		if matched == nil {
			continue
		}

		controller := matched.handlerAction
		if _, ok := controller.method[req.Method]; !ok && len(controller.method) > 0 {
			continue
		}

		controller, params, handler := registry.makeHandler(req, *matched)
		if handler == nil {
			continue
		}

		return controller, params, handler
	}
	return handlerAction{}, nil, nil
}

func (registry *RouterRegistry) makeHandler(req *http.Request, matched matchedHandler) (handlerAction, map[string]string, *Handler) {
	params := make(map[string]string)
	if len(matched.handler.params) > 0 {
		for k, param := range matched.handler.params {
			if !param.optional && param.value != "" {
				params[k] = param.value
			} else if v, ok := matched.match.Values[k]; ok {
				params[k] = v
			} else if val := req.URL.Query().Get(k); val != "" {
				params[k] = val
			} else if !param.optional && param.value == "" {
				return handlerAction{}, nil, nil
			} else {
				params[k] = param.value
			}
		}
	} else {
		params = matched.match.Values
	}
	return matched.handlerAction, params, matched.handler
}

// GetPath getter
func (handler *Handler) GetPath() string {
	return handler.path.path
}

// GetHandlerName getter
func (handler *Handler) GetHandlerName() string {
	return handler.handler
}

// Normalize enforces a normalization of passed parameters
func (handler *Handler) Normalize(params ...string) *Handler {
	if handler.path.normalize == nil {
		handler.path.normalize = make(map[string]struct{}, len(params))
	}
	for _, p := range params {
		handler.path.normalize[p] = struct{}{}
	}
	return handler
}
