package main

import (
	"net"

	"github.com/daniel840829/gameServer/msg" // 引入编译生成的包

	"google.golang.org/grpc"
	//"google.golang.org/grpc/grpclog"
	"fmt"

	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
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

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册HelloService
	msg.RegisterRpcServer(s, &msg.Rpc{UserNameMapUserInfo: make(map[string]*msg.UserInfo), Send: make(map[string](chan []byte))})

	fmt.Println("Listen on " + Address)

	s.Serve(listen)
}
