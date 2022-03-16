package main

import (
	"TcpServer/src/global"
	"TcpServer/src/initialize"
	"TcpServer/src/rpc/userRpc"
	"TcpServer/src/service"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

func main() {
	//加载配置文件
	initialize.InitConfig()
	//初始化Gorm
	global.DBHelper = initialize.CreateGorm()
	//初始化Redis
	global.RedisHelper = initialize.CreateRedis()

	l, _ := net.Listen("tcp", ":"+global.Config.Server.Port)
	s := grpc.NewServer()
	userRpc.RegisterSearchServiceServer(s, &service.UserServe{})
	err := s.Serve(l)
	if err != nil {
		fmt.Println(err)
	}
}
