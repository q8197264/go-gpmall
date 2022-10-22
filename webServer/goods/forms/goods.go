package forms

type GoodsForm struct {
	Name        string   `form:"name" json:"name" binding:"required,min=2,max=100"`
	GoodsSn     string   `form:"goods_sn" json:"goods_sn" binding:"required,min=2,lt=20"`
	Stocks      int32    `form:"stocks" json:"stocks" binding:"required,min=0"`
	CategoryId  int32    `form:"cid" json:"cid" binding:"required"`
	Brand       int32    `form:"bid" json:"bid" binding:"required"`
	MarketPrice float32  `form:"market_price" json:"market_price" binding:"required,min=0"`
	ShopPrice   float32  `form:"shop_price" json:"shop_price" binding:"required,min=0"`
	Subtitle    string   `form:"subtitle" json:"subtitle" binding:"required,min=0"`
	ShipFree    *bool    `form:"ship_free" json:"ship_free" binding:"required"`
	FrontImage  string   `form:"front_image" json:"front_image" binding:"required"`
	Images      []string `form:"images" json:"images" binding:"required"`
	DescImages  []string `form:"desc_images" json:"desc_images" binding:"required"`
}

type UpdateStatusForm struct {
	Fav_Num int32 `form:"fav_num" json:"fav_num"`
	On_Sale *bool `form:"on_sale" json:"on_sale"`
	Is_New  *bool `form:"is_new" json:"is_new"`
	Is_Hot  *bool `form:"is_hot" json:"is_hot"`
}

type GoodsByIdForm struct {
	Id int32 `uri:"id" json:"id" binding:"required"`
}
