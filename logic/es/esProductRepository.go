package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"gomall/global"
	"gomall/models/es"
)

func SaveAll(esProductList []es.EsProduct) (count int, err error) {
	// 创建一个批量操作请求
	bulkRequest := global.Es.Bulk()

	if len(esProductList) == 0 {
		return 0, nil
	}
	for _, product := range esProductList {
		// 为每个 EsProduct 创建一个批量操作请求
		request := elastic.NewBulkIndexRequest().
			Index("products").
			Id(fmt.Sprintf("%d", product.ID)).
			Doc(product)
		bulkRequest = bulkRequest.Add(request)
	}
	//执行批量操作
	response, err := bulkRequest.Do(context.Background())
	if err != nil {
		global.Logger.Errorf("将商品信息批量插入es出错：%v", err)
		return 0, fmt.Errorf("批量保存 EsProduct 失败: %v", err)
	}
	// 打印批量操作的结果
	if response.Errors {
		global.Logger.Errorf("批量操作存在错误: %v", response.Items)
		return 0, fmt.Errorf("批量操作存在错误: %v", response.Items)
	}
	global.Logger.Infof("成功保存 %d 个产品到 Elasticsearch", len(esProductList))
	return len(esProductList), nil
}

func Delete(id int64) error {
	// 通过商品 ID 构建 Elasticsearch 文档 ID
	docID := fmt.Sprintf("%d", id)

	// 执行删除请求
	_, err := global.Es.Delete().
		Index("products").       // 指定索引名，这里假设是 "products"
		Id(docID).               // 使用商品的 ID 作为文档的 ID
		Do(context.Background()) // 执行删除请求

	// 错误处理
	if err != nil {
		global.Logger.Errorf("删除ID= %d 的商品时出错：%v", id, err)
		return fmt.Errorf("删除商品失败: %v", err)
	}

	// 如果删除成功，输出日志
	global.Logger.Infof("成功删除ID= %d 的商品", id)
	return nil
}

func Save(product es.EsProduct) error {
	// 将商品 ID 转换为字符串作为文档的 ID
	docID := fmt.Sprintf("%d", product.ID)

	// 执行插入操作
	_, err := global.Es.Index().
		Index("products").       // 指定索引名称，假设是 "products"
		Id(docID).               // 商品 ID 作为文档的 ID
		BodyJson(product).       // 设置商品数据作为请求的主体
		Do(context.Background()) // 执行插入操作

	// 错误处理
	if err != nil {
		global.Logger.Errorf("插入商品（ID: %d）到 Elasticsearch 时出错：%v", product.ID, err)
		return fmt.Errorf("插入商品到 Elasticsearch 失败: %v", err)
	}

	// 插入成功后，打印日志
	global.Logger.Infof("成功插入商品（ID: %d）到 Elasticsearch", product.ID)
	return nil
}

// SearchProductsByKeyword 根据关键词和分页参数从ES中搜索商品
func SearchProductsByKeyword(keyword string, pageNum int, pageSize int) ([]es.EsProduct, error) {
	// 确保页码和页大小有效
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 计算起始位置
	from := (pageNum - 1) * pageSize

	// 构建 MultiMatchQuery
	query := elastic.NewMultiMatchQuery(keyword, "name", "subTitle", "keywords")

	// 执行搜索查询
	searchResult, err := global.Es.Search().
		Index("products").       // 指定索引名称
		Query(query).            // 设置查询条件
		From(from).              // 设置分页起始位置
		Size(pageSize).          // 设置每页大小
		Do(context.Background()) // 执行查询
	if err != nil {
		global.Logger.Errorf("搜索商品时出错: %v", err)
		return nil, fmt.Errorf("搜索商品失败: %v", err)
	}

	// 解析查询结果
	var results []es.EsProduct
	for _, hit := range searchResult.Hits.Hits {
		var product es.EsProduct
		// 反序列化 JSON 数据到 EsProduct 结构体
		if err := json.Unmarshal(hit.Source, &product); err != nil {
			global.Logger.Errorf("解析搜索结果出错: %v", err)
			continue
		}
		results = append(results, product)
	}

	global.Logger.Infof("成功获取 %d 个商品", len(results))
	return results, nil
}
