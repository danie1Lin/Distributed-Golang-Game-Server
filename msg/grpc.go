package msg

import (
	p "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"math/rand"
	"reflect"
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

type Rpc struct {
	//connection id map
	UserOnLine                   map[string]*UserInfo
	SendFuncToClientFuncToClient map[string](chan CallFuncInfo)
	RecvFuncFromClient           map[string](chan CallFuncInfo)
}

func (rpc *Rpc) SyncPos(pos Rpc_SyncPosServer) error {

	return nil
}

func (rpc *Rpc) CallFuncInfo(stream Rpc_CallMethodServer) error {
	callFuncInfo := stream.Recv()
	send := rpc.SendFuncToClient[callFuncInfo.Descript]
	for {
		select {
		case msg := <-send:
			log.Debug(msg)
		case <-time.After(2 * time.Second):
			param := make([]*any.Any, 1)
			a := AnyEecode(&TransForm{X: 3.233})
			param[0] = a
			callClient.Send(&Reply{Error: "", Param: param})
		}
	}
	return nil
}

func (rpc *Rpc) Login(ctx context.Context, in *User) (*UserInfo, error) {
	log.Info("Login:", in)
	//check pswd and get UUID
	userInfo, ok := rpc.UserNameMapUserInfo[in.UserName]
	if !ok {
		GetDbData(in)
		return &UserInfo{}, nil
	}
	userInfo.Error = ""
	//chan <- call createEntity method
	//let client start callSignal
	return userInfo, nil
}

func (rpc *Rpc) CreateAccount(ctx context.Context, userRegister *UserRegister) (*UserRegisterInfo, error) {
	log.Info("Create Account:", userRegister)
	//check userName isn't repeated
	//create a default UserInfo
	//save UserInfo
	//return
	userInfo := NewUserInfo(userRegister.UserName)
	rpc.UserNameMapUserInfo[userInfo.UserName] = userInfo
	log.Debug(userInfo)
	SaveDbData(userInfo)
	return &UserRegisterInfo{}, nil
}

func CallAllClient(entityTypeName string, id uuid.UUID, f string, args ...p.Message) {
	for _, a := range args {
		v := reflect.ValueOf(a)
		t := v.Type()
		anything, _ := ptypes.MarshalAny(a)
		x := &TransForm{}
		AnyDecode(anything, x)
		log.Debug(a, v, t, anything, x)
	}
}

func AnyEecode(in p.Message) *any.Any {
	anything, _ := ptypes.MarshalAny(in)
	return anything
}

func AnyDecode(anything *any.Any, out p.Message) {
	err := ptypes.UnmarshalAny(anything, out)
	if err != nil {
		log.Info(err)
	}
}

func NewCharacter() (c *Character) {
	c = &Character{}
	uuid, _ := uuid.NewV4()
	c.Id = uuid.String()
	c.Color = &Color{int32(rand.Intn(256)), int32(rand.Intn(256)), int32(rand.Intn(256))}
	return
}

func NewUserInfo(userName string) (u *UserInfo) {
	u = &UserInfo{IdMapCharacter: make(map[string]*Character)}
	c := NewCharacter()
	u.UserName = userName
	u.IdMapCharacter[c.Id] = c
	return
}

func FoundUser(UserName string) *UserInfo {
	return nil
}

func GetDbData(data interface{}) interface{} {
	return nil
}

func SaveDbData(data interface{}) interface{} {
	return nil
}
