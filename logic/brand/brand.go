package brand

import (
	"errors"
	"gomall/global"
	"gomall/models/brand"
	"gomall/models/home"
)

func Detail(brandId int64) (result brand.PmsBrand, err error) {
	oneBrand := &brand.PmsBrand{}
	if err := oneBrand.GetById(brandId); err != nil {
		return brand.PmsBrand{}, errors.New("查询品牌信息出错:" + err.Error())
	}
	return *oneBrand, nil
}
func RecommendList(pageNum int, pageSize int) (result []brand.PmsBrand, err error) {
	if err = global.Db.Model(&brand.PmsBrand{}).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&result).Error; err != nil {
		return nil, errors.New("分页获取推荐品牌出错" + err.Error())
	}
	return
}

func ProductList(brandId int64, pageNum int, pageSize int) (result []home.PmsProduct, err error) {
	if err = global.Db.Model(&home.PmsProduct{}).
		Where("brand_id = ? and delete_status = ? and publish_status = ?", brandId, 0, 1).
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&result).Error; err != nil {
		return nil, errors.New("获取品牌相关的商品出错" + err.Error())
	}
	return
}
