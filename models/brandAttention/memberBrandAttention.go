package brandAttention

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MemberBrandAttention struct {
	// 对应MongoDB中的_id字段，通常作为文档的唯一标识
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 成员ID，添加了索引注释
	MemberId int64 `bson:"memberId" json:"memberId"`

	// 成员昵称
	MemberNickname string `bson:"memberNickname" json:"memberNickname"`

	// 成员头像链接或路径等信息
	MemberIcon string `bson:"memberIcon" json:"memberIcon"`

	// 品牌ID，添加了索引注释
	BrandId int64 `bson:"brandId" json:"brandId"`

	// 品牌名称
	BrandName string `bson:"brandName" json:"brandName"`

	// 品牌logo链接或路径等信息
	BrandLogo string `bson:"brandLogo" json:"brandLogo"`

	// 品牌所在城市
	BrandCity string `bson:"brandCity" json:"brandCity"`

	// 创建时间
	CreateTime time.Time `bson:"createTime" json:"createTime"`
}
