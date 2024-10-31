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

func (repo *MemberBrandAttentionRepository) ClearByMemberId(ctx context.Context, memberId int64) error {
	// 设置过滤条件，仅删除 memberId 等于指定值的记录
	filter := bson.M{"memberId": memberId}

	// 执行删除操作
	_, err := repo.collection.DeleteMany(ctx, filter)
	if err != nil {
		global.Logger.Errorf("Error clearing documents with memberId %d: %v\n", memberId, err)
		return err
	}
	return nil
}

func (repo *MemberBrandAttentionRepository) Delete(ctx context.Context, brandId int64, memberId int64) error {
	// 设置过滤条件，仅删除 memberId 等于指定值的记录
	filter := bson.M{"memberId": memberId, "brandId": brandId}

	// 执行删除操作
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		global.Logger.Errorf("Error deleting documents with memberId %d: %v\n", memberId, err)
		return err
	}
	return nil
}

func (repo *MemberBrandAttentionRepository) Detail(ctx context.Context, brandId int64, memberId int64) (*brandAttention.MemberBrandAttention, error) {
	// 设置过滤条件，仅删除 memberId 等于指定值的记录
	filter := bson.M{"memberId": memberId, "brandId": brandId}

	var result brandAttention.MemberBrandAttention

	if err := repo.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			// 如果没有找到文档，返回 nil 而不是错误
			return nil, nil
		}
		// 其他查询错误
		global.Logger.Errorf("Error finding product collection with memberId %d and productId %d: %v\n", memberId, brandId, err)
		return nil, err
	}
	return &result, nil
}

func (repo *MemberBrandAttentionRepository) List(ctx context.Context, pageNum int, pageSize int, memberId int64) ([]brandAttention.MemberBrandAttention, error) {
	// 计算分页的起始位置
	skip := int64((pageNum - 1) * pageSize)
	// 设置过滤条件
	filter := bson.M{"memberId": memberId}
	// 设置查询选项：排序、跳过和限制条数
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{"createTime", -1}}) // 按创建时间降序排序

	// 查询数据库
	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []brandAttention.MemberBrandAttention
	for cursor.Next(ctx) {
		// 将当前文档解码为 bson.D 格式的原始数据
		var oneResult brandAttention.MemberBrandAttention
		if err := cursor.Decode(&oneResult); err != nil {
			return nil, err
		}
		results = append(results, oneResult)
	}
	return results, nil
}
