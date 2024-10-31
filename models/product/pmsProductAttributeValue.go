package product

import "gomall/global"

type PmsProductAttributeValue struct {
	ID                 int64  `json:"id" gorm:"column:id"`                                   // ID
	ProductID          int64  `json:"productId" gorm:"column:product_id"`                    // 产品ID
	ProductAttributeID int64  `json:"productAttributeId" gorm:"column:product_attribute_id"` // 产品属性ID
	Value              string `json:"value" gorm:"column:value"`                             // 手动添加规格或参数的值，参数单值，规格有多个时以逗号隔开
}

type PmsProductAttributeValueList []PmsProductAttributeValue

// 根据productId和attributeIds查询属性值列表
func (list *PmsProductAttributeValueList) GetByProductIdAndAttributeIds(productId int64, attributeIds []int64) (err error) {
	if err = global.Db.Debug().Model(&PmsProductAttributeValue{}).
		Where("product_id = ? AND product_attribute_id IN (?)", productId, attributeIds).Find(list).Error; err != nil {
		return err
	}
	return nil
}
func (PmsProductAttributeValue) TableName() string {
	return "pms_product_attribute_value"
}
