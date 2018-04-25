package main

import (
	"net"

	pb "github.com/daniel840829/gameServer/proto" // 引入编译生成的包

	"golang.org/x/net/context"
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
	Address = "127.0.0.1:50052"
)

type rpcService struct{}

var rpc *rpcService = &rpcService{}

func (r *rpcService) SyncPostion(ctx context.Context, in *pb.Pos) (*pb.PosReply, error) {
	log.Debug(in)
	return new(pb.PosReply), nil
}

func (r *rpcService) CallServer(ctx context.Context, in *pb.Callin) (*pb.Reply, error) {
	log.Debug(in)
	return new(pb.Reply), nil
}
func (r *rpcService) CallClient(in *pb.ClientStart, stream pb.Packet_CallClientServer) error {
	log.Debug(in)
	return nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println("failed to listen: %v", err)
	}

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册HelloService
	pb.RegisterPacketServer(s, &rpcService{})

	fmt.Println("Listen on " + Address)

	s.Serve(listen)
}
