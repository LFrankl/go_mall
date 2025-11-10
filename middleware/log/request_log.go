package log

import (
	"baby/rabbitmq"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"time"
)

//实现 r.POST("/login",log_func(LogIn,c))

var GlobalCli *rabbitmq.Client

func RegisterCli(c *rabbitmq.Client) {
	GlobalCli = c
}

//外部注册全局client实例

func LogFunc(c *gin.Context) {

	startTime := time.Now()

	//endTime := time.Now()
	//costTime := endTime.Sub(startTime)

	logData := map[string]interface{}{
		"path":   c.FullPath(),     // 请求路径
		"method": c.Request.Method, // 请求方法
		"ip":     c.ClientIP(),     // 客户端IP
		"start":  startTime,        // 开始时间
		"url":    c.Request.URL.String(),
	}

	go func(data map[string]interface{}) {

		_log, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		err = GlobalCli.Publish(_log, rabbitmq.LogExchangeName, rabbitmq.LogOfRequestRoutingKey)
		if err != nil {
			fmt.Println(err)
		}
		// 这里实现发送到消息队列的逻辑（如RabbitMQ）
		// 示例：json.Marshal(data) 后调用 Publish 方法
		// if err := rabbitmq.Publish("login.log", data); err != nil {
		//     log.Printf("日志发送失败: %v", err)
		// }
	}(logData)
	c.Next()

}
