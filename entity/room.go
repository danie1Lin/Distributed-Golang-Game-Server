package entity

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/physic"
	"github.com/daniel840829/gameServer/user"
	"github.com/gazed/vu/math/lin"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
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
	GetInfo() *RoomInfo
	EnterRoom(int64) bool
	LeaveRoom(int64) bool
	Ready(int64) bool
	GetUserInRoom() []int64
}

type Room struct {
	RoomInfo *RoomInfo
	sync.RWMutex
	EntityInRoom map[int64]IEntity
	GM           *GameManager
	EntityOfUser map[int64]int64
	World        *physic.World
	PosChans     [](chan *Position)
}

func (r *Room) Init(gm *GameManager, roomInfo *RoomInfo) {
	r.World = &physic.World{}
	r.GM = gm
	r.EntityInRoom = make(map[int64]IEntity)
	r.EntityOfUser = make(map[int64]int64)
	r.PosChans = make([](chan *Position), 0)
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
}

func (r *Room) GetInfo() *RoomInfo {
	log.Debug("[Room][GetInfo] wait get lock")
	r.RLock()
	roomInfo, _ := proto.Clone(r.RoomInfo).(*RoomInfo)
	log.Debug("[Room][GetInfo]", roomInfo, r.RoomInfo)
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
	log.Debug("[Room][EnterRoom] wait get lock")
	r.Lock()
	/*
		if _, ok := r.RoomInfo.UserInRoom[userId]; ok {
			r.roomInfoLock.Unlock()
			return false
		}
	*/
	r.PosChans = append(r.PosChans, r.GM.PosToClient[userId])
	r.RoomInfo.UserInRoom[userId] = user.Manager.GetUserInfo(userId)
	r.RoomInfo.ReadyUser[userId] = false
	log.Debug("{Room}[EnterRoom]", r.RoomInfo.UserInRoom)
	r.Unlock()
	return true
}

func (r *Room) LeaveRoom(userId int64) bool {
	r.Lock()
	delete(r.RoomInfo.UserInRoom, userId)
	delete(r.RoomInfo.ReadyUser, userId)
	r.Unlock()
	return false
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
	r.start()
	physicUpdate := time.NewTicker(PHYSIC_UPDATE_INTERVAL)
	frameUpdate := time.NewTicker(FRAME_INTERVAL)
	for {
		select {
		case <-frameUpdate.C:
			go r.Tick()
		case <-physicUpdate.C:
			go r.PhysicUpdate()
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
	pos := r.World.GetAllTransform()
	for _, posChan := range r.PosChans {
		posChan <- pos
	}
}
func (r *Room) createPlayers() {
	r.RLock()
	for _, userInfo := range r.RoomInfo.UserInRoom {
		if userInfo.UsedCharacter == int64(0) {
			for id, _ := range userInfo.OwnCharacter {
				userInfo.UsedCharacter = id
				break
			}
		}
		entity := r.GM.CreatePlayer(r, "Player", userInfo)
		if entity == nil {
			return
		}
		r.createEntity(entity)
		q := physic.EulerToQuaternion(0.0, 0.0, 0.0)
		r.World.CreateEnitity("Tank", entity.GetInfo().Uuid, *lin.NewV3S(10, 10, 10), *q)
	}
	r.RUnlock()
}

func (r *Room) createEntity(iEntity IEntity) {
	entityInfo := iEntity.GetInfo()
	//r.Lock()
	r.EntityInRoom[entityInfo.Uuid] = iEntity
	//r.Unlock()
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
	r.World.PhysicUpdate()
	for _, entity := range r.EntityInRoom {
		entity.PhysicUpdate()
	}
}
