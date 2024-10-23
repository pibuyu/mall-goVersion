package global

import (
	"github.com/go-redis/redis"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"gomall/global/config"
	"gomall/global/database/mysql"
	RedisDbFun "gomall/global/database/redis"
	log "gomall/global/logrus"
	"gomall/global/msgQueue"
	"gorm.io/gorm"
)

// 在这里执行那些实例化的init函数，然后返回预定义的对象呗
func init() {
	Logger = log.ReturnsInstance()
	RedisDb = RedisDbFun.ReturnsInstance()
	Db = mysql.ReturnsInstance()
	Config = config.ReturnsInstance()
	//普通队列的生产者
	NormalProducer = msgQueue.ReturnsNormalInstance()
	//延迟队列的生产者
	DelayProducer = msgQueue.ReturnsDelayInstance()
	//启动延时队列的消费者和普通队列的消费者
	//todo:尝试换个地方启动消费者,避免循环依赖问题
	//msgQueue.StartDelayConsumer()
	//msgQueue.StartNormalConsumer()
}

var (
	Logger         *logrus.Logger
	Config         *config.Info
	Db             *gorm.DB
	RedisDb        *redis.Client
	NormalProducer *kafka.Conn
	DelayProducer  *kafka.Conn
)
