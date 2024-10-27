package consts

const (

	/*
		RegEmailVerCode	注册验证码
		RegEmailVerCodeByForget 找回密码验证码
	*/
	RegEmailVerCode                       = "regEmailVerCode"
	RegEmailVerCodeByForget               = "regEmailVerCodeByForget"
	EmailVerificationCodeByChangePassword = "emailVerificationCodeByChangePassword"

	/*
		RegEmailVerCode	注册验证码
		RegEmailVerCodeByForget 找回密码验证码
	*/
	LiveRoomHistoricalBarrage = "liveRoomHistoricalBarrage_"

	//VideoWatchByID 观看视频
	VideoWatchByID   = "videoWatchBy_"
	ArticleWatchByID = "articleWatchBy_"

	// 推荐视频列表信息
	RecommendVideosList = "RecommendVideosList"
	HeatestVideo        = "HeatestVideo"

	//用户token后缀
	TokenString = "tokenString"

	//定时发布视频prefix

	//quartz连接池初始化调度器数
	QuartzPoolSize = 10

	//开始时间戳
	BeginTimestamp = 1640995200
	//全局增长前缀
	GlobalIdPrefix = "globalId" + "_Incr"

	//用户信息的前缀
	RedisDatabase    = "MALL"
	REDIS_KEY_MEMBER = "USERINFO"
	// REDIS_KEY_AUTH_CODE 验证码前缀
	REDIS_KEY_AUTH_CODE = "AUTH_CODE"

	REDIS_KEY_ORDER_ID = "oms:orderId"
)
