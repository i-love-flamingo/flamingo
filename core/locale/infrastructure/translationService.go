package infrastructure

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"flamingo.me/flamingo/v3/core/locale/domain"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

type (
	// TranslationService is the default TranslationService implementation
	TranslationService struct {
		translationFile  string
		translationFiles config.Slice
		logger           flamingo.Logger
		devmode          bool
		filesLoaded      bool
		i18bundle        *bundle.Bundle
	}
)

// check if translationService implements its interface
var _ domain.TranslationService = (*TranslationService)(nil)

// Inject dependencies
func (ts *TranslationService) Inject(
	logger flamingo.Logger,
	config *struct {
		DevMode          bool         `inject:"config:debug.mode"`
		TranslationFile  string       `inject:"config:locale.translationFile,optional"`
		TranslationFiles config.Slice `inject:"config:locale.translationFiles,optional"`
	},
) {
	ts.logger = logger.WithField(flamingo.LogKeyModule, "locale").WithField("category", "locale.translationService")
	if config != nil {
		ts.translationFile = config.TranslationFile
		ts.translationFiles = config.TranslationFiles
		ts.devmode = config.DevMode
	}
}

// TranslateLabel returns the result for translating a Label
func (ts *TranslationService) TranslateLabel(label domain.Label) string {
	ts.initAndLoad()
	translatedString, err := ts.translateWithLib(label.GetLocaleCode(), label.GetKey(), label.GetCount(), label.GetTranslationArguments())

	//while there is an error check fallBacks
	for _, fallbackLocale := range label.GetFallbacklocaleCodes() {
		if err != nil {
			translatedString, err = ts.translateWithLib(fallbackLocale, label.GetKey(), label.GetCount(), label.GetTranslationArguments())
		}
	}
	if err != nil {
		//default to key (=untranslated) if still an error
		translatedString = label.GetKey()
	}
	//Fallback if label was not translated
	if translatedString == label.GetKey() && label.GetDefaultLabel() != "" {
		return ts.parseDefaultLabel(label.GetDefaultLabel(), label.GetKey(), label.GetTranslationArguments())
	}
	return translatedString
}

// Translate returns the result for translating a key, with a default label for a given locale code
func (ts *TranslationService) Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string {
	ts.initAndLoad()
	label, err := ts.translateWithLib(localeCode, key, count, translationArguments)

	if err != nil {
		//default to key (=untranslated) on error
		label = key
	}

	//Fallback if label was not translated
	if label == key && defaultLabel != "" {
		return ts.parseDefaultLabel(defaultLabel, key, translationArguments)
	}
	return label

}

// AllTranslationTags returns all keys for a given locale code
func (ts *TranslationService) AllTranslationTags(localeCode string) []string {
	ts.initAndLoad()
	return ts.i18bundle.LanguageTranslationIDs(localeCode)
}

func (ts *TranslationService) parseDefaultLabel(defaultLabel string, key string, translationArguments map[string]interface{}) string {
	if translationArguments == nil {
		translationArguments = make(map[string]interface{})
	}
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

func (ts *TranslationService) translateWithLib(localeCode string, key string, count int, translationArguments map[string]interface{}) (string, error) {
	if translationArguments == nil {
		translationArguments = make(map[string]interface{})
	}
	if count < 1 {
		count = 1
	}
	T, err := ts.i18bundle.Tfunc(localeCode)
	if err != nil {
		ts.logger.Info("Error - locale.translationservice", err)
		return "", err
	}
	label := T(key, count, translationArguments)
	if key == label {
		return label, errors.New("label not found")
	}
	return label, nil
}
func (ts *TranslationService) loadFiles() {
	if ts.filesLoaded && !ts.devmode {
		return
	}

	if ts.translationFile != "" {
		err := ts.i18bundle.LoadTranslationFile(ts.translationFile)
		if err != nil {
			ts.logger.Warn(fmt.Sprintf("Load translationfile failed: %s", err))
		}
	}
	if len(ts.translationFiles) > 0 {
		for _, file := range ts.translationFiles {
			if fileName, ok := file.(string); ok {
				err := ts.i18bundle.LoadTranslationFile(fileName)
				if err != nil {
					ts.logger.Warn(fmt.Sprintf("Load translationfile failed: %s", err))
				}
			}
		}
	}
	ts.filesLoaded = true
}

func (ts *TranslationService) initAndLoad() {
	if ts.i18bundle == nil {
		ts.i18bundle = bundle.New()
	}
	ts.loadFiles()
}
