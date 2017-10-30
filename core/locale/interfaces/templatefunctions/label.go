package templatefunctions

import (
	"log"

	"github.com/nicksnyder/go-i18n/i18n"
	"fmt"
)

type (
	// GetProduct is exported as a template function
	Label struct {
	}
)

// Name alias for use in template
func (tf Label) Name() string {
	return "__"
}

func (tf Label) Func() interface{} {

	return fmt.Sprintf

	return func(key string, defaultLabel string, args interface{}) string {
		log.Printf("called with key %v", key)
		i18n.LoadTranslationFile("translations/en-US.all.json")
		T, _ := i18n.Tfunc("en-US")
		log.Printf("Some testlabel: %#v", T("formerror_email_required"))
		label := T(key)
		if label == key && defaultLabel != "" {
			label = defaultLabel
		}
		return label
	}
}
