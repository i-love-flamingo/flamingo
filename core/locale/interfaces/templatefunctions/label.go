package templatefunctions

import (
	"context"

	"flamingo.me/flamingo/v3/core/locale/application"
)

// Label is exported as a template function
type Label struct {
	translationService application.TranslationServiceInterface
}

// Inject dependencies
func (tf *Label) Inject(service application.TranslationServiceInterface) {
	tf.translationService = service
}

// Func template function factory
// todo fix
func (tf *Label) Func(context.Context) interface{} {

	// Usage:  __("key")
	// __("key","default")
	// __("key","Hello Mr {{.userName}}",{UserName: "Max"})
	// Force other than configured locale: __("switch_to_german","",{},"de-DE")
	return func(key string, params ...interface{}) string {
		localeCode := ""
		defaultLabel := key
		translationArguments := make(map[string]interface{})
		count := 1

		if len(params) > 0 {
			//if stringParam1, ok := params[0].(pugjs.String); ok {
			//	defaultLabel = string(stringParam1)
			//} else
			if stringParam1, ok := params[0].(string); ok {
				defaultLabel = string(stringParam1)
			}
		}
		if len(params) > 1 {
			//if mapParam2, ok := params[1].(*pugjs.Map); ok {
			//	for k, v := range mapParam2.Items {
			//		translationArguments[k.String()] = v
			//	}
			//	//use the special _count to set the count for T func
			//	if countArgument, ok := translationArguments["_count"]; ok {
			//		if countArgumentInt, ok := countArgument.(pugjs.Number); ok {
			//			count, _ = strconv.Atoi(countArgumentInt.String())
			//		}
			//	}
			//}
		}
		if len(params) > 2 {
			//if stringParam3, ok := params[2].(pugjs.String); ok {
			//	localeCode = string(stringParam3)
			//} else
			if stringParam3, ok := params[2].(string); ok {
				localeCode = stringParam3
			}
		}

		return tf.translationService.Translate(key, defaultLabel, localeCode, count, translationArguments)

	}
}
