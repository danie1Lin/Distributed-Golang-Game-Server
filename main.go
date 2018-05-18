package main

import (
	"net"

	"github.com/daniel840829/gameServer/entity"
	"github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/service"
	"github.com/daniel840829/gameServer/storage"
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

var MgoDb *storage.MongoDb = &storage.MongoDb{}

func init() {
	MgoDb.Init(storage.MGO_DB_NAME, storage.UserInfo_COLLECTION, storage.RegistInput_COLLECTION)
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}
	//初始化gameManager
	gm := &entity.GameManager{}
	//Regist Room
	//Regist entity
	//run gamemanager
	//Rpc handler
	rpc := service.NewRpc()
	gm.Init(rpc)
	go gm.Run()
	// 注册HelloService
	s := grpc.NewServer()
	msg.RegisterRpcServer(s, rpc)

	fmt.Println("Listen on " + Address)

	s.Serve(listen)
}
