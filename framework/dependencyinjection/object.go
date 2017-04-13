package dependencyinjection

import "reflect"

type (
	// Object represents a object which we work with
	// essentially it wraps value and tags, and some internal flags
	Object struct {
		Value interface{} // Value points to the original value
		Tags  []string    // Tags hold a list off assigned tags

		complete       bool          // complete will be set to true when the object has been properly resolved
		autocreated    bool          // autocreated signals objects which have been created, so they are not used for interface injection
		compilerpassed bool          // compilerpassed makes sure compiler pass happens only once
		reflectType    reflect.Type  // reflectType is a cache for reflect.TypeOf(Value)
		reflectValue   reflect.Value // reflectValue is a cache for reflect.ValueOf(Value)
		wrapFunc       reflect.Value // wrapFunc is generated for functions types, and wrap their argument into a new resolve call
		parameters     []string
	}
)

// AddParameters adds parameters to an object (factory only)
func (o *Object) AddParameters(params ...string) *Object {
	o.parameters = append(o.parameters, params...)
	return o
}
