package valid

import (
	"github.com/go-playground/validator/v10"
)

func Url(fl validator.FieldLevel) bool {
	url, ok := fl.Field().Interface().(string)
	if !ok || len(url) == 0 || len(url) >= 10000 {
		return false
	}
	return true
}
