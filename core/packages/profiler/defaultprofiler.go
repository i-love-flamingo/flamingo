package profiler

import (
	"bytes"
	"flamingo/core/flamingo/profiler"
	"flamingo/core/flamingo/router"
	"flamingo/core/flamingo/web"
	"fmt"
	"io/ioutil"
	"reflect"
	"time"
)

var profilestoreage map[string]*DefaultProfiler

func init() {
	profilestoreage = make(map[string]*DefaultProfiler)
}

type (
	// DefaultProfiler simply records whatever we pass into it
	DefaultProfiler struct {
		Context web.Context    `inject:""`
		Router  *router.Router `inject:""`

		key, msg string
		current  *DefaultProfiler
		childs   []*DefaultProfiler
		start    time.Time
		duration time.Duration
	}
)

// Profile something with a key and a message
func (p *DefaultProfiler) Profile(key, msg string) profiler.ProfileFinishFunc {
	if p.current == nil {
		p.current = p
	}

	var subprofiler = new(DefaultProfiler)
	subprofiler.key = key
	subprofiler.msg = msg
	subprofiler.start = time.Now()
	p.current.childs = append(p.current.childs, subprofiler)

	var parent = p.current
	p.current = subprofiler

	return func() {
		subprofiler.duration = time.Since(subprofiler.start)
		p.current = parent
	}
}

// String prints the child-logs
func (p *DefaultProfiler) String() (res string) {
	for _, c := range p.childs {
		res += c.log("")
	}
	return
}

func (p *DefaultProfiler) log(depth string) (res string) {
	res += fmt.Sprintf("%s%s: %s (%s)\n", depth, p.key, p.msg, p.duration)
	for _, c := range p.childs {
		res += c.log(depth + "    ")
	}
	return
}

// OnResponse injects the little helper into the response, and saves the profile in memory
func (p *DefaultProfiler) OnResponse(event *router.OnResponseEvent) {
	if reflect.TypeOf(event.Controller).Elem().Name() == reflect.TypeOf(ProfileController{}).Name() {
		return
	}

	if response, ok := event.Response.(*web.ContentResponse); ok {
		originalbody, _ := ioutil.ReadAll(response.Body)
		response.Body = bytes.NewBuffer(bytes.Replace(
			originalbody,
			[]byte("</body>"),
			[]byte("<div style='position:absolute;right:0;top:0;background-color:#ccc;border:solid 1px #888;'><a href='"+p.Router.URL("_profiler.view", "Profile", p.Context.ID()).String()+"'>Profile "+p.Context.ID()+"</a></div>\n</body>"),
			1,
		))
		profilestoreage[p.Context.ID()] = p
	}
}

// Events returns a list of subscribed events
func (p *DefaultProfiler) Events() []interface{} {
	return []interface{}{
		router.RESPONSE,
	}
}
