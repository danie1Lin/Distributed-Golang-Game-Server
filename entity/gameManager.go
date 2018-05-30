package entity

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/service"
	. "github.com/daniel840829/gameServer/uuid"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

type GameManager struct {
	////basic information
	Uuid int64
	////Reflection
	thisType reflect.Type
	////child component
	//rpcCallableObject map[int64]reflect.Value
	//room
	rl          sync.RWMutex
	TypeMapRoom map[string]reflect.Type
	IdMapRoom   map[int64]IRoom
	//entity
	el                sync.RWMutex
	TypeMapEntity     map[string]reflect.Type
	IdMapEntity       map[int64]IEntity
	UserIdMapEntityId map[int64]int64
	////rpc channel
	SendFuncToClient   map[int64](chan *CallFuncInfo)
	RecvFuncFromClient chan *CallFuncInfo
	PosToClient        map[int64](chan *Position)
	InputFromClient    chan *Input
	ErrToClient        map[int64](chan *Error)
	ErrFromClient      chan *Error
}

func (gm *GameManager) Init(rpc *service.Rpc) {
	gm.Uuid, _ = Uid.NewId(GM_ID)
	gm.IdMapRoom = make(map[int64]IRoom)
	gm.TypeMapRoom = make(map[string]reflect.Type)
	gm.TypeMapEntity = make(map[string]reflect.Type)
	gm.IdMapEntity = make(map[int64]IEntity)
	gm.UserIdMapEntityId = make(map[int64]int64)
	//gm.rpcCallableObject = make(map[int64]reflect.Value)
	gm.SendFuncToClient = rpc.SendFuncToClient
	gm.RecvFuncFromClient = rpc.RecvFuncFromClient
	gm.PosToClient = rpc.PosToClient
	gm.InputFromClient = rpc.InputFromClient
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
		case input := <-gm.InputFromClient:
			gm.SyncPos(input)
		}
	}
}

func (gm *GameManager) RegistRoom(roomTypeName string, iRoom IRoom) {
	if _, ok := gm.TypeMapRoom[roomTypeName]; ok {
		log.Fatal(roomTypeName, "is already registed.")
	}
	vRoom := reflect.ValueOf(iRoom)
	gm.TypeMapRoom[roomTypeName] = vRoom.Type().Elem()
	//TODO : Record method info to speed up reflection invoke.
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

func (gm *GameManager) SyncPos(input *Input) {
	//Deal with value
	entity, ok := gm.IdMapEntity[gm.UserIdMapEntityId[input.UserId]]
	if !ok {
		log.Warn("No Such Entity id", input.UserId)
	}
	//TODO
	//if timestamp is too far from now

	entity.Move(input)
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
	gm.rl.Lock()
	gm.IdMapRoom[id] = room
	gm.rl.Unlock()
	gm.getMyRoom(f.FromId, roomInfo)
}

func (gm *GameManager) RegistEnitity(EntityTypeName string, iEntity IEntity) {
	if _, ok := gm.TypeMapEntity[EntityTypeName]; ok {
		log.Fatal(EntityTypeName, "is already registed.")
	}
	vEntity := reflect.ValueOf(iEntity)
	gm.TypeMapEntity[EntityTypeName] = vEntity.Type().Elem()

}

func (gm *GameManager) CreatePlayer(room *Room, entityType string, userInfo *UserInfo) IEntity {
	roomId := room.GetInfo().Uuid
	tEntity, ok := gm.TypeMapEntity[entityType]
	if !ok {
		log.Warn(entityType, " is not registed yet. ")
		return nil
	}
	entity, ok := reflect.New(tEntity).Interface().(IEntity)
	if !ok {
		log.Warn("Something Wrong with RegistEnitity")
		return nil
	}
	entityInfo := userInfo.OwnCharacter[userInfo.UsedCharacter]
	entity.Init(gm, room, entityInfo)

	gm.rl.Lock()
	gm.IdMapEntity[entityInfo.Uuid] = entity
	gm.UserIdMapEntityId[userInfo.Uuid] = entityInfo.Uuid
	gm.rl.Unlock()

	gm.createEntity(roomId, entityInfo)
	return entity
}

//Not Done Yet
func (gm *GameManager) GetRoomStatus(f *CallFuncInfo) {
	//leftSecond := gm.IdMapRoom[f.TargetId].GetStatus()
	gm.SendFuncToClient[f.FromId] <- &CallFuncInfo{}
}

func (gm *GameManager) EnterRoom(f *CallFuncInfo) {
	log.Debug("{GM}[EnterRoom]Excute")
	if ok := gm.IdMapRoom[f.TargetId].EnterRoom(f.FromId); ok {
		gm.enterRoom(f.FromId, gm.IdMapRoom[f.TargetId].GetInfo())
	}
}

func (gm *GameManager) GetAllRoomInfo(f *CallFuncInfo) {
	gm.getAllRoomInfo(f.TargetId)
}

func (gm *GameManager) LeaveRoom(f *CallFuncInfo) {
	gm.IdMapRoom[f.TargetId].LeaveRoom(f.FromId)
}

func (gm *GameManager) GetLoginData(f *CallFuncInfo) {
	gm.getLoginData(f.FromId)
}

func (gm *GameManager) ReadyRoom(f *CallFuncInfo) {
	log.Info("{GM}[ReadyRoom]User [", f.FromId, "] is ready in [", f.TargetId, "] Room")
	gm.IdMapRoom[f.TargetId].Ready(f.FromId)
}

////Send Rpc command to client function
func (gm *GameManager) createEntity(roomId int64, entityInfo *Character) {
	f := &CallFuncInfo{}
	f.Func = "CreateEntity"
	//TODO:param
	params := make([]*any.Any, 0)
	param, _ := ptypes.MarshalAny(entityInfo)
	params = append(params, param)
	f.Param = params
	for _, id := range gm.IdMapRoom[roomId].GetUserInRoom() {
		gm.SendFuncToClient[id] <- f
	}
}

func (gm *GameManager) getLoginData(userId int64) {
	//
	//b := &BasicType{}

	//param, _ := ptypes.MarshalAny(roomInfo)

}

func (gm *GameManager) enterRoom(userId int64, roomInfo *RoomInfo) {
	params := make([]*any.Any, 0)
	param, _ := ptypes.MarshalAny(roomInfo)
	params = append(params, param)
	f := &CallFuncInfo{}
	f.Func = "EnterRoom"
	f.Param = params
	gm.SendFuncToClient[userId] <- f
}
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

func init() {

}
