package agent

import (
	"strconv"

	"github.com/daniel840829/gameServer/agent/session"
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/user"

	//. "github.com/daniel840829/gameServer/uuid"
	//"github.com/globalsign/mgo"
	//"fmt"
	//"github.com/daniel840829/gameServer/storage"
	//p "github.com/golang/protobuf/proto"
	//"github.com/golang/protobuf/ptypes"
	//any "github.com/golang/protobuf/ptypes/any"
	//log "github.com/sirupsen/logrus"

	"time"

	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	/*
		"io"
		"reflect"
		"sync"
	*/)

/*
type ClientToAgentServer interface {
	AquireSessionKey(context.Context, *Empty) (*SessionKey, error)
	AquireOtherAgent(context.Context, *Empty) (*ServerInfo, error)
	// Login
	Login(context.Context, *LoginInput) (*UserInfo, error)
	CreateAccount(context.Context, *RegistInput) (*Error, error)
	// UserSetting
	SetAccount(context.Context, *AttrSetting) (*Success, error)
	SetCharacter(context.Context, *AttrSetting) (*Success, error)
	// room
	AquireGameServer(context.Context, *Empty) (*ServerInfo, error)
	CreateRoom(context.Context, *RoomSetting) (*Success, error)
	JoinRoom(context.Context, *ID) (*Success, error)
	RoomReady(context.Context, *Empty) (*Success, error)
	// View
	UpdateHome(*Empty, ClientToAgent_UpdateHomeServer) error
	UpdateRoomList(*Empty, ClientToAgent_UpdateRoomListServer) error
	UpdateUserList(*Empty, ClientToAgent_UpdateUserListServer) error
	// rpc UpdateRoomInfo(SessionKey) returns (stream RoomInfoView) {}
	Pipe(ClientToAgent_PipeServer) error
}*/

type ErrorPipe struct {
	toClient chan (*MessageToUser)
}

func NewAgentRpc() (agent *Agent) {
	agent = &Agent{
		ErrorPipe: &ErrorPipe{
			toClient: make(chan (*MessageToUser), 10),
		},
	}
	return
}

type Agent struct {
	Uuid       int64
	ErrorPipe  *ErrorPipe
	GameServer AgentToGameClient
}

func (a *Agent) Init(ip, agentToGamePort, clientToGamePort string) {
	if os.Getenv("DONT_USE_KUBE") == "true" {
		session.RoomManager.ConnectGameServer(ip, clientToGamePort, agentToGamePort, "0")
	} else {
		session.ClusterManager.KubeClientSet()
	}
}

func (a *Agent) AquireSessionKey(c context.Context, e *Empty) (*SessionKey, error) {
	id := session.Manager.MakeSession()
	return &SessionKey{Value: strconv.FormatInt(id, 10)}, nil
}
func (a *Agent) AquireOtherAgent(c context.Context, e *Empty) (*ServerInfo, error) {
	return nil, nil
}

func (a *Agent) GetSessionCache(c context.Context, e *Empty) (*SessionCache, error) {
	s := GetSesionFromContext(c)

	if s == nil {
		log.Warn("GetSessionCache Fail")
		return &SessionCache{}, status.Errorf(codes.NotFound, "Session Not Found!")
	}

	cache := s.GetSessionCache()
	return cache, nil
}

// Login

func (a *Agent) Login(c context.Context, in *LoginInput) (*UserInfo, error) {
	s := GetSesionFromContext(c)
	if s == nil {
		return &UserInfo{}, status.Errorf(codes.NotFound, "Session Not Found!")
	}
	s.Lock()
	user := s.State.Login(in.UserName, in.Pswd)
	if user == nil {
		s.Unlock()
		return nil, nil
	}
	s.Unlock()
	return user.UserInfo, nil
}
func (a *Agent) CreateAccount(c context.Context, in *RegistInput) (*Error, error) {
	errmsg, err := user.Manager.Regist(in)
	return errmsg, err
}

// UserSetting
func (a *Agent) SetAccount(context.Context, *AttrSetting) (*Success, error) {
	return nil, nil
}

func (a *Agent) SetCharacter(c context.Context, setting *CharacterSetting) (*Success, error) {
	s := GetSesionFromContext(c)
	if s == nil {
		return &Success{
			Ok: false,
		}, status.Errorf(codes.NotFound, "Session Not Found!")
	}
	ok := s.State.SettingCharacter(setting)
	return &Success{
		Ok: ok,
	}, nil
}

// allocate room
func (a *Agent) AquireGameServer(c context.Context, e *Empty) (*ServerInfo, error) {
	s := GetSesionFromContext(c)
	log.Debug("Aquiring game server...")
	time.Sleep(5 * time.Second)
	if s == nil {
		return &ServerInfo{}, status.Errorf(codes.NotFound, "Session Not Found!")
	}
	msg := s.GetMsgChan("ServerInfo")
	if msg != nil {
		serverInfo := <-msg.DataCh
		return serverInfo.(*ServerInfo), nil
	}
	return &ServerInfo{}, nil
}

// View
func (a *Agent) UpdateHome(*Empty, ClientToAgent_UpdateHomeServer) error {
	return nil
}

func (a *Agent) UpdateRoomList(e *Empty, stream ClientToAgent_UpdateRoomListServer) error {
	s := GetSesionFromContext(stream.Context())
	if s == nil {
		return status.Errorf(codes.NotFound, "Session Not Found!")
	}
	msgChan := s.GetMsgChan("RoomList")
	data := msgChan.DataCh
	stop := msgChan.StopSignal
	for {
		select {
		case msg := <-data:
			stream.Send(msg.(*RoomList))
		case <-stop:
			break
		}
	}
	return nil
}

func (a *Agent) UpdateUserList(*Empty, ClientToAgent_UpdateUserListServer) error {
	return nil
}

// rpc UpdateRoomInfo(SessionKey) returns (stream RoomInfoView) {}
func (a *Agent) Pipe(ClientToAgent_PipeServer) error {
	return nil
}
func (a *Agent) CreateRoom(c context.Context, roomSetting *RoomSetting) (*Success, error) {
	s := GetSesionFromContext(c)
	if s == nil {
		return &Success{}, status.Errorf(codes.NotFound, "Session Not Found!")
	}
	ok := s.State.CreateRoom(roomSetting)
	return &Success{
		Ok: ok,
	}, nil
}

func (a *Agent) JoinRoom(c context.Context, id *ID) (*Success, error) {
	s := GetSesionFromContext(c)

	if s == nil {
		return &Success{
			Ok: false,
		}, status.Errorf(codes.NotFound, "Session Not Found!")
	}
	ok := s.State.EnterRoom(id.Value)
	return &Success{
		Ok: ok,
	}, nil
}

func (a *Agent) UpdateRoomContent(e *Empty, stream ClientToAgent_UpdateRoomContentServer) error {
	s := GetSesionFromContext(stream.Context())
	if s == nil {
		return status.Errorf(codes.NotFound, "Session Not Found!")
	}
	msgChan := s.GetMsgChan("RoomContent")
	data := msgChan.DataCh
	stop := msgChan.StopSignal
	for {
		select {
		case msg := <-data:
			stream.Send(msg.(*RoomContent))
		case <-stop:
			break
		}
	}
	return nil
}

func (a *Agent) RoomReady(c context.Context, e *Empty) (*Success, error) {
	s := GetSesionFromContext(c)
	if s == nil {
		return &Success{
			Ok: false,
		}, status.Errorf(codes.NotFound, "Session Not Found!")
	}

	if s.IsReady {
		if s.State.CancelReady() {
			return &Success{
				Ok: s.IsReady,
			}, nil
		}
	} else {
		if s.State.ReadyRoom() {
			return &Success{
				Ok: s.IsReady,
			}, nil
		}
	}
	return &Success{}, status.Errorf(codes.Internal, "Somethig Wrong")
}

func GetSesionFromContext(c context.Context) *session.Session {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil
	}
	s := session.Manager.GetSession(md)
	return s
}
