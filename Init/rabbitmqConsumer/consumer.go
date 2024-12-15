package rabbitmqConsumer

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"gomall/consts"
	"gomall/global"
	orderLogic "gomall/logic/order"
	"gomall/models/users"
	"gorm.io/gorm"
	"log"
	"strconv"
)

func connectRabbitMQ() (*amqp.Connection, error) {
	return amqp.Dial("amqp://mall:mall@101.126.144.39:5672/mall")
}

// 这个消费者不能和测试类里的消费者一起启动，不然会导致测试类的消费者把消息提前消费了
func StartDelayConsumer() {
	conn, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("连接 RabbitMQ 失败: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法打开通道: %v", err)
	}
	defer ch.Close()

	// 创建死信队列
	dlqName := "dlq"
	_, err = ch.QueueDeclare(
		dlqName, // 队列名称
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Fatalf("声明死信队列失败: %v", err)
	}

	// 消费死信队列中的消息
	go func() {
		dlqMsgs, err := ch.Consume(
			dlqName, // 死信队列名称
			"",      // 消费者名称
			false,   // auto-ack
			false,   // exclusive
			false,   // no-local
			false,   // no-wait
			nil,     // args
		)
		if err != nil {
			log.Fatalf("注册死信队列消费者失败: %v", err)
		}

		for msg := range dlqMsgs {
			strId := string(msg.Body)
			orderId, _ := strconv.ParseInt(strId, 10, 64)
			if err = orderLogic.CancelOrder(orderId); err != nil {
				global.Logger.Errorf("取消id=%d的订单失败:%v", orderId, err)
				if err = msg.Nack(false, true); err != nil {
					global.Logger.Errorf("消息重回dlqName队列失败:%v", orderId, err)
				}
			}
			global.Logger.Infof("取消了id为%d的订单，因为其超时未支付", orderId)

			// 确认消息
			if err = msg.Ack(false); err != nil {
				global.Logger.Errorf("消费消息出错:%v", err)
			}
		}
	}()

	// 保持程序运行
	select {}
}

func StartUpdateuserinfoConsumer() {
	conn, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("连接 RabbitMQ 失败: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("无法打开通道: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		consts.UPDATE_USER_INFO.Name, // 队列名称
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		log.Fatalf("声明UPDATE_USER_INFO队列失败: %v", err)
	}

	// 消费死信队列中的消息
	go func() {
		msgs, err := ch.Consume(
			consts.UPDATE_USER_INFO.Name, // 队列名称
			"",                           // 消费者名称
			false,                        // auto-ack
			false,                        // exclusive
			false,                        // no-local
			false,                        // no-wait
			nil,                          // args
		)
		if err != nil {
			global.Logger.Errorf("注册UPDATE_USER_INFO队列消费者失败: %v", err)
		}

		for msg := range msgs {
			var userInfo UserInfo
			if err = json.Unmarshal(msg.Body, &userInfo); err != nil {
				global.Logger.Errorf("解析消息体出错:%v", err)
				if err = msg.Nack(false, true); err != nil {
					global.Logger.Errorf("使用 Nack 将消息返回队列重新消费失败:%v", err)
				}
				continue
			}
			if err := global.Db.Model(&users.User{}).Debug().Where("id = ?", userInfo.MemberId).
				UpdateColumn("integration", gorm.Expr("integration + ?", userInfo.Integration)).
				UpdateColumn("growth", gorm.Expr("growth + ?", userInfo.Growth)).Error; err != nil {
				global.Logger.Errorf("更新用户积分和成长值出错:%v", err)
				// 使用 Nack 将消息返回队列重新消费
				if err := msg.Nack(false, true); err != nil {
					global.Logger.Errorf("使用 Nack 将消息返回队列重新消费失败:%v", err)
				}
				continue
			}
			// 确认消息
			if err = msg.Ack(false); err != nil {
				global.Logger.Errorf("ack UPDATE_USER_INFO消息出错:%v", err)
				continue
			}
			global.Logger.Infof("将用户%d的积分增加了%d,成长值增加了%d", userInfo.MemberId, userInfo.Integration, userInfo.Growth)
		}
	}()

	// 保持程序运行
	select {}
}

type UserInfo struct {
	MemberId    int64 `json:"memberId"`
	Integration int   `json:"integration"`
	Growth      int   `json:"growth"`
}
