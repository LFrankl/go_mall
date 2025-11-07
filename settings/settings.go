package settings

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Database struct {
	User     string
	Password string
	Host     string
	Name     string
}

var MySQLSetting = &Database{
	User:     "root",
	Password: "123456",
	Host:     "127.0.0.1",
	Name:     "go_mall",
}

// gin.DebugMode
//gin.ReleaseMode
//TestMode

var Mode = gin.DebugMode

//jwt有效时间

var TokenExpireDuration = time.Minute * 30

//jwt加密盐

var Secret = []byte("你好")

//分页

var PageSize = 6

//支付信息

var AppId = ""
var AlipayPublicKeyString = ``
var AppPrivateKeyString = ``
