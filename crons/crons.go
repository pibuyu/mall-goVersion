package crons

import (
	"gomall/Init/kafkaConsumer"
	"gomall/global"

	"errors"
	"github.com/robfig/cron/v3"
)

var job *cron.Cron

func InitCrons() {

	//在这里启动消费者
	kafkaConsumer.StartNormalConsumer()
	kafkaConsumer.StartDelayConsumer()

}

// AddTask 添加定时任务
func AddTask(spec string, task func()) error {
	_, err := job.AddFunc(spec, task)
	if err != nil {
		global.Logger.Errorf("添加定时任务出错：%s", err.Error())
		return errors.New("添加定时任务出错：" + err.Error())
	}
	return nil
}
