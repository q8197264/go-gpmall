package initialize

import transaltor "webServer/userop/utils/translator"

func InitTranslator() {
	var jsonConf = map[string]map[string]string{
		"mobile": {
			"pattern": "^1([38][0-9]|14[579]|16[6]|5[^4]|7[1-35-8|9[189]])\\d{8}",
			"tip":     "手机号不合法",
		},
		"username": {
			"pattern": "^[a-zA-Z][0-9a-zA-Z_]{6,18}",
			"tip":     "用户名不合法",
		},
	}
	transaltor.ValidateTranslator("zh", jsonConf)
}
