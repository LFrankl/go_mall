package log

import (
	"github.com/gin-gonic/gin"
	"time"
)

//实现 r.POST("/login",log_func(LogIn,c))

// 修正参数顺序，符合gin中间件规范（先接收handler，返回包装后的handler）
func logFunc(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 请求前记录开始时间
		startTime := time.Now()

		// 2. 执行原始处理函数（如登录逻辑）
		handler(c)

		// 3. 请求后收集日志信息
		endTime := time.Now()
		costTime := endTime.Sub(startTime) // 计算耗时

		// 收集日志数据（可根据需要扩展）
		logData := map[string]interface{}{
			"path":   c.FullPath(),      // 请求路径
			"method": c.Request.Method,  // 请求方法
			"ip":     c.ClientIP(),      // 客户端IP
			"start":  startTime,         // 开始时间
			"cost":   costTime.String(), // 耗时
			"status": c.Writer.Status(), // 响应状态码
			// 可添加用户ID、请求参数等信息
		}

		// 4. 异步发送到消息队列
		go func(data map[string]interface{}) {
			// 这里实现发送到消息队列的逻辑（如RabbitMQ）
			// 示例：json.Marshal(data) 后调用 Publish 方法
			// if err := rabbitmq.Publish("login.log", data); err != nil {
			//     log.Printf("日志发送失败: %v", err)
			// }
		}(logData)
	}
}
