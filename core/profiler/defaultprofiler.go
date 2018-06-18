package profiler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"
	"time"

	"flamingo.me/flamingo/framework/profiler"
	"flamingo.me/flamingo/framework/profiler/collector"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	profilemap struct {
		sync.Map
	}
)

func (pm *profilemap) Load(key string) (*defaultProfiler, bool) {
	p, ok := pm.Map.Load(key)
	if !ok {
		return nil, false
	}
	pp, ok := p.(*defaultProfiler)
	return pp, ok
}

func (pm *profilemap) Store(key string, val *defaultProfiler) {
	pm.Map.Store(key, val)
}

var (
	profilestorage = new(profilemap)
)

type (
	// defaultProfiler simply records whatever we pass into it
	defaultProfiler struct {
		Router    *router.Router            `inject:""`
		Collector []collector.DataCollector `inject:""`

		Fnc, Msg         string
		File             string
		Startpos, Endpos int
		current          *defaultProfiler
		Childs           []*defaultProfiler
		Start            time.Time
		Duration         time.Duration
		Depth            int
		Link             string
		Data             []string
	}
)

// Profile something with a Fnc and a message
func (p *defaultProfiler) Profile(key, msg string) profiler.ProfileFinishFunc {
	if p.current == nil {
		p.current = p
		p.Start = time.Now()
	}

	pc, _, _, _ := runtime.Caller(2)
	fnc := runtime.FuncForPC(pc)
	file, line := fnc.FileLine(pc)

	subprofiler := &defaultProfiler{
		Fnc:      fnc.Name(),
		File:     file,
		Startpos: line,
		Msg:      key + ": " + msg,
		Start:    time.Now(),
		Depth:    p.current.Depth + 1,
	}
	p.current.Childs = append(p.current.Childs, subprofiler)

	parent := p.current
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
func (p *defaultProfiler) Collect(ctx web.Context) {
	for _, c := range p.Collector {
		p.Data = append(p.Data, c.Collect(ctx))
	}
}

// ProfileOffline is called for profiling events not directly in flamingo, such as the browser
func (p *defaultProfiler) ProfileOffline(key, msg string, duration time.Duration) {
	p.current.Childs = append(p.current.Childs, &defaultProfiler{
		Msg:      key + ": " + msg,
		Duration: duration,
	})
}

// ProfileExternal connects another profile to the current, e.g. for Ajax requests
func (p *defaultProfiler) ProfileExternal(key, id string, duration time.Duration) {
	p.current.Childs = append(p.current.Childs, &defaultProfiler{
		Msg:      key,
		Duration: duration,
		Link:     id,
	})
}

// Filehint gives the source file's content hint
func (p *defaultProfiler) Filehint() string {
	c, err := ioutil.ReadFile(p.File)
	if err != nil {
		return err.Error()
	}
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
