package productCollection

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gomall/global"
	"gomall/models/productCollection"
	"log"
)

type MemberProductCollectionRepository struct {
	collection *mongo.Collection
}

func NewMemberProductCollectionRepository(db *mongo.Database, collectionName string) *MemberProductCollectionRepository {
	return &MemberProductCollectionRepository{
		collection: db.Collection(collectionName),
	}
}

func (repo *MemberProductCollectionRepository) FindByMemberIdAndProductId(ctx context.Context, memberId int64, productId int64) (*productCollection.MemberProductCollection, error) {
	filter := bson.M{"memberId": memberId, "productId": productId}
	var result productCollection.MemberProductCollection
	if err := repo.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		global.Logger.Errorf("从mongodb查询记录出错:%v", err)
		return nil, err
	}
	return &result, nil
}

func (repo *MemberProductCollectionRepository) Save(ctx context.Context, collection *productCollection.MemberProductCollection) (err error) {
	// 设置过滤条件，确保同一 (memberId, productId) 组合不会重复插入
	filter := bson.M{
		"memberId":  collection.MemberId,
		"productId": collection.ProductId,
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
