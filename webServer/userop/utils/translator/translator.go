package transaltor

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translate_en "github.com/go-playground/validator/v10/translations/en"
	translate_zh "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"
)

var translate ut.Translator

func ValidateTranslator(locale string, tags map[string]map[string]string) (ut.Translator, error) {
	en := en.New()
	zh := zh.New()
	un := ut.New(en, zh, en)

	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("")
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tagName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if tagName == "-" {
			return ""
		}
		return tagName
	})

	trans, ok := un.GetTranslator("zh")
	if !ok {
		return nil, errors.New("")
	}
	switch locale {
	case "en":
		translate_en.RegisterDefaultTranslations(validate, trans)
	case "zh":
		translate_zh.RegisterDefaultTranslations(validate, trans)
	default:
		translate_en.RegisterDefaultTranslations(validate, trans)
	}

	for tag, item := range tags {
		fmt.Printf("%+v\n", tag)
		validateRuleOverride(tag, item, validate, trans)
	}
	translate = trans
	return trans, nil
}

func validateRuleOverride(tag string, item map[string]string, validate *validator.Validate, trans ut.Translator) {
	validate.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		ok, err := regexp.MatchString(item["pattern"], fl.FieldName())
		if err != nil {
			zap.S().Warnf("验证错误:%s", err.Error())
		}
		return ok
	})
	validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, item["tip"], true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		s, err := ut.T(fe.Field())
		if err != nil {
			zap.S().Warnf("字段验证信息翻译失败:%s", err.Error())
		}
		return s
	})
}

func TranslateValidateErrors(err validator.ValidationErrors) map[string]string {
	va := err.Translate(translate)
	res := map[string]string{}
	for k, v := range va {
		st := strings.Split(k, ".")
		res[st[len(st)-1]] = v
	}

	return res
}
