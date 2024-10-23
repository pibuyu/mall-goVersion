package common

import "time"

type PublicModel struct {
	ID         int64     `json:"id" gorm:"column:id"`            // 主键ID
	CreateTime time.Time `json:"create_time" gorm:"create_time"` // 创建时间
}
