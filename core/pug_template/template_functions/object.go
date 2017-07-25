package template_functions

import "flamingo/core/pug_template/pugast"

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

func (o Object) Assign(target *pugast.Map, sources ...*pugast.Map) pugast.Object {
	for _, source := range sources {
		if source != nil {
			for k, v := range source.Items {
				target.Items[k] = v
			}
		}
	}

	return target
}
