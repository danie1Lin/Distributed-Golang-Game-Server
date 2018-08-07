package main

import (
	"net"

	//"github.com/daniel840829/gameServer2/entity"
	"github.com/daniel840829/gameServer2/agent"
	"github.com/daniel840829/gameServer2/msg"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/grpclog"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	serverType       string = *flag.String("type", "agent", "choose server Type")
	configFile       string = *flag.String("config", "", "config file's path")
	AgentPort        string = *flag.String("agentPort", "8080", "ClientToAgent Port")
	AgentToGamePort  string = *flag.String("agentToGamePort", "3000", "AgentToGame Port")
	ClientToGamePort string = *flag.String("clientToGamePort", "8000", "ClientToGame Port")
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	flag.Parse()

}

func main() {
	if serverType == "agent" {
		RunAgent()
	} else if serverType == "game" {
		RunGame()
	} else {
		go RunGame()
		RunAgent()
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

func RunAgent() {

	listen, err := net.Listen("tcp", ":"+AgentPort)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}
	agentRpc := agent.NewAgentRpc()
	s := grpc.NewServer()
	msg.RegisterClientToAgentServer(s, agentRpc)

	fmt.Println("Listen on " + port)

	s.Serve(listen)
}

func RunGame() {
	listen, err := net.Listen("tcp", ":"+AgentPort)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}
	agentRpc := agent.NewAgentRpc()
	agentRpc.Init(AgentToGamePort)
	s := grpc.NewServer()
	msg.RegisterClientToAgentServer(s, agentRpc)

	fmt.Println("Listen on " + port)

	s.Serve(listen)
}
