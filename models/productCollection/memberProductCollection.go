package productCollection

import "time"

type MemberProductCollection struct {
	ID              string    `bson:"_id,omitempty"`
	MemberId        int64     `bson:"memberId"`
	MemberNickname  string    `bson:"memberNickname"`
	MemberIcon      string    `bson:"memberIcon"`
	ProductId       int64     `bson:"productId"`
	ProductName     string    `bson:"productName"`
	ProductPic      string    `bson:"productPic"`
	ProductSubTitle string    `bson:"productSubTitle"`
	ProductPrice    float32   `bson:"productPrice"`
	CreateTime      time.Time `bson:"createTime"`
}
