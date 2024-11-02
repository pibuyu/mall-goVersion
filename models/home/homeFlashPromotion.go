package home

type HomeFlashPromotion struct {
	StartDate     string                    `json:"startTime"` //本场开始时间
	EndDate       string                    `json:"endTime"`
	NextStartDate string                    `json:"nextStartTime"`
	NextEndDate   string                    `json:"nextEndTime"`
	ProductList   FlashPromotionProductList `json:"productList"` //属于该秒杀活动的商品
}
