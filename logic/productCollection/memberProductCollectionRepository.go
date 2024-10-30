package productCollection

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gomall/global"
	"gomall/models/productCollection"
	"log"
	"math"
	"strconv"
	"time"
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

func (repo *MemberProductCollectionRepository) ClearByMemberId(ctx context.Context, memberId int64) error {
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

func (repo *MemberProductCollectionRepository) Delete(ctx context.Context, productId int64, memberId int64) error {
	// 设置过滤条件，仅删除 memberId 等于指定值的记录
	filter := bson.M{"memberId": memberId, "productId": productId}

	// 执行删除操作
	_, err := repo.collection.DeleteOne(ctx, filter)
	if err != nil {
		global.Logger.Errorf("Error deleting documents with memberId %d: %v\n", memberId, err)
		return err
	}
	return nil
}

func (repo *MemberProductCollectionRepository) Detail(ctx context.Context, productId int64, memberId int64) (*productCollection.MemberProductCollection, error) {
	// 设置过滤条件，仅删除 memberId 等于指定值的记录
	filter := bson.M{"memberId": memberId, "productId": productId}

	// 存储原始的查询结果
	var rawResult bson.D

	// 执行查询
	if err := repo.collection.FindOne(ctx, filter).Decode(&rawResult); err != nil {
		if err == mongo.ErrNoDocuments {
			// 如果没有找到文档，返回 nil 而不是错误
			return nil, nil
		}
		// 其他查询错误
		global.Logger.Errorf("Error finding product collection with memberId %d and productId %d: %v\n", memberId, productId, err)
		return nil, err
	}

	global.Logger.Infof("查询收藏的detail信息为:%v", rawResult)

	// todo:This is a special comment.
	//  自动映射的话会把productPrice映射成string类型，但我们定义的是float32类型，所以需要手动映射。
	//手动映射
	result := &productCollection.MemberProductCollection{}
	for _, elem := range rawResult {
		switch elem.Key {
		case "memberId":
			result.MemberId, _ = elem.Value.(int64)
		case "memberNickname":
			result.MemberNickname, _ = elem.Value.(string)
		case "memberIcon":
			result.MemberIcon, _ = elem.Value.(string)
		case "productId":
			result.ProductId, _ = elem.Value.(int64)
		case "productName":
			result.ProductName, _ = elem.Value.(string)
		case "productPic":
			result.ProductPic, _ = elem.Value.(string)
		case "productSubTitle":
			result.ProductSubTitle, _ = elem.Value.(string)
		case "productPrice":
			if v, ok := elem.Value.(string); ok {
				// 将productPrice转换为float32
				//todo:这里需要保留2位小数，改来改去都不对
				if price, err := strconv.ParseFloat(v, 32); err == nil {
					price = math.Round(price*100) / 100
					result.ProductPrice = float32(price)
				} else {
					global.Logger.Errorf("price转换为float32解析出错:%v", err)
				}
			}
		case "createTime":
			result.CreateTime, _ = elem.Value.(time.Time)
		}
	}
	// 返回查询结果
	return result, nil
}

func (repo *MemberProductCollectionRepository) List(ctx context.Context, pageNum int, pageSize int, memberId int64) ([]productCollection.MemberProductCollection, error) {
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

	// 解析查询结果
	var results []productCollection.MemberProductCollection
	for cursor.Next(ctx) {
		// 将当前文档解码为 bson.D 格式的原始数据
		var rawResult bson.D
		if err := cursor.Decode(&rawResult); err != nil {
			return nil, err
		}
		// 手动映射字段
		var collection productCollection.MemberProductCollection
		for _, elem := range rawResult {
			switch elem.Key {
			case "memberId":
				if v, ok := elem.Value.(int64); ok {
					collection.MemberId = v
				}
			case "memberNickname":
				if v, ok := elem.Value.(string); ok {
					collection.MemberNickname = v
				}
			case "memberIcon":
				if v, ok := elem.Value.(string); ok {
					collection.MemberIcon = v
				}
			case "productId":
				if v, ok := elem.Value.(int64); ok {
					collection.ProductId = v
				}
			case "productName":
				if v, ok := elem.Value.(string); ok {
					collection.ProductName = v
				}
			case "productPic":
				if v, ok := elem.Value.(string); ok {
					collection.ProductPic = v
				}
			case "productSubTitle":
				if v, ok := elem.Value.(string); ok {
					collection.ProductSubTitle = v
				}
			case "productPrice":
				if v, ok := elem.Value.(string); ok {
					// 将字符串转换为 float32 类型，保留两位小数
					if price, err := strconv.ParseFloat(v, 32); err == nil {
						collection.ProductPrice = float32(math.Round(price*100) / 100) // 保留两位小数
					} else {
						return nil, fmt.Errorf("productPrice 字段解析出错: %v", err)
					}
				}
			case "createTime":
				if v, ok := elem.Value.(primitive.DateTime); ok {
					collection.CreateTime = v.Time()
				}
			}
		}
		results = append(results, collection)
	}

	// 检查游标是否遇到错误
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
