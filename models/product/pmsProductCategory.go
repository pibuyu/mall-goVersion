package product

type PmsProductCategory struct {
	ID           int64  `json:"id" gorm:"column:id"`                      // ID
	ParentID     int64  `json:"parentId" gorm:"column:parent_id"`         // 上级分类的编号：0表示一级分类
	Name         string `json:"name" gorm:"column:name"`                  // 分类名称
	Level        int    `json:"level" gorm:"column:level"`                // 分类级别：0->1级；1->2级
	ProductCount int    `json:"productCount" gorm:"column:product_count"` // 产品数量
	ProductUnit  string `json:"productUnit" gorm:"column:product_unit"`   // 产品单位
	NavStatus    int    `json:"navStatus" gorm:"column:nav_status"`       // 是否显示在导航栏：0->不显示；1->显示
	ShowStatus   int    `json:"showStatus" gorm:"column:show_status"`     // 显示状态：0->不显示；1->显示
	Sort         int    `json:"sort" gorm:"column:sort"`                  // 排序字段
	Icon         string `json:"icon" gorm:"column:icon"`                  // 图标
	Keywords     string `json:"keywords" gorm:"column:keywords"`          // 关键词
	Description  string `json:"description" gorm:"column:description"`    // 描述
}

func (PmsProductCategory) TableName() string {
	return "pms_product_category"
}
