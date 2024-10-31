package brandAttention

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gomall/global"
	"gomall/models/brandAttention"
	"log"
)

type MemberBrandAttentionRepository struct {
	collection *mongo.Collection
}

func NewMemberBrandAttentionRepository(db *mongo.Database, collectionName string) *MemberBrandAttentionRepository {
	return &MemberBrandAttentionRepository{
		collection: db.Collection(collectionName),
	}
}

func (repo *MemberBrandAttentionRepository) FindByMemberIdAndProductId(ctx context.Context, memberId int64, brandId int64) (*brandAttention.MemberBrandAttention, error) {
	filter := bson.M{"memberId": memberId, "brandId": brandId}
	var result brandAttention.MemberBrandAttention
	if err := repo.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		global.Logger.Errorf("从mongodb查询用户关注品牌的记录出错:%v", err)
		return nil, err
	}
	return &result, nil
}

func (repo *MemberBrandAttentionRepository) Save(ctx context.Context, collection *brandAttention.MemberBrandAttention) (err error) {
	// 设置过滤条件，确保同一 (memberId, productId) 组合不会重复插入
	filter := bson.M{
		"brandId": collection.BrandId,
	}
	// 使用 upsert 选项：如果存在则更新，不存在则插入
	update := bson.M{
		"$set": collection,
	}
	opts := options.Update().SetUpsert(true)
	//执行插入
	_, err = repo.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error saving document: %v\n", err)
		return err
	}

	return nil
}
