package profiler

import (
	"bytes"
	"io/ioutil"
	"time"

	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

// eventSubscriber for the profiler
type eventSubscriber struct {
	Router *router.Router `inject:""`
}

// Notify on events
func (e *eventSubscriber) Notify(ev event.Event) {
	switch ev := ev.(type) {
	case *router.OnResponseEvent:
		e.onResponse(ev)
	}
}

// onResponse injects the little helper into the response, and saves the profile in memory
func (e *eventSubscriber) onResponse(event *router.OnResponseEvent) {
	// ensure we are not profiling ourself
	if _, ok := event.Controller.(*profileController); ok {
		return
	}

	context := event.Request.Context().Value(web.CONTEXT).(web.Context)
	p := context.Profiler().(*defaultProfiler)

	event.ResponseWriter.Header().Set("X-Correlation-ID", context.ID())

	if context.Session() != nil {
		if _, ok := event.Response.(*web.RedirectResponse); ok {
			context.Session().Values["context.id"] = context.ID()
		} else {
			delete(context.Session().Values, "context.id")
		}
	}

	p.Collect(context)

	if response, ok := event.Response.(*web.ContentResponse); ok {
		p.Duration = time.Since(p.Start)
		originalbody, _ := ioutil.ReadAll(response.Body)
		response.Body = bytes.NewBuffer(bytes.Replace(
			bytes.Replace(
				originalbody,
				[]byte("</head>"),
				[]byte(`
<script type='text/javascript'>
var __start = Date.now(), __open = XMLHttpRequest.prototype.open;

XMLHttpRequest.prototype.open = function(a, b) {
	r = __open.call(this, a, b);
	this.setRequestHeader("X-Request-Id", "`+context.ID()+`");
	return r;
}

function __profileStatic(key, message, duration) {
	var r = new XMLHttpRequest();
	r.open("POST", "`+e.Router.URL("_profiler.view", router.P{"profile": context.ID()}).String()+`");
	r.setRequestHeader("Content-Type", "application/json");
	r.send(JSON.stringify({"key": key, "message": message, "duration": duration.toString()}));
}

function __profile(key, message) {
	start = Date.now();
	return function(){
		__profileStatic(key, message, Date.now() - start);
	}
}

window.addEventListener("error", function (e) {
    __profileStatic("browser.error", e.error ? e.error.stack : e.message, Date.now() - __start);
});

window.addEventListener("load", function load(e) {
    window.removeEventListener("load", load);
    var t = window.performance.timing,
          dcl = t.domContentLoadedEventStart - t.domLoading,
          complete = t.domComplete - t.domLoading;
    __profileStatic("browser", "DOMContentLoaded", dcl);
    __profileStatic("browser", "Load", complete);
});

</script>
</head>`),
				1,
			),
			[]byte("</body>"),
			[]byte(`<div style='position:absolute;right:0;bottom:0;background-color:#ccc;border:1px solid #888;z-index:10000;padding:3px 2px 0;'>
	<a href='`+p.Router.URL("_profiler.view", router.P{"profile": context.ID()}).String()+`' style='text-decoration:none;'>⚙</a>
</div>
</body>`),
			1,
		),
		)
	}

	if existing, ok := profilestorage.Load(context.ID()); ok {
		p.Childs = append(existing.Childs, p.Childs...)
	}
	if existing, ok := profilestorage.Load(context.Request().Header.Get("X-Correlation-Id")); ok {
		existing.ProfileExternal(context.Request().RequestURI, context.ID(), p.Duration)
	}
	profilestorage.Store(context.ID(), p)
}
