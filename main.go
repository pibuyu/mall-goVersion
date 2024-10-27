package main

import (
	"gomall/router"
)

func main() {
	//启动定时器和kafka的消费者
	//crons.InitCrons()//todo:先把消费者关掉，没啥用

	//启动路由器
	router.InitRouter()
}
