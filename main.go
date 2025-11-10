package main

import (
	"baby/models"
	"baby/rabbitmq"
	"baby/routers"
	"fmt"
	"net/http"
	"time"
)

var cli *rabbitmq.Client

func init() {

	fmt.Println("rabbitmq init")

	err := models.Setup()
	if err != nil {
		panic(err)
	}

	cli, err = rabbitmq.NewClient("localhost:5672", "guest", "guest")
	if err != nil {
		panic(err)
	}
	err = cli.DeclareOrderExchangeQueue()
	if err != nil {
		panic(err)
	}

	fmt.Println("rabbitmq init success")
}

func main() {

	msg := models.LoginLogMessage{
		OrderID:     999,
		UserID:      998,
		CommodityID: 3432,
		Quantity:    9999,
		Price:       9.9,
	}

	errr := cli.PublishOrder(msg)
	if errr != nil {
		panic(errr)
	}

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
