package rabbitmq

import (
	"baby/models"
	"github.com/goccy/go-json"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

const (
	OrderExchangeName = "order.exchange" // 订单交换机
	OrderQueueName    = "order.queue"    // 订单队列
	OrderRoutingKey   = "order.routing"  // 路由键
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

// DeclareOrderExchangeQueue 声明订单相关的交换机和队列（direct 类型，确保消息可靠投递）
func (c *Client) DeclareOrderExchangeQueue() error {
	// 声明交换机（持久化）
	if err := c.channel.ExchangeDeclare(
		OrderExchangeName,
		amqp091.ExchangeDirect, // 直接交换机，精确匹配路由键
		true,                   // 持久化
		false,                  // 不自动删除
		false,                  // 非排他
		false,                  // 非阻塞
		nil,
	); err != nil {
		return err
	}

	// 声明队列（持久化，避免消息丢失）
	_, err := c.channel.QueueDeclare(
		OrderQueueName,
		true,  // 持久化队列
		false, // 不自动删除
		false, // 非排他
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 绑定队列到交换机（指定路由键）
	return c.channel.QueueBind(
		OrderQueueName,
		OrderRoutingKey,
		OrderExchangeName,
		false,
		nil,
	)
}

// PublishOrder 发送订单消息（持久化消息）
func (c *Client) PublishOrder(msg models.LoginLogMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.channel.Publish(
		OrderExchangeName,
		OrderRoutingKey,
		false, // 非 mandatory
		false, // 非 immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp091.Persistent, // 持久化消息（避免 RabbitMQ 重启丢失）
			Timestamp:    time.Now(),
		},
	)
}

// ConsumeOrder 消费订单消息（手动确认）
func (c *Client) ConsumeOrder() (<-chan amqp091.Delivery, error) {
	return c.channel.Consume(
		OrderQueueName,
		"order_consumer", // 消费者标签
		false,            // 关闭自动确认（手动确认确保任务完成）
		false,
		false,
		false,
		nil,
	)
}

// Close 关闭连接
func (c *Client) Close() {
	c.channel.Close()
	c.conn.Close()
}
