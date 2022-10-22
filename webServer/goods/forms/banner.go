package forms

type BannerForm struct {
	Index int32  `form:"index" json:"index"`
	Image string `form:"image" json:"image" binding:"required"`
	Url   string `form:"url" json:"url" binding:"required"`
}
