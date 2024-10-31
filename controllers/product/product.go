package product

import (
	"github.com/gin-gonic/gin"
	controller "gomall/controllers"
	"gomall/global"
	receive "gomall/interaction/receive/product"
	productLogic "gomall/logic/product"
)

type ProductController struct {
	controller.BaseControllers
}

// 获取前台商品详情
func (c *ProductController) Detail(ctx *gin.Context) {
	var rec receive.DetailReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("Detail请求传入参绑定失败: %v", err)
		c.Response(ctx, "请求参数错误", nil, err)
		return
	}
	result, err := productLogic.Detail(rec.Id)
	if err != nil {
		c.Response(ctx, "获取商品详情失败", nil, err)
	}
	c.Response(ctx, "获取商品详情成功", result, nil)
}

// CategoryTreeList 以树形结构获取所有商品分类
func (c *ProductController) CategoryTreeList(ctx *gin.Context) {
	result, err := productLogic.CategoryTreeList()
	if err != nil {
		c.Response(ctx, "以树形结构获取所有商品分类失败", nil, err)
	}
	c.Response(ctx, "以树形结构获取所有商品分类成功", result, nil)
}

// Search 综合搜索、筛选、排序
func (c *ProductController) Search(ctx *gin.Context) {
	var rec receive.SearchReqStruct
	if err := ctx.ShouldBindJSON(&rec); err != nil {
		global.Logger.Errorf("综合搜索、筛选、排序请求传入参绑定失败: %v", err)
		c.Response(ctx, "综合搜索、筛选、排序请求参数错误", nil, err)
		return
	}
	result, err := productLogic.Search(&rec)
	if err != nil {
		c.Response(ctx, "综合搜索、筛选、排序失败", nil, err)
	}
	c.Response(ctx, "综合搜索、筛选、排序成功", result, nil)

}
