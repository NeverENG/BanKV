package main

import (
	"fmt"

	"github.com/NeverENG/BanKV/internal/network/znet"
	"github.com/NeverENG/BanKV/internal/service"
)

func main() {
	// 初始化 FSM
	fsm := service.NewFSM()

	// 启动 FSM
	go fsm.Run()

	// 初始化 HA
	ha := service.NewHA(fsm)

	// 初始化网络服务
	server := znet.NewServer()

	// 创建路由
	router := service.NewRouter(fsm)

	// 注册路由
	server.AddRouter(1, router) // PUT 操作
	server.AddRouter(2, router) // GET 操作
	server.AddRouter(3, router) // DELETE 操作

	// 启动服务
	fmt.Println("Starting server...")
	fmt.Printf("HA initialized, initial health status: %v\n", ha.IsHealthy())
	server.Serve()
}
