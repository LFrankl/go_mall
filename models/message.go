package models

import "time"

// LogRequestMessage api请求日志消息结构体（用于 RabbitMQ 传输）
type LogRequestMessage struct {
	Path   string
	Ip     string
	Start  time.Time
	Url    string
	Method string
}
