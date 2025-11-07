package v1

import (
	"baby/models"
	"github.com/gin-gonic/gin"
)

//商品加购

func ShopperShopCart(c *gin.Context) {
	context := gin.H{"state": "success", "msg": "获取成功"}
	userId, _ := c.Get("userId")
	if c.Request.Method == "GET" {
		if userId != 0 {
			var carts models.Carts
			models.DB.Preload("Commodities").Where("user_id = ?", userId.(int64)).Order("id desc").Find(&carts)
			//找这个人的所有订单
			context["data"] = carts
		}
	}

	if c.Request.Method == "POST" {
		context = gin.H{"state": "fail", "msg": "加载失败"}
		var body struct {
			Id       int64 `json:"id"`
			Quantity int64 `json:"quantity"`
		}
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(200, gin.H{"state": "fail", "msg": err.Error()})
		}
		id := body.Id
		quantity := body.Quantity
		var commodity models.Commodities
		models.DB.Where("id = ?", id).First(&commodity)
		//查找商品是否存在
		if commodity.ID > 0 {
			//购物车同一商品，只增加购买数量
			var cart models.Carts
			models.DB.Where("commodity_id = ? and user_id = ?", id, userId).First(&cart)
			if cart.ID > 0 {
				cart.Quantity = quantity + cart.Quantity
				models.DB.Save(&cart)
			} else {
				carts := models.Carts{UserId: userId.(int64), CommodityId: id, Quantity: quantity}
				models.DB.Create(&carts)
			}
			context = gin.H{"state": "success", "msg": "加购成功"}
		}
	}
	c.JSON(200, context)
}

func ShopperDelete(c *gin.Context) {
	var body struct {
		CartId int64 `json:"cartId"`
	}
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(200, gin.H{"state": "fail", "msg": err.Error()})
	}
	var cart []models.Carts
	cartId := body.CartId

	if cartId != 0 {
		models.DB.Where("id = ?", cartId).Find(&cart)
		//查cart表的主键
	} else {
		userId, _ := c.Get("userId")
		models.DB.Where("user_id = ?", userId.(int64)).Find(&cart)
	}
	models.DB.Unscoped().Delete(&cart)
	context := gin.H{"state": "success", "msg": "删除成功"}
	c.JSON(200, context)

}

func ShopperHome(c *gin.Context) {
	context := gin.H{"state": "success", "msg": "获取成功"}
	data := gin.H{}
	userId, _ := c.Get("userId")
	payInfo := c.DefaultQuery("out_trade_no", "")

	if payInfo != "" {
		models.DB.Model(&models.Orders{}).Where("pay_info = ?", payInfo).Update("state", 1)
	}
	if userId != 0 {
		var orders []models.Orders
		models.DB.Where("user_id = ?", userId).Order("id desc").Find(&orders)
		data["orders"] = data
	}
	context["data"] = data
	c.JSON(200, context)
}
