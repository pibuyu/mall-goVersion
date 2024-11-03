package crons

import (
	"github.com/sirupsen/logrus"
	"gomall/consts"
	"gomall/global"
	orderLogic "gomall/logic/order"
	"log"

	"errors"
	"github.com/robfig/cron/v3"
)

var job *cron.Cron

func InitCrons() {
	job = cron.New(cron.WithSeconds())

	//每10分钟查找并删除超时订单
	_, err := job.AddFunc(consts.CRON_EVERY_10MINS, func() {
		_, err := orderLogic.CancelTimeOutOrder(12)
		if err != nil {
			logrus.Errorf("执行CancelTimeOutOrder定时任务出错：%s", err.Error())
		}
	})
	log.Println("启动CancelTimeOutOrder定时任务")
	if err != nil {
		logrus.Errorf("添加定时任务failed：%s", err.Error())
	}

	job.Start()
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
