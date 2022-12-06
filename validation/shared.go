package validation

import (
	val "github.com/go-ozzo/ozzo-validation"
)

type SharedValFunc = func(value interface{}) error

// some helpers for ozzo-validation
func CreateSaredValidation(rules ...val.Rule) SharedValFunc {
	return func(value interface{}) error {
		return val.Validate(value, rules...)
	}
}
