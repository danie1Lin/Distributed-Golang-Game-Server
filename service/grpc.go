package service

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/user"
	//. "github.com/daniel840829/gameServer/uuid"
	//"github.com/globalsign/mgo"
	"fmt"
	//"github.com/daniel840829/gameServer/storage"
	p "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"io"
	"reflect"
	//"time"
)

/*
type RpcServer interface {
	SyncPos(Rpc_SyncPosServer) error
	CallMethod(Rpc_CallMethodServer) error
	ErrorPipLine(Rpc_ErrorPipLineServer) error
	Login(context.Context, *LoginInput) (*UserInfo, error)
	CreateAccount(context.Context, *RegistInput) (*Error, error)
}
*/
func NewRpc() *Rpc {
	return &Rpc{
		UserOnLine:         make(map[int64]*UserInfo),
		SendFuncToClient:   make(map[int64](chan *CallFuncInfo)),
		RecvFuncFromClient: make(chan *CallFuncInfo, 100),
		PosToClient:        make(map[int64](chan *Position)),
		PosFromClient:      make(chan *Position, 100),
		ErrToClient:        make(map[int64](chan *Error)),
		ErrFromClient:      make(chan *Error, 100),
	}
}

type Rpc struct {
	//connection id map
	UserOnLine         map[int64]*UserInfo
	SendFuncToClient   map[int64](chan *CallFuncInfo)
	RecvFuncFromClient chan *CallFuncInfo
	PosToClient        map[int64](chan *Position)
	PosFromClient      chan *Position
	ErrToClient        map[int64](chan *Error)
	ErrFromClient      chan *Error
}

func (rpc *Rpc) ErrorPipLine(stream Rpc_ErrorPipLineServer) error {
	err, _ := stream.Recv()
	id := err.FromId
	send := rpc.ErrToClient[id]
	recv := rpc.ErrFromClient
	SendEnd := make(chan struct{})
	fmt.Println("Client[", id, "] start postion sync")
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(SendEnd)
				break
			}
			if err != nil {
				close(SendEnd)
				fmt.Println("Recv Error: ", err)
				break
			}
			recv <- in
		}
	}()
	for {
		select {
		case msg := <-send:
			err := stream.Send(msg)
			if err != nil {
				fmt.Println("Send erro: r", err)
				break
			}
		case <-SendEnd:
			break
		}
	}
	delete(rpc.ErrToClient, id)
	return nil
}

func (rpc *Rpc) SyncPos(stream Rpc_SyncPosServer) error {
	position, _ := stream.Recv()
	id := position.FromId
	send := rpc.PosToClient[id]
	recv := rpc.PosFromClient
	SendEnd := make(chan struct{})
	fmt.Println("Client[", id, "] start postion sync")
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(SendEnd)
				break
			}
			if err != nil {
				close(SendEnd)
				fmt.Println("Recv Error: ", err)
				break
			}
			recv <- in
		}
	}()
	for {
		select {
		case msg := <-send:
			err := stream.Send(msg)
			if err != nil {
				fmt.Println("Send erro: r", err)
				break
			}
		case <-SendEnd:
			break
		}
	}
	delete(rpc.PosToClient, id)
	return nil
}

func (rpc *Rpc) CallMethod(stream Rpc_CallMethodServer) error {
	callFuncInfo, _ := stream.Recv()
	id := callFuncInfo.FromId
	send := rpc.SendFuncToClient[id]
	recv := rpc.RecvFuncFromClient
	SendEnd := make(chan struct{})
	fmt.Println("Client[", id, "] start CallMethod")
	//recv := rpc.RecvFuncFromClient[id]
	go func() {
		for {
			in, err := stream.Recv()
			log.Debug("[RECIEVE]", in)
			if err == io.EOF {
				close(SendEnd)
				log.Info("[RPC Disconnect] ID:", id)
				break
			}
			if err != nil {
				close(SendEnd)
				fmt.Println("Recv Error: ", err)
				log.Info("[RPC Disconnect] ID:", id, "[Error]", err)
				break
			}
			recv <- in
		}
	}()
	for {
		select {
		case out := <-send:
			log.Debug("[SEND]", out)
			err := stream.Send(out)
			if err != nil {
				fmt.Println("Send error", err)
				log.Info("[RPC Disconnect] ID:", id, "[Error]", err)
				break
			}
		case <-SendEnd:
			break
		}
	}
	delete(rpc.SendFuncToClient, id)
	return nil
}

func (rpc *Rpc) Login(ctx context.Context, in *LoginInput) (*UserInfo, error) {
	log.Info("Login:", in)
	userInfo, err := user.Manager.Login(in)
	if err != nil {
		log.Warn("login", err)
		return &UserInfo{}, nil
	}
	rpc.SendFuncToClient[userInfo.Uuid] = make(chan *CallFuncInfo, 100)
	rpc.PosToClient[userInfo.Uuid] = make(chan *Position, 100)
	rpc.ErrToClient[userInfo.Uuid] = make(chan *Error, 100)
	//TEST send a default Room
	/*
		send := rpc.SendFuncToClient[userInfo.Uuid]
		roomId, _ := Uid.NewId(ROOM_ID)
		param := make([]*any.Any, 0)
		room := &RoomInfo{
			Name: "DefaultRoom",
			Uuid: roomId,
		}
		param = append(param, AnyEncode(room))
		send <- &CallFuncInfo{
			Func:     "AddRoomInfo",
			TargetId: roomId,
			Param:    param,
		}
	*/
	//Test Done
	//TODO: create Real room
	rpc.RecvFuncFromClient <- &CallFuncInfo{
		Func:     "GetAllRoomInfo",
		TargetId: userInfo.Uuid,
	}
	return userInfo, nil

}

func (rpc *Rpc) CreateAccount(ctx context.Context, userRegister *RegistInput) (*Error, error) {
	log.Info("Create Account:", userRegister)
	user.Manager.Regist(userRegister)
	return &Error{}, nil
}

func CallAllClient(entityTypeName string, id int64, f string, args ...p.Message) {
	for _, a := range args {
		v := reflect.ValueOf(a)
		t := v.Type()
		anything, _ := ptypes.MarshalAny(a)
		x := &TransForm{}
		AnyDecode(anything, x)
		log.Debug(a, v, t, anything, x)
	}
}

func AnyEncode(in p.Message) *any.Any {
	anything, _ := ptypes.MarshalAny(in)
	return anything
}

func AnyDecode(anything *any.Any, out p.Message) {
	err := ptypes.UnmarshalAny(anything, out)
	if err != nil {
		log.Info(err)
	}
}
