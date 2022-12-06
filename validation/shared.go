package validation

import (
	val "github.com/go-ozzo/ozzo-validation"
)

// some helpers for ozzo-validation
func CreateSaredValidation(rules ...val.Rule) func(value interface{}) error {
	return func(value interface{}) error {
		return val.Validate(value, rules...)
	}
}
