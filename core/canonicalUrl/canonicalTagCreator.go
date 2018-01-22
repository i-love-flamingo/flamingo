package canonicalUrl

import (
	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
	"go.aoe.com/flamingo/framework/event"
)

type (
	canonicalTagCreator struct{}
)

func (c *canonicalTagCreator) Notify(e event.Event) {
	switch e := e.(type) {
	case *pugjs.OnRenderHTMLBlockEvent:
		c.onRenderHTMLBlockEvent(e)
	}
}

func (c *canonicalTagCreator) onRenderHTMLBlockEvent(event *pugjs.OnRenderHTMLBlockEvent) {
	switch event.BlockName {
	case "head":
		event.Buffer.WriteString(`<link rel="canonical" href="{{canonicalUrl}}" />`)
	}
}
