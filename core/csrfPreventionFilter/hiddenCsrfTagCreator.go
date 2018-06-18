package csrfPreventionFilter

import (
	"flamingo.me/flamingo/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/framework/event"
)

type (
	hiddenCsrfTagCreator struct{}
)

// Notify is called on events
func (c *hiddenCsrfTagCreator) Notify(e event.Event) {
	switch e := e.(type) {
	case *pugjs.OnRenderHTMLBlockEvent:
		c.onRenderHTMLBlockEvent(e)
	}
}

func (c *hiddenCsrfTagCreator) onRenderHTMLBlockEvent(event *pugjs.OnRenderHTMLBlockEvent) {
	switch event.BlockName {
	case "form":
		event.Buffer.WriteString(`<input type="hidden" name="csrf_token" value="{{ csrftoken }}" />`)
	}
}
