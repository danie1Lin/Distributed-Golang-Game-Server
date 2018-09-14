package service

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/session"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"strconv"
)

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
	s := GetSesionFromContext(c)
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
func (a *Agent) UpdateUserList(*Empty, ClientToAgent_UpdateUserListServer) error {
	return nil
}

// rpc UpdateRoomInfo(SessionKey) returns (stream RoomInfoView) {}
func (a *Agent) Pipe(ClientToAgent_PipeServer) error {
	return nil
}
func (a *Agent) CreateRoom(context.Context, *RoomSetting) (*Success, error) {
	return nil, nil
}
func (a *Agent) JoinRoom(context.Context, *ID) (*Success, error) {
	return nil, nil
}
func (a *Agent) RoomReady(context.Context, *Empty) (*Success, error) {
	return nil, nil
}

func GetSesionFromContext(c context.Context) *Session {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil
	}
	s := session.Manager.GetSession(md)

	return s
}
