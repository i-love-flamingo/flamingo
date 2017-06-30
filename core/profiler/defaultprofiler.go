package profiler

import (
	"bytes"
	"flamingo/framework/profiler"
	"flamingo/framework/profiler/collector"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
	"sync"
)

var profilestorage map[string]*DefaultProfiler
var profilelock *sync.Mutex

func init() {
	profilestorage = make(map[string]*DefaultProfiler)
	profilelock = new(sync.Mutex)
}

type (
	// DefaultProfiler simply records whatever we pass into it
	DefaultProfiler struct {
		Router    *router.Router            `inject:""`
		Collector []collector.DataCollector `inject:""`

		Fnc, Msg         string
		File             string
		Startpos, Endpos int
		current          *DefaultProfiler
		Childs           []*DefaultProfiler
		Start            time.Time
		Duration         time.Duration
		Depth            int
		Link             string
		Data             []string
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

// Collect data
func (p *DefaultProfiler) Collect(ctx web.Context) {
	for _, c := range p.Collector {
		p.Data = append(p.Data, c.Collect(ctx))
	}
}

// ProfileOffline is called for profiling events not directly in flamingo, such as the browser
func (p *DefaultProfiler) ProfileOffline(key, msg string, duration time.Duration) {
	var subprofiler = new(DefaultProfiler)
	subprofiler.Msg = key + ": " + msg
	subprofiler.Duration = duration
	p.current.Childs = append(p.current.Childs, subprofiler)
}

// ProfileExternal connects another profile to the current, e.g. for Ajax requests
func (p *DefaultProfiler) ProfileExternal(key, id string, duration time.Duration) {
	var subprofiler = new(DefaultProfiler)
	subprofiler.Msg = key
	subprofiler.Duration = duration
	subprofiler.Link = id
	p.current.Childs = append(p.current.Childs, subprofiler)
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
