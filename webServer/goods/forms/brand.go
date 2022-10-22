package forms

type BrandForm struct {
	Name string `form:"name" json:"name" binding:"required"`
	Logo string `form:"logo" json:"logo" binding:"required"`
}

type CategoryBrandForm struct {
	CategoryId int32 `form:"cid" json:"cid" binding:"required"`
	BrandId    int32 `form:"bid" json:"bid" binding:"required"`
}
