package forms

type AddGoodsForm struct {
	GoodsId int `json:"gid" binding:"required,gte=1"`
	Nums    int `json:"nums" binding:"required,gte=1"`
}
