package model

import "github.com/go-playground/validator/v10"

const (
	male   = "MALE"
	female = "FEMALE"
)

func ValidateSex(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str != male && str != female {
		return false
	}
	return true
}
