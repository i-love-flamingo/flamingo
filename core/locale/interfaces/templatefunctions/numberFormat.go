package templatefunctions

import (
	"github.com/leekchan/accounting"
	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/framework/config"
)

type (
	// NumberFormatFunc for formatting numbers
	NumberFormatFunc struct {
		Config             config.Map                      `inject:"config:locale.numbers"`
		TranslationService *application.TranslationService `inject:""`
	}
)

type NumberConf struct {
	Precision int
	Thousand string
	Decimal string
}

// Name alias for use in template
func (pff NumberFormatFunc) Name() string {
	return "numberFormat"
}

// Func as implementation of debug method
func (pff NumberFormatFunc) Func() interface{} {
	return func(value interface{}) string {

		numberConfig := NumberConf {}

		// read values from config if they are set
		precision, ok := pff.Config["precision"].(int)
		if ok {
			numberConfig.Precision = precision
		}
		thousand, ok := pff.Config["thousand"].(string)
		if ok {
			numberConfig.Thousand = thousand
		}
		decimal, ok := pff.Config["decimal"].(string)
		if ok {
			numberConfig.Decimal = decimal
		}

		return accounting.FormatNumber(value, numberConfig.Precision, numberConfig.Thousand, numberConfig.Decimal)
	}
}
