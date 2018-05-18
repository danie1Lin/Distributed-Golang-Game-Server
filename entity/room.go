package entity

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/user"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	FRAME_INTERVAL         = 100 * time.Millisecond
	PHYSIC_UPDATE_INTERVAL = 10 * time.Millisecond
)

type IGameBehavier interface {
	Init(*GameManager, *RoomInfo)
	Tick()
	PhysicUpdate()
	Destroy()
	Run()
}

type IRoom interface {
	IGameBehavier
	CreateEnitity()
	GetInfo() *RoomInfo
	EnterRoom(int64)
	LeaveRoom(int64)
}

type Room struct {
	RoomInfo     *RoomInfo
	roomInfoLock sync.RWMutex
	EntityInRoom map[string]IEntity
	GM           *GameManager
}

func (r *Room) Init(gm *GameManager, roomInfo *RoomInfo) {
	r.EntityInRoom = make(map[string]IEntity)
	r.roomInfoLock.Lock()
	r.RoomInfo = roomInfo
	r.RoomInfo.ReadyUser = make(map[int64]bool)
	r.roomInfoLock.Unlock()
	log.Info("[", roomInfo.Uuid, "]Room is Create: ", roomInfo)
	go r.Run()
}

func (r *Room) Tick() {
	//Syncpostion
	//callfuncinfo
	//Game Logic
	for _, entity := range r.EntityInRoom {
		entity.Tick()
	}
}
func (r *Room) Destroy() {
}
func (r *Room) GetInfo() *RoomInfo {
	r.roomInfoLock.RLock()
	roomInfo, _ := proto.Clone(r.RoomInfo).(*RoomInfo)
	r.roomInfoLock.RUnlock()
	return roomInfo
}

func (r *Room) EnterRoom(userId int64) {
	r.roomInfoLock.Lock()
	r.RoomInfo.UserInRoom[userId] = user.Manager.GetUserInfo(userId)
	r.RoomInfo.ReadyUser[userId] = false
	r.roomInfoLock.Unlock()
}

func (r *Room) LeaveRoom(userId int64) {
	r.roomInfoLock.Lock()
	delete(r.RoomInfo.UserInRoom, userId)
	delete(r.RoomInfo.ReadyUser, userId)
	r.roomInfoLock.Unlock()
}
func (r *Room) Run() {
	r.roomInfoLock.RLock()
	//Read Info
	allReady := false
	for !allReady {
		for _, ready := range r.RoomInfo.ReadyUser {
			if !ready {
				allReady = false
				break
			} else {
				allReady = true
			}
		}
		<-time.After(time.Millisecond)
	}
	r.roomInfoLock.RUnlock()
	physicUpdate := time.NewTicker(PHYSIC_UPDATE_INTERVAL)
	frameUpdate := time.NewTicker(FRAME_INTERVAL)
	for {
		select {
		case <-frameUpdate.C:
			r.Tick()
		case <-physicUpdate.C:
			r.PhysicUpdate()
		}
	}
}
func (r *Room) PhysicUpdate() {
	for _, entity := range r.EntityInRoom {
		entity.PhysicUpdate()
	}
}
