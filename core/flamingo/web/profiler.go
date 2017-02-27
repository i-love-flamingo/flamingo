package web

import (
	"log"
	"time"
)

type (
	// ProfileFinishFunc is used to finish a profiling-block
	ProfileFinishFunc func()

	// Profiler is used to measure the time certain things need
	Profiler interface {
		Profile(string, string) ProfileFinishFunc
		Init(Context)
		Log()
	}

	// DefaultProfiler simply records whatever we pass into it
	DefaultProfiler struct {
		ctx      Context
		key, msg string
		current  *DefaultProfiler
		childs   []*DefaultProfiler
		start    time.Time
		duration time.Duration
	}
)

// Profile something with a key and a message
func (p *DefaultProfiler) Profile(key, msg string) ProfileFinishFunc {
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

// Log prints the child-logs
func (p *DefaultProfiler) Log() {
	for _, c := range p.childs {
		c.log("")
	}
}

func (p *DefaultProfiler) log(depth string) {
	log.Printf("%s%s: %s (%s)\n", depth, p.key, p.msg, p.duration)
	for _, c := range p.childs {
		c.log(depth + "    ")
	}
}

// Init the profiler with current and context
func (p *DefaultProfiler) Init(ctx Context) {
	p.current = p
	p.ctx = ctx
}
