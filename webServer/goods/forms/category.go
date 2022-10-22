package forms

type CategoryForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	ParentId int32  `form:"parent_id" json:"parent_id"`
	IsTab    *bool  `form:"is_tab" json:"is_tab"`
	Level    int32  `form:"level" json:"level"`
}

// type CategoryFilterForm struct {
// 	Id int32 `uri:"cid" default:"1"`
// }
