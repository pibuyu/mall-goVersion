package cart

// 向购物车添加商品请求体
type AddCartItemRequestStruct struct {
	CreateDate        string  `json:"create_date"`
	DeleteStatus      int     `json:"delete_status"`
	Id                int     `json:"id"`
	MemberID          int64   `json:"memberId"`
	MemberNickname    string  `json:"memberNickname"`
	ModifyDate        string  `json:"modifyDate"`
	Price             float32 `json:"price" binding:"required"`             //必须
	ProductAttr       string  `json:"productAttr" binding:"required"`       //必须
	ProductBrand      string  `json:"productBrand" binding:"required"`      //必须
	ProductCategoryID int64   `json:"productCategoryId" binding:"required"` //必须
	ProductID         int     `json:"productId" binding:"required"`         //必须
	ProductName       string  `json:"productName" binding:"required"`       //必须
	ProductPic        string  `json:"productPic" binding:"required"`        //必须
	ProductSkuCode    string  `json:"productSkuCode" binding:"required"`    //必须
	ProductSkuID      int     `json:"productSkuId" binding:"required"`      //必须
	ProductSN         string  `json:"productSn" binding:"required"`         //必须
	ProductSubTitle   string  `json:"productSubTitle" binding:"required"`   //必须
	Quantity          int     `json:"quantity" binding:"required"`          //必须
}

// ClearCartRequestStruct 清空购物车请求体,不需要参数，从request.header取出userId
type ClearCartRequestStruct struct {
}

// DeleteCartItemsByIdsRequestStruct 根据商品id删除购物车中的商品请求体
type DeleteCartItemsByIdsRequestStruct struct {
	Ids []int64 `json:"ids"`
}

type GetProductByIdRequestStruct struct {
	ProductId int `json:"productId"`
}

type CartListPromotionRequestStruct struct {
	CartIds []int64 `json:"cartIds"`
}

type UpdateAttrRequestStruct struct {
	Id                int64   `json:"id"`
	ProductId         int64   `json:"product_id"`
	ProductSkuId      int64   `json:"product_sku_id"`
	MemberId          int64   `json:"member_id"`
	Quantity          int     `json:"quantity"`
	Price             float32 `json:"price"`
	ProductPic        string  `json:"product_pic"`
	ProductName       string  `json:"product_name"`
	ProductSubTitle   string  `json:"product_sub_title"`
	ProductSkuCode    string  `json:"product_sku_code"`
	MemberNickname    string  `json:"member_nickname"`
	CreateDate        string  `json:"create_date"`
	ModifyDate        string  `json:"modify_date"`
	DeleteStatus      int     `json:"delete_status"`
	ProductCategoryId int64   `json:"product_category_id"`
	ProductBrand      string  `json:"product_brand"`
	ProductSn         string  `json:"product_sn"`
	ProductAttr       string  `json:"product_attr"`
}

type UpdateQuantityRequestStruct struct {
	Id       int64 `json:"id"`
	Quantity int   `json:"quantity"`
}
