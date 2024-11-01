package history

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// mongodb查询的返回类
type MemberReadHistoryMongoStruct struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`                // 使用 ObjectID 类型作为 MongoDB 的 ID
	MemberID        int64              `bson:"memberId" json:"memberId"`               // 会员 ID
	MemberNickname  string             `bson:"memberNickname" json:"memberNickname"`   // 会员昵称
	MemberIcon      string             `bson:"memberIcon" json:"memberIcon"`           // 会员头像
	ProductID       int64              `bson:"productId" json:"productId"`             // 产品 ID
	ProductName     string             `bson:"productName" json:"productName"`         // 产品名称
	ProductPic      string             `bson:"productPic" json:"productPic"`           // 产品图片
	ProductSubTitle string             `bson:"productSubTitle" json:"productSubTitle"` // 产品副标题
	ProductPrice    string             `bson:"productPrice" json:"productPrice"`       // 产品价格
	CreateTime      time.Time          `bson:"createTime" json:"createTime"`           // 创建时间
}
