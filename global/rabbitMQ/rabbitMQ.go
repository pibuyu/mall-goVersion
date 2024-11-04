package rabbitMQ

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sethvargo/go-retry"
	"gomall/global/config"
	"log"
	"time"
)

var rabbitMQConn *amqp.Connection

func ReturnInstance() *amqp.Connection {
	var err error
	var rabbitMQConfig = config.Config.RabbitMQConfig // 确保你的配置文件包含 RabbitMQ 的相关配置

	b := retry.NewFibonacci(10 * time.Second) //重试的斐波那契机制，最大重试间隔时间为10秒
	ctx := context.Background()
	if err := retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		rabbitMQConn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
			rabbitMQConfig.User,
			rabbitMQConfig.Password,
			rabbitMQConfig.Host,
			rabbitMQConfig.Port,
			rabbitMQConfig.Vhost, // 确保你的配置中有这个字段
		))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		// 多次重连仍旧失败
		log.Fatalf("重试后仍然无法连接RabbitMQ，请排查RabbitMQ服务端是否启动/配置信息是否正确，错误信息为： %v\n", err)
	}

	return rabbitMQConn
}
