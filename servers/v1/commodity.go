package v1

import (
	"baby/middleware"
	"baby/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"strconv"
)

// Home 商城首页
func Home(c *gin.Context) {
	context := gin.H{"state": "success", "msg": "获取成功"}
	data := gin.H{}

	//今日必抢商品信息
	var commodity []models.Commodities
	models.DB.Order("sold DESC").Find(&commodity)
	//按 sold 字段（销量）从高到低排序，并将结果存入 commodity 变量中
	//ASC：升序排序（默认，可省略），从小到大排列。
	//示例：Order("price ASC") → 按 price 字段从低到高排序。
	//DESC：降序排序，从大到小排列。
	//示例：Order("created_at DESC") → 按 created_at（创建时间）从新到旧排序。

	data["commodityInfos"] = [][]models.Commodities{commodity[0:4], commodity[4:8]}

	//分类商品信息
	var classification = map[string]string{"clothes": "儿童服饰", "food": "奶粉辅食", "goods": "儿童用品"}

	for k, v := range classification {
		var types = []string{}
		var temp []models.Commodities

		models.DB.Model(&models.Types{}).Where("firsts = ?", v).Select("seconds").Find(&types)
		models.DB.Where("types in ?", types).Order("sold DESC").Find(&temp)
		data[k] = temp[0:5]
	}
	context["data"] = data
	c.JSON(200, context)
}

func CommodityList(c *gin.Context) {
	context := gin.H{"state": "success", "msg": "success"}
	data := gin.H{} //待返回的结构体
	//http://xxx.com/goods?type=clothes&search=T恤&sort=sold&page=2
	types := c.DefaultQuery("type", "")
	search := c.DefaultQuery("search", "")
	sort := c.DefaultQuery("sort", "")
	page := c.DefaultQuery("page", "1")
	p, _ := strconv.Atoi(page)

	//商品分类列表
	var firsts = []string{}
	//models.DB.Model(&models.Types{}).Where("firsts = ?", types).Select("firsts").Find(&firsts)
	models.DB.Model(&models.Types{}).Distinct("firsts").Find(&firsts)
	var res []map[string]interface{}

	for _, f := range firsts {
		var seconds = []string{}
		models.DB.Model(&models.Types{}).Where("firsts = ?", f).Select("seconds").Find(&seconds)
		res = append(res, map[string]interface{}{"name": f, "value": seconds})
		//对每一个一级分类下，进行批量查询
	}
	data["types"] = res
	//商品列表信息

	var commodity []models.Commodities
	querys := models.DB.Model(&models.Commodities{})
	if types != "" {
		querys = querys.Where("types = ?", types)
	}
	if sort != "" {
		//这里指的是sort 可能是按价格，按销量
		querys = querys.Order(sort + " DESC")
	}
	if search != "" {
		querys = querys.Where("name like ?", "%"+search+"%")
	}
	//根据是否有搜索关键词（search）来动态添加模糊查询条件

	querys, previous, next, count, pageCount := models.Paginate(querys, p)
	querys = querys.Find(&commodity)
	data["commodityInfos"] = map[string]interface{}{
		"data":      commodity,
		"previous":  previous,
		"next":      next,
		"count":     count,
		"pageCount": pageCount,
	}
	context["data"] = data
	c.JSON(200, context)

}

func CommodityDetail(c *gin.Context) {

	//api/v1/commodity/detail/:id/

	context := gin.H{"state": "success", "msg": "获取成功"}
	data := gin.H{}
	id := c.Param("id")
	//获取商品详细信息
	var commodity models.Commodities
	models.DB.Where("id = ?", id).First(&commodity)
	data["commodity"] = commodity

	//获取推荐商品
	var recommend []models.Commodities
	models.DB.Where("id != ?", id).Order("sold DESC").Limit(5).Find(&recommend)
	data["recommend"] = recommend

	//收藏状态
	data["likes"] = false
	//获取请求头的Authorization
	//获取用户信息，根据用户信息和商品id查找收藏记录
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "" {
		mc, _ := middleware.ParseToken(authHeader)
		if mc != nil {
			UserId := mc.UserId
			if UserId != 0 {
				var records []models.Records
				models.DB.Where("user_id = ? and commodity_id = ?", UserId, id).Find(&records)
				if len(records) > 0 {
					data["likes"] = true
				}
			}
		}
	}
	context["data"] = data
	c.JSON(200, context)
}

func CommodityCollect(c *gin.Context) {
	context := gin.H{"state": "fail", "msg": "收藏失败"}
	data, _ := c.GetRawData()
	var body map[string]int64
	err := json.Unmarshal(data, &body)
	if err != nil {
		fmt.Println("error:", err)
	}
	id := body["id"]
	userId, _ := c.Get("userId")
	var records []models.Records
	models.DB.Where("user_id = ? and commodity_id = ?", userId.(int64), id).Find(&records)

	if len(records) > 0 {
		context["msg"] = "收藏取消"
		context["state"] = "success"
		models.DB.Unscoped().Delete(&records)
	} else {
		models.DB.Model(&models.Commodities{}).Where("id = ?", id).Update("likes", 1)
		r := models.Records{UserId: userId.(int64), CommodityId: id}
		models.DB.Create(&r)
		context["msg"] = "收藏成功"
		context["state"] = "success"
	}
	c.JSON(200, context)
}
