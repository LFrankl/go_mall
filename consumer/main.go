package main

import (
	"baby/rabbitmq"
	"fmt"
)

/*

写一个处理所有业务
复杂度上来了再解耦

先全局初始化一个Client，各个模块都能使用

消费者进程也要写
DeclareOrderExchangeQueue()
......等

*/

var cli *rabbitmq.Client

func init() {
	fmt.Println("connecting to rabbitmq")
	var err error
	cli, err = rabbitmq.NewClient("localhost:5672", "guest", "guest")
	if err != nil {
		panic(err)
	}
	fmt.Println("connected ")
	return

}

func main() {

	err := cli.DeclareExchangeQueue(rabbitmq.LogExchangeName,
		rabbitmq.LogOfRequestQueueName,
		rabbitmq.LogOfRequestRoutingKey,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := cli.Consume(rabbitmq.LogOfRequestQueueName, "consumer_1")
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		// msg 是 amqp091.Delivery 类型
		fmt.Printf("收到消息：路由键=%s，内容=%s\n", msg.RoutingKey, string(msg.Body))
		err := msg.Ack(false)
		if err != nil {
			panic(err)
		}
	}

}
