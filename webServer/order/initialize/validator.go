package initialize

import (
	"encoding/json"
	"path/filepath"
	"webServer/order/utils/translator"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"webServer/order/global"
)

type config struct{}

func InitTranslate(locale string) {
	content := []byte(`
		{
			"mobile":{"pattern": "^1([38][0-9]|14[579]|16[6]|5[^4]|7[1-35-8|9[189]])\\d{8}", "tip": "手机号不合法"},
			"username":{"pattern":"^[a-zA-Z][0-9a-zA-Z_]{6,18}","tip":"用户名不合法"}
		}
	`)
	var data map[string]map[string]string
	err := json.Unmarshal(content, &data)
	if err != nil {
		zap.S().DPanic(err.Error())
	}

	// readYamlConf()

	cfg := translator.NewDefaultConfig(locale)
	if global.Trans, err = cfg.InitValidator(data); err != nil {
		zap.S().Warnf(err.Error())
	}
}

func readYamlConf() map[string]map[string]string {
	root, _ := filepath.Abs(".")
	path := filepath.Join(root, "utils", "validator", "rules.yaml")
	vp := viper.New()
	vp.SetConfigFile(path)
	if err := vp.ReadInConfig(); err != nil {
		zap.S().DPanic(err.Error())
	}
	var data map[string]map[string]string
	vp.Unmarshal(&data)

	return data
}
