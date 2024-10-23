package home

import "gomall/models/home"

// GetHomeContentResponseStruct 响应主页内容
type GetHomeContentResponseStruct struct {
	//轮播广告
	AdvertiseList home.AdvertiseList `json:"advertiseList"`
	//推荐品牌
	BrandList home.BrandList `json:"brandList"`
	//当前秒杀场次
	HomeFlashPromotion home.HomeFlashPromotion `json:"homeFlashPromotion"`
	//新品推荐
	NewProductList home.PmsProductList `json:"newProductList"`
	//人气推荐
	HotProductList home.PmsProductList `json:"hotProductList"`
	//推荐专题
	SubjectList home.RecommendSubjectList `json:"subjectList"`
}
