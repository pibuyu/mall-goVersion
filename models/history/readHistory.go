package history

type MemberReadHistory struct {
	Id              int64   `json:"id" gorm:"id"`
	MemberId        int64   `json:"member_id" gorm:"member_id"`
	MemberNickname  string  `json:"member_nickname" gorm:"member_nickname"`
	MemberIcon      string  `json:"member_icon" gorm:"member_icon"`
	ProductId       int64   `json:"product_id" gorm:"product_id"`
	ProductName     string  `json:"product_name" gorm:"product_name"`
	ProductIcon     string  `json:"product_icon" gorm:"product_icon"`
	ProductSubTitle string  `json:"product_sub_title" gorm:"product_sub_title"`
	ProductPrice    float32 `json:"product_price" gorm:"product_price"`
	CreateTime      string  `json:"create_time" gorm:"create_time"`
	ProductPic      string  `json:"product_pic" gorm:"product_pic"`
}

func (MemberReadHistory) TableName() string {
	return "ums_member_read_history"
}
