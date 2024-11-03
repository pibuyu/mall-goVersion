package kafkaConsumer

import (
	"context"
	"gomall/global/config"

	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var msgConfig = config.Config.KafkaConfig

func StartDelayConsumer() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("启动消费者出现错误:%v", r)
		}
	}()

	delayReader := kafka.NewReader(kafka.ReaderConfig{

		Brokers:     []string{msgConfig.Server},
		Topic:       msgConfig.DelayTopic,
		StartOffset: kafka.FirstOffset,
		GroupID:     "delay-consumer",
	})

	log.Println("创建监听延时队列的DelayConsumer成功")

	//normalWriter, err := kafka.DialLeader(context.Background(), "tcp", consts.KafkaServerAddr, consts.KafkaTopic, 0)
	normalWriter, err := kafka.DialLeader(context.Background(), "tcp", msgConfig.Server, msgConfig.NormalTopic, 0)
	if err != nil {
		log.Fatalf("创建写入到普通队列的生产者失败：%v", err)
	}
	//delayWriter, err := kafka.DialLeader(context.Background(), "tcp", consts.KafkaServerAddr, consts.DelayQueue, 0)
	delayWriter, err := kafka.DialLeader(context.Background(), "tcp", msgConfig.Server, msgConfig.DelayTopic, 0)
	if err != nil {
		log.Fatalf("创建写入到延时队列的生产者失败：%v", err)
	}
	//监听延时队列的消息，到达处理时间则写入正常队列，否则回写到延时队列中去
	go func() {
		for {
			message, err := delayReader.ReadMessage(context.Background())
			if err != nil {
				//打印一下就行，然后继续读取消息
				log.Printf("读取延时队列的消息出错：%v", err)
				continue
			}
			log.Printf("监听到延时队列的消息：%s", message.Value)

			if time.Now().After(message.Time) {
				//写入到正常队列中去
				_, err := normalWriter.WriteMessages(kafka.Message{Value: message.Value})
				if err != nil {
					log.Printf("延时队列中的消息转写入普通队列出错：%v", err)
					//此时跳过提交偏移量，让剩下的消费者继续消费
					continue
				}
				//提交偏移量，kafka才知道消息被消费过了
				if err := delayReader.CommitMessages(context.Background(), message); err != nil {
					log.Printf("延时队列提交偏移量失败：%v", err)
				}
			} else {
				//写回到延时队列中去

				_, err := delayWriter.WriteMessages(message)
				if err != nil {
					log.Printf("延时队列消息回写失败：%v", err)
				}
			}
		}
	}()
}

// StartNormalConsumer 启动普通队列(即时处理)的消费者
func StartNormalConsumer() {
	normalReader := kafka.NewReader(kafka.ReaderConfig{
		//Brokers:     []string{consts.KafkaServerAddr},
		//Topic:       consts.KafkaTopic,
		Brokers:     []string{msgConfig.Server},
		Topic:       msgConfig.NormalTopic,
		StartOffset: kafka.LastOffset,
		GroupID:     "normal-consumer",
	})

	log.Println("创建监听普通队列的NormalConsumer成功")

	go func() {
		for {
			message, err := normalReader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("读取普通队列的消息错误：%v", err)
				continue
			}
			//messsageBody := string(message.Value)

			//处理定时发布视频的消息,此类消息形如:publishVideo_1234

			//处理正常队列的消息
			log.Printf("监听到了普通队列的消息：%s", string(message.Value))
			//提交偏移量
			if err := normalReader.CommitMessages(context.Background(), message); err != nil {
				log.Printf("普通队列提交偏移量失败：%v", err)
			}
		}
	}()
}
