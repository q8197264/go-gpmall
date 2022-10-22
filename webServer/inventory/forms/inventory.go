package forms

type InventoryForm struct {
	GoodsId int `json:"gid" binding:"required"`
	Nums    int `json:"nums" binding:"required"`
}

type InvDetailForm struct {
	GoodsId int `uri:"id" binding:"required"`
}
