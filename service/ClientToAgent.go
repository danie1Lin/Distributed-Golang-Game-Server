package service

import (
	. "github.com/daniel840829/gameServer2/msg"
	"github.com/daniel840829/gameServer2/session"
	//"github.com/daniel840829/gameServer2/user"
	"strconv"
	//. "github.com/daniel840829/gameServer/uuid"
	//"github.com/globalsign/mgo"
	//"fmt"
	//"github.com/daniel840829/gameServer/storage"
	//p "github.com/golang/protobuf/proto"
	//"github.com/golang/protobuf/ptypes"
	//any "github.com/golang/protobuf/ptypes/any"
	//log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	/*
		"io"
		"reflect"
		"sync"
		"time"
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
	// allocate room
	AquireGameServer(context.Context, *Empty) (*ServerInfo, error)
	// View
	UpdateHome(*Empty, ClientToAgent_UpdateHomeServer) error
	UpdateRoomList(*Empty, ClientToAgent_UpdateRoomListServer) error
}*/

func NewAgentRpc() (agent *Agent) {
	agent = &Agent{}
	return
}

type Agent struct {
	Uuid int64
}

func (a *Agent) Init() {

}

func (a *Agent) AquireSessionKey(c context.Context, e *Empty) (*SessionKey, error) {
	id := session.Manager.MakeSession()
	return &SessionKey{Value: strconv.FormatInt(id, 10)}, nil
}
func (a *Agent) AquireOtherAgent(c context.Context, e *Empty) (*ServerInfo, error) {
	return nil, nil
}

// Login
func (a *Agent) Login(c context.Context, in *LoginInput) (*UserInfo, error) {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, nil
	}
	s := session.Manager.GetSession(md)
	if s == nil {
		return nil, nil
	}

	s.Lock()
	user := s.State.Login(in.UserName, in.Pswd)
	s.Unlock()
	return user.UserInfo, nil
}
func (a *Agent) CreateAccount(context.Context, *RegistInput) (*Error, error) {
	return nil, nil
}

// UserSetting
func (a *Agent) SetAccount(context.Context, *AttrSetting) (*Success, error) {
	return nil, nil
}
func (a *Agent) SetCharacter(context.Context, *AttrSetting) (*Success, error) {
	return nil, nil
}

// allocate room
func (a *Agent) AquireGameServer(context.Context, *Empty) (*ServerInfo, error) {
	return nil, nil
}

// View
func (a *Agent) UpdateHome(*Empty, ClientToAgent_UpdateHomeServer) error {
	return nil
}
func (a *Agent) UpdateRoomList(*Empty, ClientToAgent_UpdateRoomListServer) error {
	return nil
}
