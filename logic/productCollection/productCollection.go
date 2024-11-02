package productCollection

import (
	"context"
	"errors"
	"gomall/global"
	receive "gomall/interaction/receive/productCollection"
	"gomall/models/home"
	productCollectionModels "gomall/models/productCollection"
	"strconv"
	"time"
)

func Add(data *receive.AddReqStruct) (count int, err error) {
	//将data的数据转换为productCollection对象
	productCollection := &productCollectionModels.MemberProductCollection{
		MemberID:        data.MemberId,
		MemberNickname:  data.MemberNickname,
		MemberIcon:      data.MemberIcon,
		ProductID:       data.ProductId,
		ProductName:     data.ProductName,
		ProductPic:      data.ProductPic,
		ProductPrice:    strconv.FormatFloat(float64(data.ProductPrice), 'f', -1, 32),
		ProductSubTitle: data.ProductSubTitle,
		CreateTime:      time.Now(),
	}

	//先根据memberId和productId在mongodb中找到对应的记录
	db := global.MongoDb.Database("mall-port")
	repo := NewMemberProductCollectionRepository(db, "memberProductCollection")
	//执行查询
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findCollection, err := repo.FindByMemberIdAndProductId(ctx, data.MemberId, data.ProductId)
	if err != nil {
		return 0, errors.New("添加收藏的过程中，根据memberId和productId查询记录出错: " + err.Error())
	}

	if findCollection == nil {
		product := &home.PmsProduct{}
		if err = product.GetById(data.ProductId); err != nil {
			return 0, errors.New("添加收藏的过程中，查询商品时出错: " + err.Error())
		}
		//
		if product.Id == 0 || product.DeleteStatus == 1 {
			return 0, nil
		}
		productCollection.ProductName = product.Name
		productCollection.ProductSubTitle = product.SubTitle
		productCollection.ProductPrice = strconv.FormatFloat(float64(product.Price), 'f', -1, 32)
		productCollection.ProductPic = product.Pic
		//保存
		if err = repo.Save(ctx, productCollection); err != nil {
			return 0, errors.New("添加收藏的过程中，插入商品时出错: " + err.Error())
		}
	}
	return 1, nil
}

func Clear(memberId int64) (err error) {

	db := global.MongoDb.Database("mall-port")
	repo := NewMemberProductCollectionRepository(db, "memberProductCollection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//执行查询
	if err = repo.ClearByMemberId(ctx, memberId); err != nil {
		return err
	}
	return nil
}

func Delete(productId int64, memberId int64) (count int, err error) {
	db := global.MongoDb.Database("mall-port")
	repo := NewMemberProductCollectionRepository(db, "memberProductCollection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//执行查询
	if err = repo.Delete(ctx, productId, memberId); err != nil {
		return 0, errors.New("删除收藏的商品失败: " + err.Error())
	}
	return 1, nil
}

func Detail(productId int64, memberId int64) (result *productCollectionModels.MemberProductCollection, err error) {
	db := global.MongoDb.Database("mall-port")
	repo := NewMemberProductCollectionRepository(db, "memberProductCollection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//执行查询
	result, err = repo.Detail(ctx, productId, memberId)
	if err != nil {
		return nil, errors.New("获取收藏商品详情失败: " + err.Error())
	}
	return result, nil
}

func List(pageNum int, pageSize int, memberId int64) (result []productCollectionModels.MemberProductCollection, err error) {
	db := global.MongoDb.Database("mall-port")
	repo := NewMemberProductCollectionRepository(db, "memberProductCollection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//执行查询
	result, err = repo.List(ctx, pageNum, pageSize, memberId)
	if err != nil {
		return nil, errors.New("获取收藏商品列表失败: " + err.Error())
	}
	return result, nil
}
