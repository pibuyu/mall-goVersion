package main

import (
	"gomall/crons"
	"gomall/router"
)

func main() {
	//启动定时器和kafka的消费者
	crons.InitCrons()

	//启动路由器
	router.InitRouter()
}
