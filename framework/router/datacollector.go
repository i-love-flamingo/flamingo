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
		handler *Handler
	}
)

// Collect data
func (dc *DataCollector) Collect(ctx web.Context) string {
	data, ok := ctx.Value("Handler").(handlerdata)
	if !ok {
		return "no Handler data"
	}
	var params []string
	for k, v := range data.match {
		ps := k
		ps += `="` + v + `"`

		params = append(params, ps)
	}
	if data.handler != nil && data.handler.path != nil {
		return fmt.Sprintf("Router: %s: %s(%s)", data.handler.path.path, data.handler.handler, strings.Join(params, ", "))
	}
	return "-"
}
