package rabbitmqConsumer

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gomall/global"
	orderLogic "gomall/logic/order"
	"log"
	"strconv"
)

func connectRabbitMQ() (*amqp.Connection, error) {
	return amqp.Dial("amqp://mall:mall@101.126.144.39:5672/mall")
}

// todo:This is a special comment.这个消费者不能和测试类里的消费者一起启动，不然会导致测试类的消费者把消息提前消费了
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
			if err := orderLogic.CancelOrder(orderId); err != nil {
				global.Logger.Errorf("取消id=%d的订单失败")
			}
			global.Logger.Infof("取消了id为%d的订单，因为其超时未支付")

			// 确认消息
			msg.Ack(false)
		}
	}()

	// 保持程序运行
	select {}
}
