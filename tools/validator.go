package tools

import (
	v "github.com/go-playground/validator/v10"
)


var validator *v.Validate
func NewValidator() {
	validate := v.New()
	validator = validate
}

func Validate(i interface{}) error {
	return validator.Struct(i)
}