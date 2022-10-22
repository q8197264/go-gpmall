package forms

type FavoritesForm struct {
	UserId  int32 `uri:"uid" json:"uid" binding:""`
	GoodsId int32 `uri:"gid" json:"gid" binding:"required"`
}
