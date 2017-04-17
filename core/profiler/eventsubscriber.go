package profiler

import (
	"bytes"
	"flamingo/framework/event"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"io/ioutil"
	"reflect"
	"time"
)

type EventSubscriber struct {
	Router *router.Router `inject:""`
}

func (e *EventSubscriber) Notify(ev event.Event) {
	switch ev := ev.(type) {
	case *router.OnResponseEvent:
		e.OnResponse(ev)
	}
}

// OnResponse injects the little helper into the response, and saves the profile in memory
func (e *EventSubscriber) OnResponse(event *router.OnResponseEvent) {
	if reflect.TypeOf(event.Controller).Kind() != reflect.Ptr {
		return
	}
	if reflect.TypeOf(event.Controller).Elem().Name() == reflect.TypeOf(ProfileController{}).Name() {
		return
	}

	context := event.Request.Context().Value(web.CONTEXT).(web.Context)
	p := context.Profiler().(*DefaultProfiler)

	if response, ok := event.Response.(*web.ContentResponse); ok {
		p.Duration = time.Since(p.Start)
		originalbody, _ := ioutil.ReadAll(response.Body)
		response.Body = bytes.NewBuffer(bytes.Replace(
			originalbody,
			[]byte("</body>"),
			[]byte(`
<script type='text/javascript'>
var __start = 0;

function __profileStatic(key, message, duration) {
	var r = new XMLHttpRequest();
	r.open("POST", "`+e.Router.URL("_profiler.view", "profile", context.ID()).String()+`");
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
    __profileStatic("browser.error", e.error.stack, Date.now() - __start);
});

window.addEventListener("DOMContentLoaded", function(e){
	__start = Date.now() - e.timeStamp;
	__profileStatic("browser", "DOMContentLoaded", e.timeStamp);
});

window.addEventListener("load", function load(e) {
    window.removeEventListener("load", load);
    __profileStatic("browser", "Load", e.timeStamp);
});

</script>
<div style='position:absolute;right:0;top:0;background-color:#ccc;border:solid 1px #888;'>
	<a href='`+p.Router.URL("_profiler.view", "profile", context.ID()).String()+`'>`+p.Duration.String()+`: `+context.ID()+`</a>
</div>
</body>`),
			1,
		))
		profilestorage[context.ID()] = p
	}
}
