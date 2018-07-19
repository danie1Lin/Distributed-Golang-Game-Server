package main

import (
	"net"

	//"github.com/daniel840829/gameServer2/entity"
	"github.com/daniel840829/gameServer2/msg"
	"github.com/daniel840829/gameServer2/service"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/grpclog"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	serverType string = *flag.String("type", "agent", "choose server Type")
	configFile string = *flag.String("config", "", "config file's path")
	port       string = *flag.String("port", "8080", "port")
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	flag.Parse()

}

func main() {
	if serverType == "agent" {
		listen, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("failed to listen: %v", err)
		}
		agentRpc := service.NewAgentRpc()
		s := grpc.NewServer()
		msg.RegisterClientToAgentServer(s, agentRpc)

		fmt.Println("Listen on " + port)

		s.Serve(listen)
	} else if serverType == "game" {
		listen, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Println("failed to listen: %v", err)
		}
		agentRpc := service.NewAgentRpc()
		s := grpc.NewServer()
		msg.RegisterClientToAgentServer(s, agentRpc)

		fmt.Println("Listen on " + port)

		s.Serve(listen)
	}
	/*
		//初始化gameManager
		rpc := service.NewRpc()
		gm := &entity.GameManager{}
		gm.Init(rpc)
		gm.RegistRoom("room", &entity.Room{})
		gm.RegistEnitity("Player", &entity.Player{})
		gm.RegistEnitity("Shell", &entity.Shell{})
		gm.RegistEnitity("Enemy", &entity.Enemy{})
		go gm.Run()
		// 注册HelloService
	*/
}
