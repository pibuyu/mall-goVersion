package order

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"gomall/consts"
	"gomall/global"
	receive "gomall/interaction/receive/order"
	cartLogic "gomall/logic/cart"
	"gomall/models/cart"
	"gomall/models/coupon"
	"gomall/models/integration"
	"gomall/models/order"
	orderModel "gomall/models/order"
	"gomall/models/users"
	"gomall/utils/jwt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// Detail 根据订单id获取订单详情
func Detail(orderId int64) (result *order.OmsOrderDetail, err error) {
	omsOrder := &order.OmsOrder{}
	if err = global.Db.Model(&order.OmsOrder{}).
		Where("id = ?", orderId).Find(&omsOrder).Error; err != nil {
		return nil, errors.New("查询订单详情时，获取omsOrder出错:" + err.Error())
	}

	//根据orderId再去oms_order_item表里查询出OmsOrderItem列表
	orderItemList := make([]order.OmsOrderItem, 0)
	if err = global.Db.Model(&order.OmsOrderItem{}).
		Where("order_id = ?", orderId).Find(&orderItemList).Error; err != nil {
		return nil, errors.New("查询订单详情时，获取orderItemList出错:" + err.Error())
	}
	//然后可以构造omsOrderDetail对象并返回了
	orderDetail := &order.OmsOrderDetail{
		OrderItemList: orderItemList,
	}
	if err := copyProperties(omsOrder, orderDetail); err != nil {
		return nil, errors.New("字段赋值错误:" + err.Error())
	}
	return orderDetail, nil
}

func CancelOrder(orderId int64) (err error) {
	//查询未付款的取消订单
	cancelOrderList := make([]order.OmsOrder, 0)
	if err = global.Db.Model(&order.OmsOrder{}).
		Where("id = ?", orderId).Where("status = ?", 0).Where("delete_status = ?", 0).
		Find(&cancelOrderList).Error; err != nil {
		return errors.New("取消订单时，查询未付款订单出错:" + err.Error())
	}
	if len(cancelOrderList) == 0 {
		return nil
	}
	cancelOrder := cancelOrderList[0]
	if cancelOrder.ID != 0 {
		//修改订单状态为取消
		cancelOrder.Status = 4
		if err = global.Db.Model(&order.OmsOrder{}).Where("id", cancelOrder.ID).Updates(cancelOrder).Error; err != nil {
			return errors.New("取消订单时，更新订单状态出错:" + err.Error())
		}
		orderItemList := make([]order.OmsOrderItem, 0)
		if err = global.Db.Model(&order.OmsOrderItem{}).Where("order_id = ?", cancelOrder.ID).
			Find(&orderItemList).Error; err != nil {
			return errors.New("取消订单时，获取订单列表出错:" + err.Error())
		}
		//解除订单商品库存锁定
		if len(orderItemList) != 0 {
			for _, item := range orderItemList {
				if err := releaseStockBySkuID(item.ProductSkuId, item.ProductQuantity); err != nil {
					//这里貌似不应该返回报错
					global.Logger.Errorf("库存不足，无法释放")
				}
			}
		}
		//修改优惠券的使用状态（退还已用的优惠券）
		if err := updateCouponStatus(cancelOrder.CouponID, cancelOrder.MemberID, 0); err != nil {
			global.Logger.Errorf("退还优惠券时,修改优惠券状态失败")
		}
		//返还使用积分
		if cancelOrder.Integration != 0 {
			user := &users.User{}
			if err := user.GetMemberById(cancelOrder.MemberID); err != nil {
				global.Logger.Errorf("退还积分时，查询用户信息failed:%v", err.Error())
			}
			user.Integration += cancelOrder.Integration
			if err := user.Update(); err != nil {
				global.Logger.Errorf("退还积分时，更新用户积分failed:%v", err.Error())
			}
		}
	}
	return

}

func releaseStockBySkuID(productSkuId int64, productQuantity int) (err error) {
	if err = global.Db.Model(&cart.PmsSkuStock{}).
		Where("id = ?", productSkuId).Where("loc_stock - ? >=0", productQuantity).
		Update("lock_stock", gorm.Expr("lock_stock - ?", productQuantity)).Error; err != nil {
		return errors.New("释放库存失败:" + err.Error())
	}
	return nil
}

// 将优惠券信息更改为指定状态
func updateCouponStatus(couponID int64, memberId int64, useStatus int) (err error) {
	if couponID == 0 {
		return nil
	}

	//找到这张优惠券
	couponHistory := &coupon.SmsCouponHistory{}
	status := 1
	if useStatus == 1 {
		useStatus = 0
	}
	if err = global.Db.Model(&coupon.SmsCouponHistory{}).
		Where("member_id = ?", memberId).Where("coupon_id = ?", couponID).Where("use_status = ?", status).
		First(&couponHistory).Error; err != nil {
		return errors.New("查询优惠券出错:" + err.Error())
	}
	//将其使用时间和使用状态修改一下
	now := time.Now()
	couponHistory.UseTime = &now
	couponHistory.UseStatus = useStatus
	if err = global.Db.Model(&coupon.SmsCouponHistory{}).Where("id", couponHistory.ID).Updates(&couponHistory).Error; err != nil {
		return errors.New("更新优惠券状态出错:" + err.Error())
	}
	return nil
}

func ConfirmReceiveOrder(ctx *gin.Context, rec *receive.ConfirmReceiveOrderReqStruct) (err error) {
	//先验证一下是否是当前用户的订单
	//根据订单id查询订单
	curOrder := &order.OmsOrder{}
	if err = global.Db.Model(&order.OmsOrder{}).Where("id = ?", rec.OrderId).Find(&curOrder).Error; err != nil {
		global.Logger.Errorf("查询订单信息出错:%v", err)
	}
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	if curOrder.MemberID != memberId {
		return errors.New("没有操作权限，这是别人的订单")
	}
	if curOrder.Status != 2 {
		return errors.New("订单还未发货")
	}
	curOrder.Status = 3
	curOrder.ConfirmStatus = 1
	now := time.Now()
	curOrder.ReceiveTime = &now
	if err = global.Db.Model(&order.OmsOrder{}).Where("id = ?", curOrder.ID).Updates(curOrder).Error; err != nil {
		return errors.New("更新订单信息failed:" + err.Error())
	}
	return nil
}

func DeleteOrder(ctx *gin.Context, rec *receive.DeleteOrderReqStruct) (err error) {
	//先验证一下是否是当前用户的订单
	//根据订单id查询订单
	curOrder := &order.OmsOrder{}
	if err = global.Db.Model(&order.OmsOrder{}).Where("id = ?", rec.OrderId).Find(&curOrder).Error; err != nil {
		global.Logger.Errorf("查询订单信息出错:%v", err)
	}
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	if curOrder.MemberID != memberId {
		return errors.New("没有操作权限，这是别人的订单")
	}

	if curOrder.Status == 3 || curOrder.Status == 4 {
		curOrder.DeleteStatus = 1
		if err = global.Db.Model(&order.OmsOrder{}).Where("id = ?", curOrder.ID).Updates(&curOrder).Error; err != nil {
			return errors.New("删除订单时，更新订单状态失败" + err.Error())
		}
		return nil
	} else {
		return errors.New("只能删除已完成或已关闭的订单")
	}
}

func GenerateConfirmOrder(cartIds []int64, ctx *gin.Context) (result order.ConfirmOrderResult, err error) {
	//1.获取当前用户的购物车信息
	memberId, _ := jwt.GetMemberIdFromCtx(ctx)
	user := &users.User{}
	if err := user.GetMemberById(memberId); err != nil {
		return order.ConfirmOrderResult{}, errors.New("获取当前用户信息出错:" + err.Error())
	}
	cartPromotionItemList, err := cartLogic.CartListPromotion(cartIds, memberId)
	if err != nil {
		return order.ConfirmOrderResult{}, errors.New("获取用户的购物车信息出错:" + err.Error())
	}
	result.CartPromotionItemList = cartPromotionItemList
	//2.获取用户的收货地址列表
	memberReceiveAddressList := make([]order.UmsMemberReceiveAddress, 0)
	if err = global.Db.Model(&order.UmsMemberReceiveAddress{}).Where("member_id = ?", memberId).
		Find(&memberReceiveAddressList).Error; err != nil {
		return order.ConfirmOrderResult{}, errors.New("获取用户收货地址出错:" + err.Error())
	}
	result.MemberReceiveAddressList = memberReceiveAddressList
	//3.获取用户可用优惠券列表
	//todo:couponHistoryDetailList的结构还有问题
	couponHistoryDetailList, err := ListCart(cartPromotionItemList, 1, memberId)
	if err != nil {
		return order.ConfirmOrderResult{}, errors.New("获取用户可用优惠券列表出错:" + err.Error())
	}
	result.CouponHistoryDetailList = couponHistoryDetailList
	//4.补充其他的信息
	//获取用户积分
	result.MemberIntegration = user.Integration
	//获取积分使用规则
	integrationConsumeSetting := &integration.UmsIntegrationConsumeSetting{}
	integrationConsumeSetting.GetById(1)
	result.IntegrationConsumeSetting = *integrationConsumeSetting
	//计算总金额、活动优惠、应付金额
	calcAmount := calcCartAmount(cartPromotionItemList)
	result.CalcAmount = calcAmount
	return
}

// 根据购物车信息获取可用优惠券
func ListCart(cartItemList cart.CartPromotionItemList, couponType int, memberId int64) (result []coupon.SmsCouponHistoryDetail, err error) {
	//获取当前用户的信息
	user := &users.User{}
	if err := user.GetMemberById(memberId); err != nil {
		return nil, errors.New("listCart时，查询用户信息failed:" + err.Error())
	}
	//1.获取该用户所有优惠券
	allList, err := getDetailList(memberId)
	if err != nil {
		return nil, errors.New("获取用户的优惠券failed:" + err.Error())
	}

	//2.根据优惠券使用类型判断优惠券是否可用
	enableList := make([]coupon.SmsCouponHistoryDetail, 0)
	disableList := make([]coupon.SmsCouponHistoryDetail, 0)
	for _, couponHistoryDetail := range allList {
		useType := couponHistoryDetail.Coupon.UseType
		minPoint := couponHistoryDetail.Coupon.MinPoint
		endTime := couponHistoryDetail.Coupon.EndTime
		if useType == 0 {
			//0->全场通用
			//判断是否满足优惠起点
			//计算购物车商品的总价
			totalAmount := calcTotalAmount(cartItemList)
			if time.Now().Before(endTime) && totalAmount-minPoint >= 0 {
				enableList = append(enableList, couponHistoryDetail)
			} else {
				disableList = append(disableList, couponHistoryDetail)
			}
		} else if useType == 1 {
			//1->指定分类
			//计算指定分类商品的总价
			var productCategoryIds []int64
			for _, categoryRelation := range couponHistoryDetail.CategoryRelationList {
				productCategoryIds = append(productCategoryIds, categoryRelation.ProductCategoryId)
			}
			totalAmount := calcTotalAmountByproductCategoryId(cartItemList, productCategoryIds)
			if time.Now().Before(endTime) && totalAmount-minPoint >= 0 {
				enableList = append(enableList, couponHistoryDetail)
			} else {
				disableList = append(disableList, couponHistoryDetail)
			}
		} else if useType == 2 {
			//2->指定商品
			//计算指定商品的总价
			var productIds []int64
			for _, productRelation := range couponHistoryDetail.ProductRelationList {
				productIds = append(productIds, productRelation.ProductId)
			}
			totalAmount := calcTotalAmountByProductId(cartItemList, productIds)
			if time.Now().Before(endTime) && totalAmount-minPoint >= 0 {
				enableList = append(enableList, couponHistoryDetail)
			} else {
				disableList = append(disableList, couponHistoryDetail)
			}
		}
	}
	if couponType == 1 {
		return enableList, nil
	} else {
		return disableList, nil
	}
}

// 获取优惠券历史详情
func getDetailList(memberId int64) (result []coupon.SmsCouponHistoryDetail, err error) {
	// 查询优惠券历史
	var couponHistories []coupon.SmsCouponHistory
	if err = global.Db.Table("sms_coupon_history").
		Where("member_id = ? AND use_status = 0", memberId).Find(&couponHistories).Error; err != nil {
		return nil, errors.New("获取优惠券历史记录出错: " + err.Error())
	}
	// 初始化结果切片
	result = make([]coupon.SmsCouponHistoryDetail, 0, len(couponHistories))
	for _, history := range couponHistories {
		var detail coupon.SmsCouponHistoryDetail
		detail.ID = history.ID
		detail.CouponID = history.ID
		detail.MemberID = history.MemberID
		detail.CouponCode = history.CouponCode
		detail.MemberNickname = history.MemberNickname
		detail.GetType = history.GetType
		detail.CreateTime = history.CreateTime
		detail.UseStatus = history.UseStatus
		detail.OrderID = history.OrderID
		detail.OrderSn = history.OrderSn
		detail.Coupon = coupon.SmsCoupon{}
		//查询优惠券关联的商品信息
		if err = global.Db.Table("sms_coupon").
			Where("id=?", history.CouponID).First(&detail.Coupon).Error; err != nil {
			return nil, errors.New("获取优惠券信息出错: " + err.Error())
		}
		//查询关联的商品
		if err = global.Db.Table("sms_coupon_product_relation").
			Where("coupon_id = ?", history.CouponID).Find(&detail.ProductRelationList).Error; err != nil {
			return nil, errors.New("查询优惠券关联的商品信息出错:" + err.Error())
		}
		//查询优惠券关联的商品分类
		if err = global.Db.Table("sms_coupon_product_category_relation").
			Where("coupon_id = ?", history.CouponID).Find(&detail.CategoryRelationList).Error; err != nil {
			return nil, errors.New("查询优惠券关联的商品分类出错:" + err.Error())
		}
		result = append(result, detail)
	}
	return

	//if err = global.Db.Table("sms_coupon_history ch").
	//	Select("ch.*, c.id as c_id, c.name as c_name, c.amount as c_amount, c.min_point as c_min_point, c.platform as c_platform, c.start_time as c_start_time, c.end_time as c_end_time, c.note as c_note, c.use_type as c_use_type, c.type as c_type, cpr.id as cpr_id, cpr.product_id as cpr_product_id, cpcr.id as cpcr_id, cpcr.product_category_id as cpcr_product_category_id").
	//	Joins("LEFT JOIN sms_coupon c ON ch.coupon_id = c.id").
	//	Joins("LEFT JOIN sms_coupon_product_relation cpr ON cpr.coupon_id = c.id").
	//	Joins("LEFT JOIN sms_coupon_product_category_relation cpcr ON cpcr.coupon_id = c.id").
	//	Where("ch.member_id =?", memberId).
	//	Where("ch.use_status = 0").
	//	Scan(&result).Error; err != nil {
	//	return nil, errors.New("获取优惠券历史详情查表出错:" + err.Error())
	//}
	//return
}

// 计算购物车商品的总价
func calcTotalAmount(cartItemList cart.CartPromotionItemList) (result float32) {
	for _, item := range cartItemList {
		realPrice := item.Price - item.ReduceAmount
		result += realPrice * float32(item.Quantity)
	}
	return result
}

// 计算指定分类商品的总价
func calcTotalAmountByproductCategoryId(cartItemList cart.CartPromotionItemList, productCategoryIds []int64) (result float32) {
	for _, item := range cartItemList {
		if contain(item.ProductCategoryId, productCategoryIds) {
			realPrice := item.Price - item.ReduceAmount
			result += realPrice * float32(item.Quantity)
		}
	}
	return
}

func contain(id int64, ids []int64) bool {
	for _, value := range ids {
		if id == value {
			return true
		}
	}
	return false
}

// 计算指定商品的总价
func calcTotalAmountByProductId(cartItemList cart.CartPromotionItemList, productIds []int64) (result float32) {
	for _, item := range cartItemList {
		if contain(item.ProductId, productIds) {
			realPrice := item.Price - item.ReduceAmount
			result += realPrice * float32(item.Quantity)
		}
	}
	return
}

func getIntegrationConsumeSettingById(id int64) (result integration.UmsIntegrationConsumeSetting) {
	global.Db.Model(&integration.UmsIntegrationConsumeSetting{}).Where("id = ?", id).First(&result)
	return
}

// 计算购物车中商品的价格
func calcCartAmount(cartPromotionItemList cart.CartPromotionItemList) (result order.CalcAmount) {
	calcAmount := &order.CalcAmount{}
	totalAmount := float32(0)
	promotionAmount := float32(0)
	for _, cartPromotionItem := range cartPromotionItemList {
		totalAmount += cartPromotionItem.Price * float32(cartPromotionItem.Quantity)
		promotionAmount += cartPromotionItem.ReduceAmount * float32(cartPromotionItem.Quantity)
	}
	calcAmount.TotalAmount = totalAmount
	calcAmount.PromotionAmount = promotionAmount
	calcAmount.PayAmount = totalAmount - promotionAmount
	return *calcAmount
}

func GenerateOrder(data *receive.GenerateOrderReqStruct, ctx *gin.Context) (result map[string]interface{}, err error) {
	orderItemList := make([]order.OmsOrderItem, 0)
	//校验收货地址
	if data.MemberReceiveAddressId == 0 {
		return nil, errors.New("请选择收货地址！")
	}
	//获取购物车及优惠信息
	memberId, err := jwt.GetMemberIdFromCtx(ctx)
	currentMember := &users.User{}
	if err := currentMember.GetMemberById(memberId); err != nil {
		return nil, errors.New("获取用户身份信息出错:" + err.Error())
	}

	cartPromotionItemList, err := cartLogic.CartListPromotion(data.CartIds, memberId)
	if err != nil {
		return nil, errors.New("获取购物车及优惠信息failed:" + err.Error())
	}

	for _, cartPromotionItem := range cartPromotionItemList {
		//生成下单商品信息
		orderItem := &order.OmsOrderItem{
			ProductId:         cartPromotionItem.ProductId,
			ProductName:       cartPromotionItem.ProductName,
			ProductPic:        cartPromotionItem.ProductPic,
			ProductAttr:       cartPromotionItem.ProductAttr,
			ProductBrand:      cartPromotionItem.ProductBrand,
			ProductSn:         cartPromotionItem.ProductSn,
			ProductPrice:      cartPromotionItem.Price,
			ProductQuantity:   cartPromotionItem.Quantity,
			ProductSkuId:      cartPromotionItem.ProductSkuId,
			ProductSkuCode:    cartPromotionItem.ProductSkuCode,
			ProductCategoryId: cartPromotionItem.ProductCategoryId,
			PromotionAmount:   cartPromotionItem.ReduceAmount,
			PromotionName:     cartPromotionItem.PromotionMessage,
			GiftIntegration:   cartPromotionItem.Integration,
			GiftGrowth:        cartPromotionItem.Growth,
		}
		orderItemList = append(orderItemList, *orderItem)
	}

	//判断购物车中商品是否都有库存
	if !hasStock(cartPromotionItemList) {
		return nil, errors.New("生成订单时，某些商品的库存不足")
	}
	//判断使用使用了优惠券
	if data.CouponId == 0 {
		//不用优惠券
		for _, orderItem := range orderItemList {
			orderItem.CouponAmount = float32(0)
		}
	} else {
		//使用优惠券
		couponHistoryDetail, err := getUseCoupon(cartPromotionItemList, data.CouponId, memberId)
		if err != nil {
			return nil, errors.New("该优惠券不可用")
		}
		//对下单商品的优惠券进行处理
		handleCouponAmount(orderItemList, couponHistoryDetail)
	}

	//判断是否使用积分
	if data.UseIntegration == 0 {
		//不使用积分
		for _, orderItem := range orderItemList {
			orderItem.IntegrationAmount = float32(0)
		}
	} else {
		//使用积分
		totalAmount := calcOrderTotalAmount(orderItemList)
		integrationAmount := getUseIntegrationAmount(data.UseIntegration, totalAmount, currentMember, data.CouponId != 0)
		if integrationAmount == 0 {
			return nil, errors.New("积分不可用")
		} else {
			//可用情况下分摊到可用商品中
			for _, orderItem := range orderItemList {
				perAmount := orderItem.ProductPrice / totalAmount
				orderItem.IntegrationAmount = perAmount
			}
		}
	}
	//计算order_item的实付金额
	pointerOrderItemList := make([]*order.OmsOrderItem, len(orderItemList))
	for index := range orderItemList {
		pointerOrderItemList[index] = &orderItemList[index]
	}
	if err := handleRealAmount(pointerOrderItemList); err != nil {
		return nil, errors.New("计算订单的实际支付金额出错:" + err.Error())
	}

	//进行库存锁定
	//todo:This is a special comment.
	// 这个地方有点坑爹：make([]*cart.CartPromotionItem, 0)初始化内存时，指定长度应该为0。如果指定长度为len(cartPromotionItemList)：
	// 当cartPromotionItemList的长度=1时，pointerCartPromotionItemList赋值完毕自动扩容变为长度为2的切片，会出现一个nil指针在里面，后续处理的时候会遇到空指针错误
	pointerCartPromotionItemList := make([]*cart.CartPromotionItem, 0)
	for _, item := range cartPromotionItemList {
		pointerCartPromotionItemList = append(pointerCartPromotionItemList, &item)
	}

	if err := lockStock(pointerCartPromotionItemList); err != nil {
		return nil, err
	}
	//根据商品合计、运费、活动优惠、优惠券、积分计算应付金额
	order := &order.OmsOrder{
		DiscountAmount:  float32(0),
		TotalAmount:     calcOrderTotalAmount(orderItemList),
		FreightAmount:   float32(0),
		PromotionAmount: calcPromotionAmount(orderItemList),
		PromotionInfo:   getOrderPromotionInfo(orderItemList),
	}
	if data.CouponId == 0 {
		order.CouponAmount = 0
	} else {
		order.CouponID = data.CouponId
		order.CouponAmount = calcCouponAmount(orderItemList)
	}
	if data.UseIntegration == 0 {
		order.Integration = 0
		order.IntegrationAmount = 0
	} else {
		order.Integration = data.UseIntegration
		order.IntegrationAmount = calcIntegrationAmount(orderItemList)
	}
	order.PayAmount = calcPayAmount(*order)

	//转化为订单信息并插入数据库
	order.MemberID = currentMember.Id
	now := time.Now()
	order.CreateTime = &now
	order.MemberUsername = currentMember.Username
	//支付方式：0->未支付；1->支付宝；2->微信
	order.PayType = data.PayType
	//订单来源：0->PC订单；1->app订单
	order.SourceType = 1
	//订单状态：0->待付款；1->待发货；2->已发货；3->已完成；4->已关闭；5->无效订单
	order.Status = 0
	//订单类型：0->正常订单；1->秒杀订单
	order.OrderType = 0
	//收货人信息：姓名、电话、邮编、地址
	address := &orderModel.UmsMemberReceiveAddress{}
	address.GetAddressByid(data.MemberReceiveAddressId)
	order.ReceiverName = address.Name
	order.ReceiverPhone = address.PhoneNumber
	order.ReceiverPostCode = address.PostCode
	order.ReceiverProvince = address.Province
	order.ReceiverCity = address.City
	order.ReceiverRegion = address.Region
	order.ReceiverDetailAddress = address.DetailAddress
	//0->未确认；1->已确认
	order.ConfirmStatus = 0
	order.DeleteStatus = 0
	//计算赠送积分
	order.Integration = calcGifIntegration(orderItemList)
	//计算赠送成长值
	order.Growth = calcGiftGrowth(orderItemList)
	//生成订单号
	order.OrderSn = generateOrderSn(*order)
	//设置自动收货天数
	orderSettings := &orderModel.OmsOrderSettingList{}
	if err := orderSettings.GetAll(); err != nil {
		return nil, err
	}
	if len(*orderSettings) != 0 {
		order.AutoConfirmDay = (*orderSettings)[0].ConfirmOvertime
	}
	//插入订单表
	//todo:this is a special comment.
	// 凡是涉及到日期的，怎么全都是零值（应该是nil空值而不是00:00:00的零值），前台查询order表的时候不允许数据库的时间字段出现零值.
	// 解决办法：把time.Time类型的字段全都改成*time.Time类型，这样未赋值的字段就是nil类型了
	if err := order.Insert(); err != nil {
		return nil, errors.New("创建订单时，插入订单表failed:" + err.Error())
	}
	//todo:This is a special comment.
	// 这里是容易错的地方：如果用for _,item :=range(orderItemList)的方式遍历orderItemList，循环变量 orderItem 会得到元素的一个副本。
	// 解决办法：1.转换为指针切片然后修改；2.通过索引进行访问.这里选择更简单的索引访问
	for index := range orderItemList {
		orderItemList[index].OrderId = order.ID
		orderItemList[index].OrderSn = order.OrderSn
	}
	if err := insertOrderItemList(orderItemList); err != nil {
		return nil, err
	}
	//如使用优惠券更新优惠券使用状态
	if data.CouponId != 0 {
		if err := updateCouponStatus(data.CouponId, currentMember.Id, 1); err != nil {
			return nil, err
		}
	}
	//如使用积分需要扣除积分
	if data.UseIntegration != 0 {
		order.UseIntegration = data.UseIntegration
		if err := currentMember.Update(); err != nil {
			return nil, err
		}
	}
	//删除购物车中的下单商品
	if err := deleteCartItemList(cartPromotionItemList, *currentMember); err != nil {
		return nil, err
	}
	//发送延迟消息取消订单
	//todo:完善这个消息队列的逻辑
	if err := sendDelayMessageCancelOrder(order.ID); err != nil {
		global.Logger.Errorf("向延迟队列发送消息失败:%v", err)
	}
	result = map[string]interface{}{
		"order":         order,
		"orderItemList": orderItemList,
	}
	return result, nil
}

func sendDelayMessageCancelOrder(orderId int64) error {
	orderSetting := orderModel.OmsOrderSetting{}
	if err := orderSetting.GetById(1); err != nil {
		return fmt.Errorf("查询普通订单的过期时间出错: %v", err)
	}
	delayOrderTime := orderSetting.NormalOrderOvertime * 60 * 1000
	ch, err := global.RabbitMQ.Channel()
	if err != nil {
		return fmt.Errorf("获取rabbitmq的channel出错: %v", err)
	}
	defer ch.Close() // 确保在函数结束时关闭通道

	// 创建死信队列
	dlqName := "dlq"
	_, err = ch.QueueDeclare(
		dlqName, // 队列名称
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("声明死信队列失败: %v", err)
	}

	// 创建主队列，设置 TTL 和死信队列
	mainQueueName := "main_queue"
	_, err = ch.QueueDeclare(
		mainQueueName, // 队列名称
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		amqp.Table{
			"x-message-ttl":             delayOrderTime, // 消息的 TTL，单位为毫秒
			"x-dead-letter-exchange":    "",             // 使用默认交换机
			"x-dead-letter-routing-key": dlqName,        // 死信队列的路由键
		},
	)
	if err != nil {
		return fmt.Errorf("声明主队列失败: %v", err)
	}

	// 发送消息到主队列
	if err = ch.Publish(
		"",            // 使用默认交换机
		mainQueueName, // 路由键（主队列名称）
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(strconv.FormatInt(orderId, 10)),
		},
	); err != nil {
		return fmt.Errorf("发送消息到主队列失败: %v", err)
	}

	return nil
}

func hasStock(cartPromotionItemList cart.CartPromotionItemList) bool {
	for _, cartPromotionItem := range cartPromotionItemList {
		if cartPromotionItem.RealStock <= 0 || cartPromotionItem.RealStock < cartPromotionItem.Quantity {
			return false
		}
	}
	return true
}

// 获取该用户可以使用的优惠券
func getUseCoupon(cartPromotionItemList cart.CartPromotionItemList, couponId int64, memberId int64) (result coupon.SmsCouponHistoryDetail, err error) {
	couponHistoryDetailList, err := ListCart(cartPromotionItemList, 1, memberId)
	if err != nil {
		return coupon.SmsCouponHistoryDetail{}, err
	}
	for _, couponHistoryDetail := range couponHistoryDetailList {
		if couponHistoryDetail.Coupon.ID == couponId {
			return couponHistoryDetail, nil
		}
	}
	return coupon.SmsCouponHistoryDetail{}, nil
}

// 对优惠券优惠进行处理
func handleCouponAmount(orderItemList []order.OmsOrderItem, couponHistoryDetail coupon.SmsCouponHistoryDetail) {
	//先转换为指针切片
	coupon := couponHistoryDetail.Coupon
	pointerOrderItemList := make([]*order.OmsOrderItem, len(orderItemList))
	for index := range orderItemList {
		pointerOrderItemList[index] = &orderItemList[index]
	}
	if coupon.UseType == 0 {
		//全场通用
		calcPerCouponAmount(pointerOrderItemList, coupon)
	} else if coupon.UseType == 1 {
		//指定分类
		couponOrderItemList, err := getCouponOrderItemByRelation(couponHistoryDetail, orderItemList, 0)
		if err != nil {
			global.Logger.Errorf("获取couponOrderItemListfaield:%v", err)
		}
		//将couponOrderItemList转换为指针切片
		pointerCouponOrderItemList := make([]*order.OmsOrderItem, len(couponOrderItemList))
		for index := range couponOrderItemList {
			pointerCouponOrderItemList[index] = &couponOrderItemList[index]
		}
		calcPerCouponAmount(pointerCouponOrderItemList, coupon)
	} else if coupon.UseType == 2 {
		//指定商品
		couponOrderItemList, err := getCouponOrderItemByRelation(couponHistoryDetail, orderItemList, 1)
		if err != nil {
			global.Logger.Errorf("获取couponOrderItemList  faield:%v", err)
		}
		//couponOrderItemList转换为指针切片
		pointerCouponOrderItemList := make([]*order.OmsOrderItem, len(couponOrderItemList))
		for index := range couponOrderItemList {
			pointerCouponOrderItemList[index] = &couponOrderItemList[index]
		}
		calcPerCouponAmount(pointerCouponOrderItemList, coupon)
	}
}

// 计算总金额
func calcOrderTotalAmount(orderItemList []order.OmsOrderItem) (result float32) {
	for _, orderItem := range orderItemList {
		result += orderItem.ProductPrice * float32(orderItem.ProductQuantity)
	}
	return
}

// 获取可用积分抵扣金额
func getUseIntegrationAmount(useIntegration int, totalAmount float32, currentMember *users.User, hasCoupon bool) (result float32) {
	//判断用户是否有这么多积分
	if currentMember.Integration < useIntegration {
		return 0
	}
	//根据积分使用规则判断是否可以同时使用
	integrationConsumeSetting := &integration.UmsIntegrationConsumeSetting{}
	integrationConsumeSetting.GetById(1)
	if hasCoupon && integrationConsumeSetting.CouponStatus == 0 {
		//不可与优惠券共用
		return 0
	}
	//是否达到最低使用积分门槛
	if useIntegration < integrationConsumeSetting.UseUnit {
		return 0
	}
	//是否超过订单抵用最高百分比
	integrationAmount := float32(useIntegration) / float32(integrationConsumeSetting.UseUnit)
	maxPercent := float32(integrationConsumeSetting.MaxPercentPerOrder / 100)
	if integrationAmount > maxPercent {
		return 0
	}
	return integrationAmount
}

// 计算订单的实际支付金额
func handleRealAmount(orderItemList []*order.OmsOrderItem) (err error) {
	//原价-促销优惠-优惠券抵扣-积分抵扣
	for _, orderItem := range orderItemList {
		realAmount := orderItem.ProductPrice - orderItem.PromotionAmount - orderItem.CouponAmount - orderItem.IntegrationAmount
		orderItem.RealAmount = realAmount
	}
	return nil
}

// 锁定下单商品的所有库存
func lockStock(cartPromotionItemList []*cart.CartPromotionItem) (err error) {
	for _, cartPromotionItem := range cartPromotionItemList {
		skuStock := &cart.PmsSkuStock{}
		skuStock.GetSkuStockById(cartPromotionItem.ProductSkuId)
		skuStock.LockStock += cartPromotionItem.Quantity
		if err := lockStockBySkuId(cartPromotionItem.ProductSkuId, cartPromotionItem.Quantity); err != nil {
			return errors.New("锁定库存出错:" + err.Error())
		}
	}
	return nil
}

func lockStockBySkuId(productSkusId int64, quantity int) (err error) {
	if err = global.Db.Model(&cart.PmsSkuStock{}).
		Where("id = ?", productSkusId).Update("lock_stock", gorm.Expr("lock_stock - ?", quantity)).Error; err != nil {
		return errors.New("锁定库存时，更新库存failed:" + err.Error())
	}
	return nil
}

func calcPromotionAmount(orderItemList []order.OmsOrderItem) (result float32) {
	for _, orderItem := range orderItemList {
		if orderItem.PromotionAmount != 0 {
			result += orderItem.PromotionAmount * float32(orderItem.ProductQuantity)
		}
	}
	return
}

// 获取订单促销信息
func getOrderPromotionInfo(orderItemList []order.OmsOrderItem) (result string) {
	var sb strings.Builder
	for _, orderItem := range orderItemList {
		sb.WriteString(orderItem.PromotionName)
		sb.WriteString(";")
	}
	result = sb.String()
	if strings.HasSuffix(result, ";") {
		result = result[:len(result)-1]
	}
	return
}

// 计算订单优惠券金额
func calcCouponAmount(orderItemList []order.OmsOrderItem) (result float32) {
	for _, orderItem := range orderItemList {
		if orderItem.CouponAmount != 0 {
			result += orderItem.CouponAmount * float32(orderItem.ProductQuantity)
		}
	}
	return result
}

// 计算订单优惠券金额
func calcIntegrationAmount(orderItemList []order.OmsOrderItem) (result float32) {
	for _, orderItem := range orderItemList {
		if orderItem.IntegrationAmount != 0 {
			result += orderItem.IntegrationAmount * float32(orderItem.ProductQuantity)
		}
	}
	return
}

// 计算订单应付金额
func calcPayAmount(order order.OmsOrder) (result float32) {
	//总金额+运费-促销优惠-优惠券优惠-积分抵扣
	payAmount := order.TotalAmount + order.FreightAmount - order.PromotionAmount - order.CouponAmount - order.IntegrationAmount
	return payAmount
}

// 计算该订单赠送的积分
func calcGifIntegration(orderItemList []orderModel.OmsOrderItem) (result int) {
	for _, orderItem := range orderItemList {
		result += orderItem.GiftIntegration * orderItem.ProductQuantity
	}
	return
}

// 计算该订单赠送的成长值
func calcGiftGrowth(orderItemList []orderModel.OmsOrderItem) (result int) {
	for _, orderItem := range orderItemList {
		result += orderItem.GiftGrowth * orderItem.ProductQuantity
	}
	return
}

// 生成18位订单编号:8位日期+2位平台号码+2位支付方式+6位以上自增id
func generateOrderSn(order orderModel.OmsOrder) string {
	var sb strings.Builder
	date := time.Now().Format("2006-01-02 15:04:05")
	key := fmt.Sprintf("%s:%s%s", consts.RedisDatabase, consts.REDIS_KEY_ORDER_ID, date)
	increment, _ := global.RedisDb.IncrBy(key, 1).Result()
	sb.WriteString(date)
	sb.WriteString(fmt.Sprintf("%2d", order.SourceType))
	sb.WriteString(fmt.Sprintf("%2d", order.PayType))
	incrementStr := strconv.FormatInt(increment, 10)
	if len(incrementStr) < 6 {
		sb.WriteString(fmt.Sprintf("%6d", increment))
	} else {
		sb.WriteString(incrementStr)
	}
	return sb.String()
}

func insertOrderItemList(orderItemList []orderModel.OmsOrderItem) error {
	//orderItemList中的元素挨个插入
	for _, orderItem := range orderItemList {
		if err := global.Db.Create(&orderItem).Error; err != nil {
			return errors.New("插入到orderItem表出错:" + err.Error())
		}
	}
	return nil
}

// 删除下单商品的购物车信息
func deleteCartItemList(cartPromotionItemList []cart.CartPromotionItem, currentMember users.User) error {
	var ids []int64
	for _, cartPromotionItem := range cartPromotionItemList {
		ids = append(ids, cartPromotionItem.Id)
	}
	if err := global.Db.Model(&cart.OmsCartItem{}).
		Where("id in ?", ids).Where("member_id = ?", currentMember.Id).
		Update("delete_status", 1).Error; err != nil {
		return errors.New("删除下单商品的购物车信息出错:" + err.Error())
	}
	return nil
}

// calcPerCouponAmount 对每个下单商品进行优惠券金额分摊的计算
// 因为我们不需要返回值，直接操作每一个orderItem，所以传递过来的应该是指针切片
func calcPerCouponAmount(orderItemList []*orderModel.OmsOrderItem, coupon coupon.SmsCoupon) {
	//指针切片转换为值切片,作为参数传递给calcOrderTotalAmount
	valueOrderItemList := make([]orderModel.OmsOrderItem, len(orderItemList))
	for index, item := range orderItemList {
		valueOrderItemList[index] = *item
	}
	totalAmount := calcOrderTotalAmount(valueOrderItemList)
	for _, orderItem := range orderItemList {
		//(商品价格/可用商品总价)*优惠券面额
		couponAmount := orderItem.ProductPrice / totalAmount * coupon.Amount
		orderItem.CouponAmount = couponAmount
	}
}

// 获取与优惠券有关系的下单商品
func getCouponOrderItemByRelation(couponHistoryDetail coupon.SmsCouponHistoryDetail, orderItemList []orderModel.OmsOrderItem, useType int) (result []orderModel.OmsOrderItem, err error) {
	result = make([]orderModel.OmsOrderItem, 0)

	if useType == 0 {
		var categoryIdList []int64
		for _, productCategoryRelation := range couponHistoryDetail.CategoryRelationList {
			categoryIdList = append(categoryIdList, productCategoryRelation.ProductCategoryId)
		}
		for _, orderItem := range orderItemList {
			if contain(orderItem.ProductCategoryId, categoryIdList) {
				result = append(result, orderItem)
			} else {
				orderItem.CouponAmount = float32(0)
			}
		}
	} else if useType == 1 {
		var productIdList []int64
		for _, productRelation := range couponHistoryDetail.ProductRelationList {
			productIdList = append(productIdList, productRelation.ProductId)
		}
		for _, orderItem := range orderItemList {
			if contain(orderItem.ProductId, productIdList) {
				result = append(result, orderItem)
			} else {
				orderItem.CouponAmount = float32(0)
			}
		}
	}
	return result, nil
}

// List 按状态分页获取用户订单列表
func List(data *receive.ListReqStruct, memberId int64) (result []orderModel.OmsOrderDetail, err error) {
	//result的orderList部分
	var orderList []*orderModel.OmsOrder
	query := global.Db.Model(&orderModel.OmsOrder{}).
		Where("member_id", memberId).Where("delete_status = ?", 0)
	if data.Status != -1 {
		query = query.Where("status = ?", data.Status)
	}
	//todo:This is a special comment.
	// order需要放在find之前，否则不生效。
	if err = query.Offset((data.PageNum - 1) * data.PageSize).
		Limit(data.PageSize).
		Order("create_time desc").
		Find(&orderList).
		Error; err != nil {
		return nil, errors.New("分页查询订单列表failed:" + err.Error())
	}

	//设置数据信息
	orderIds := make([]int64, 0) //收集所有order的ids
	for _, item := range orderList {
		orderIds = append(orderIds, item.ID)
	}

	orderItemList := make([]orderModel.OmsOrderItem, 0)
	if err = global.Db.Model(&orderModel.OmsOrderItem{}).Where("order_id in ?", orderIds).
		Find(&orderItemList).Error; err != nil {
		return nil, errors.New("根据orderIds查询OmsOrderItem表失败:" + err.Error())
	}

	var orderDetailList []orderModel.OmsOrderDetail
	for _, item := range orderList {
		orderDetail := &orderModel.OmsOrderDetail{}
		//这里需要将order的所有字段全部赋值给orderDetail
		if err := copyProperties(item, orderDetail); err != nil {
			return nil, errors.New("赋值失败")
		}

		//然后在orderItemList里找到和当前orderDetail对应的几个item
		for _, item2 := range orderItemList {
			if item2.OrderId == item.ID {
				orderDetail.OrderItemList = append(orderDetail.OrderItemList, item2)
			}
		}
		orderDetailList = append(orderDetailList, *orderDetail)
	}
	return orderDetailList, nil
}

// copyProperties 将源结构体的字段值复制到目标结构体
func copyProperties(item *orderModel.OmsOrder, detail *orderModel.OmsOrderDetail) error {
	// 手动赋值字段
	detail.ID = item.ID
	detail.MemberID = item.MemberID
	detail.CouponID = item.CouponID
	detail.OrderSn = item.OrderSn
	detail.CreateTime = item.CreateTime
	detail.MemberUsername = item.MemberUsername
	detail.TotalAmount = item.TotalAmount
	detail.PayAmount = item.PayAmount
	detail.FreightAmount = item.FreightAmount
	detail.PromotionAmount = item.PromotionAmount
	detail.IntegrationAmount = item.IntegrationAmount
	detail.CouponAmount = item.CouponAmount
	detail.DiscountAmount = item.DiscountAmount
	detail.PayType = item.PayType
	detail.SourceType = item.SourceType
	detail.Status = item.Status
	detail.OrderType = item.OrderType
	detail.DeliveryCompany = item.DeliveryCompany
	detail.DeliverySn = item.DeliverySn
	detail.AutoConfirmDay = item.AutoConfirmDay
	detail.Integration = item.Integration
	detail.Growth = item.Growth
	detail.PromotionInfo = item.PromotionInfo
	detail.BillType = item.BillType
	detail.BillHeader = item.BillHeader
	detail.BillContent = item.BillContent
	detail.BillReceiverPhone = item.BillReceiverPhone
	detail.BillReceiverEmail = item.BillReceiverEmail
	detail.ReceiverName = item.ReceiverName
	detail.ReceiverPhone = item.ReceiverPhone
	detail.ReceiverPostCode = item.ReceiverPostCode
	detail.ReceiverProvince = item.ReceiverProvince
	detail.ReceiverCity = item.ReceiverCity
	detail.ReceiverRegion = item.ReceiverRegion
	detail.ReceiverDetailAddress = item.ReceiverDetailAddress
	detail.Note = item.Note
	detail.ConfirmStatus = item.ConfirmStatus
	detail.DeleteStatus = item.DeleteStatus
	detail.UseIntegration = item.UseIntegration
	detail.PaymentTime = item.PaymentTime
	detail.DeliveryTime = item.DeliveryTime
	detail.ReceiveTime = item.ReceiveTime
	detail.CommentTime = item.CommentTime
	detail.ModifyTime = item.ModifyTime

	return nil
}

func PaySuccess(data *receive.PaySuccessReqStruct) (count int, err error) {
	//修改订单的支付状态
	oneOrder := &orderModel.OmsOrder{}
	if err = global.Db.Model(&orderModel.OmsOrder{}).Where("id", data.OrderId).
		Where("delete_status", 0).Where("status", 0).First(oneOrder).Error; err != nil {
		return 0, errors.New("按照orderId和状态查询订单failed:" + err.Error())
	}
	oneOrder.PayType = data.PayType
	oneOrder.Status = 1
	now := time.Now()
	oneOrder.PaymentTime = &now
	if err = global.Db.Model(&orderModel.OmsOrder{}).Where("id", oneOrder.ID).Updates(oneOrder).Error; err != nil {
		return 0, errors.New("修改订单的支付状态和支付时间failed:" + err.Error())
	}
	//恢复商品的锁定库存，扣减真实库存
	orderDetail, err := getDetail(data.OrderId)
	if err != nil {
		return 0, err
	}
	if orderDetail.ID == 0 {
		return 0, errors.New("订单详情为空")
	}
	totalCount := 0
	for _, item := range orderDetail.OrderItemList {
		err := reduceSkuStock(item.ProductSkuId, item.ProductQuantity)
		if err != nil {
			return 0, errors.New("库存不足，无法扣减")
		}
		totalCount++
	}
	return totalCount, nil
}

func getDetail(orderId int64) (result orderModel.OmsOrderDetail, err error) {
	oneOrder := &orderModel.OmsOrder{}
	orderItem := &orderModel.OmsOrderItem{}
	if err = global.Db.Model(&orderModel.OmsOrder{}).Where("id", orderId).First(&oneOrder).Error; err != nil {
		return orderModel.OmsOrderDetail{}, errors.New("getDetail时，查询OmsOrder表failed:" + err.Error())
	}
	if err = global.Db.Model(&orderModel.OmsOrderItem{}).Where("order_id", orderId).First(&orderItem).Error; err != nil {
		return orderModel.OmsOrderDetail{}, errors.New("getDetail时，查询OmsOrderItem表failed:" + err.Error())
	}
	if err := copyProperties(oneOrder, &result); err != nil {
		return orderModel.OmsOrderDetail{}, errors.New("字段赋值错误:" + err.Error())
	}
	result.OrderItemList = append(result.OrderItemList, *orderItem)

	return
}

func reduceSkuStock(productSkuId int64, productQuantity int) (err error) {
	if err = global.Db.Model(&cart.PmsSkuStock{}).
		Where("id", productSkuId).
		Where("stock - ? >=0", productQuantity).
		Where("lock_stock - ? >=0", productQuantity).
		Update("lock_stock", gorm.Expr("lock_stock -?", productQuantity)).
		Update("stock", gorm.Expr("stock - ?", productQuantity)).
		Error; err != nil {
		return errors.New("更新锁定库存和真实库存failed:" + err.Error())
	}
	return nil
}

func CancelTimeOutOrder(memberId int64) (count int, err error) {
	orderSetting := &orderModel.OmsOrderSetting{}
	if err := orderSetting.GetById(1); err != nil {
		return 0, err
	}
	//查询超时、未支付的订单及订单详情
	timeOutOrders, err := getTimeOutOrders(orderSetting.NormalOrderOvertime)

	if err != nil {
		return 0, err
	}
	//修改订单状态为交易取消
	var ids []int64
	for _, timeOutOrder := range timeOutOrders {
		ids = append(ids, timeOutOrder.ID)
	}
	if err := updateOrderStatus(ids, 4); err != nil {
		return 0, err
	}
	for _, item := range timeOutOrders {
		//解除订单商品库存锁定
		if err := releaseSkuStockLock(item.OrderItemList); err != nil {
			return 0, errors.New("解除订单商品的库存锁定failed:" + err.Error())
		}
		//修改优惠券使用状态
		if err := updateCouponStatus(item.CouponID, memberId, 0); err != nil {
			return 0, err
		}
		//返还使用积分
		if item.UseIntegration != 0 {
			//找到这个用户，更新其积分
			curUser := &users.User{}
			if err := curUser.GetMemberById(memberId); err != nil {
				return 0, errors.New("查询用户信息出错:" + err.Error())
			}
			curUser.Integration += item.UseIntegration
			if err := curUser.Update(); err != nil {
				return 0, errors.New("返还用户积分出错:" + err.Error())
			}
		}
	}
	return len(timeOutOrders), nil
}

// getTimeOutOrders 查询已经超时的订单
func getTimeOutOrders(normalOverTime int) (result []orderModel.OmsOrderDetail, err error) {
	//todo:This is a special comment.
	// 这个地方有个雷，计算时间差值是否大于normalOverTime来判断订单是否超时，需要注意数据库和代码中的时区保持一致，不然可能查不到已经超时的订单。
	orders := make([]orderModel.OmsOrder, 0)
	orderItems := make([]orderModel.OmsOrderItem, 0)
	if err = global.Db.Model(&orderModel.OmsOrder{}).
		Where("status", 0).
		Where("TIMESTAMPDIFF(MINUTE,create_time,?) > ?", time.Now().Local(), normalOverTime).
		Find(&orders).Error; err != nil {
		return nil, errors.New("查询超时订单时，查询OmsOrder表failed:" + err.Error())
	}

	//将已经超时的订单的ids整合起来
	var ids []int64
	for _, oneOrder := range orders {
		ids = append(ids, oneOrder.ID)
	}
	//根据上述ids查询对应的OmsOrderItem项
	if err = global.Db.Model(&orderModel.OmsOrderItem{}).Where("order_id in ?", ids).
		Find(&orderItems).Error; err != nil {
		return nil, errors.New("查询超时订单时，查询OmsOrderItem表failed:" + err.Error())
	}
	//整合为DTO类型并返回
	result = make([]orderModel.OmsOrderDetail, 0, len(orders))
	for _, oneOrder := range orders {
		oneOrderDetail := &orderModel.OmsOrderDetail{}
		if err := copyProperties(&oneOrder, oneOrderDetail); err != nil {
			return nil, errors.New("字段赋值错误:" + err.Error())
		}
		for _, oneOrderItem := range orderItems {
			if oneOrderItem.OrderId == oneOrder.ID {
				oneOrderDetail.OrderItemList = append(oneOrderDetail.OrderItemList, oneOrderItem)
			}
		}
		result = append(result, *oneOrderDetail)
	}
	return result, nil
}

// 更新订单状态
func updateOrderStatus(ids []int64, newStatus int) (err error) {
	if err = global.Db.Model(&orderModel.OmsOrder{}).Where("id in (?)", ids).
		Update("status", newStatus).Error; err != nil {
		return errors.New("更新订单状态failed:" + err.Error())
	}
	return nil
}

// 解除订单商品库存锁定
func releaseSkuStockLock(orderItemList []orderModel.OmsOrderItem) (err error) {
	for _, orderItem := range orderItemList {
		if err = global.Db.Model(&cart.PmsSkuStock{}).Where("id", orderItem.ProductSkuId).
			Update("lock_stock", gorm.Expr("lock_stock -?", orderItem.ProductQuantity)).
			Where("id", orderItem.ProductSkuId).Error; err != nil {
			return errors.New("解除订单商品库存锁定failed:" + err.Error())
		}
	}
	return nil
}

func CreateReturnApply(data *receive.CreateReturnApplyReqStruct) error {
	realApply := &orderModel.OmsOrderReturnApply{
		OrderId:          data.OrderId,
		ProductId:        data.ProductId,
		OrderSn:          data.OrderSn,
		MemberUsername:   data.MemberUsername,
		ReturnName:       data.ReturnName,
		ReturnPhone:      data.ReturnPhone,
		ProductPic:       data.ProductPic,
		ProductName:      data.ProductName,
		ProductBrand:     data.ProductBrand,
		ProductAttr:      data.ProductAttr,
		ProductCount:     data.ProductCount,
		ProductPrice:     data.ProductPrice,
		ProductRealPrice: data.ProductRealPrice,
		Reason:           data.Reason,
		Description:      data.Description,
		ProofPics:        data.ProofPics,
	}
	now := time.Now()
	realApply.CreateTime = &now
	realApply.Status = 0
	return global.Db.Create(&realApply).Error
}
