package readHistory

import (
	"errors"
	"gomall/global"
	receive "gomall/interaction/receive/readHistory"
	"gomall/models/history"
	"gomall/models/home"
	"gomall/models/users"
	"time"
)

func CreateReadHistory(data *receive.CreateReadHistoryReqStruct, memberId int64) (err error) {
	readHistory := &history.MemberReadHistory{}
	user := &users.User{}
	product := &home.PmsProduct{}
	//需要根据memberId倒查用户信息
	if err = global.Db.Model(&users.User{}).Where("id = ?", memberId).First(&user).Error; err != nil {
		return errors.New("创建浏览记录时，查询用户信息failed:" + err.Error())
	}
	//根据data.productId倒查商品信息
	if err = global.Db.Model(&home.PmsProduct{}).
		Where("id = ?", data.ProductId).Where("delete_status = ?", 0).First(&product).Error; err != nil {
		return errors.New("创建浏览记录时，查询商品信息failed:" + err.Error())
	}
	//然后开始构造readHistory
	readHistory.MemberId = memberId
	readHistory.MemberNickname = user.Nickname
	readHistory.MemberIcon = user.Icon
	readHistory.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	if product.Id == 0 {
		return errors.New("商品不存在")
	}
	readHistory.ProductId = data.ProductId
	readHistory.ProductName = product.Name
	readHistory.ProductSubTitle = product.SubTitle
	readHistory.ProductPrice = product.Price
	readHistory.ProductPic = product.Pic
	if err = global.Db.Create(&readHistory).Error; err != nil {
		return errors.New("创建浏览记录时，插入浏览记录failed:" + err.Error())
	}
	return
}

func ClearReadHistory(memberId int64) (err error) {
	if err = global.Db.Model(&history.MemberReadHistory{}).Where("member_id = ?", memberId).
		Delete(&history.MemberReadHistory{}).Error; err != nil {
		return errors.New("清空浏览记录出错:" + err.Error())
	}
	return nil
}

func DeleteReadHistoryByIds(data *receive.DeleteReadHistoryReqStruct, memberId int64) (err error) {
	if err = global.Db.Model(&history.MemberReadHistory{}).
		Where("member_id = ?", memberId).Where("id in ?", data.Ids).
		Delete(&history.MemberReadHistory{}).Error; err != nil {
		return errors.New("删除浏览记录出错：" + err.Error())
	}
	return nil
}

func ListReadHistory(offset int64, limit int64, memberId int64) (results []history.MemberReadHistory, err error) {
	if err = global.Db.Model(&history.MemberReadHistory{}).Where("member_id = ?", memberId).
		Offset(int(offset)).Limit(int(limit)).Order("create_time desc").Find(&results).Error; err != nil {
		return nil, errors.New("分页查找浏览记录出错:%v" + err.Error())
	}
	return
}
