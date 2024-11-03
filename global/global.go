package global

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gomall/global/config"
	MongoDbFun "gomall/global/database/mongodb"
	"gomall/global/database/mysql"
	RedisDbFun "gomall/global/database/redis"
	log "gomall/global/logrus"

	"gorm.io/gorm"
)

// 在这里执行那些实例化的init函数，然后返回预定义的对象呗
func init() {
	Logger = log.ReturnsInstance()
	RedisDb = RedisDbFun.ReturnsInstance()
	MongoDb = MongoDbFun.ReturnsInstance()
	Db = mysql.ReturnsInstance()
	Config = config.ReturnsInstance()
}

var (
	Logger  *logrus.Logger
	Config  *config.Info
	Db      *gorm.DB
	RedisDb *redis.Client
	MongoDb *mongo.Client
)
