package order

import "time"

type OmsOrderReturnApply struct {
	ID               int64      `gorm:"column:id;primaryKey" json:"id"  // 记录的唯一标识，对应数据库表中的主键列`
	OrderId          int64      `gorm:"column:order_id" json:"orderId"  // 订单id`
	CompanyAddressId *int64     `gorm:"column:company_address_id" json:"companyAddressId"  // 收货地址表id`
	ProductId        int64      `gorm:"column:product_id" json:"productId"  // 退货商品id`
	OrderSn          string     `gorm:"column:order_sn" json:"orderSn"  // 订单编号`
	CreateTime       *time.Time `gorm:"column:create_time" json:"createTime"  // 申请时间`
	MemberUsername   string     `gorm:"column:member_username" json:"memberUsername"  // 会员用户名`
	ReturnAmount     *float64   `gorm:"column:return_amount" json:"returnAmount"  // 退款金额`
	ReturnName       string     `gorm:"column:return_name" json:"returnName"  // 退货人姓名`
	ReturnPhone      string     `gorm:"column:return_phone" json:"returnPhone"  // 退货人电话`
	Status           int        `gorm:"column:status" json:"status"  // 申请状态：0->待处理；1->退货中；2->已完成；3->已拒绝`
	HandleTime       *time.Time `gorm:"column:handle_time" json:"handleTime"  // 处理时间`
	ProductPic       string     `gorm:"column:product_pic" json:"productPic"  // 商品图片`
	ProductName      string     `gorm:"column:product_name" json:"productName"  // 商品名称`
	ProductBrand     string     `gorm:"column:product_brand" json:"productBrand"  // 商品品牌`
	ProductAttr      string     `gorm:"column:product_attr" json:"productAttr"  // 商品销售属性：颜色：红色；尺 码：xl;`
	ProductCount     int        `gorm:"column:product_count" json:"productCount"  // 退货数量`
	ProductPrice     float64    `gorm:"column:product_price" json:"productPrice"  // 商品单价`
	ProductRealPrice float64    `gorm:"column:product_real_price" json:"productRealPrice"  // 商品实际支付单价`
	Reason           string     `gorm:"column:reason" json:"reason"  // 原因`
	Description      string     `gorm:"column:description" json:"description"  // 描述`
	ProofPics        string     `gorm:"column:proof_pics" json:"proofPics"  // 凭证图片，以逗号隔开`
	HandleNote       string     `gorm:"column:handle_note" json:"handleNote"  // 处理备注`
	HandleMan        string     `gorm:"column:handle_man" json:"handleMan"  // 处理人员`
	ReceiveMan       string     `gorm:"column:receive_man" json:"receiveMan"  // 收货人`
	ReceiveTime      *time.Time `gorm:"column:receive_time" json:"receiveTime"  // 收货时间`
	ReceiveNote      string     `gorm:"column:receive_note" json:"receiveNote"  // 收货备注`
}

func (o *OmsOrderReturnApply) TableName() string {
	return "oms_order_return_apply"
}
