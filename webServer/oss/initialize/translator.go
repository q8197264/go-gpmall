package initialize

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"

	"webServer/oss/global"
)

func InitTranslator(locale string) {
	en := en.New()
	zh := zh.New()
	uni := ut.New(en, zh, en)

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		global.Trans, _ = uni.GetTranslator(locale)
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(validate, global.Trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(validate, global.Trans)
		}

		translateOverride(validate, global.Trans)
	}
}

// 自定义汉化 gin框架的 英文验证(validator)提示
func translateOverride(validate *validator.Validate, trans ut.Translator) {
	if err := validate.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		ok, err := regexp.MatchString(
			`^1([38][0-9]|14[579]|16[6]|5[^4]|7[1-35-8|9[189]])\d{8}$`,
			fl.Field().String(),
		)
		if err != nil {
			zap.S().Warn(err.Error())
		}
		return ok
	}); err != nil {
		zap.S().Info("mobile register validate fail")
	}
	err := validate.RegisterTranslation("mobile", trans, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0} 必填手机号码", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field())

		return t
	})
	if err != nil {
		zap.S().Warn(err.Error())
	}
}
