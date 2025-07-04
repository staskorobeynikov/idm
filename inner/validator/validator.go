package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"unicode"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	validate := validator.New()
	err := validate.RegisterValidation("minnows3", minNonWhitespaceCount)
	if err != nil {
		return nil
	}
	return &Validator{validate: validate}
}

func (v Validator) Validate(request any) (err error) {
	err = v.validate.Struct(request)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			return validateErrs
		}
	}
	return err
}

func minNonWhitespaceCount(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	count := 0
	for _, r := range text {
		if !unicode.IsSpace(r) {
			count++
		}
	}
	return count >= 3
}
