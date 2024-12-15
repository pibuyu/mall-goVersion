package es

import (
	"fmt"
	"gomall/global"
	receive "gomall/interaction/receive/es"
	esModels "gomall/models/es"
	"gomall/models/home"
	"gomall/models/product"
)

func ImportAll() (int, error) {
	esProductList, err := GetAllEsProductList(0)
	if err != nil {
		return 0, fmt.Errorf("es导入所有商品过程中，查询商品列表错误：%v", err)
	}
	count, err := SaveAll(esProductList)
	if err != nil {
		return 0, fmt.Errorf("es导入所有商品过程中，调用SaveAll方法出错：%v", err)
	}
	global.Logger.Infof("成功导入%d个商品到es中", count)
	return count, nil
}

// queryAttr := global.Db.Table("pms_product_attribute_value pav").
//
//	Select("pav.product_id, pav.id as attr_id, pav.value as value, pav.product_attribute_id as product_attribute_id, pa.type as type, pa.name as name").
//	Joins("LEFT JOIN pms_product_attribute pa ON pav.product_attribute_id = pa.id").
//	Where("pav.product_id IN (?)", getProductIds(productList))
func GetAllEsProductList(id int64) (result []esModels.EsProduct, err error) {
	var esProductList []esModels.EsProduct

	// Step 1: 查询 pms_product 表
	var productList []home.PmsProduct
	query := global.Db.Model(&home.PmsProduct{}).
		Where("delete_status = ? AND publish_status = ?", 0, 1)

	// 如果 id 不为 nil，则添加条件
	if id != 0 {
		query = query.Where("id = ?", id)
	}

	// 执行查询获取所有产品信息
	if err = query.Find(&productList).Error; err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	// Step 2: 查询与产品相关的属性和属性值.对每个pms_product，其id对应pms_product_attribute_value表中的多个项;然后pms_product_attribute_value.product_attribute_id又一一对应pms_product_attribute的表项
	productIds := getProductIds(productList)

	//根据product_id获取PmsProductAttributeValue表中的多个表项
	var productAttributeValueList []product.PmsProductAttributeValue
	if err = global.Db.Model(&product.PmsProductAttributeValue{}).
		Where("product_id IN (?)", productIds).
		Find(&productAttributeValueList).Error; err != nil {
		return nil, fmt.Errorf("查询属性值失败: %v", err)
	}
	productAttributeValueListIds := getProductAttributeValueListIds(productAttributeValueList)

	// 再根据pms_product_attribute_value.product_attribute_id去找pms_product_attribute表项
	var productAttributeList []product.PmsProductAttribute
	if err = global.Db.Model(&product.PmsProductAttribute{}).
		Where("id IN (?)", productAttributeValueListIds).
		Find(&productAttributeList).Error; err != nil {
		return nil, fmt.Errorf("查询属性失败: %v", err)
	}

	// Step 3: 将属性信息拼接到对应的产品信息中.每一个pms_product都对应一条esProduct
	for _, product := range productList {
		oneEsProduct := &esModels.EsProduct{
			ProductSn:           product.ProductSn,
			BrandId:             product.BrandId,
			BrandName:           product.BrandName,
			ProductCategoryId:   product.ProductCategoryId,
			ProductCategoryName: product.ProductCategoryName,
			Pic:                 product.Pic,
			Name:                product.Name,
			SubTitle:            product.SubTitle,
			Keywords:            product.Keywords,
			Price:               product.Price,
			Sale:                product.Sale,
			NewStatus:           product.NewStatus,
			RecommandStatus:     product.RecommandStatus,
			Stock:               product.Stock,
			PromotionType:       product.PromotionType,
			Sort:                product.Sort,
		}
		//再拼接AttrValueList
		oneAttrValue := &esModels.EsProductAttributeValue{}
		for _, value := range productAttributeValueList {
			if value.ProductID == product.Id {
				oneAttrValue.ProductAttributeId = value.ProductAttributeID
				oneAttrValue.Name = value.Value
				for _, attribute := range productAttributeList {
					if attribute.ID == value.ProductAttributeID {
						oneAttrValue.Type = attribute.Type
						oneAttrValue.Name = attribute.Name
					}
				}
				oneEsProduct.AttrValueList = append(oneEsProduct.AttrValueList, *oneAttrValue)
			}
		}
		esProductList = append(esProductList, *oneEsProduct)
	}

	return esProductList, nil
}

// 获取所有产品的 ID 列表，用于查询相关属性
func getProductIds(products []home.PmsProduct) []int64 {
	var ids []int64
	for _, product := range products {
		ids = append(ids, product.Id)
	}
	return ids
}

func getProductAttributeValueListIds(productAttributeValueList []product.PmsProductAttributeValue) []int64 {
	var ids []int64
	for _, value := range productAttributeValueList {
		ids = append(ids, value.ProductAttributeID)
	}
	return ids
}

func DeleteById(id int64) error {
	return Delete(id)
}

func Create(id int64) error {
	//先从数据库找到这个商品的信息以及attribute等信息；然后插入es中去
	result, err := GetAllEsProductList(id)
	if err != nil {
		return fmt.Errorf("根据id查询数据库错误：%v", err)
	}

	if len(result) == 0 {
		return nil
	}
	esProduct := result[0]
	return Save(esProduct)
}

func SimpleSearch(data *receive.SimpleSearchReqStruct) (result []esModels.EsProduct, err error) {
	global.Logger.Infof("es中关键字查找接受的关键字为：%s,pagenum=%d,pageSize=%d", data.Keyword, data.PageNum, data.PageSize)
	return SearchProductsByKeyword(data.Keyword, data.PageNum, data.PageSize)
}
