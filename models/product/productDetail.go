package product

import (
	"gomall/models/brand"
	"gomall/models/cart"
	"gomall/models/coupon"
	"gomall/models/home"
)

type PmsPortalProductDetail struct {
	Product                   home.PmsProduct                `json:"product" description:"商品信息"`                            // 商品信息
	Brand                     brand.PmsBrand                 `json:"brand" description:"商品品牌"`                              // 商品品牌
	ProductAttributeList      PmsProductAttributeList        `json:"productAttributeList" description:"商品属性与参数"`            // 商品属性与参数
	ProductAttributeValueList PmsProductAttributeValueList   `json:"productAttributeValueList" description:"手动录入的商品属性与参数值"` // 手动录入的商品属性与参数值
	SkuStockList              cart.PmsSkuStockList           `json:"skuStockList" description:"商品的sku库存信息"`                 // 商品的sku库存信息
	ProductLadderList         []cart.PmsProductLadder        `json:"productLadderList" description:"商品阶梯价格设置"`              // 商品阶梯价格设置
	ProductFullReductionList  []cart.PmsProductFullReduction `json:"productFullReductionList" description:"商品满减价格设置"`       // 商品满减价格设置
	CouponList                []coupon.SmsCoupon             `json:"couponList" description:"商品可用优惠券"`                      // 商品可用优惠券
}
