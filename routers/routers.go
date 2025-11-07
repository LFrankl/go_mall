package routers

import (
	"baby/middleware"
	v1 "baby/servers/v1"
	"baby/settings"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {

	gin.SetMode(settings.Mode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.StaticFS("/static", http.Dir("static"))
	//注册静态文件服务，将本地文件系统中的 static 目录映射到 HTTP 路径 /static，使客户端可以通过 URL 访问该目录下的静态资源（如图片、CSS、JS 文件等）。

	//配置跨域访问
	config := cors.DefaultConfig()

	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"}

	r.Use(cors.New(config))

	apiv1 := r.Group("/api/v1/")

	commodity := apiv1.Group("")
	{
		//首页
		commodity.GET("home/", v1.Home)
		//商品列表
		commodity.GET("commodity/list/", v1.CommodityList)
		//商品详细
		commodity.GET("commodity/detail/:id/", v1.CommodityDetail)

	}

	auth := apiv1.Group("")
	{
		//用户注册登录
		auth.POST("auth/login/", v1.ShopperLogin)
		//退出登录
		auth.POST("auth/logout/", v1.ShopperLogout)
		//用户注销
		auth.POST("auth/cancel/", v1.ShopperCancel)
	}

	shopper := apiv1.Group("", middleware.JWTAuthMiddleware)
	{
		//商品收藏
		shopper.POST("commodity/collect/", v1.CommodityCollect)

		//个人主页
		shopper.GET("shopper/home/", v1.ShopperHome)
		//加入购物车
		shopper.POST("shopper/shopcart/", v1.ShopperShopCart)
		//购物车列表
		shopper.GET("shopper/shopcart/", v1.ShopperShopCart)
		//在线支付
		//shopper.POST("shopper/pays/", v1.ShopperPays)
		//删除购物车商品
		shopper.POST("shopper/delete/", v1.ShopperDelete)
	}
	return r

}
