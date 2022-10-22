package response

type OrderGoodsDetail struct {
	GoodsId     int32
	GoodsName   string
	MarketPrice float32
	ShopPrice   float32
	Nums        int32
}

type OrderDetail struct {
	OrderId      int32
	OrderSn      string
	UserId       int32
	SignerMobile string
	Amount       float32
	PayAmount    float32
	PayTime      string
	GoodsList    []OrderGoodsDetail
}
