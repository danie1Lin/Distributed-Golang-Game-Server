package main

import (
	"net"

	//"github.com/daniel840829/gameServer/entity"
	"runtime/pprof"

	"github.com/daniel840829/gameServer/agent"
	"github.com/daniel840829/gameServer/game"
	"github.com/daniel840829/gameServer/msg"
	"google.golang.org/grpc"

	//"google.golang.org/grpc/grpclog"
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

var (
	serverType       *string
	configFile       *string
	AgentPort        *string
	AgentToGamePort  *string
	ClientToGamePort *string
	cpuprofile       *string
)

func main() {

	serverType = flag.String("type", "game", "choose server Type")
	configFile = flag.String("config", "", "config file's path")
	AgentPort = flag.String("agentPort", "50051", "ClientToAgent Port")
	AgentToGamePort = flag.String("agentToGamePort", "3000", "AgentToGame Port")
	ClientToGamePort = flag.String("clientToGamePort", "8080", "ClientToGame Port")
	cpuprofile = flag.String("cpuprofile", "./cpu.prof", "write cpu profile to file,set blank to close profile function")
	log.Debug("config :", "type :", serverType)
	flag.Parse()
	//SERVER_TYPE ID CLIENT_TO_AGENT_PORT CLIENT_TO_GAME_PORT AGENT_TO_GAME_PORT
	ReadEnv(serverType, "SERVER_TYPE")
	ReadEnv(AgentToGamePort, "AGENT_TO_GAME_PORT")
	ReadEnv(ClientToGamePort, "CLIENT_TO_GAME_PORT")
	ReadEnv(AgentPort, "CLIENT_TO_AGENT_PORT")

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *serverType == "agent" {
		RunAgent()
	} else if *serverType == "game" {
		go RunGame()
		RunAgentToGame()
	} else {
		go RunGame()
		go RunAgentToGame()
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

//ReadEnv if use kubernete Read para from env
func ReadEnv(para *string, envName string) {
	v := os.Getenv(envName)
	if v != "" {
		log.Info(envName, " change from ", *para, " to ", v)
		*para = v
	}
}

func RunAgent() {
	listen, err := net.Listen("tcp", ":"+*AgentPort)
	if err != nil {
		fmt.Println("AgentServer failed to listen: %v", err)
	}
	agentRpc := agent.NewAgentRpc()
	s := grpc.NewServer()
	msg.RegisterClientToAgentServer(s, agentRpc)
	fmt.Println("AgentServer Listen on " + *AgentPort)
	agentRpc.Init("127.0.0.1", *AgentToGamePort, *ClientToGamePort)
	s.Serve(listen)
}

func RunGame() {
	listen, err := net.Listen("tcp", ":"+*ClientToGamePort)
	if err != nil {
		fmt.Println("GameServer failed to listen: %v", err)
	}
	s := grpc.NewServer()
	msg.RegisterClientToGameServer(s, &game.CTGServer{})
	fmt.Println("GameServer Listen on " + *ClientToGamePort)
	s.Serve(listen)
}

func RunAgentToGame() {
	listen, err := net.Listen("tcp", ":"+*AgentToGamePort)
	if err != nil {
		fmt.Println("AgentToGameServer failed to listen: %v", err)
	}
	s := grpc.NewServer()
	msg.RegisterAgentToGameServer(s, &game.ATGServer{})
	fmt.Println("AgentToGameServer Listen on " + *AgentToGamePort)
	s.Serve(listen)
}
