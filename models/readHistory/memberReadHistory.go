package readHistory

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MemberReadHistory struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MemberID        int64              `bson:"memberId" json:"memberId"`
	MemberNickname  string             `bson:"memberNickname" json:"memberNickname"`
	MemberIcon      string             `bson:"memberIcon" json:"memberIcon"`
	ProductID       int64              `bson:"productId" json:"productId"`
	ProductName     string             `bson:"productName" json:"productName"`
	ProductPic      string             `bson:"productPic" json:"productPic"`
	ProductSubTitle string             `bson:"productSubTitle" json:"productSubTitle"`
	ProductPrice    string             `bson:"productPrice" json:"productPrice"`
	CreateTime      time.Time          `bson:"createTime" json:"createTime"`
}
