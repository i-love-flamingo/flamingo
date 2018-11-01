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

		if r.Values[ctxKey] == nil {
			r.Values[ctxKey] = make(map[string]interface{})
		}

		r.Values[ctxKey].(map[string]interface{})[key] = val

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

		if r.Values[ctxKey] == nil {
			return nil
		}

		return r.Values[ctxKey].(map[string]interface{})
	}
}
