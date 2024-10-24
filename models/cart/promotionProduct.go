package cart

import "gomall/models/home"

type PromotionProduct struct {
	Product home.PmsProduct `json:"product"`

	SkuStockList         []PmsSkuStock             `json:"sku_stock_list"`         //商品库存信息
	ProductLadderList    []PmsProductLadder        `json:"product_ladder_list"`    //商品打折信息
	ProductFullReduction []PmsProductFullReduction `json:"product_full_reduction"` //商品满减信息
}
