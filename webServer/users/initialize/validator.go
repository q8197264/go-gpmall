package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"

	"webServer/users/global"
	"webServer/users/validator_match"
)

func TranslateInit(locale string) (err error) {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法(错误提示中的字段用请求中使用的字段替代)
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		en := en.New()
		zh := zh.New()

		// 备用语言环境
		uni := ut.New(en, zh, en)

		global.Trans, ok = uni.GetTranslator(locale)

		if !ok {
			zap.S().Infof("翻译初始化失败：uni.GetTranslator(%s)", locale)
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(validate, global.Trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(validate, global.Trans)
		default:
			zh_translations.RegisterDefaultTranslations(validate, global.Trans)
		}

		customValidators(validate)
	}

	return nil
}

/*
// 格式化错误信息
func ErrToJson(errs validator.ValidationErrors) map[string]string {
	tips := errs.Translate(global.Trans)
	data := make(map[string]string)
	var t []string
	for k, v := range tips {
		t = strings.Split(k, ".")
		data[t[len(t)-1]] = v
	}

	return data
}
*/

// 注册自定义指定字段验证
func customValidators(validate *validator.Validate) {
	if err := validate.RegisterValidation("mobile", validator_match.ValidateMobile); err != nil {
		zap.S().Info("mobile register validate fail")
	}
	err := validate.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
		return ut.Add("mobile", "非法手机号码", true)
	}, func(trans ut.Translator, fe validator.FieldError) string {
		t, err := trans.T("mobile", fe.Field())
		if err != nil {
			panic(fe.(error).Error())
		}
		return t
	})
	if err != nil {
		zap.S().Info("mobile register validate translation fail")
	}
}
