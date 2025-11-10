package main

import (
	"baby/middleware/log"
	"baby/models"
	"baby/rabbitmq"
	"baby/routers"
	"fmt"
	"net/http"
	"time"
)

var cli *rabbitmq.Client

func initRabbitMQ() {
	fmt.Println("rabbitmq init")

	err := models.Setup()
	if err != nil {
		panic(err)
	}

	cli, err = rabbitmq.NewClient("localhost:5672", "guest", "guest")
	if err != nil {
		panic(err)
	}
	err = cli.DeclareExchangeQueue(rabbitmq.LogExchangeName,
		rabbitmq.LogOfRequestQueueName,
		rabbitmq.LogOfRequestRoutingKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("rabbitmq init success")
}

func init() {
	initRabbitMQ()
	log.RegisterCli(cli)
}

func main() {

	server := &http.Server{
		Addr:         ":8080",
		Handler:      routers.InitRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//可以使用fvbock/endless替换http的ListenAndServe实现平滑重启

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
