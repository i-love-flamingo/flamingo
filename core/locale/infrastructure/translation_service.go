package infrastructure

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"text/template"
	"time"

	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

// TranslationService is the default TranslationService implementation
type TranslationService struct {
	mutex            sync.Mutex
	lastReload       time.Time
	translationFiles []string
	logger           flamingo.Logger
	devmode          bool
	i18bundle        *bundle.Bundle
}

// check if translationService implements its interface
var _ domain.TranslationService = (*TranslationService)(nil)

// Inject dependencies
func (ts *TranslationService) Inject(
	logger flamingo.Logger,
	config *struct {
		DevMode          bool         `inject:"config:flamingo.debug.mode"`
		TranslationFile  string       `inject:"config:core.locale.translationFile,optional"`
		TranslationFiles config.Slice `inject:"config:core.locale.translationFiles,optional"`
	},
) {
	ts.logger = logger.WithField(flamingo.LogKeyModule, "locale").WithField(flamingo.LogKeyCategory, "locale.translationService")
	if config != nil {
		err := config.TranslationFiles.MapInto(&ts.translationFiles)
		if err != nil {
			ts.logger.Warn(fmt.Sprintf("could not map core.locale.translationFiles: %v", err))
		}

		if config.TranslationFile != "" {
			ts.translationFiles = append(ts.translationFiles, config.TranslationFile)
		}
		ts.devmode = config.DevMode
	}

	ts.i18bundle = bundle.New()
	ts.mutex.Lock()
	ts.loadFiles()
	ts.mutex.Unlock()
}

// TranslateLabel returns the result for translating a Label
func (ts *TranslationService) TranslateLabel(label domain.Label) string {
	ts.reloadFilesIfNecessary()
	translatedString, err := ts.translateWithLib(label.GetLocaleCode(), label.GetKey(), label.GetCount(), label.GetTranslationArguments())

	//while there is an error check fallBacks
	for _, fallbackLocale := range label.GetFallbackLocaleCodes() {
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
	ts.reloadFilesIfNecessary()
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

// AllTranslationKeys returns all keys for a given locale code
func (ts *TranslationService) AllTranslationKeys(localeCode string) []string {
	ts.reloadFilesIfNecessary()
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

func (ts *TranslationService) reloadFilesIfNecessary() {
	if !ts.devmode {
		return
	}

	var lastFileChange time.Time
	for _, fileName := range ts.translationFiles {
		stat, err := os.Stat(fileName)
		if err != nil {
			continue
		}

		if stat.ModTime().After(lastFileChange) {
			lastFileChange = stat.ModTime()
		}
	}

	ts.mutex.Lock()
	if lastFileChange.After(ts.lastReload) {
		ts.loadFiles()
	}
	ts.mutex.Unlock()
}

// loadFiles must only be called when mutex is locked
func (ts *TranslationService) loadFiles() {
	for _, fileName := range ts.translationFiles {
		err := ts.i18bundle.LoadTranslationFile(fileName)
		if err != nil {
			ts.logger.Warn(fmt.Sprintf("loading of translationfile %s failed: %s", fileName, err))
		}
	}

	ts.lastReload = time.Now()
}
