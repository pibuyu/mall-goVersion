package mongodb

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-retry"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gomall/global/config"
	"log"
	"time"
)

var Db *mongo.Client

func ReturnsInstance() *mongo.Client {
	var mongoConfig = config.Config.MongoDBConfig
	b := retry.NewFibonacci(10 * time.Second)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var url = fmt.Sprintf("mongodb://%s:%s@%s:%d", mongoConfig.User, mongoConfig.Password, mongoConfig.Host, mongoConfig.Port)

	if err := retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		clientOptions := options.Client().ApplyURI(url)
		Db, _ = mongo.NewClient(clientOptions)
		//尝试连接到MongoDB
		if err := Db.Connect(ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("重试 5 次后仍然无法连接 MongoDB，请排查 MongoDB 服务端是否启动/配置信息是否正确，错误信息为： %v \n", err)
	}
	return Db
}
