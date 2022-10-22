package forms

type AddressIdForm struct {
	Id int `uri:"id" binding:"required,gte=1"`
}

type AddressDetailForm struct {
	Province     string `json:"province" binding:"required"`
	City         string `json:"city" binding:"required"`
	District     string `json:"district" binding:"required"`
	Address      string `json:"address" binding:"required"`
	SignerName   string `json:"signer_name" binding:"required"`
	SignerMobile string `json:"signer_mobile" binding:"required"`
	IsDefault    *bool  `json:"is_default" binding:""`
}
