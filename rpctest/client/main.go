package main

import (
	pb "github.com/daniel840829/gameServer/proto" // 引入proto包

	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	Address = "35.185.75.79:8080"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	// 初始化客户端
	c := pb.NewPacketClient(conn)

	// 调用方法
	reqBody := &pb.Pos{Id: "test", Vector3: &pb.Vector3{1.2, 1.3, 1.4}, Rotation: &pb.Rotation{}}
	r, err := c.SyncPostion(context.Background(), reqBody)
	if err != nil {
		fmt.Println(err)
	}
	log.Debug(r)
}
