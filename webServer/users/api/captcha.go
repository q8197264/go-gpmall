package api

import (
	"image/color"
	"net/http"
	"webServer/users/forms"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	// var captchaForm forms.CaptchaForm
	// if err := c.ShouldBind(&captchaForm); err != nil {
	// 	HandleValidateHttp(err, c)
	// 	return
	// }

	var driver base64Captcha.Driver
	CaptchaType := "math"
	switch CaptchaType {
	case "audio":
		var driverAudio *base64Captcha.DriverAudio
		driver = driverAudio
	case "string": //数字混合字母
		driver = base64Captcha.NewDriverString(80, 240, 15, 7, 5, "1234567890qwertyuioplkjhgfdsazxcvbnm", nil, nil, nil)
	case "math": //加减乘除
		driver = base64Captcha.NewDriverMath(80, 240, 15, 7, &color.RGBA{0, 0, 0, 0}, nil, nil)
	case "chinese": //汉字
		driver = base64Captcha.NewDriverChinese(80, 240, 5, 7, 2, "消费,没声,在无,任体是,何出能", nil, nil, nil)
	default: //数字
		driver = base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	}

	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   1,
			"errmsg": "generate captcha fail",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       0,
		"captchaId":  id,
		"captchaImg": b64s,
		"message":    "success",
	})
}

// base64Captcha verify http handler
func CaptchaVerify(c *gin.Context) {
	var captchaForm forms.CaptchaForm
	if err := c.ShouldBind(&captchaForm); err != nil {
		HandleValidateHttp(err, c)
		return
	}

	//verify the captcha
	// arg three: Do you want to delete captcahId after verify?
	if !store.Verify(captchaForm.CaptchaId, captchaForm.VerifyValue, false) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   1,
			"errmsg": "fail",
		})
		return
	}

	//set json response
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    "验证码正确",
		"message": "success",
	})
}
