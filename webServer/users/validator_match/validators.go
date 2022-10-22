package validator_match

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|16[6]|5[^4]|7[1-35-8|9[189]])\d{8}$`, mobile)
	if ok {
		return true
	}
	return false
}
