package forms

type commonFields struct {
	Mobile    string `json:"mobile" binding:"required,len=11,mobile"`
	Nick_Name string `json:"nick_name"`
	Gender    string
	Birthday  int32 `json:"birthday"`
	Role      int32 `json:"role" binding:"required"`
	Avatar    string
	Desc      string
	Country   string
	Provice   string
	City      string
	Area      string
	Address   string
}

type RegisterForm struct {
	CaptchaForm
	commonFields
	Password   string `json:"password" binding:"required,gte=6,lte=20"`
	Repassword string `json:"repassword" binding:"required,gte=6,lte=20,eqfield=Password"`
}

type UpdateUserForm struct {
	commonFields
	Uid int32 `json:"uid" binding:"required"`
}

type CheckLoginForm struct {
	Username string `json:"username" binding:"required,gte=2,lte=15"`
	Password string `json:"password" binding:"required,gte=6,lte=20"`
}

type UidForm struct {
	Uid int32 `uri:"uid" json:"uid" binding:"gte=0,required"`
}

type MobileForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,len=11,mobile"`
}

type UserListForm struct {
	Page  int32 `form:"page" binding:"gte=0"`
	Limit int32 `form:"limit" binding:"required,gte=0,lte=100"`
}
