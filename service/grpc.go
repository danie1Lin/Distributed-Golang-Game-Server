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
	"sync"
	"time"
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
		SendFuncToClient:   make(map[int64](chan *CallFuncInfo)),
		RecvFuncFromClient: make(chan *CallFuncInfo, 100),
		PosToClient:        make(map[int64](chan *Position)),
		InputFromClient:    make(chan *Input, 100),
		ErrToClient:        make(map[int64](chan *Error)),
		ErrFromClient:      make(chan *Error, 100),
	}
}

type Rpc struct {
	//connection id map
	sync.RWMutex
	IsDisconnect       sync.Map
	SendFuncToClient   map[int64](chan *CallFuncInfo)
	RecvFuncFromClient chan *CallFuncInfo
	PosToClient        map[int64](chan *Position)
	InputFromClient    chan *Input
	ErrToClient        map[int64](chan *Error)
	ErrFromClient      chan *Error
}

func (rpc *Rpc) ErrorPipLine(stream Rpc_ErrorPipLineServer) error {
	err, _ := stream.Recv()
	id := err.FromId
	send := rpc.ErrToClient[id]
	recv := rpc.ErrFromClient
	c, _ := rpc.IsDisconnect.Load(id)
	closeRpc := c.(chan struct{})
	fmt.Println("Client[", id, "] start postion sync")
	go func() {
		for {
			select {
			case <-closeRpc:
			default:
				in, err := stream.Recv()
				if err == io.EOF {
					rpc.Disconnect(id)
					break
				}
				if err != nil {
					fmt.Println("Recv Error: ", err)
					rpc.Disconnect(id)
					break
				}
				recv <- in
			}
		}
		close(recv)
	}()
	for {
		select {
		case msg := <-send:
			err := stream.Send(msg)
			if err != nil {
				fmt.Println("Send erro: r", err)
				rpc.Disconnect(id)
				break
			}
		case <-closeRpc:
			break
		}
	}
	return nil
}

func (rpc *Rpc) SyncPos(stream Rpc_SyncPosServer) error {
	input, _ := stream.Recv()
	id := input.UserId
	send := rpc.PosToClient[id]
	recv := rpc.InputFromClient
	c, _ := rpc.IsDisconnect.Load(id)
	closeRpc := c.(chan struct{})
	fmt.Println("Client[", id, "] start postion sync")
	go func() {
		for {
			select {
			case <-closeRpc:
			default:
				in, err := stream.Recv()
				if err == io.EOF {
					rpc.Disconnect(id)
					break
				}
				if err != nil {
					fmt.Println("Recv Error: ", err)
					rpc.Disconnect(id)
					break
				}
				recv <- in
				t := time.Now().Sub(time.Unix(in.TimeStamp * 1000000)).Second()
				log.Debug("[grpc]{SyncPos}Recv:", in, " [Delay]", t, "(s)")
			}
		}
		close(recv)
	}()
	for {
		select {
		case data := <-send:
			data.TimeStamp = time.Now().UnixNano()
			err := stream.Send(data)
			log.Debug("[grpc]{SyncPos}Send", data)
			if err != nil {
				fmt.Println("Send erro: r", err)
				rpc.Disconnect(id)
				break
			}
		case <-closeRpc:
			break
		}
	}
	return nil
}

func (rpc *Rpc) CallMethod(stream Rpc_CallMethodServer) error {
	callFuncInfo, _ := stream.Recv()
	id := callFuncInfo.FromId
	send := rpc.SendFuncToClient[id]
	recv := rpc.RecvFuncFromClient
	c, _ := rpc.IsDisconnect.Load(id)
	closeRpc := c.(chan struct{})
	fmt.Println("Client[", id, "] start CallMethod")
	//recv := rpc.RecvFuncFromClient[id]
	go func() {
		for {
			select {

			case <-closeRpc:
				break
			default:
				in, err := stream.Recv()
				log.Debug("[RECIEVE]", in)
				if err == io.EOF {
					log.Info("[RPC Disconnect] ID:", id)
					rpc.Disconnect(id)
					break
				}
				if err != nil {
					fmt.Println("Recv Error: ", err)
					rpc.Disconnect(id)
					log.Info("[RPC Disconnect] ID:", id, "[Error]", err)
					break
				}
				recv <- in
			}
		}
		close(recv)
	}()
	for {
		select {
		case out := <-send:
			log.Debug("[SEND]", out)
			err := stream.Send(out)
			if err != nil {
				fmt.Println("Send error", err)
				rpc.Disconnect(id)
				log.Info("[RPC Disconnect] ID:", id, "[Error]", err)
				break
			}
		case <-closeRpc:
			break
		}
	}
	return nil
}

func (rpc *Rpc) Disconnect(userId int64) bool {
	if closeCh, ok := rpc.IsDisconnect.Load(userId); ok {
		rpc.Lock()
		user.Manager.Logout(userId)
		close(closeCh.(chan struct{}))
		delete(rpc.SendFuncToClient, userId)
		delete(rpc.ErrToClient, userId)
		delete(rpc.PosToClient, userId)
		rpc.IsDisconnect.Delete(userId)
		rpc.Unlock()
		return true
	} else {
		return false
	}
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
	rpc.IsDisconnect.Store(userInfo.Uuid, make(chan struct{}))
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
