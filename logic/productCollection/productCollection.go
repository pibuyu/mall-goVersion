package productCollection

import (
	"context"
	"errors"
	"gomall/global"
	receive "gomall/interaction/receive/productCollection"
	"gomall/models/home"
	productCollectionModels "gomall/models/productCollection"
	"time"
)

func Add(data *receive.AddReqStruct) (count int, err error) {
	//将data的数据转换为productCollection对象
	productCollection := &productCollectionModels.MemberProductCollection{
		MemberId:        data.MemberId,
		MemberNickname:  data.MemberNickname,
		MemberIcon:      data.MemberIcon,
		ProductId:       data.ProductId,
		ProductName:     data.ProductName,
		ProductPic:      data.ProductPic,
		ProductPrice:    data.ProductPrice,
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
		productCollection.ProductPrice = product.Price
		productCollection.ProductPic = product.Pic
		//保存
		if err = repo.Save(ctx, productCollection); err != nil {
			return 0, errors.New("添加收藏的过程中，插入商品时出错: " + err.Error())
		}
	}
	return 1, nil
}