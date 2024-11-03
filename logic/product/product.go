package product

import (
	"errors"
	"fmt"
	"gomall/global"
	receive "gomall/interaction/receive/product"
	pmsBrandModels "gomall/models/brand"
	"gomall/models/cart"
	"gomall/models/coupon"
	"gomall/models/home"
	productModels "gomall/models/product"
)

// Detail 获取前台商品详情
func Detail(id int64) (result productModels.PmsPortalProductDetail, err error) {

	//根据id获取商品信息
	product := &home.PmsProduct{}
	if err := product.GetById(id); err != nil {
		return productModels.PmsPortalProductDetail{}, errors.New("获取前台商品详情时,查询product表出错:" + err.Error())
	}
	result.Product = *product
	//获取品牌信息
	brand := &pmsBrandModels.PmsBrand{}
	if err := brand.GetById(product.BrandId); err != nil {
		return productModels.PmsPortalProductDetail{}, errors.New("获取前台商品详情时,查询brand表出错:" + err.Error())
	}
	result.Brand = *brand
	//获取商品属性信息
	productAttributeList := &productModels.PmsProductAttributeList{}
	if err := productAttributeList.GetById(product.ProductAttributeCategoryId); err != nil {
		return productModels.PmsPortalProductDetail{}, errors.New("获取前台商品详情时,查询PmsProductAttribute表出错:" + err.Error())
	}
	result.ProductAttributeList = *productAttributeList
	//获取商品属性值信息
	if len(*productAttributeList) != 0 {
		attributeIds := make([]int64, 0)
		for _, v := range *productAttributeList {
			attributeIds = append(attributeIds, v.ID)
		}

		productAttributeValueList := &productModels.PmsProductAttributeValueList{}
		if err := productAttributeValueList.GetByProductIdAndAttributeIds(product.Id, attributeIds); err != nil {
			return productModels.PmsPortalProductDetail{}, errors.New("获取前台商品详情时,查询PmsProductAttributeValue表出错:" + err.Error())
		}
		result.ProductAttributeValueList = *productAttributeValueList
	}
	//获取商品SKU库存信息
	skuStockList := &cart.PmsSkuStockList{}
	if err := skuStockList.GetByProductId(product.Id); err != nil {
		return productModels.PmsPortalProductDetail{}, errors.New("获取前台商品详情时,查询PmsSkuStock表出错:" + err.Error())
	}
	result.SkuStockList = *skuStockList
	//商品阶梯价格设置
	if product.PromotionType == 3 {
		productLadderList := &cart.PmsProductLadderList{}
		if err := productLadderList.GetByProductId(product.Id); err != nil {
			return productModels.PmsPortalProductDetail{}, errors.New("获取前台商品详情时,查询PmsProductLadder表出错:" + err.Error())
		}
		result.ProductLadderList = *productLadderList
	}
	//商品满减价格设置
	if product.PromotionType == 4 {
		productFullReductionList := &cart.PmsProductFullReductionList{}
		if err := productFullReductionList.GetByProductId(product.Id); err != nil {
			return productModels.PmsPortalProductDetail{}, errors.New("获取前台商品详情时,查询PmsProductFullReduction表出错:" + err.Error())
		}
		result.ProductFullReductionList = *productFullReductionList
	}
	//商品可用优惠券
	availableCouponList, err := getAvailableCouponList(product.Id, product.ProductCategoryId)
	if err != nil {
		return productModels.PmsPortalProductDetail{}, errors.New("获取前台商品详情时,获取可用优惠券出错:" + err.Error())
	}
	result.CouponList = availableCouponList
	return
}

// todo:This is a special comment.不知道怎么写sql的时候，不妨先写出原生sql语句，然后用这种方式来搞定
func getAvailableCouponList(productId int64, productCategoryId int64) (result []coupon.SmsCoupon, err error) {
	query := `
        SELECT *
        FROM sms_coupon
        WHERE use_type = 0
          AND start_time < NOW()
          AND end_time > NOW()
        UNION
        SELECT c.*
        FROM sms_coupon_product_category_relation cpc
        LEFT JOIN sms_coupon c ON cpc.coupon_id = c.id
        WHERE c.use_type = 1
          AND c.start_time < NOW()
          AND c.end_time > NOW()
          AND cpc.product_category_id = ?
        UNION
        SELECT c.*
        FROM sms_coupon_product_relation cp
        LEFT JOIN sms_coupon c ON cp.coupon_id = c.id
        WHERE c.use_type = 2
          AND c.start_time < NOW()
          AND c.end_time > NOW()
          AND cp.product_id = ?
    `
	if err := global.Db.Raw(query, productCategoryId, productId).Scan(&result).Error; err != nil {
		return nil, err
	}
	return
}

// CategoryTreeList 树形获取所有商品
func CategoryTreeList() (result []productModels.PmsProductCategoryNode, err error) {
	var allList []productModels.PmsProductCategory
	if err := global.Db.Model(&productModels.PmsProductCategory{}).Find(&allList).Error; err != nil {
		return nil, errors.New("树形获取所有商品时,查询PmsProductCategory表出错:" + err.Error())
	}

	for _, item := range allList {
		if item.ParentID == 0 { //当前节点是个根节点
			node := convert(item, allList)
			result = append(result, node)
		}
	}
	return result, nil
}

// 递归方法，组成树形结构
func convert(item productModels.PmsProductCategory, allList []productModels.PmsProductCategory) productModels.PmsProductCategoryNode {
	node := productModels.PmsProductCategoryNode{
		PmsProductCategory: item,
	}
	for _, subItem := range allList {
		if subItem.ParentID == item.ID {
			childNode := convert(subItem, allList)
			node.Children = append(node.Children, childNode)
		}
	}
	return node
}

func Search(data *receive.SearchReqStruct) (result []home.PmsProduct, err error) {
	query := global.Db.Model(&home.PmsProduct{}).
		Where("delete_status = ? AND publish_status = ?", 0, 1)

	//模糊匹配
	if data.Keyword != "" {
		//sql里的通配符% 表示匹配任意多个字符；前面再多加一个%是拼接时的转义字符，表示向keyword拼接一个真实的%
		keyword := fmt.Sprintf("%%%s%%", data.Keyword) // equivalent to "%keyword%"
		query = query.Where("name LIKE ?", keyword)
	}

	if data.BrandId != 0 {
		query = query.Where("brand_id = ?", data.BrandId)
	}
	if data.ProductCategoryId != 0 {
		query = query.Where("product_category_id  = ?", data.ProductCategoryId)
	}

	//1->按新品；2->按销量；3->价格从低到高；4->价格从高到低
	switch data.Sort {
	case 1:
		query = query.Order("id desc")
	case 2:
		query = query.Order("sale desc")
	case 3:
		query = query.Order("price asc")
	case 4:
		query = query.Order("price desc")
	}
	if err := query.Offset((data.PageNum - 1) * data.PageSize).Limit(data.PageSize).
		Find(&result).Error; err != nil {
		return nil, errors.New("分页查询商品列表失败:" + err.Error())
	}
	return result, nil
}
