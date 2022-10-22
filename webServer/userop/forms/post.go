package forms

type PostListForm struct {
	Page  int32 `form:"page" binding:""`
	Limit int32 `form:"nums" binding:""`
}

type PostForm struct {
	Id      int32  `json:"id" binding:""`
	Type    int32  `json:"type" binding:"required"`
	Subject string `json:"subject" binding:"required,gte=5,lte=55"`
	Message string `json:"message" binding:"required,gte=5,lte=255"`
	File    string `json:"file" binding:""`
}
