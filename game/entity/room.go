package entity

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/physic"
	"github.com/daniel840829/gameServer/user"
	. "github.com/daniel840829/gameServer/uuid"
	"github.com/gazed/vu/math/lin"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

const (
	FRAME_INTERVAL         = 60 * time.Millisecond
	PHYSIC_UPDATE_INTERVAL = 10 * time.Millisecond
)

type IGameBehavier interface {
	Tick()
	PhysicUpdate()
	Destroy()
	Run()
}

type IRoom interface {
	IGameBehavier
	Init(*GameManager, *RoomInfo)
	//CreateEnitity()
	UserDisconnect(userId int64)
	UserReconnect(userId int64)
	GetInfo() *RoomInfo
	EnterRoom(int64) bool
	LeaveRoom(int64) int64
	Ready(int64) bool
	GetUserInRoom() []int64
}

type Room struct {
	RoomInfo *RoomInfo
	sync.RWMutex
	EntityInRoom map[int64]IEntity
	GM           *GameManager
	EntityOfUser map[int64]int64 //map[id]id
	World        *physic.World
	PosChans     map[int64](chan *Position)
	FuncChans    map[int64](chan *CallFuncInfo)
	roomEndChan  chan (struct{})
}

func (r *Room) UserDisconnect(userId int64) {
	entity, ok := r.EntityInRoom[r.GM.UserIdMapEntityId[userId]]
	if !ok {
		log.Info("Entity of User[", userId, "] is not created yet.")

	} else {

		entity.Destroy()
	}
	r.LeaveRoom(userId)
	if len(r.GetInfo().UserInRoom) == 0 {
		r.Destroy()
	}
}

func (r *Room) UserReconnect(userId int64) {
	//TODO
}

func (r *Room) Init(gm *GameManager, roomInfo *RoomInfo) {
	r.World = &physic.World{}
	r.roomEndChan = make(chan (struct{}))
	r.GM = gm
	r.EntityInRoom = make(map[int64]IEntity)
	r.EntityOfUser = make(map[int64]int64)
	r.PosChans = make(map[int64](chan *Position), 0)
	r.FuncChans = make(map[int64](chan *CallFuncInfo), 0)
	r.Lock()
	r.RoomInfo = roomInfo
	r.World.Init(r.RoomInfo.Uuid)
	r.RoomInfo.ReadyUser = make(map[int64]bool)
	r.RoomInfo.UserInRoom = make(map[int64]*UserInfo)
	r.Unlock()
	log.Info("[", roomInfo.Uuid, "]Room is Create: ", roomInfo)
	go r.Run()
}

func (r *Room) Tick() {
	//Syncpostion
	//callfuncinfo
	//Game Logic
	r.GetAllTransform()
	for _, entity := range r.EntityInRoom {
		entity.Tick()
	}
}
func (r *Room) Destroy() {
	log.Debug("[Room]{Destroy}")
	r.Lock()
	r.roomEndChan <- struct{}{}
	r.World.Destroy()
	r.Unlock()
	r = nil
}

func (r *Room) DestroyEntity(id int64) {
	r.Lock()
	delete(r.EntityInRoom, id)
	if _, ok := r.EntityOfUser[id]; ok {
		delete(r.EntityOfUser, id)
	}
	f := &CallFuncInfo{}
	f.Func = "DestroyEntity"
	f.TargetId = id
	r.SendFuncToAll(f)
	r.Unlock()
}
func (r *Room) GetInfo() *RoomInfo {
	r.RLock()
	roomInfo, _ := proto.Clone(r.RoomInfo).(*RoomInfo)
	r.RUnlock()
	return roomInfo
}

func (r *Room) Ready(userId int64) bool {
	log.Debug("{Room}[Ready]:", userId, " is ready ")
	r.Lock()
	r.RoomInfo.ReadyUser[userId] = true
	r.Unlock()
	return true
}
func (r *Room) EnterRoom(userId int64) bool {
	r.Lock()
	if _, ok := r.RoomInfo.UserInRoom[userId]; ok {
		r.Unlock()
		log.Debug("User", userId, " is already in room")
		return false
	}
	r.PosChans[userId] = r.GM.PosToClient[userId]
	r.FuncChans[userId] = r.GM.SendFuncToClient[userId]
	r.RoomInfo.UserInRoom[userId] = user.Manager.GetUserInfo(userId)
	r.RoomInfo.ReadyUser[userId] = false
	log.Debug("{Room}[EnterRoom]", r.RoomInfo.UserInRoom)
	r.Unlock()
	return true
}

func (r *Room) LeaveRoom(userId int64) int64 {
	r.Lock()
	log.Debug("[room]{LeaveRoom}")
	delete(r.RoomInfo.UserInRoom, userId)
	delete(r.RoomInfo.ReadyUser, userId)
	delete(r.FuncChans, userId)
	delete(r.PosChans, userId)
	r.Unlock()
	return r.RoomInfo.Uuid
}
func (r *Room) Run() {
	//Read Info
	allReady := false
	for !allReady {
		r.RLock()
		for _, ready := range r.RoomInfo.ReadyUser {
			if !ready {
				allReady = false
				break
			} else {
				allReady = true
			}
		}
		r.RUnlock()
		<-time.After(time.Millisecond)
	}
	log.Debug("{room}[Run]:start")
	r.createPlayers()
	r.CreateEnemies()
	r.start()
	physicUpdate := time.NewTicker(PHYSIC_UPDATE_INTERVAL)
	frameUpdate := time.NewTicker(FRAME_INTERVAL)
	for {
		select {
		case <-frameUpdate.C:
			r.Tick()
		case <-physicUpdate.C:
			r.PhysicUpdate()
		case <-r.roomEndChan:
			return
		}
	}
}

func (r *Room) GetUserInRoom() (ids []int64) {
	r.RLock()
	for id, _ := range r.RoomInfo.UserInRoom {
		ids = append(ids, id)
	}
	r.RUnlock()
	return
}

func (r *Room) GetAllTransform() {
	//t1 := time.Now().UnixNano()
	pos := r.World.GetAllTransform()
	//log.Debug("[room]{GetAllTransform}GetAllTransform time:", (time.Now().UnixNano() - t1))
	pos.TimeStamp = int64(time.Now().UnixNano() / 1000000)
	go func() {
		for _, posChan := range r.PosChans {
			posChan <- pos
		}
	}()
}

func (r *Room) SendFuncToAll(f *CallFuncInfo) {

	go func() {
		for _, funcChan := range r.FuncChans {
			funcChan <- f
		}
	}()
}

func (r *Room) CreateEnemies() {
	log.Debug("[Room]{CreateEnemies}")
	for i := 0; i < 1; i++ {
		r.CreateEnemy()
	}
}

func (r *Room) CreateEnemy() {
	r.Lock()
	log.Debug("[Room]{CreateEnemy}")
	entityInfo := &Character{
		MaxHealth:     100.0,
		CharacterType: "Enemy",
		Color: &Color{
			B: 255.0,
			R: 0.0,
			G: 0.0,
		},
		Ability: &Ability{
			SPD:  0.5,
			TSPD: 1.0,
		},
	}
	entityInfo.Uuid, _ = Uid.NewId(ENTITY_ID)
	q := physic.EulerToQuaternion(0.0, 0.0, 0.0)
	p := lin.NewV3S((rand.Float64()*64 - 32), (rand.Float64()*64 - 32), 1.0)
	r.World.CreateEntity("Tank", entityInfo.Uuid, *p, *q)
	entity := r.GM.CreateEntity(r, entityInfo, "Enemy")
	if entity == nil {
		return
	}
	r.createEntity(entity, p, q)
	r.Unlock()
}
func (r *Room) createPlayers() {
	r.RLock()
	for id, userInfo := range r.RoomInfo.UserInRoom {
		log.Debug("[RoomInfo.UserInRoom]", id, userInfo)
		if userInfo.UsedCharacter == int64(0) {
			for id, _ := range userInfo.OwnCharacter {
				userInfo.UsedCharacter = id
				break
			}
		}
		entityInfo := userInfo.OwnCharacter[userInfo.UsedCharacter]
		q := physic.EulerToQuaternion(0.0, 0.0, 0.0)
		p := lin.NewV3S((rand.Float64()*64 - 32), (rand.Float64()*64 - 32), 1.0)
		r.World.CreateEntity("Tank", entityInfo.Uuid, *p, *q)
		entity := r.GM.CreatePlayer(r, "Player", userInfo)
		if entity == nil {
			return
		}
		r.EntityOfUser[entityInfo.Uuid] = userInfo.Uuid
		r.createEntity(entity, p, q)
	}
	r.RUnlock()
}

func (r *Room) createEntity(iEntity IEntity, position *lin.V3, quaternion *lin.Q) {
	entityInfo := iEntity.GetInfo()
	//r.Lock()
	r.EntityInRoom[entityInfo.Uuid] = iEntity

	f := &CallFuncInfo{}
	f.Func = "CreateEntity"
	f.FromPos = &TransForm{Position: physic.V3_LinToMsg(position), Rotation: physic.Q_LinToMsg(quaternion)}
	params := make([]*any.Any, 0)
	param, _ := ptypes.MarshalAny(entityInfo)
	params = append(params, param)
	f.Param = params
	r.SendFuncToAll(f) //r.Unlock()
}

//TODO : These two function should be combined
func (r *Room) CreateShell(iEntity IEntity, entityInfo *Character, p *Vector3, q *Quaternion) {
	//send msg to all
	r.EntityInRoom[entityInfo.Uuid] = iEntity

	f := &CallFuncInfo{}
	f.Func = "CreateShell"
	f.FromPos = &TransForm{Position: p, Rotation: q}
	param, _ := ptypes.MarshalAny(entityInfo)
	params := make([]*any.Any, 0)
	params = append(params, param)
	f.Param = params
	r.SendFuncToAll(f)
}

func (r *Room) start() {
	f := &CallFuncInfo{}
	f.Func = "StartRoom"
	r.RLock()
	for id, _ := range r.RoomInfo.UserInRoom {
		r.GM.SendFuncToClient[id] <- f
	}
	r.RUnlock()
}
func (r *Room) PhysicUpdate() {
	//r.Lock()
	r.World.PhysicUpdate()
	for _, entity := range r.EntityInRoom {
		entity.PhysicUpdate()
	}
	//r.Lock()
}
