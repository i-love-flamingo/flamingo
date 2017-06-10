package router

import (
	"flamingo/framework/web"
	"fmt"
	"strings"
)

type (
	// DataCollector for core/profiler
	DataCollector struct{}

	handlerdata struct {
		match   map[string]string
		handler *handler
	}
)

// Collect data
func (dc *DataCollector) Collect(ctx web.Context) string {
	data, ok := ctx.Value("handler").(handlerdata)
	if !ok {
		return "no handler data"
	}
	var params []string
	for k, v := range data.match {
		ps := k
		ps += `="` + v + `"`

		params = append(params, ps)
	}
	return fmt.Sprintf("Router: %s: %s(%s)", data.handler.path.path, data.handler.handler, strings.Join(params, ", "))
}
