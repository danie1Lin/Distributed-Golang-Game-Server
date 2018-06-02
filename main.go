package main

import (
	"net"

	"github.com/daniel840829/gameServer/entity"
	"github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/service"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/grpclog"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

const (
	// Address gRPC服务地址
	Address = ":8080"
)

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}
	rpc := service.NewRpc()
	//初始化gameManager
	gm := &entity.GameManager{}
	gm.Init(rpc)
	gm.RegistRoom("room", &entity.Room{})
	gm.RegistEnitity("Player", &entity.Player{})
	gm.RegistEnitity("Shell", &entity.Shell{})
	gm.RegistEnitity("Enemy", &entity.Enemy{})
	go gm.Run()
	// 注册HelloService
	s := grpc.NewServer()
	msg.RegisterRpcServer(s, rpc)

	fmt.Println("Listen on " + Address)

	s.Serve(listen)
}
