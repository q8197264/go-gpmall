package forms

type OrderInfoForm struct {
	UserId  int32  `json:"uid" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Mobile  string `json:"mobile" binding:"required"`
	Address string `json:"address" binding:"required"`
	Post    string `json:"post" binding:""`
}
