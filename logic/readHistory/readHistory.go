package readHistory

import (
	"context"
	"errors"
	"fmt"
	"gomall/global"
	receive "gomall/interaction/receive/readHistory"
	"gomall/models/history"
	"gomall/models/home"
	readHistoryCollection "gomall/models/readHistory"
	"gomall/models/users"
	"time"
)

func CreateReadHistory(data *receive.CreateReadHistoryReqStruct, memberId int64) (err error) {

	readHistory := &readHistoryCollection.MemberReadHistory{}
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
	if product.Id == 0 {
		return errors.New("商品不存在")
	}
	//然后开始构造readHistory的用户信息和商品信息
	readHistory.MemberID = memberId
	readHistory.MemberNickname = user.Nickname
	readHistory.MemberIcon = user.Icon
	readHistory.CreateTime = time.Now()
	readHistory.ProductID = data.ProductId
	readHistory.ProductName = product.Name
	readHistory.ProductSubTitle = product.SubTitle
	readHistory.ProductPrice = fmt.Sprintf("%.2f", product.Price)
	readHistory.ProductPic = product.Pic

	db := global.MongoDb.Database("mall-port")
	repo := NewMemberReadHistoryRepository(db, "memberReadHistory")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//执行查询
	if err = repo.CreateReadHistory(ctx, readHistory, memberId); err != nil {
		return errors.New("创建浏览记录时出错，向mongodb插入数据失败: " + err.Error())
	}
	return nil
}

func ClearReadHistory(memberId int64) (err error) {
	db := global.MongoDb.Database("mall-port")
	repo := NewMemberReadHistoryRepository(db, "memberReadHistory")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err = repo.Clear(ctx, memberId); err != nil {
		return err
	}
	return nil
}

func DeleteReadHistoryByIds(data *receive.DeleteReadHistoryReqStruct, memberId int64) (err error) {
	db := global.MongoDb.Database("mall-port")
	repo := NewMemberReadHistoryRepository(db, "memberReadHistory")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err = repo.DeleteByIds(ctx, data.Ids, memberId); err != nil {
		return err
	}

	return nil
}

func ListReadHistory(pageNum int64, pageSize int64, memberId int64) (results []history.MemberReadHistoryMongoStruct, err error) {
	db := global.MongoDb.Database("mall-port")
	repo := NewMemberReadHistoryRepository(db, "memberReadHistory")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//执行查询
	result, err := repo.FindByMemberIdOrderByCreateTimeDesc(ctx, int(pageNum), int(pageSize), memberId)
	if err != nil {
		return nil, errors.New("获取收藏商品列表失败: " + err.Error())
	}
	global.Logger.Infof("分页查询浏览记录的result为:%v", result)

	return result, nil
}
