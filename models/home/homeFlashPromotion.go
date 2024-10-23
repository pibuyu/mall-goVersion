package home

type HomeFlashPromotion struct {
	StartDate     string                    `json:"start_time"` //本场开始时间
	EndDate       string                    `json:"end_time"`
	NextStartDate string                    `json:"next_start_time"`
	NextEndDate   string                    `json:"next_end_time"`
	ProductList   FlashPromotionProductList `json:"product_list"` //属于该秒杀活动的商品
}
