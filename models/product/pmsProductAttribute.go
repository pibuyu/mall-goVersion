package product

import "gomall/global"

type PmsProductAttribute struct {
	ID                         int64  `json:"id" gorm:"column:id"`                                                    // ID
	ProductAttributeCategoryID int64  `json:"productAttributeCategoryId" gorm:"column:product_attribute_category_id"` // 产品属性分类ID
	Name                       string `json:"name" gorm:"column:name"`                                                // 属性名称
	SelectType                 int    `json:"selectType" gorm:"column:select_type"`                                   // 属性选择类型：0->唯一；1->单选；2->多选
	InputType                  int    `json:"inputType" gorm:"column:input_type"`                                     // 属性录入方式：0->手工录入；1->从列表中选取
	InputList                  string `json:"inputList" gorm:"column:input_list"`                                     // 可选值列表，以逗号隔开
	Sort                       int    `json:"sort" gorm:"column:sort"`                                                // 排序字段：最高的可以单独上传图片
	FilterType                 int    `json:"filterType" gorm:"column:filter_type"`                                   // 分类筛选样式：1->普通；2->颜色
	SearchType                 int    `json:"searchType" gorm:"column:search_type"`                                   // 检索类型；0->不需要进行检索；1->关键字检索；2->范围检索
	RelatedStatus              int    `json:"relatedStatus" gorm:"column:related_status"`                             // 相同属性产品是否关联；0->不关联；1->关联
	HandAddStatus              int    `json:"handAddStatus" gorm:"column:hand_add_status"`                            // 是否支持手动新增；0->不支持；1->支持
	Type                       int    `json:"type" gorm:"column:type"`                                                // 属性的类型；0->规格；1->参数
}

type PmsProductAttributeList []PmsProductAttribute

func (list *PmsProductAttributeList) GetById(id int64) error {
	return global.Db.Model(PmsProductAttribute{}).
		Where("product_attribute_category_id = ?", id).Find(list).Error
}

func (PmsProductAttribute) TableName() string {
	return "pms_product_attribute"
}
