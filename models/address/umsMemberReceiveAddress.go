package address

type UmsMemberReceiveAddress struct {
	ID            int64  `gorm:"column:id;primaryKey" json:"id"  // 地址记录的唯一标识，对应数据库表中的主键列`
	MemberID      int64  `gorm:"column:member_id" json:"member_id"  // 会员的唯一标识，关联到会员表中的会员ID列`
	Name          string `gorm:"column:name" json:"name"  // 收货人名称`
	PhoneNumber   string `gorm:"column:phone_number" json:"phone_number"  // 收货人电话号码`
	DefaultStatus int    `gorm:"column:default_status" json:"default_status"  // 是否为默认收货地址，1表示是，0表示否`
	PostCode      string `gorm:"column:post_code" json:"post_code"  // 邮政编码`
	Province      string `gorm:"column:province" json:"province"  // 省份/直辖市名称`
	City          string `gorm:"column:city" json:"city"  // 城市名称`
	Region        string `gorm:"column:region" json:"region"  // 区名称`
	DetailAddress string `gorm:"column:detail_address" json:"detail_address"  // 详细地址(街道)信息`
}

func (u UmsMemberReceiveAddress) TableName() string {
	return "ums_member_receive_address"
}
