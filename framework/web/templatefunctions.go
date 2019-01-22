package web

import (
	"context"
)

type (
	partialDataContextKey string

	// SetPartialDataFunc allows to set partial data
	SetPartialDataFunc struct{}

	// GetPartialDataFunc allows to get partial data
	GetPartialDataFunc struct{}
)

const ctxKey partialDataContextKey = "partialData"

// Func getter to bind the context
func (*SetPartialDataFunc) Func(c context.Context) interface{} {
	return func(key string, val interface{}) interface{} {
		r := RequestFromContext(c)
		if r == nil {
			return nil
		}

		data, ok := r.Values.Load(ctxKey)
		if !ok || data == nil {
			data = make(map[string]interface{})
		}

		data.(map[string]interface{})[key] = val

		r.Values.Store(ctxKey, data)

		return nil
	}
}

// Func getter to bind the context
func (*GetPartialDataFunc) Func(c context.Context) interface{} {
	return func() map[string]interface{} {
		r := RequestFromContext(c)
		if r == nil {
			return nil
		}

		data, ok := r.Values.Load(ctxKey)
		if !ok || data == nil {
			return nil
		}

		return data.(map[string]interface{})
	}
}
