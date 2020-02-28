package main

import (
	user "elementser/src/auth"
	middleware "elementser/src/middleware"
	transaction "elementser/src/transaction"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/web"
)

func main() {
	serviceName := "dev_test"
	serviceVersion := "latest"

	fmt.Printf("start service:%s version:%s", serviceName, serviceVersion)

	service := web.NewService(
		web.Name(serviceName),
		web.Address(":8012"),
		web.RegisterTTL(time.Second*10),
		web.RegisterInterval(time.Second*5))
	_ = service.Init()

	engine := gin.Default()
	route := engine.Group("/dev-api")
	route.Use(middleware.CheckToken())

	user := new(user.User)
	transaction := new(transaction.TransactionList)
	route.POST("/vue-element-admin/user/login", user.Login)
	route.GET("/vue-element-admin/user/info", user.Info)
	route.POST("/vue-element-admin/user/logout", user.Logout)
	route.GET("/vue-element-admin/transaction/list", transaction.List)

	service.Handle("/", engine)

	if err := service.Run(); err != nil {
		println(err)
	}
}
