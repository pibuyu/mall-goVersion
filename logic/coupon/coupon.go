package coupon

import (
	"errors"
	"fmt"
	"gomall/global"
	orderLogic "gomall/logic/order"
	"gomall/models/cart"
	couponModels "gomall/models/coupon"
	"gomall/models/home"
	"gomall/models/users"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 领取优惠券
func AddCoupon(couponId int64, memberId int64) (err error) {
	//先获取当前用户
	curUser := &users.User{}
	if err = curUser.GetMemberById(memberId); err != nil {
		return errors.New("领取优惠券过程中，获取用户信息失败:" + err.Error())
	}
	//然后找到对应的优惠券信息，判断数量
	oneCoupon := &couponModels.SmsCoupon{}
	if err := oneCoupon.GetCouponById(couponId); err != nil {
		return errors.New("领取优惠券过程中，获取优惠券信息失败:" + err.Error())
	}
	global.Logger.Infof("查询出来的优惠券信息为:%v", oneCoupon)
	if oneCoupon.Count <= 0 {
		return errors.New("优惠券已经领完了")
	}
	if time.Now().Before(oneCoupon.EnableTime) {
		return errors.New("还没到优惠券发放时间")
	}
	//判断用户领取的优惠券数量是否超过限制
	couponHistoryExample := &couponModels.SmsCouponHistoryList{}
	if err = couponHistoryExample.GetByCouponIdAndMemberId(couponId, memberId); err != nil {
		return errors.New("领取优惠券过程中，获取优惠券领取历史记录失败:" + err.Error())
	}
	//解引用，获取到指向的切片
	global.Logger.Infof("当前id=%d的用户已经领取了id=%d的优惠券%d张", memberId, couponId, len(*couponHistoryExample))
	global.Logger.Infof("当前优惠券限制领取%d张", oneCoupon.PerLimit)
	if len(*couponHistoryExample) >= oneCoupon.PerLimit {
		return errors.New("优惠券领取数量已经超过限制")
	}
	//生成优惠券领取历史项
	couponHistory := &couponModels.SmsCouponHistory{
		CouponID:       couponId,
		CouponCode:     GenerateCouponCode(memberId),
		CreateTime:     time.Now(),
		MemberID:       memberId,
		MemberNickname: curUser.Nickname,
		GetType:        1,
		UseStatus:      0,
	}
	if err = couponHistory.Create(); err != nil {
		return errors.New("领取优惠券过程中，创建优惠券领取历史记录失败:" + err.Error())
	}
	oneCoupon.Count -= 1
	if oneCoupon.ReceiveCount == 0 {
		oneCoupon.ReceiveCount = 1
	} else {
		oneCoupon.ReceiveCount += 1
	}
	if err := oneCoupon.Update(); err != nil {
		return errors.New("修改优惠券表的数量和领取数量时出错:" + err.Error())
	}
	return nil
}

func GetCouponList(memberId int64, useStatus int) (result []couponModels.SmsCoupon, err error) {
	query := global.Db.Table("sms_coupon_history ch").
		Select("c.*").
		Joins("LEFT JOIN sms_coupon c ON ch.coupon_id = c.id").
		Where("ch.member_id = ?", memberId)
	if useStatus != 2 {
		query = query.Where("ch.use_status = ?", useStatus).
			Where("NOW() > c.start_time").
			Where("c.end_time > NOW()")
	} else {
		query = query.Where("NOW() > c.end_time")
	}
	if err = query.Debug().Find(&result).Error; err != nil {
		return nil, errors.New("查询对应的优惠券列表:" + err.Error())
	}
	return result, nil
}

// generateCouponCode 16位优惠码生成：时间戳后8位+4位随机数+用户id后4位
func GenerateCouponCode(memberId int64) string {
	var sb strings.Builder
	currentTimeMillis := time.Now().UnixNano() / 1e6
	timeMillisStr := strconv.FormatInt(currentTimeMillis, 10)
	sb.WriteString(timeMillisStr[len(timeMillisStr)-8:])
	//生成4位随机数字
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 4; i++ {
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(10)))
	}
	//添加会员id的后4位
	memberIdStr := strconv.FormatInt(memberId, 10)
	if len(memberIdStr) <= 4 {
		sb.WriteString(fmt.Sprintf("%04d", memberId))
	} else {
		sb.WriteString(memberIdStr[len(memberIdStr)-4:])
	}
	return sb.String()
}

// ListCartCoupon 根据购物车信息获取可用的优惠券
func ListCartCoupon(cartPromotionItemList cart.CartPromotionItemList, couponType int, memberId int64) (result []couponModels.SmsCouponHistoryDetail, err error) {
	return orderLogic.ListCart(cartPromotionItemList, couponType, memberId)
}

func ListByProductId(productId int64, memberId int64) (result couponModels.SmsCouponList, err error) {
	allCouponIds := make([]int64, 0)
	//获取指定商品优惠券
	cprList := &couponModels.SmsCouponProductRelationList{}
	if err = cprList.GetByProductId(productId); err != nil {
		return nil, errors.New("查询商品与优惠券关系出错:" + err.Error())
	}
	//拿到所有的couponId
	if len(*cprList) != 0 {
		for _, cpr := range *cprList {
			allCouponIds = append(allCouponIds, cpr.CouponId)
		}
	}
	//获取指定分类优惠券
	product := &home.PmsProduct{}
	if err = product.GetById(productId); err != nil {
		return nil, errors.New("根据productId获取商品信息出错：" + err.Error())
	}

	cpcrList := &couponModels.SmsCouponProductCategoryRelationList{}
	if err = cpcrList.GetByProductCategoryId(product.ProductCategoryId); err != nil {
		return nil, errors.New("根据productCategoryId获取SmsCouponProductCategoryRelationList出错：" + err.Error())
	}
	if len(*cpcrList) != 0 {
		for _, cpcr := range *cpcrList {
			allCouponIds = append(allCouponIds, cpcr.CouponId)
		}
	}
	//所有优惠券
	couponLists := &couponModels.SmsCouponList{}
	if err = couponLists.GetListByIds(allCouponIds); err != nil {
		return nil, errors.New("查询所有优惠券出错:" + err.Error())
	}
	return *couponLists, nil
}

func ListHistory(memberId int64, useStatus int) (result []couponModels.SmsCouponHistory, err error) {
	var couponHistoryList []couponModels.SmsCouponHistory
	if err = global.Db.Model(&couponModels.SmsCouponHistory{}).
		Where("member_id = ? and use_status = ?", memberId, useStatus).
		Find(&couponHistoryList).Error; err != nil {
		return nil, errors.New("根据member_id和use_status查询优惠券历史列表出错:" + err.Error())
	}
	return couponHistoryList, nil
}
