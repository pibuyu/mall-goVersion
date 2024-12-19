package payment

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
	"gomall/consts"
	controller "gomall/controllers"
	"gomall/global"
	"gomall/models/order"
	"gomall/utils/uniqueid"
	"net/http"
	"strconv"
	"time"
)

type PaymentController struct {
	controller.BaseControllers
}

var (
	avoid_repeat_payment_script = `
        local lock_key=KEYS[1]
		local lock_value=ARGV[1] --锁的持有者
		local ttl=tonumber(ARGV[2])

		--若锁不存在，则新增锁，并设置锁重入计数为1、设置锁的过期时间
		if redis.call('exists',lock_key)==0 then
			redis.call('set',lock_key,lock_value,'PX',ttl)
			return 'LOCK_ACQUIRED'
		else
			--否则表明锁被其他线程占用，当前线程无权解锁，返回未获取到锁即可
			return 'LOCK_UNACQUIRED'
		end
	`

	release_dis_lock_script = `
		local lock_key=KEYS[1]
		local lock_holder=ARGV[1] --锁的持有者
			
		--若锁不存在，则返回释放成功
		if redis.call('exists',lock_key)==0 then
			return 'LOCK_RELEASED'
		else
			--检查锁的持有者和当前试图解锁的进程是否一致
			if redis.call('get',lock_key)~=lock_holder then
				return 'LOCK_UNLEASED'
			else
				redis.call('del',lock_key)
				return 'LOCK_RELEASED'
			end
		end
	`
)

var (
	ORDERSN_TO_OUTTRADENO_MAP             = "ORDERSN_TO_OUTTRADENO_MAP_"
	SERVER_BUSY_ERROR                     = "服务器繁忙，请稍后再试"
	INSERT_ORDERSN_TO_OUTTRADENO_MAP_LOCK = "INSERT_ORDERSN_TO_OUTTRADENO_MAP_LOCK" //插入orderSn到outTradeNo的分布式锁前缀
)

func (c *PaymentController) WebPay(ctx *gin.Context) {
	var err error
	orderSn := ctx.Query("outTradeNo")      //订单号,可以拿订单号和ctx获取的memberId倒查订单的其他信息
	subject := ctx.Query("subject")         //订单注释字段
	totalAmount := ctx.Query("totalAmount") //应付金额

	//1.先尝试获取到缓存的支付网址
	result, err := global.RedisDb.Get(fmt.Sprintf("%s%s", ORDERSN_TO_OUTTRADENO_MAP, orderSn)).Result()
	if err != nil && err != redis.Nil {
		global.Logger.Errorf("从redis中获取orderSn到outTradeNo的映射关系failed:%v", err)
		c.Response(ctx, SERVER_BUSY_ERROR, nil, err)
		return
	}
	if len(result) != 0 {
		global.Logger.Infof("订单 %s 的支付请求命中缓存，直接返回支付链接,无需生成outTradeNo", orderSn)
		ctx.Redirect(http.StatusTemporaryRedirect, result)
		return
	}

	//2.如果缓存中没有支付网址的话，竞争分布式锁，然后将orderSn到outTradeNo的映射关系写入到redis中去
	lockKey := fmt.Sprintf("%s%s", INSERT_ORDERSN_TO_OUTTRADENO_MAP_LOCK, orderSn)
	lockValue := uuid.New().String() //锁的value应该为表示锁持有者的的唯一标识，这里用uuid
	ttl := 1000                      //锁的过期时间，单位为ms
	acquireLockRes, err := global.RedisDb.Eval(avoid_repeat_payment_script, []string{lockKey}, lockValue, ttl).Result()
	if err != nil {
		global.Logger.Errorf("订单: %s 尝试支付，获取分布式锁出错: %v", orderSn, err)
		c.Response(ctx, "获取分布式锁失败", nil, err)
		return
	}
	if acquireLockRes != "LOCK_ACQUIRED" {
		global.Logger.Infof("订单 %s 尝试请求支付，当前竞争分布式锁的唯一id为: %s ,竞争分布式锁失败，直接返回了", orderSn, lockValue)
		c.Response(ctx, "该订单正在支付中，请勿重复支付", nil, errors.New(SERVER_BUSY_ERROR))
		return
	}

	//3.获取到了分布式锁
	global.Logger.Infof("当前线程为: %s ,成功获取到了分布式锁: %s ", lockValue, lockKey)
	outTradeNo := xid.Next() //生成支付流水单号
	//3.1.生成预支付单
	//todo:支付状态发生变化时需要考虑将这个hash里订单号与memberId的映射删除,避免空间无限膨胀
	memberIdFromCtx, err := global.RedisDb.HGet(consts.Order2MemberIdMap, orderSn).Result()
	memberIdInt64, _ := strconv.ParseInt(memberIdFromCtx, 10, 64)
	if err != nil {
		global.Logger.Errorf("预支付时，查询memberId failed:%v", err)
		c.Response(ctx, "获取用户信息失败", nil, err)
		return
	}
	orderInfo := &order.OmsOrder{}
	if err := global.Db.Model(&order.OmsOrder{}).Where("member_id = ? and order_sn = ?", memberIdInt64, orderSn).Find(orderInfo).Error; err != nil {
		global.Logger.Errorf("预支付时，查询订单详情failed : %v", err)
		c.Response(ctx, "查询订单详情失败", nil, err)
		return
	}
	float64Price, _ := strconv.ParseFloat(totalAmount, 64)
	payment := &order.OmsPayment{
		OrderId:     orderInfo.ID,
		PaymentSn:   uniqueid.GenSn(uniqueid.SN_PREFIX_THIRD_PAYMENT),
		UserId:      memberIdInt64,
		OrderSn:     orderSn,
		Subject:     subject,
		OutTradeNo:  strconv.FormatInt(outTradeNo, 10),
		TotalAmount: float64Price,
	}
	if err = payment.Insert(); err != nil {
		global.Logger.Errorf("预支付时，插入预支付单失败:%v", err)
		c.Response(ctx, "预支付时，插入预支付单失败", nil, err)
		return
	}

	//3.2.构造并发起支付请求
	var p = alipay.TradePagePay{}
	//dev:这是ngrok配置的9090端口的内网穿透地址，每次重启ngrok都要修改这里的回调地址;prod:填服务器的回调地址
	p.NotifyURL = "http://101.126.144.39:9090/alipay/callback"
	//支付完成之后跳转页面，直接填前端的订单页即可
	p.ReturnURL = "http://localhost:8060/#/pages/order/order?state=0"
	p.Subject = subject
	p.OutTradeNo = strconv.FormatInt(outTradeNo, 10)
	p.TotalAmount = totalAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY" //网页支付,这个字段是固定的，目前支付沙箱仅支持这一个模式

	url, _ := global.AliPay.TradePagePay(p)
	//3.3释放分布式锁并将支付链接放入redis中去，有效期为15分钟;同时释放分布式锁
	global.RedisDb.Set(fmt.Sprintf("%s%s", ORDERSN_TO_OUTTRADENO_MAP, orderSn), url.String(), 15*time.Minute)
	releaseLockRes, err := global.RedisDb.Eval(release_dis_lock_script, []string{lockKey}, lockValue).Result()
	if err != nil || releaseLockRes != "LOCK_RELEASED" {
		global.Logger.Errorf("释放分布式锁: %s failed", lockKey)
	}
	global.Logger.Infof("当前线程为: %s ,成功写入了缓存: %s ", lockValue, url.String())
	//跳转到支付宝支付页面
	ctx.Redirect(http.StatusTemporaryRedirect, url.String())
}

func (c *PaymentController) AliPayCallback(ctx *gin.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		global.Logger.Errorf("解析支付回调参数failed:%v", err)
		c.Response(ctx, "解析支付回调参数失败", nil, err)
		return
	}

	// 1.获取回调中的支付流水单号
	outTradeNo := ctx.Request.Form.Get("out_trade_no")
	// 2.验证签名，确保这个回调的安全性
	if err := global.AliPay.VerifySign(ctx.Request.Form); err != nil {
		global.Logger.Errorf("验证支付回调签名未通过：%v", err)
		c.Response(ctx, "验证支付回调签名未通过", nil, err)
		return
	}
	// 3.主动查询支付状态
	var p = alipay.TradeQuery{
		OutTradeNo: outTradeNo,
	}
	resp, err := global.AliPay.TradeQuery(ctx, p)
	if err != nil {
		global.Logger.Errorf("查询订单支付状态失败：%v", err)
		c.Response(ctx, "查询订单支付状态失败", nil, err)
		return
	}

	// 3.1 如果支付失败，返回告知用户结果
	if resp.TradeStatus != "TRADE_SUCCESS" {
		global.Logger.Errorf("支付流水单：%s 支付失败", outTradeNo)
		c.Response(ctx, "订单支付失败", nil, errors.New("订单支付失败"))
		return
	}

	// 4.根据outTradeNo查询本地的支付记录，比对支付的金额是否一致
	payment := &order.OmsPayment{}
	if err := global.Db.Model(&order.OmsPayment{}).Where("out_trade_no = ?", outTradeNo).Find(&payment).Error; err != nil {
		global.Logger.Errorf("根据payment_sn查询订单详情failed : %v", err)
		c.Response(ctx, "查询订单详情失败", nil, err)
		return
	}
	payAmount, _ := strconv.ParseFloat(resp.TotalAmount, 10)
	if payAmount != payment.TotalAmount {
		global.Logger.Errorf("支付流水%s 的实际支付金额= %v ,应付金额为= %v ", outTradeNo, payAmount, payment.TotalAmount)
		c.Response(ctx, "实付金额与应付金额不符", nil, err)
		return
	}

	// 5.1 支付成功,修改订单状态为已支付
	if err = global.Db.Model(&order.OmsOrder{}).Where("id = ?", payment.OrderId).UpdateColumn("status", 1).Error; err != nil {
		global.Logger.Errorf("修改订单 %d 为已支付状态失败 %v", payment.OrderId, err)
	}

	// 5.2 修改支付状态、实付金额
	payment.PayStatus = 1
	payment.ActualPay = payAmount
	if err = global.Db.Model(&order.OmsPayment{}).Where("out_trade_no = ?", payment.OutTradeNo).Updates(payment).Error; err != nil {
		global.Logger.Errorf("修改支付流水 %s 出错 %v", outTradeNo, err)
	}
	c.Response(ctx, "订单支付成功", nil, nil)
}

// 不要等结果，如果真的修改失败了通知用户即可
func updatePaymentStatus(outTradeNo string, orderId int64, payAmount float64, newStatus int) {
	var err error
	tx := global.Db.Begin()

	if err = tx.Table("oms_order").Debug().Where("id = ?", orderId).UpdateColumn("status", 1).Error; err != nil {
		global.Logger.Errorf("修改订单 : %d 状态failed:%v", orderId, err)
		tx.Rollback()
	}
	if err = tx.Table("oms_payment").Debug().Where("out_trade_no = ?", outTradeNo).UpdateColumn("pay_status", 1).UpdateColumn("actual_pay", payAmount).Error; err != nil {
		global.Logger.Errorf("修改订单 : %s 状态failed:%v", outTradeNo, err)
		tx.Rollback()
		//todo:应该丢进消息队列一直重试
	}
}
