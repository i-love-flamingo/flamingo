package domain

import (
	"bytes"
	"html/template"
	"log"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

type (
	TranslationService struct {
		DefaultLocaleCode string `inject:"config:locale.locale"`
		TranslationFile   string `inject:"config:locale.translationFile"`
	}
)

var i18bundle *bundle.Bundle
var filesLoaded bool

func init() {
	i18bundle = bundle.New()
	filesLoaded = false
}

func (ts *TranslationService) Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string {
	if count < 1 {
		count = 1
	}
	if translationArguments == nil {
		translationArguments = make(map[string]interface{})
	}
	if !filesLoaded {
		i18bundle.LoadTranslationFile(ts.TranslationFile)
		filesLoaded = true
	}
	label := ""
	//Use default configured localeCode if nothing is given explicitly
	if localeCode == "" {
		localeCode = ts.DefaultLocaleCode
	}
	T, err := i18bundle.Tfunc(localeCode)
	if err != nil {
		log.Printf("Error - locale.translationservice %v", err)
	} else {
		log.Printf("called with key %v  default: %v  Code: %v translationArguments: %#v Count %v", key, defaultLabel, localeCode, translationArguments, count)
		label = T(key, count, translationArguments)
	}

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
