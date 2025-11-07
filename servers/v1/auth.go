package v1

import (
	"baby/middleware"
	"baby/models"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"time"
)

// ShopperLogin 注册/登录
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

// ShopperLogout 退出账号
func ShopperLogout(c *gin.Context) {
	context := gin.H{"state": "fail", "msg": "退出失败"}
	userId, _ := c.Get("userId")

	if userId != 0 {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			var jwts models.Jwts
			models.DB.Where("token = ?", authHeader).First(&jwts)
			models.DB.Where("token = ?", authHeader).Unscoped().Delete(&jwts)
			context["state"] = "success"
			context["msg"] = "退出成功"
		}
	}
	c.JSON(200, context)
}

// ShopperCancel 注销账号
func ShopperCancel(c *gin.Context) {}
