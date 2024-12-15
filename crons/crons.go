package crons

import (
	"github.com/robfig/cron/v3"
	"gomall/Init/rabbitmqConsumer"
)

var job *cron.Cron

// 启动定时任务和rabbitmq的消费者
func InitCrons() {
	//1.启动消费者
	go rabbitmqConsumer.StartDelayConsumer()
	go rabbitmqConsumer.StartUpdateuserinfoConsumer()
	//2.启动定时任务.貌似没有必要定时取消了，因为这个任务已经交给了rabbitmq的消费者去做
	//job = cron.New(cron.WithSeconds())
	//
	////每10分钟查找并删除超时订单
	//_, err := job.AddFunc(consts.CRON_EVERY_10MINS, func() {
	//	_, err := orderLogic.CancelTimeOutOrder(12)
	//	if err != nil {
	//		logrus.Errorf("执行CancelTimeOutOrder定时任务出错：%s", err.Error())
	//	}
	//})
	//log.Println("启动CancelTimeOutOrder定时任务")
	//if err != nil {
	//	logrus.Errorf("添加定时任务failed：%s", err.Error())
	//}
	//
	//job.Start()
}
