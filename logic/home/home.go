package home

import (
	"errors"
	"gomall/global"
	homeReceive "gomall/interaction/receive/home"
	"gomall/models/home"
	"time"
)

func GetHomeAdvertiseList() (results home.AdvertiseList, err error) {
	if err := global.Db.Where("type = ?", 1).Where("status = ?", 1).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func GetRecommendBrandList(offset int, limit int) (results home.BrandList, err error) {
	if err = global.Db.Model(&home.Brand{}).Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, err
	}
	return
}

func GetHomeFlashPromotion() (results *home.HomeFlashPromotion, err error) {
	homeFlashPromotion := &home.HomeFlashPromotion{}
	//获取当前正在进行的秒杀活动
	curDate := time.Now().Format("2006-01-02")
	flashPromotion, err := getFlashPromotion(curDate)
	if err != nil {
		return nil, err
	}
	if flashPromotion == nil {
		return
	}

	//根据时间获取秒杀场次
	promotionSession, err := getFlashPromotionSession(curDate)
	if err != nil {
		return nil, err
	}
	if promotionSession == nil {
		return
	}
	//获取下一个秒杀场次
	homeFlashPromotion.StartDate = promotionSession.StartDate
	homeFlashPromotion.EndDate = promotionSession.EndDate
	nextSession, err := getNextFlashPromotionSession(homeFlashPromotion.StartDate)
	if err != nil {
		return nil, err
	}
	if nextSession == nil {
		return
	}
	homeFlashPromotion.NextStartDate = nextSession.StartDate
	homeFlashPromotion.NextEndDate = nextSession.EndDate
	//获取秒杀商品
	flashProductList, err := getFlashProductList(flashPromotion.Id, promotionSession.Id)
	if err != nil {
		return nil, err
	}
	homeFlashPromotion.ProductList = flashProductList

	return homeFlashPromotion, nil
}

// getFlashPromotion 获取当前时间正在进行的秒杀活动
func getFlashPromotion(curDate string) (results *home.FlashPromotion, err error) {
	if err := global.Db.Model(&home.FlashPromotion{}).Where("start_date <= ?", curDate).Where("end_date >= ?", curDate).Find(&results).Error; err != nil {
		return nil, errors.New("查询FlashPromotion表错误:" + err.Error())
	}
	return
}

// getFlashPromotionSession 获取当前时间正在进行的秒杀场次
func getFlashPromotionSession(curDate string) (results *home.HomeFlashPromotionSession, err error) {
	if err = global.Db.Model(&home.HomeFlashPromotionSession{}).Where("start_time <= ?", curDate).Where("end_time >= ?", curDate).Find(&results).Error; err != nil {
		return nil, errors.New("查询PromotionSession表错误:" + err.Error())
	}
	if results == nil {
		return nil, errors.New("当前没有进行的秒杀的场次")
	}
	return
}

// getNextFlashPromotionSession 获取下一个秒杀场次
func getNextFlashPromotionSession(curDate string) (results *home.HomeFlashPromotionSession, err error) {
	if err = global.Db.Model(&home.HomeFlashPromotionSession{}).Where("start_time <= ?", curDate).Where("end_time >= ?", curDate).Find(&results).Error; err != nil {
		return nil, errors.New("获取nextSession，查询PromotionSession表错误:" + err.Error())
	}
	if results == nil {
		return nil, errors.New("指定时间内没有下一场秒杀的场次")
	}
	return
}

func getFlashProductList(flashPromotionId int64, promotionSessionId int64) (productList home.FlashPromotionProductList, err error) {
	//从relation表，根据flashPromotionId和promotionSessionId字段查出ProductId，再根据查出的ProductId和product表进行关联，从而获取到商品的所有信息以及与秒杀活动有关的秒杀价格、秒杀数量和用于秒杀的库存三个字段
	//将上述结果封装为HomeFlashPromotionProduct的列表，返回即可
	var reletions []home.FlashPromotionProductRelation
	if err := global.Db.Model(&home.FlashPromotionProductRelation{}).Where("flash_promotion_id = ?", flashPromotionId).
		Where("flash_promotion_session_id = ?", promotionSessionId).Find(&reletions).Error; err != nil {
		return nil, errors.New("查询关系表错误:" + err.Error())
	}
	productList = make([]home.HomeFlashPromotionProduct, 0)
	for _, relation := range reletions {
		var product home.PmsProduct
		if err := global.Db.First(&product, relation.ProductId).Error; err != nil {
			return nil, errors.New("查询PmsProduct表出错:" + err.Error())
		}
		//将结果封装到HomeFlashPromotionProduct
		flashProduct := home.HomeFlashPromotionProduct{
			Product:             product,
			FlashPromotionPrice: relation.FlashPromotionPrice,
			FlashPromotionCount: relation.FlashPromotionCount,
			FlashPromotionLimit: relation.FlashPromotionLimit,
		}
		productList = append(productList, flashProduct)
	}
	return productList, nil

}

func GetNewProductList(offset int, limit int) (productList home.PmsProductList, err error) {
	if err := global.Db.Model(&home.PmsProduct{}).Joins("LEFT JOIN sms_home_new_product ON pms_product.id = sms_home_new_product.product_id").
		Order("sms_home_new_product.sort desc").Limit(limit).Offset(offset).Find(&productList).Error; err != nil {
		return nil, errors.New("查询sms_home_new_product出错:" + err.Error())
	}
	return productList, nil
}

func GetHotProductList(offset int, limit int) (productList home.PmsProductList, err error) {
	if err := global.Db.Model(&home.PmsProduct{}).Joins("LEFT JOIN sms_home_recommend_product ON pms_product.id=sms_home_recommend_product.product_id").
		Where("pms_product.delete_status = ? and pms_product.publish_status = ?", 0, 1).Order("sms_home_recommend_product.sort desc").Limit(limit).Offset(offset).Find(&productList).Error; err != nil {
		return nil, errors.New("查询sms_home_recommend_product出错:" + err.Error())
	}
	return productList, nil
}

func GetRecommendSubjectList(offset int, limit int) (recommendSubjectList home.RecommendSubjectList, err error) {
	if err := global.Db.Model(&home.SmsHomeRecommendSubject{}).Joins("LEFT JOIN cms_subject ON cms_subject.id=sms_home_recommend_subject.subject_id").
		Where("sms_home_recommend_subject.recommend_status = ?", 1).
		Where("cms_subject.show_status = ?", 1).
		Order("sms_home_recommend_subject.sort desc").
		Limit(limit).Offset(offset).Find(&recommendSubjectList).Error; err != nil {
		return nil, errors.New("查询sms_home_recommend_subject表出错:" + err.Error())
	}
	return
}

func GetRecommendProductList(offset int, limit int) (productList home.PmsProductList, err error) {
	// 查询所有商品的总数
	var total int64
	if err := global.Db.Model(&home.PmsProduct{}).Where("delete_status = ?", 0).Where("publish_status = ?", 1).Count(&total).Error; err != nil {
		return nil, errors.New("获取商品总数出错:" + err.Error())
	}
	// 判断 offset 和 limit 是否超出商品数量,如果超出了商品数量，需要返回空值
	if offset >= int(total) {
		return productList, nil // 返回空结果
	}
	//暂时默认推荐所有商品
	if err := global.Db.Model(&home.PmsProduct{}).Where("delete_status = ?", 0).
		Where("publish_status = ?", 1).
		Offset(offset).Limit(limit).
		Find(&productList).Error; err != nil {
		return nil, errors.New("GetRecommendProductList查询出错:" + err.Error())
	}
	return
}

func GetSubjectListByCategoryId(data *homeReceive.GetSubjectListRequestStruct) (results home.CmsSubjectList, err error) {
	//构造基本查询
	query := global.Db.Model(&home.CmsSubject{}).Where("show_status = ?", 1)
	//如果CateId字段不为0，添加条件
	if data.CateId != 0 {
		query = query.Where("category_id = ?", data.CateId)
	}
	//执行查询
	if err := query.Offset((data.PageNum - 1) * data.PageSize).
		Limit(data.PageSize).Find(&results).Error; err != nil {
		return nil, errors.New("查询CmsSubject表出错:" + err.Error())
	}
	return
}

func GetProductCateList(parentId int64) (results home.ProductCategoryList, err error) {
	if err := global.Db.Model(&home.PmsProductCategory{}).Where("show_status = ?", 1).
		Where("parent_id = ?", parentId).Order("sort desc").Find(&results).Error; err != nil {
		return nil, errors.New("查询PmsProductCategory表出错:" + err.Error())
	}
	return
}
