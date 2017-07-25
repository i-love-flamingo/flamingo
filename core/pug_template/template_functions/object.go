package template_functions

import "flamingo/core/pug_template/pugjs"

type (
	ObjectLib struct{}
	Object    struct{}
)

// Name alias for use in template
func (ol ObjectLib) Name() string {
	return "Object"
}

// Func as implementation of debug method
func (ol ObjectLib) Func() interface{} {
	return func() Object {
		return Object{}
	}
}

func (o Object) NoConvert() {}

func (o Object) Assign(target *pugjs.Map, sources ...*pugjs.Map) pugjs.Object {
	for _, source := range sources {
		if source != nil {
			for k, v := range source.Items {
				target.Items[k] = v
			}
		}
	}

	return target
}
