package flamingo

import (
	"context"

	"flamingo.me/dingo"
)

// TemplateFunc defines an interface for a custom function to be used in gotemplates/pug templates
type TemplateFunc interface {
	Func(ctx context.Context) interface{}
}

// BindTemplateFunc makes sure a template function is correctly bound via dingo
func BindTemplateFunc(injector *dingo.Injector, name string, fnc TemplateFunc) {
	injector.BindMap(new(TemplateFunc), name).To(fnc)
}
