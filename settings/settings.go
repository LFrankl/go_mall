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

// MQConfig 消息队列配置（支持多队列）
type MQConfig struct {
	Type     string            // 消息队列类型（rabbitmq/kafka等）
	Address  string            // 连接地址
	QueueMap map[string]string // 队列名映射（业务类型 → 实际队列名）
}

// MQSetting 初始化多队列配置
var MQSetting = &MQConfig{
	Type:    "rabbitmq",
	Address: "amqp://guest:guest@localhost:5672/",
	QueueMap: map[string]string{
		"login":   "login_logs",
		"order":   "order_logs",
		"payment": "payment_logs",
		"user":    "user_operate",
	},
}
