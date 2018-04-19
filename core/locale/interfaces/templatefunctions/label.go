package templatefunctions

import (
	"strconv"

	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
)

type (
	// Label is exported as a template function
	Label struct {
		TranslationService application.TranslationServiceInterface `inject:""`
	}
)

// Name alias for use in template
func (tf Label) Name() string {
	return "__"
}

func (tf Label) Func() interface{} {

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
			if stringParam1, ok := params[0].(pugjs.String); ok {
				defaultLabel = string(stringParam1)
			} else if stringParam1, ok := params[0].(string); ok {
				defaultLabel = string(stringParam1)
			}
		}
		if len(params) > 1 {
			if mapParam2, ok := params[1].(*pugjs.Map); ok {
				for k, v := range mapParam2.Items {
					translationArguments[k.String()] = v
				}
				//use the special _count to set the count for T func
				if countArgument, ok := translationArguments["_count"]; ok {
					if countArgumentInt, ok := countArgument.(pugjs.Number); ok {
						count, _ = strconv.Atoi(countArgumentInt.String())
					}
				}
			}
		}
		if len(params) > 2 {
			if stringParam3, ok := params[2].(pugjs.String); ok {
				localeCode = string(stringParam3)
			} else if stringParam3, ok := params[2].(string); ok {
				localeCode = stringParam3
			}
		}

		return tf.TranslationService.Translate(key, defaultLabel, localeCode, count, translationArguments)

	}
}
