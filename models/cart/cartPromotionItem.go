package cart

type CartPromotionItem struct {
	OmsCartItem

	PromotionMessage string  `json:"promotionMessage"` //促销活动信息
	ReduceAmount     float32 `json:"reduceAmount"`     //促销活动减去的金额，针对每个商品
	RealStock        int     `json:"realStock"`        //锁定库存-剩余库存=真实库存
	Integration      int     `json:"integration"`      //购买商品赠送的积分
	Growth           int     `json:"growth"`           //购买商品赠送的成长值
}

type CartPromotionItemList []CartPromotionItem
