package readHistory

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gomall/global"
	readHistoryModels "gomall/models/history"
	"time"
)

type MemberReadHistoryRepository struct {
	collection *mongo.Collection
}

func NewMemberReadHistoryRepository(db *mongo.Database, collectionName string) *MemberReadHistoryRepository {
	return &MemberReadHistoryRepository{
		collection: db.Collection(collectionName),
	}
}

func (repo *MemberReadHistoryRepository) FindByMemberIdOrderByCreateTimeDesc(ctx context.Context, pageNum int, pageSize int, memberId int64) (result []readHistoryModels.MemberReadHistoryMongoStruct, err error) {
	filter := bson.M{"memberId": memberId}

	// 使用 Find 方法进行查询
	cursor, err := repo.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{"createTime", -1}}).SetSkip(int64((pageNum-1)*pageSize)).SetLimit(int64(pageSize)))
	if err != nil {
		global.Logger.Errorf("从mongodb查询用户关注品牌的记录出错:%v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// 解码结果到切片中
	for cursor.Next(ctx) {
		var history readHistoryModels.MemberReadHistoryMongoStruct
		if err := cursor.Decode(&history); err != nil {
			global.Logger.Errorf("解码浏览记录出错:%v", err)
			return nil, err
		}
		result = append(result, history)
	}

	if err := cursor.Err(); err != nil {
		global.Logger.Errorf("游标出错:%v", err)
		return nil, err
	}

	return result, nil
}

func (repo *MemberReadHistoryRepository) Clear(ctx context.Context, memberId int64) error {
	// 设置过滤条件，仅删除 memberId 等于指定值的记录
	filter := bson.M{"memberId": memberId}

	// 执行删除操作
	_, err := repo.collection.DeleteMany(ctx, filter)
	if err != nil {
		global.Logger.Errorf("清除用户id=%d的浏览记录失败: %v\n", memberId, err)
		return err
	}
	return nil
}

func (repo *MemberReadHistoryRepository) DeleteByIds(ctx context.Context, ids []int64, memberId int64) error {
	// 设置过滤条件，删除 memberId 等于指定值，并且 ID 在 ids 列表中的记录
	filter := bson.M{
		"memberId": memberId,
		"_id":      bson.M{"$in": convertIdsToObjectIDs(ids)},
	}

	// 执行删除操作
	_, err := repo.collection.DeleteMany(ctx, filter)
	if err != nil {
		global.Logger.Errorf("根据 IDs 清除用户 id=%d 的浏览记录失败: %v\n", memberId, err)
		return err
	}
	return nil
}

func convertIdsToObjectIDs(ids []int64) []primitive.ObjectID {
	objectIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectIDs[i] = primitive.NewObjectIDFromTimestamp(time.Unix(id, 0))
	}
	return objectIDs
}
