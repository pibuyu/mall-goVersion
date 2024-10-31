package brandAttention

import (
	"context"
	"errors"
	"gomall/global"
	receive "gomall/interaction/receive/brandAttention"
	pmsBrandModels "gomall/models/brand"
	brandAttentionModels "gomall/models/brandAttention"
	"time"
)

func Add(data *receive.AddBrandAttentionReqStruct) (count int, err error) {
	//将data数据映射到brandAttention对象中
	memberBrandAttention := &brandAttentionModels.MemberBrandAttention{
		//member
		MemberId:       data.MemberId,
		MemberNickname: data.MemberNickname,
		MemberIcon:     data.MemberIcon,
		//brand
		BrandId:    data.BrandId,
		BrandName:  data.BrandName,
		BrandLogo:  data.BrandLogo,
		BrandCity:  data.BrandCity,
		CreateTime: time.Now(),
	}
	//先根据memberId和productId在mongodb中找到对应的记录
	db := global.MongoDb.Database("mall-port")
	repo := NewMemberBrandAttentionRepository(db, "memberBrandAttention")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//执行查询
	findCollection, err := repo.FindByMemberIdAndProductId(ctx, data.MemberId, data.BrandId)
	if err != nil {
		return 0, errors.New("添加关注的品牌的过程中，根据memberId和brandId在mongodb查询记录出错: " + err.Error())
	}

	if findCollection == nil {
		brand := &pmsBrandModels.PmsBrand{}
		if err = brand.GetById(data.BrandId); err != nil {
			return 0, errors.New("添加关注的品牌的过程中，在pmsBrand表查询记录出错: " + err.Error())
		}
		if brand.ID == 0 {
			return 0, errors.New("添加关注的品牌的过程中，在pmsBrand表中未查询到该品牌相关信息: " + err.Error())
		}

		//保存
		memberBrandAttention.BrandName = brand.Name
		memberBrandAttention.BrandLogo = brand.Logo
		if err = repo.Save(ctx, memberBrandAttention); err != nil {
			return 0, errors.New("添加关注的品牌的过程中，在mongodb中插入记录出错: " + err.Error())
		}
	}
	return 1, nil
}