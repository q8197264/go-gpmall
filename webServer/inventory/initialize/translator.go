package initialize

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func InitTranslator(locale string) {
	en := en.New()
	zh := zh.New()

	uni := ut.New(en, zh, en)
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		trans, _ := uni.GetTranslator(locale)
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(validate, trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(validate, trans)
		}
		// 自定义规则
		translateOverride(validate, trans)
	}
}

func translateOverride(validate *validator.Validate, trans ut.Translator) {
	// 注册指定tag的验证规则
	validate.RegisterValidation("mobile", func(fe validator.FieldLevel) bool {
		ok, err := regexp.MatchString(
			`^1([38][0-9]|14[579]|16[6]|5[^4]|7[1-35-8|9[189]])\d{8}$`,
			fe.FieldName(),
		)
		if err != nil {
			zap.S().Warn(err.Error())
		}
		return ok
	})

	// 注册指定tag的错误翻译
	validate.RegisterTranslation("mobile", trans, func(ut ut.Translator) error {
		return ut.Add("mobile", "不是合法的手机号码", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field())
		return t
	})
}
