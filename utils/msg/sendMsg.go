package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gomall/global"
)

func SendToQueue(queueName string, messageBody map[string]interface{}) error {
	ch, err := global.RabbitMQ.Channel()
	if err != nil {
		global.Logger.Errorf("无法创建channel：%v", err)
		return errors.New(fmt.Sprintf("无法创建channel：%v", err))
	}
	defer ch.Close()

	// 声明队列 (确保队列存在)
	_, err = ch.QueueDeclare(
		queueName, // 队列名称
		true,      // durable: 队列持久化
		false,     // auto-delete: 不自动删除
		false,     // exclusive: 非独占
		false,     // no-wait: 无等待
		nil,       // 额外参数
	)
	if err != nil {
		global.Logger.Errorf("声明队列失败: %v", err)
		return errors.New(fmt.Sprintf("声明队列失败: %v", err))
	}

	// 将 map 转换为 JSON 字符串（[]byte）
	messageBytes, err := json.Marshal(messageBody)
	if err != nil {
		global.Logger.Errorf("消息转换为JSON失败: %v", err)
		return errors.New(fmt.Sprintf("消息转换为JSON失败: %v", err))
	}
	// 将消息发送到队列
	err = ch.Publish(
		"",        // 默认交换机
		queueName, // 队列名称
		false,     // mandatory: 如果队列不可路由则返回消息
		false,     // immediate: 如果消息不能立即被处理，则返回
		amqp.Publishing{
			ContentType:  "text/plain", // 消息类型
			Body:         messageBytes, // 消息内容
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		global.Logger.Error("发送消息到队列失败: %v", err)
		return errors.New(fmt.Sprintf("发送消息到队列失败: %v", err))
	}

	global.Logger.Infof("成功将消息：%s 发送到队列 → %s", string(messageBytes), queueName)
	return nil
}
