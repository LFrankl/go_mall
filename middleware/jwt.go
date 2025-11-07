package middleware

import (
	"baby/models"
	"baby/settings"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

/*
在我的逻辑中，每次触发登录操作，都会加一个jwts记录
理论上用户只能退出后再登录，所以能保证一直都只有一个记录
但是如果api地址被人用脚本恶意重复登录，这样jwts表会不会爆炸？

*/

type CustomClaims struct {
	Username string `json:"username"`
	UserId   int64  `json:"userId"`
	jwt.RegisteredClaims
}

// GenToken 生成jwt
func GenToken(username string, userId int64) (string, error) {
	expire := time.Now().Add(settings.TokenExpireDuration)

	claims := CustomClaims{
		username, //自定义用户名字段
		userId,
		jwt.RegisteredClaims{
			Issuer: "奥里给",
		},
	}
	//使用指定的签名方法确定签名对象
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(settings.Secret)
	if err != nil {
		return "", err
	}

	models.DB.Where("token = ?", token).Unscoped().Delete(&models.Jwts{})
	//使用指定的secret签名并获得完整编码后的字符串token
	//token写入数据库
	j := models.Jwts{Token: token, Expire: expire}
	models.DB.Create(&j)
	return token, nil
}

// ParseToken 解析jwt
func ParseToken(tokenString string) (*CustomClaims, error) {
	//验证 Token 的签名合法性，确保 Token 是由持有相同密钥的服务器签发，且传输过程中未被篡改
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return settings.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	//校验
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	//token.Claims：JWT 解析后得到的载荷数据（类型为 jwt.Claims 接口）。
	//.(*CustomClaims)：类型断言操作符，声明 “我认为 token.Claims 实际是 *CustomClaims 类型”。
	//返回值 claims：如果断言成功，claims 就是 *CustomClaims 类型的指针，可直接访问其字段（如 claims.Username）。
	//返回值 ok：布尔值，true 表示断言成功（token.Claims 确实是 *CustomClaims 类型），false 表示断言失败（类型不匹配）。
	return nil, errors.New("invalid token")
}

func JWTAuthMiddleware(c *gin.Context) {
	//客户端携带token有三种方式：1.放在请求头2.放在请求体 3.放在URI
	//把Token放在Header的Authorization中
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(200, gin.H{
			"state": "fail",
			"msg":   "请求头的Authorization为空",
		})
		c.Abort()
		return
	}
	mc, err := ParseToken(authHeader)
	if err != nil {
		c.JSON(200, gin.H{
			"state": "fail",
			"msg":   "无效的token",
		})
		c.Abort()
		return
	}
	var jwts models.Jwts
	models.DB.Where("token = ?", authHeader).First(&jwts)
	if jwts.Token != "" {
		if jwts.Expire.After(time.Now()) {
			jwts.Expire = time.Now().Add(settings.TokenExpireDuration)
			models.DB.Save(&jwts)
		} else {
			//强制删除表数据
			models.DB.Unscoped().Delete(&jwts)
			//默认 Delete 是软删除（更新 DeletedAt），结合 Unscoped() 可执行硬删除（真正从数据库中删除记录）：
		}
	} else {
		c.JSON(200, gin.H{
			"state": "fail",
			"msg":   "无效的token",
		})
		c.Abort()
		return
	}
	//把当前请求的username信息保存在请求的上下文c上
	//路由函数通过c.get获取当前用户的信息
	c.Set("username", mc.Username)
	c.Set("userId", mc.UserId)
	c.Next()

}
