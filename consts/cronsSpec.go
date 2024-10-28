package consts

const (
	//每10秒一次
	CRON_EVERY_10S = "*/10 * * * * *"
	//每1分钟一次
	CRON_EVERY_1MIN = "@every 1m"
	//每天零点
	CRON_EVERYDAY_MIDNIGHT = "@midnight"
	//每10分钟一次，扫描超时订单并取消
	CRON_EVERY_10MINS = "0 */10 * * * *"
)
