package profiler

import (
	"bytes"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"time"
)

var profilestorage map[string]*DefaultProfiler

func init() {
	profilestorage = make(map[string]*DefaultProfiler)
}

type (
	// DefaultProfiler simply records whatever we pass into it
	DefaultProfiler struct {
		Context web.Context    `inject:""`
		Router  *router.Router `inject:""`

		Fnc, Msg         string
		File             string
		Startpos, Endpos int
		current          *DefaultProfiler
		Childs           []*DefaultProfiler
		Start            time.Time
		Duration         time.Duration
		Depth            int
	}
)

// Profile something with a Fnc and a message
func (p *DefaultProfiler) Profile(key, msg string) profiler.ProfileFinishFunc {
	if p.current == nil {
		p.current = p
		p.Start = time.Now()
	}

	var subprofiler = new(DefaultProfiler)

	pc, _, _, _ := runtime.Caller(2)
	fnc := runtime.FuncForPC(pc)
	file, line := fnc.FileLine(pc)
	subprofiler.Fnc = fnc.Name()
	subprofiler.File = file
	subprofiler.Startpos = line

	subprofiler.Msg = key + ": " + msg
	subprofiler.Start = time.Now()
	subprofiler.Depth = p.current.Depth + 1
	p.current.Childs = append(p.current.Childs, subprofiler)

	var parent = p.current
	p.current = subprofiler

	return func() {
		subprofiler.Duration = time.Since(subprofiler.Start)
		pc, _, _, _ := runtime.Caller(1)
		fnc := runtime.FuncForPC(pc)
		_, line := fnc.FileLine(pc)
		subprofiler.Endpos = line
		p.current = parent
	}
}

func (p *DefaultProfiler) ProfileOffline(key, msg string, duration time.Duration) {
	var subprofiler = new(DefaultProfiler)
	subprofiler.Msg = key + ": " + msg
	subprofiler.Duration = duration
	p.current.Childs = append(p.current.Childs, subprofiler)
}

// OnResponse injects the little helper into the response, and saves the profile in memory
func (p *DefaultProfiler) OnResponse(event *router.OnResponseEvent) {
	if reflect.TypeOf(event.Controller).Kind() != reflect.Ptr {
		return
	}
	if reflect.TypeOf(event.Controller).Elem().Name() == reflect.TypeOf(ProfileController{}).Name() {
		return
	}

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
	r.open("POST", "`+p.Router.URL("_profiler.view", "profile", p.Context.ID()).String()+`");
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
	<a href='`+p.Router.URL("_profiler.view", "profile", p.Context.ID()).String()+`'>`+p.Duration.String()+`: `+p.Context.ID()+`</a>
</div>
</body>`),
			1,
		))
		profilestorage[p.Context.ID()] = p
	}
}

// Events returns a list of subscribed events
func (p *DefaultProfiler) Events() []interface{} {
	return []interface{}{
		router.RESPONSE,
	}
}

// Filehint gives the source file's content hint
func (p *DefaultProfiler) Filehint() string {
	c, err := ioutil.ReadFile(p.File)
	if err != nil {
		return err.Error()
	}
	//log.Println(bytes.Split(c, []byte("\n")))
	//log.Println(p.Startpos, p.Endpos, len(bytes.Split(c, []byte("\n"))))
	//os.Exit(1)
	if p.Endpos < p.Startpos {
		p.Endpos = p.Startpos + 1
	}
	if len(bytes.Split(c, []byte("\n"))) < p.Endpos || len(bytes.Split(c, []byte("\n"))) < p.Startpos {
		return "--"
	}
	lines := bytes.Split(c, []byte("\n"))[p.Startpos-1 : p.Endpos]
	res := ""
	for i, l := range lines {
		res += fmt.Sprintf("%03d: %s\n", i+p.Startpos, string(l))
	}
	return res
}
