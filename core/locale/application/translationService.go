package application

import (
	"bytes"
	"fmt"
	"text/template"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

type (
	// TranslationServiceInterface defines the translation service
	TranslationServiceInterface interface {
		Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string
	}

	// TranslationService is the default TranslationServiceInterface implementation
	TranslationService struct {
		defaultLocaleCode string
		translationFile   string
		translationFiles  config.Slice
		logger            flamingo.Logger
		devmode           bool
	}
)

// check if translationService implements its interface
var _ TranslationServiceInterface = (*TranslationService)(nil)

var i18bundle *bundle.Bundle
var filesLoaded bool

func init() {
	i18bundle = bundle.New()
	filesLoaded = false
}

// Inject dependencies
func (ts *TranslationService) Inject(
	logger flamingo.Logger,
	config *struct {
		DefaultLocaleCode string       `inject:"config:locale.locale"`
		DevMode           bool         `inject:"config:debug.mode"`
		TranslationFile   string       `inject:"config:locale.translationFile,optional"`
		TranslationFiles  config.Slice `inject:"config:locale.translationFiles,optional"`
	},
) {
	ts.logger = logger
	ts.defaultLocaleCode = config.DefaultLocaleCode
	ts.translationFile = config.TranslationFile
	ts.translationFiles = config.TranslationFiles
	ts.devmode = config.DevMode
}

// Translate returns the result for translating a key, with a default label for a given locale code
func (ts *TranslationService) Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string {
	if count < 1 {
		count = 1
	}
	if translationArguments == nil {
		translationArguments = make(map[string]interface{})
	}
	if !filesLoaded || ts.devmode {
		ts.loadFiles()
		filesLoaded = true
	}
	label := ""
	//Use default configured localeCode if nothing is given explicitly
	if localeCode == "" {
		localeCode = ts.defaultLocaleCode
	}
	T, err := i18bundle.Tfunc(localeCode)
	if err != nil {
		ts.logger.Info("Error - locale.translationservice", err)
		label = defaultLabel
	} else {
		//ts.Logger.Debug("called with key %v  default: %v  localeCode: %v translationArguments: %#v Count %v", key, defaultLabel, localeCode, translationArguments, count)
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
func (ts *TranslationService) loadFiles() {
	if ts.translationFile != "" {
		err := i18bundle.LoadTranslationFile(ts.translationFile)
		if err != nil {
			ts.logger.Warn(fmt.Sprintf("Load translationfile failed: %s", err))
		}
	}
	if len(ts.translationFiles) > 0 {
		for _, file := range ts.translationFiles {
			if fileName, ok := file.(string); ok {
				err := i18bundle.LoadTranslationFile(fileName)
				if err != nil {
					ts.logger.Warn(fmt.Sprintf("Load translationfile failed: %s", err))
				}
			}
		}
	}
}
