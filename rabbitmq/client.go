package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"time"
)

const (
	OrderExchangeName = "order.exchange" // 订单交换机
	OrderQueueName    = "order.queue"    // 订单队列
	OrderRoutingKey   = "order.routing"  // 路由键

	LogExchangeName        = "log.exchange" // 日志交换机，内部对应多个队列和路由键
	LogOfRequestQueueName  = "log.queue"
	LogOfRequestRoutingKey = "log.routing"
)

type Client struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// NewClient 创建 RabbitMQ 客户端
func NewClient(addr, user, pass string) (*Client, error) {
	conn, err := amqp091.Dial("amqp://" + user + ":" + pass + "@" + addr + "/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{conn: conn, channel: ch}, nil
}

// DeclareExchangeQueue 声明订单相关的交换机和队列（direct 类型，确保消息可靠投递）
func (c *Client) DeclareExchangeQueue(exchangeName string,
	queueName string,
	routeKey string) error {
	// 声明交换机（持久化）
	if err := c.channel.ExchangeDeclare(
		exchangeName,
		amqp091.ExchangeDirect, // 直接交换机，精确匹配路由键
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// 声明队列
	_, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 绑定队列到交换机（指定路由键）
	return c.channel.QueueBind(
		queueName,
		routeKey,
		exchangeName,
		false,
		nil,
	)
}

// Consume 消费消息（手动确认）
func (c *Client) Consume(queueName string, consumerTag string) (<-chan amqp091.Delivery, error) {
	return c.channel.Consume(
		queueName,
		consumerTag, // 消费者标签
		false,       // 关闭自动确认
		false,
		false,
		false,
		nil,
	)
}

/*
	LogExchangeName        = "log.exchange" // 日志交换机，内部对应多个队列和路由键
	LogOfRequestQueueName  = "log.queue"
	LogOfRequestRoutingKey = "log.routing"
*/
//Publish 发送消息，持久化
func (c *Client) Publish(msg []byte, exChange string, routeKey string) error {
	/*
		传参的时候传入确定类型，和一个字节流
		然后根据类型，确定要发到哪个管道
		交换机 + 路由键

	*/
	return c.channel.Publish(
		exChange,
		routeKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         msg,
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

// Close 关闭连接
func (c *Client) Close() {
	c.channel.Close()
	c.conn.Close()
}
