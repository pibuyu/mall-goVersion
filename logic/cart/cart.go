package cart

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gomall/global"
	receive "gomall/interaction/receive/cart"
	response "gomall/interaction/response/cart"
	"gomall/models/cart"
	"gomall/utils/jwt"
	"sort"
	"time"
)

// AddCartItem 向购物车添加物品
func AddCartItem(data *receive.AddCartItemRequestStruct, claims *jwt.Claims) (err error) {
	//先从context的Auth中获取当前连接的用户,这一步在中间件里已经完成了，未经过认证的请求会被拦截；但是需要传递过来用户的memberId
	//将cartItem的memberId、memberNickName=用户信息、deleteStatus=0
	cartItem := ConvertToOmsCartItem(data)
	cartItem.MemberId = claims.UserID
	cartItem.MemberNickname = claims.User.Nickname
	cartItem.DeleteStatus = 0
	//根据MemberId、ProductId、DeleteStatus=0、ProductSkuId查询购物车中是否已经存在该物品
	existInCart, err := ItemIsExistInCart(cartItem)
	if err != nil {
		return err
	}
	//不存在则新插入一条数据
	if existInCart == nil {
		cartItem.CreateDate = time.Now().Format("2006-01-02 15:04:06")
		if err := InsertCartItem(cartItem); err != nil {
			return err
		}
		return nil
	}
	//存在将物品的修改日期和数量进行调整
	cartItem.ModifyDate = time.Now().Format("2006-01-02 15:04:06")
	cartItem.Quantity += existInCart.Quantity
	if err := UpdateCartItem(cartItem); err != nil {
		return err
	}
	return nil
}

func InsertCartItem(item cart.OmsCartItem) error {
	item.ModifyDate = time.Now().Format("2006-01-02 15:04:06")
	if err := global.Db.Create(&item).Error; err != nil {
		return errors.New("插入CartItem到OmsCartItem表出错:" + err.Error())
	}
	return nil
}

func UpdateCartItem(item cart.OmsCartItem) (err error) {
	if err := global.Db.Model(&cart.OmsCartItem{}).Where("id = ?", item.Id).Updates(map[string]interface{}{
		"modify_date": item.ModifyDate,
		"quantity":    item.Quantity,
	}).Error; err != nil {
		return errors.New("更新OmsCartItem表出错:" + err.Error())
	}
	return nil
}

func ConvertToOmsCartItem(request *receive.AddCartItemRequestStruct) cart.OmsCartItem {
	return cart.OmsCartItem{
		Id:                int64(request.Id), // 如果Id需要从AddCartItemRequestStruct转换
		ProductId:         int64(request.ProductID),
		ProductSkuId:      int64(request.ProductSkuID),
		MemberId:          int64(request.MemberID),
		Quantity:          request.Quantity,
		Price:             request.Price,
		ProductPic:        request.ProductPic,
		ProductName:       request.ProductName,
		ProductSubTitle:   request.ProductSubTitle,
		ProductSkuCode:    request.ProductSkuCode,
		MemberNickname:    request.MemberNickname,
		CreateDate:        request.CreateDate,
		ModifyDate:        request.ModifyDate,
		DeleteStatus:      request.DeleteStatus,
		ProductCategoryId: request.ProductCategoryID,
		ProductBrand:      request.ProductBrand,
		ProductSn:         request.ProductSN,
		ProductAttr:       request.ProductAttr,
	}
}

func ItemIsExistInCart(item cart.OmsCartItem) (*cart.OmsCartItem, error) {
	var cartItem cart.OmsCartItem
	if err := global.Db.Model(&cart.OmsCartItem{}).
		Where("member_id = ?", item.MemberId).Where("product_id = ?", item.ProductId).
		Where("delete_status = ?", 0).Find(&cartItem).Error; err != nil {
		return nil, errors.New("查询OmsCartItem表出错" + err.Error())
	}
	//todo:这一步是在特判什么？？
	if cartItem.Id > 0 {
		if cartItem.ProductSkuId == item.ProductSkuId {
			return &cartItem, nil
		}
	}
	return nil, nil
}

// todo:应该是软删除！！！
func Clear(memberId int64) error {
	if err := global.Db.Model(&cart.OmsCartItem{}).Where("member_id = ?", memberId).
		UpdateColumn("delete_status", 1).Error; err != nil {
		return errors.New("从OmsCartItem表中删除购物车物品出错:" + err.Error())
	}
	return nil
}

func DeleteCartItemsByIds(memberId int64, Ids []int64) error {
	query := global.Db.Model(&cart.OmsCartItem{}).Where("member_id = ?", memberId)
	if len(Ids) > 0 {
		query.Where("id IN (?)", Ids)
	}
	if err := query.UpdateColumn("delete_status", 1).Error; err != nil {
		return err
	}
	return nil
}

func GetProductById(data *receive.GetProductByIdRequestStruct) (results *response.GetCartProductResponseStruct, err error) {
	//oms_product表连接pms_product_attribute,再连接pms_sku_stock表
	if err := global.Db.Table("pms_product p").Debug().
		Select("p.id AS id, p.name AS name, p.sub_title AS sub_title, p.price AS price, p.pic AS pic, p.product_attribute_category_id AS productAttributeCategoryId, p.stock AS stock, pa.id AS attr_id, pa.name AS attr_name, ps.id AS sku_id, ps.sku_code AS sku_code, ps.price AS sku_price, ps.stock AS sku_stock, ps.pic AS sku_pic").
		Joins("LEFT JOIN pms_product_attribute pa ON pa.product_attribute_category_id = p.product_attribute_category_id").
		Joins("LEFT JOIN pms_sku_stock ps ON ps.product_id = p.id").
		Where("p.id = ?", data.ProductId).
		Where("pa.type = ?", 0).
		Order("pa.sort DESC").
		Find(&results).Error; err != nil {
		return nil, errors.New("获取购物车商品详情时查表出错:" + err.Error())
	}
	return results, nil
}

func List(memberId int64) (results cart.OmsCartItemList, err error) {
	//delete_status=0  &&  memberId=param
	if err := global.Db.Model(&cart.OmsCartItem{}).Where("delete_status = ?", 0).
		Where("member_id = ?", memberId).Find(&results).Error; err != nil {
		return nil, errors.New("获取购物车全部信息时查表出错:" + err.Error())
	}
	return
}

// CartListPromotion 获取包含促销活动信息的购物车列表
func CartListPromotion(data *receive.CartListPromotionRequestStruct, memberId int64) (results cart.CartPromotionItemList, err error) {
	//先找到当前会员的购物车列表
	cartItemList, err := List(memberId)
	if err != nil {
		return nil, errors.New("查询会员的购物车信息出错:" + err.Error())
	}
	//过滤一下data.cartIds，确保其存在于当前会员的cartItemList中

	filteredItemList := make([]cart.OmsCartItem, 0)
	//将会员购物车里的itemId转化为map
	cartIdSet := make(map[int64]struct{})
	for _, item := range cartItemList {
		cartIdSet[item.Id] = struct{}{}
	}

	//收集合法的ids
	for _, cartId := range data.CartIds {
		if _, exist := cartIdSet[cartId]; exist {
			for _, item := range cartItemList {
				if item.Id == cartId {
					filteredItemList = append(filteredItemList, item)
				}
			}
		}
	}
	//然后计算购物车促销信息
	if len(filteredItemList) != 0 {
		results, err = calcCartPromotion(filteredItemList)
		if err != nil {
			return nil, errors.New("计算购物车促销信息出错:" + err.Error())
		}
		return results, nil
	}
	return
}

func calcCartPromotion(cartItemList []cart.OmsCartItem) (cartPromotionItemList cart.CartPromotionItemList, err error) {
	//1.现根据productId对cartItemList进行分组，以spu为单位计算优惠
	productCartMap := groupCartItemBySpu(cartItemList)
	//2.查询所有商品的优惠相关信息
	promotionProductList, err := getPromotionProductList(cartItemList)
	if err != nil {
		return nil, err
	}
	//3.根据商品促销类型计算商品促销优惠价格
	cartPromotionItemList = make([]cart.CartPromotionItem, 0)
	for productId, itemList := range productCartMap {
		promotionProduct := getPromotionProductById(productId, promotionProductList)
		promotionType := promotionProduct.Product.PromotionType
		if promotionType == 1 {
			for _, item := range itemList {
				cartPromotionItem := copyFromOmsCartItem(item)
				cartPromotionItem.PromotionMessage = "单品促销"
				//商品原价-促销价
				skuStock := getOriginalPrice(promotionProduct, item.ProductSkuId)
				originalPrice := skuStock.Price
				//单品促销使用原价
				cartPromotionItem.CartItem.Price = originalPrice
				cartPromotionItem.ReduceAmount = originalPrice - skuStock.PromotionPrice
				cartPromotionItem.RealStock = skuStock.Stock - skuStock.LockStock
				cartPromotionItem.Integration = promotionProduct.Product.GiftPoint
				cartPromotionItem.Growth = promotionProduct.Product.GiftGrowth
				cartPromotionItemList = append(cartPromotionItemList, cartPromotionItem)
			}
		} else if promotionType == 3 {
			//打折优惠
			count := getCartItemCount(itemList)
			ladder := getProductLadder(count, promotionProduct.ProductLadderList)
			if ladder.Id != 0 {
				for _, item := range itemList {
					cartPromotionItem := copyFromOmsCartItem(item)
					message := getLadderPromotionMessage(ladder)
					cartPromotionItem.PromotionMessage = message
					//商品原价-折扣*商品原价
					skuStock := getOriginalPrice(promotionProduct, item.ProductSkuId)
					originalPrice := skuStock.Price
					reduceAmount := originalPrice - ladder.Discount*originalPrice
					cartPromotionItem.ReduceAmount = reduceAmount
					cartPromotionItem.RealStock = skuStock.Stock - skuStock.LockStock
					cartPromotionItem.Integration = promotionProduct.Product.GiftPoint
					cartPromotionItem.Growth = promotionProduct.Product.GiftGrowth
					cartPromotionItemList = append(cartPromotionItemList, cartPromotionItem)
				}
			}
		} else {
			handleNoReduce(cartPromotionItemList, itemList, promotionProduct)
		}
	}
	return cartPromotionItemList, nil
}

// 以spu为单位对购物车中商品进行分组
func groupCartItemBySpu(cartItemList []cart.OmsCartItem) (results map[int64][]cart.OmsCartItem) {
	results = make(map[int64][]cart.OmsCartItem)
	for _, cartItem := range cartItemList {
		//获得当前商品id对应的购物车项列表
		productCartItemList, exist := results[cartItem.ProductId]

		if !exist {
			//如果不存在，创建新的列表并添加当前项
			productCartItemList = []cart.OmsCartItem{cartItem}
			results[cartItem.ProductId] = productCartItemList
		} else {
			//如果已经存在,将当前项添加到列表中
			productCartItemList = append(productCartItemList, cartItem)
			results[cartItem.ProductId] = productCartItemList
		}
	}
	return results
}

// 查询所有商品的优惠相关信息
func getPromotionProductList(cartItemList []cart.OmsCartItem) (results []cart.PromotionProduct, err error) {
	productIdList := make([]int64, 0)
	for _, cartItem := range cartItemList {
		productIdList = append(productIdList, cartItem.ProductId)
	}
	//然后根据这些ids去查询优惠信息
	results, err = getProductPromotionByIds(productIdList)
	if err != nil {
		return nil, err
	}
	return
}

func getProductPromotionByIds(productIdList []int64) (results []cart.PromotionProduct, err error) {
	if err = global.Db.Table("pms_product p").
		Select(`
			p.id,
			p.name,
			p.promotion_type,
			p.gift_growth,
			p.gift_point,
			sku.id AS sku_id,
			sku.price AS sku_price,
			sku.sku_code AS sku_sku_code,
			sku.promotion_price AS sku_promotion_price,
			sku.stock AS sku_stock,
			sku.lock_stock AS sku_lock_stock,
			ladder.id AS ladder_id,
			ladder.count AS ladder_count,
			ladder.discount AS ladder_discount,
			full_re.id AS full_id,
			full_re.full_price AS full_full_price,
			full_re.reduce_price AS full_reduce_price`).
		Joins("LEFT JOIN pms_sku_stock sku ON p.id = sku.product_id").
		Joins("LEFT JOIN pms_product_ladder ladder ON p.id = ladder.product_id").
		Joins("LEFT JOIN pms_product_full_reduction full_re ON p.id = full_re.product_id").
		Where("p.id IN ?", productIdList).Find(&results).Error; err != nil {
		return nil, errors.New("根据商品ids查询优惠信息出错:" + err.Error())
	}
	return
}

// 根据商品id获取商品的促销信息
func getPromotionProductById(productId int64, promotionProductList []cart.PromotionProduct) (results cart.PromotionProduct) {
	for _, promotionProduct := range promotionProductList {
		if promotionProduct.Product.Id == productId {
			return promotionProduct
		}
	}
	return
}

func copyFromOmsCartItem(item cart.OmsCartItem) (result cart.CartPromotionItem) {
	var promotionItem cart.CartPromotionItem
	promotionItem.CartItem = item
	return promotionItem
}

func getOriginalPrice(promotionProduct cart.PromotionProduct, skuId int64) (result cart.PmsSkuStock) {
	for _, skuStock := range promotionProduct.SkuStockList {
		if skuId == skuStock.Id {
			return skuStock
		}
	}
	return cart.PmsSkuStock{}
}

func getCartItemCount(itemList []cart.OmsCartItem) (count int) {
	result := 0
	for _, item := range itemList {
		result += item.Quantity
	}
	return result
}

// 根据购买商品数量获取满足条件的打折优惠策略
func getProductLadder(count int, productLadderList []cart.PmsProductLadder) (productLadder cart.PmsProductLadder) {
	//按数量从大到小排序
	sort.Slice(productLadderList, func(i, j int) bool {
		return productLadderList[i].Count > productLadderList[j].Count
	})
	for _, productLadder := range productLadderList {
		if count >= productLadder.Count {
			return productLadder
		}
	}
	//返回零值
	return cart.PmsProductLadder{}
}

func getLadderPromotionMessage(ladder cart.PmsProductLadder) string {
	return fmt.Sprintf("打折优惠：满%d件，打%d折", ladder.Count, int(ladder.Discount*10))
}

// 对没满足优惠条件的商品进行处理
func handleNoReduce(cartPromotionItemList []cart.CartPromotionItem, itemList []cart.OmsCartItem, promotionProduct cart.PromotionProduct) {
	for _, item := range itemList {
		cartPromotionItem := copyFromOmsCartItem(item)
		cartPromotionItem.PromotionMessage = "无优惠"
		cartPromotionItem.ReduceAmount = 0
		skuStock := getOriginalPrice(promotionProduct, item.ProductSkuId)
		if skuStock.Id != 0 {
			cartPromotionItem.RealStock = skuStock.Stock - skuStock.LockStock
		}
		cartPromotionItem.Integration = promotionProduct.Product.GiftPoint
		cartPromotionItem.Growth = promotionProduct.Product.GiftGrowth
		cartPromotionItemList = append(cartPromotionItemList, cartPromotionItem)
	}
}

// todo:没转化过来，除了手动赋值的两个字段，插入的数据都是空的
func UpdateAttr(data *receive.UpdateAttrRequestStruct) (err error) {
	//从data里结构出item对象，然后删除原本的数据，插入新的数据
	cartItem := &cart.OmsCartItem{}
	jsonData, _ := json.Marshal(data)
	if err = json.Unmarshal(jsonData, cartItem); err != nil {
		//如果这种映射方法不可行，采取别的映射方法，这里不应该错误返回,因为至少可以手动逐字段映射保证成功
		cartItem = mapRequestDataToCartItem(data)
	}
	//先删除旧的表项
	if err := global.Db.Model(&cart.OmsCartItem{}).Where("id = ?", cartItem.Id).UpdateColumn("delete_status", 1).Error; err != nil {
		return errors.New("从OmsCartItem表中删除旧数据出错:" + err.Error())
	}
	//用cartItem去更新表项
	cartItem.CreateDate = time.Now().Format("2006-01-02 15:04:05")
	cartItem.ModifyDate = time.Now().Format("2006-01-02 15:04:05")
	global.Logger.Infof("即将插入的新数据为%v", cartItem)
	if err := global.Db.Model(&cart.OmsCartItem{}).Create(cartItem).Error; err != nil {
		return errors.New("在OmsCartItem插入新数据出错:" + err.Error())
	}
	return nil
}

func mapRequestDataToCartItem(data *receive.UpdateAttrRequestStruct) (cartItem *cart.OmsCartItem) {
	if err := mapstructure.Decode(data, &cartItem); err != nil {
		//如果还是不对，那就手动映射
		result := &cart.OmsCartItem{
			Id:                data.Id,
			ProductId:         data.ProductId,
			ProductSkuId:      data.ProductSkuId,
			MemberId:          data.MemberId,
			Quantity:          data.Quantity,
			Price:             data.Price,
			ProductPic:        data.ProductPic,
			ProductName:       data.ProductName,
			ProductSubTitle:   data.ProductSubTitle,
			ProductSkuCode:    data.ProductSkuCode,
			MemberNickname:    data.MemberNickname,
			CreateDate:        data.CreateDate,
			ModifyDate:        data.ModifyDate,
			DeleteStatus:      data.DeleteStatus,
			ProductCategoryId: data.ProductCategoryId,
			ProductBrand:      data.ProductBrand,
			ProductSn:         data.ProductSn,
			ProductAttr:       data.ProductAttr,
		}
		return result
	}
	return
}

func UpdateQuantity(data *receive.UpdateQuantityRequestStruct, memberId int64) (err error) {
	//三个条件：删除状态=0，memberId对上，Id对上
	if err = global.Db.Model(&cart.OmsCartItem{}).
		Where("id = ?", data.Id).Where("member_id = ?", memberId).Where("delete_status = ?", 0).
		UpdateColumn("quantity", data.Quantity).Error; err != nil {
		return errors.New("更新OmsCartItem表的quantity出错:" + err.Error())
	}
	return nil
}
