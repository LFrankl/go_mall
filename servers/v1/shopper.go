package v1

import (
	"baby/middleware"
	"baby/models"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"time"
)

func ShopperLogin(c *gin.Context) {
	//map[string]interface{}{}
	context := gin.H{
		"state": "fail",
		"msg":   "注册或登录失败"}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(200, gin.H{"state": "fail", "msg": err.Error()})
	}
	username := body.Username
	p := body.Password
	//fmt.Printf("原始密码: %q, 长度: %d\n", p, len(p)) // %q 会显示引号和转义字符（如 \n）
	//fmt.Println("username:", username)
	//fmt.Println("password:", p)
	if username != "" && p != "" {
		context["state"] = "success"
		//生成登陆时间
		lastLogin := time.Now()
		context["last_login"] = lastLogin.Format("2006-01-02 15:04:05")
		//密码加密
		m := md5.New()
		m.Write([]byte(p))
		password := hex.EncodeToString(m.Sum(nil))
		//查找用户，用户存在则登陆成功，不存在则创建
		//fmt.Println("加密后password:", password)

		var userID uint
		var users models.Users
		models.DB.Where("username = ?", username).First(&users)
		//fmt.Println("取到的users密码:", users.Password)
		if users.ID > 0 {
			if users.Password == password {
				userID = users.ID
				users.LastLogin = lastLogin
				models.DB.Save(&users)
				context["msg"] = "登录成功"
			} else {
				context["msg"] = "请输入正确密码"
				context["state"] = "fail"
			}
		} else {
			context["msg"] = "注册成功"
			r := models.Users{
				Username:  username,
				Password:  password,
				IsStaff:   1,
				LastLogin: lastLogin,
			}
			//fmt.Println("即将存入数据库的密码:", r.Password) // 若不是 d8e7...，说明赋值错误
			models.DB.Create(&r)
			if r.ID > 0 {
				userID = r.ID
			} else {
				context["state"] = "fail"
				context["msg"] = "注册失败"
			}
		}
		//创建token
		token := ""
		if userID > 0 {
			token, _ = middleware.GenToken(username, int64(userID))
		}

		context["token"] = token

	}
	c.JSON(200, context)

}

func ShopperLogout(c *gin.Context) {
	context := gin.H{"state": "fail", "msg": "退出失败"}
	userId, _ := c.Get("userId")

	if userId != 0 {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			var jwts models.Jwts
			models.DB.Where("token = ?", authHeader).First(&jwts)
			models.DB.Unscoped().Delete(&jwts)
			context["state"] = "success"
			context["msg"] = "退出成功"
		}
	}
	c.JSON(200, context)
}

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
