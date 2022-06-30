package validatator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(level validator.FieldLevel) bool {
	mobile := level.Field().String()
	// 正则表达式判断是否合法
	regRuler := "^1[345789]{1}\\d{9}$"
	ok, _ := regexp.MatchString(regRuler, mobile)
	if !ok{
		return false
	}
	return true
}
