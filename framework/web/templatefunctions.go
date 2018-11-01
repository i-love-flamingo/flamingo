package web

import (
	"context"
)

type partialDataContextKey string

const ctxKey partialDataContextKey = "partialData"

type SetPartialDataFunc struct{}

func (*SetPartialDataFunc) Func(c context.Context) interface{} {
	return func(key string, val interface{}) interface{} {
		r, ok := FromContext(c)
		if !ok {
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

type GetPartialDataFunc struct{}

func (*GetPartialDataFunc) Func(c context.Context) interface{} {
	return func() map[string]interface{} {
		r, ok := FromContext(c)
		if !ok {
			return nil
		}

		data, ok := r.Values.Load(ctxKey)
		if !ok || data == nil {
			return nil
		}

		return data.(map[string]interface{})
	}
}
