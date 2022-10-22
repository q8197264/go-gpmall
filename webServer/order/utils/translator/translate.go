package translator

import (
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

	"webServer/order/global"
)

/*
 */
type validate struct {
	locale string
}

func NewDefaultConfig(locale string) *validate {
	return &validate{
		locale: locale,
	}
}

// Gin 绑定 validator
func (v *validate) InitValidator(tags map[string]map[string]string) (ut.Translator, error) {
	en := en.New()
	zh := zh.New()
	un := ut.New(en, zh, en)

	validate, ok := binding.Validator.Engine().(*validator.Validate)
	trans, ok2 := un.GetTranslator(v.locale)
	if !ok || !ok2 {
		return nil, fmt.Errorf("翻译初始化失败：un.GetTranslator(%s)", v.locale)
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		// 如果有多个标签名：则只分割两个，并返回取首个
		tagName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if tagName == "-" {
			return ""
		}
		return tagName
	})
	switch v.locale {
	case "en":
		translate_en.RegisterDefaultTranslations(validate, trans)
	case "zh":
		translate_zh.RegisterDefaultTranslations(validate, trans)
	default:
		translate_en.RegisterDefaultTranslations(validate, trans)
	}

	for tag, item := range tags {
		v.validateRuleOverride(validate, trans, tag, item)
	}

	return trans, nil
}

func (v *validate) validateRuleOverride(validate *validator.Validate, trans ut.Translator, tag string, item map[string]string) {
	validate.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		b, err := regexp.MatchString(item["pattern"], fl.FieldName())
		if err != nil {
			zap.S().Warnf("%s 字段验证失败：%s", tag, err.Error())
		}
		return b
	})

	validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, item["tip"], true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		s, err := ut.T(fe.Field())
		if err != nil {
			zap.S().Warnf("%s 字段验证信息翻译失败：", tag, err.Error())
		}
		return s
	})
}

// 格式化错误提示 【剔除多余字符】
func FormatValidationError(errs validator.ValidationErrors) map[string]string {
	tips := errs.Translate(global.Trans)
	msg := make(map[string]string)
	for k, v := range tips {
		ss := strings.Split(k, ".")
		msg[ss[len(ss)-1]] = v
	}

	return msg
}
