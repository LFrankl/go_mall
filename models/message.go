package models

// LoginLogMessage 登录日志消息结构体（用于 RabbitMQ 传输）
type LoginLogMessage struct {
	OrderID     int64   `json:"order_id"`
	UserID      int64   `json:"user_id"`
	CommodityID int64   `json:"commodity_id"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
}
