package home

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/home"
	homeResponse "gomall/interaction/response/home"
	"gomall/logic/home"
	"strconv"
)

type HomeController struct {
	controller.BaseControllers
}

func (c *HomeController) GetHomeContent(ctx *gin.Context) {
	if _, err := controller.ShouldBind(ctx, new(receive.GetHomeContentRequestStruct)); err == nil {
		//1.获取首页广告
		advertiseList, err := home.GetHomeAdvertiseList()
		if err != nil {
			c.Response(ctx, "获取首页广告出错", nil, err)
			return
		}
		//2.获取推荐品牌
		brandList, err := home.GetRecommendBrandList(0, 6)
		if err != nil {
			c.Response(ctx, "获取推荐品牌出错", nil, err)
			return
		}
		//3.获取秒杀信息
		flashPromotion, err := home.GetHomeFlashPromotion()
		if err != nil {
			c.Response(ctx, "获取秒杀信息出错", nil, err)
			return
		}
		responseData := &homeResponse.GetHomeContentResponseStruct{
			AdvertiseList:      advertiseList,
			BrandList:          brandList,
			HomeFlashPromotion: *flashPromotion,
		}
		//4.获取新品推荐
		newProductList, err := home.GetNewProductList(0, 4)
		if err != nil {
			c.Response(ctx, "获取新品推荐出错", nil, err)
			return
		}
		responseData.NewProductList = newProductList
		//5.获取人气推荐
		hotProductList, err := home.GetHotProductList(0, 4)
		if err != nil {
			c.Response(ctx, "获取热门推荐出错", nil, err)
			return
		}
		responseData.HotProductList = hotProductList
		//6.获取专题推荐
		recommendSubjectList, err := home.GetRecommendSubjectList(0, 4)
		if err != nil {
			c.Response(ctx, "获取推荐专题出错", nil, err)
			return
		}
		responseData.SubjectList = recommendSubjectList
		//全部封装完毕，返回数据
		c.Response(ctx, "获取首页内容成功", responseData, nil)
	}
}

// 分页获取热门商品推荐
func (c *HomeController) GetHotProductList(ctx *gin.Context) {
	//绑定参数
	var rec receive.GetHotProductListRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("GetHotProductList请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	//请求并返回
	hotProductList, err := home.GetHotProductList((rec.PageNum-1)*rec.PageSize, rec.PageSize)
	if err != nil {
		c.Response(ctx, "获取热门商品出错", nil, err)
		return
	}
	c.Response(ctx, "获取热门商品成功", hotProductList, nil)
}

func (c *HomeController) GetNewProductList(ctx *gin.Context) {
	var rec receive.GetNewProductListRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("GetNewProductList请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	newProductList, err := home.GetNewProductList((rec.PageNum-1)*rec.PageSize, rec.PageSize)
	if err != nil {
		c.Response(ctx, "获取最新商品出错", nil, err)
		return
	}
	c.Response(ctx, "获取最新商品成功", newProductList, nil)

}

// 分页获取推荐商品
func (c *HomeController) GetRecommendProductList(ctx *gin.Context) {
	var rec receive.GetRecommendProductListRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("GetRecommendProductList请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	newProductList, err := home.GetRecommendProductList((rec.PageNum-1)*rec.PageSize, rec.PageSize)
	if err != nil {
		c.Response(ctx, "获取推荐商品出错", nil, err)
		return
	}
	c.Response(ctx, "获取推荐商品成功", newProductList, nil)
}

// 分类分页获取专题
func (c *HomeController) GetSubjectList(ctx *gin.Context) {
	var rec receive.GetSubjectListRequestStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("GetSubjectList请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	newProductList, err := home.GetSubjectListByCategoryId(&rec)
	if err != nil {
		c.Response(ctx, "获取专题出错", nil, err)
		return
	}
	c.Response(ctx, "获取专题成功", newProductList, nil)

}

// GetProductCateList 获取首页商品分类
func (c *HomeController) GetProductCateList(ctx *gin.Context) {
	//var rec receive.GetProductCateListRequestStruct
	//if err := ctx.ShouldBindJSON(&rec); err != nil {
	//	global.Logger.Errorf("GetProductCateList请求传入参绑定失败: %v", err)
	//	c.Response(ctx, "请求参数错误", nil, err)
	//	return
	//}
	parentId, _ := strconv.ParseInt(ctx.Param("parentId"), 10, 64)
	productCateList, err := home.GetProductCateList(parentId)
	if err != nil {
		c.Response(ctx, "获取专题出错", nil, err)
		return
	}
	c.Response(ctx, "获取专题成功", productCateList, nil)
}
