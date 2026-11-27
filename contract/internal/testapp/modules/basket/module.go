// Package basket is a test module with routes that read and write the session.
package basket

import (
	"context"
	"slices"
	"strings"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/web"
)

const keyPrefix = "basket."

type (
	// Module of the basket routes.
	Module struct{}

	// Holder lets a test control when the route "basket.hold" writes.
	// A test binds it in the root area. Without a binding the route does not wait.
	Holder struct {
		Loaded  chan string              // gets the item when the request has loaded its session
		Release map[string]chan struct{} // the request of an item writes after the channel of the item is closed
	}

	routes struct {
		responder *web.Responder
		holder    *Holder
	}
)

// Configure DI.
func (*Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

// Inject dependencies.
func (r *routes) Inject(responder *web.Responder, cfg *struct {
	Holder *Holder `inject:",optional"`
}) {
	r.responder = responder

	if cfg != nil {
		r.holder = cfg.Holder
	}
}

// Routes of the module.
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/basket/add/:item", "basket.add")
	registry.HandleGet("basket.add", func(_ context.Context, req *web.Request) web.Result {
		return r.add(req)
	})

	registry.MustRoute("/basket", "basket.show")
	registry.HandleGet("basket.show", func(_ context.Context, req *web.Request) web.Result {
		return r.responder.Data(keys(req.Session()))
	})

	registry.MustRoute("/basket/count", "basket.count")
	registry.HandleGet("basket.count", func(_ context.Context, req *web.Request) web.Result {
		count := 0

		// Load each item, do not only count the keys: the save mode test needs a request that reads values.
		for _, key := range keys(req.Session()) {
			if _, ok := req.Session().Load(key); ok {
				count++
			}
		}

		return r.responder.Data(map[string]int{"count": count})
	})

	// The route cancels its context first, to show that the detached context stays alive and still has the
	// request and the session of the caller.
	registry.MustRoute("/basket/detached", "basket.detached")
	registry.HandleGet("basket.detached", func(ctx context.Context, req *web.Request) web.Result {
		canceled, cancel := context.WithCancel(ctx)
		cancel()

		var state map[string]bool

		web.RunWithDetachedContext(canceled, func(detached context.Context) {
			state = map[string]bool{
				"request": web.RequestFromContext(detached) == req,
				"session": web.SessionFromContext(detached) == req.Session(),
				"live":    detached.Err() == nil,
			}
		})

		return r.responder.Data(state)
	})

	registry.MustRoute("/basket/hold/:item", "basket.hold")
	registry.HandleGet("basket.hold", func(ctx context.Context, req *web.Request) web.Result {
		r.hold(ctx, req.Params["item"])

		return r.add(req)
	})
}

// add stores the item in the session.
func (r *routes) add(req *web.Request) web.Result {
	item := req.Params["item"]

	req.Session().Store(keyPrefix+item, "1")

	return r.responder.Data(map[string]string{"added": item})
}

// hold tells the holder that the session is loaded, and waits for the release of the item.
func (r *routes) hold(ctx context.Context, item string) {
	if r.holder == nil || r.holder.Loaded == nil {
		return
	}

	select {
	case r.holder.Loaded <- item:
	case <-ctx.Done():
		return
	}

	select {
	case <-r.holder.Release[item]:
	case <-ctx.Done():
	}
}

// keys gives the sorted keys of the session that start with the basket prefix.
func keys(session *web.Session) []string {
	result := []string{}

	for _, key := range session.Keys() {
		if name, ok := key.(string); ok && strings.HasPrefix(name, keyPrefix) {
			result = append(result, name)
		}
	}

	slices.Sort(result)

	return result
}
