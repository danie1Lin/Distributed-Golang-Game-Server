package entity

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/service"
	. "github.com/daniel840829/gameServer/uuid"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type GameManager struct {
	////basic information
	Uuid int64
	////Reflection
	thisType reflect.Type
	////child component
	//rpcCallableObject map[int64]reflect.Value
	//room
	TypeMapRoom map[string]reflect.Type
	IdMapRoom   map[int64]IRoom
	//entity
	TypeMapEntity map[string]reflect.Type
	IdMapEntity   map[int64]IEntity

	////rpc channel
	SendFuncToClient   map[int64](chan *CallFuncInfo)
	RecvFuncFromClient chan *CallFuncInfo
	PosToClient        map[int64](chan *Position)
	PosFromClient      chan *Position
	ErrToClient        map[int64](chan *Error)
	ErrFromClient      chan *Error
}

func (gm *GameManager) Init(rpc *service.Rpc) {
	gm.Uuid, _ = Uid.NewId(GM_ID)
	gm.IdMapRoom = make(map[int64]IRoom)
	gm.TypeMapRoom = make(map[string]reflect.Type)
	//gm.rpcCallableObject = make(map[int64]reflect.Value)
	gm.SendFuncToClient = rpc.SendFuncToClient
	gm.RecvFuncFromClient = rpc.RecvFuncFromClient
	gm.PosToClient = rpc.PosToClient
	gm.PosFromClient = rpc.PosFromClient
	gm.ErrFromClient = rpc.ErrFromClient
	gm.ErrToClient = rpc.ErrToClient
	gm.thisType = reflect.TypeOf(gm)
}

func (gm *GameManager) Run() {
	for {
		select {
		case f := <-gm.RecvFuncFromClient:
			log.Debug("[GM Call]", f)
			go gm.Call(f)
		case err := <-gm.ErrFromClient:
			gm.ErrorHandle(err)
		case pos := <-gm.PosFromClient:
			gm.SyncPos(pos)
		}
	}
}

func (gm *GameManager) RegisterRoom(roomTypeName string, iRoom IRoom) {
	if _, ok := gm.TypeMapRoom[roomTypeName]; ok {
		log.Fatal(roomTypeName, "is already registed.")
	}
	vRoom := reflect.ValueOf(iRoom)
	gm.TypeMapRoom[roomTypeName] = vRoom.Type().Elem()
	//TODO : Record method info to speed up reflection invoke.
}

func (gm *GameManager) RegisterEnitity(iEntity IEntity) {

}

func (gm *GameManager) CreateEnitity(entityType string, inRoomId string) {

}

func (gm *GameManager) Call(f *CallFuncInfo) {
	log.Debug("Function INFO :", f)

	method, ok := gm.thisType.MethodByName(f.Func)
	log.Debug("reflect.Method:", method)
	log.Debug("Is OK? :", ok)
	param := make([]reflect.Value, 0)
	param = append(param, reflect.ValueOf(gm))
	param = append(param, reflect.ValueOf(f))
	method.Func.Call(param)
	/*
		targetType, _ := Uid.ParseId(f.TargetId)
		switch targetType {
		case ENTITY_ID:
		case EQUIP_ID:
		case GM_ID:
		case ROOM_ID:
		default:
			if targetType == "" {
				log.Warn(f.TargetId, " is not existed !")
				return
			}
			log.Warn(targetType, " is not callable!")
		}
	*/
}

func (gm *GameManager) ErrorHandle(err *Error) {
	log.Warn("Something Wrong", err)
}

func (gm *GameManager) SyncPos(pos *Position) {
	//pos.RoomId
}

/*

Reflection callable method

*/
func (gm *GameManager) RoomReady(f *CallFuncInfo) {
	log.Info("User [", f.FromId, "] is ready in [", f.TargetId, "] Room")
}

func (gm *GameManager) CreateNewRoom(f *CallFuncInfo) {
	roomInfo := &RoomInfo{}
	err := ptypes.UnmarshalAny(f.Param[0], roomInfo)
	if err != nil {
		log.Warn("[*any Unmarshal Error]", f.Param[0])
		return
	}
	tRoom, ok := gm.TypeMapRoom[roomInfo.GameType]
	if !ok {
		log.Warn(roomInfo.GameType, " is not registed yet. ")
		return
	}
	room, ok := reflect.New(tRoom).Interface().(IRoom)
	if !ok {
		log.Warn("Something Wrong with RegisterRoom")
		return
	}
	id, err := Uid.NewId(ROOM_ID)
	if err != nil {
		log.Fatal("Id generator error:", err)
		return
	}
	roomInfo.Uuid = id
	room.Init(gm, roomInfo)

	gm.getMyRoom(f.FromId, roomInfo)
}

//Not Done Yet
func (gm *GameManager) GetRoomStatus(f *CallFuncInfo) {
	//leftSecond := gm.IdMapRoom[f.TargetId].GetStatus()
	gm.SendFuncToClient[f.FromId] <- &CallFuncInfo{}
}

func (gm *GameManager) EnterRoom(f *CallFuncInfo) {
	gm.IdMapRoom[f.TargetId].EnterRoom(f.FromId)
}

func (gm *GameManager) GetAllRoomInfo(f *CallFuncInfo) {
	gm.getAllRoomInfo(f.FromId)
}

func (gm *GameManager) LeaveRoom(f *CallFuncInfo) {
	gm.IdMapRoom[f.TargetId].LeaveRoom(f.FromId)
}

////Send Rpc command to client function

func (gm *GameManager) getAllRoomInfo(userId int64) {
	params := make([]*any.Any, 0)
	for _, room := range gm.IdMapRoom {
		roomInfo := room.GetInfo()
		param, _ := ptypes.MarshalAny(roomInfo)
		params = append(params, param)
	}
	gm.SendFuncToClient[userId] <- &CallFuncInfo{
		Func:  "GetAllRoomInfo",
		Param: params,
	}
}

func (gm *GameManager) getMyRoom(userId int64, roomInfo *RoomInfo) {
	param, _ := ptypes.MarshalAny(roomInfo)
	gm.SendFuncToClient[userId] <- &CallFuncInfo{
		Func:  "GetMyRoom",
		Param: append(make([]*any.Any, 0), param),
	}
}
