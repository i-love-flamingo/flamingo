package templatefunctions

import (
	"text/template"

	"bytes"

	"strconv"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
)

type (
	// Label is exported as a template function
	Label struct {
		LocaleCode      string `inject:"config:locale.locale"`
		TranslationFile string `inject:"config:locale.translationFile"`
	}
)

// Name alias for use in template
func (tf Label) Name() string {
	return "__"
}

func (tf Label) Func() interface{} {

	i18bundle := bundle.New()
	i18bundle.LoadTranslationFile(tf.TranslationFile)

	// Usage:  __("key")
	// __("key","default")
	// __("key","Hello Mr {{.userName}}",{UserName: "Max"})
	// Force other than configured locale: __("switch_to_german","",{},"de-DE")
	return func(key string, params ...interface{}) string {
		localeCode := tf.LocaleCode
		defaultLabel := key
		translationArguments := make(map[string]interface{})
		count := 1
		if len(params) > 0 {
			if stringParam1, ok := params[0].(string); ok && stringParam1 != "" {
				defaultLabel = stringParam1
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
			if stringParam3, ok := params[2].(string); ok {
				localeCode = stringParam3
			}
		}
		T, _ := i18bundle.Tfunc(localeCode)

		//log.Printf("called with key %v param: %#v  default: %v  Code: %v translationArguments: %#v Count %v", key, params, defaultLabel, localeCode, translationArguments, count)
		label := T(key, count, translationArguments)

		//Fallback if label was not translated
		if label == key && defaultLabel != "" {
			tmpl, err := template.New(key).Parse(defaultLabel)
			if err != nil {
				return defaultLabel
			}
			var doc bytes.Buffer
			err = tmpl.Execute(&doc, translationArguments)
			if err != nil {
				return defaultLabel
			}
			return doc.String()
		}
		return label
	}
}
