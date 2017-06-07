package router

import (
	"flamingo/framework/web"
	"fmt"
	"strings"
)

// DataCollector for core/profiler
type DataCollector struct{}

// Collect data
func (dc *DataCollector) Collect(ctx web.Context) string {
	handler := ctx.Value("handler").(*handler)
	var params []string
	for n, p := range handler.params {
		ps := n + " "
		if p.optional {
			ps += "?"
		}
		if p.value != "" {
			ps += `="` + p.value + `"`
		}
		params = append(params, ps)
	}
	return fmt.Sprintf("Router: %s: %s(%s)", handler.path.path, handler.handler, strings.Join(params, ", "))
}
