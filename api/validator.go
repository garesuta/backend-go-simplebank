package api

import (
	"github.com/backendproduction-2/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check if the currency is valid
		return util.IsSupportedCurrency(currency)
	}
	return false
}
