package forms

type CaptchaForm struct {
	CaptchaId   string `form:"captcha_id" json:"captcha_id" binding:"required"`
	VerifyValue string `form:"verify_value" json:"verify_value" binding:"required"`
	CaptchaType string
}
