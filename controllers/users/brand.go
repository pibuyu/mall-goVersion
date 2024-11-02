package users

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/brand"
	brandLogic "gomall/logic/brand"
	"gomall/utils/response"
	"strconv"
)

type BrandController struct {
	controller.BaseControllers
}

func (c *BrandController) Detail(ctx *gin.Context) {
	//var rec receive.DetailReqStruct
	//if err := ctx.ShouldBindJSON(&rec); err != nil {
	//	c.Response(ctx, "获取品牌详情时，参数绑定错误", nil, err)
	//	global.Logger.Errorf("获取品牌详情时，参数绑定错误:%v", err)
	//	return
	//}
	brandId, _ := strconv.ParseInt(ctx.Param("brandId"), 10, 64)

	brand, err := brandLogic.Detail(brandId)
	if err != nil {
		c.Response(ctx, "获取品牌详情失败", nil, err)
		return
	}
	c.Response(ctx, "获取品牌详情成功", brand, nil)
}

// 分页获取推荐品牌
func (c *BrandController) RecommendList(ctx *gin.Context) {
	var rec receive.RecommendListReqStruct
	if err := ctx.ShouldBind(&rec); err != nil {
		c.Response(ctx, "分页获取推荐品牌时，参数绑定错误", nil, err)
		global.Logger.Errorf("分页获取推荐品牌时，参数绑定错误:%v", err)
		return
	}

	brand, err := brandLogic.RecommendList(rec.PageNum, rec.PageSize)
	if err != nil {
		c.Response(ctx, "分页获取推荐品牌失败", nil, err)
		return
	}
	//先分页再返回
	//pageResult := response.ResetPage(brand, int64(len(brand)), rec.PageNum, rec.PageSize)
	c.Response(ctx, "分页获取推荐品牌成功", brand, nil)
}

// 分页获取品牌相关的商品
func (c *BrandController) ProductList(ctx *gin.Context) {
	var rec receive.ProductListReqStruct
	if err := ctx.ShouldBind(&rec); err != nil {
		c.Response(ctx, "分页获取品牌相关的商品时，参数绑定错误", nil, err)
		global.Logger.Errorf("分页获取品牌相关的商品时，参数绑定错误:%v", err)
		return
	}

	result, err := brandLogic.ProductList(rec.BrandId, rec.PageNum, rec.PageSize)
	if err != nil {
		c.Response(ctx, "分页获取品牌相关的商品失败", nil, err)
		return
	}
	//先分页再返回
	pageResult := response.ResetPage(result, int64(len(result)), rec.PageNum, rec.PageSize)
	c.Response(ctx, "分页获取品牌相关的商品成功", pageResult, nil)
}
