package session

import (
	. "github.com/daniel840829/gameServer2/msg"
	. "github.com/daniel840829/gameServer2/uuid"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"sync"
)

type ChatMessage struct {
	SpeakerId   int32
	SpeakerName string
	Content     string
}

type ChatRoom struct {
	ReadingBuffer []*ChatMessage
}

//TODO : Chat function

var limitReadyTime = 1000 //wait ms to start a room

var RoomManager *roomManager = &roomManager{
	Rooms:                make(map[int64]*Room),
	UserIdleRoomListChan: make(map[*MsgChannel]struct{}),
}

type roomManager struct {
	UserIdleRoomListChan map[*MsgChannel]struct{}
	Rooms                map[int64]*Room
	sync.RWMutex
	RoomList *RoomList
}

func (rm *roomManager) UpdateRoomList() {
	roomList := &RoomList{
		Item: make([]*RoomReview, 0),
	}
	for _, room := range rm.Rooms {
		roomList.Item = append(roomList.Item, room.GetReview())
	}
	for c, _ := range rm.UserIdleRoomListChan {
		c.DataCh <- roomList
	}
	log.Debug("Update room list", roomList)
}

func (rm *roomManager) AddIdleUserMsgChan(m *MsgChannel) {
	if _, ok := rm.UserIdleRoomListChan[m]; ok {
		return
	}
	rm.UserIdleRoomListChan[m] = struct{}{}
}

func (rm *roomManager) RemoveIdleUserMsgChan(m *MsgChannel) {
	if _, ok := rm.UserIdleRoomListChan[m]; ok {
		delete(rm.UserIdleRoomListChan, m)
	}
}

func (rm *roomManager) Run() {

}

func (rm *roomManager) DeletRoom() {

}

func (rm *roomManager) LeaveRoom() {

}

func (rm *roomManager) CreateRoom(master *Session, setting *RoomSetting) *Room {
	id, _ := Uid.NewId(ROOM_ID)
	room := NewRoom(master, id, setting)
	rm.Rooms[id] = room
	rm.UpdateRoomList()
	return room
}

func NewRoom(master *Session, roomId int64, setting *RoomSetting) *Room {
	room := &Room{
		Master:       master,
		Client:       make(map[*Session]struct{}),
		Uuid:         roomId,
		Name:         setting.Name,
		GameType:     setting.GameType,
		MaxPlyer:     setting.MaxPlayer,
		PlayerInRoom: 0,
		Review: &RoomReview{
			Uuid: roomId,
		},
	}
	room.UpdateReview()
	RoomManager.UpdateRoomList()
	room.UpdateRoomContent()
	return room
}

type Room struct {
	Name         string
	GameType     string
	Master       *Session
	Client       map[*Session]struct{}
	Uuid         int64
	IsFull       bool
	MaxPlyer     int32
	PlayerInRoom int32
	Review       *RoomReview
	sync.RWMutex
}

func (r *Room) GetReview() (Review *RoomReview) {
	r.RLock()
	Review = proto.Clone(r.Review).(*RoomReview)
	r.RUnlock()
	return
}

func (r *Room) UpdateReview() {
	r.Lock()
	r.Review.Name = r.Name
	r.Review.GameType = r.GameType
	r.Review.InRoomPlayer = r.PlayerInRoom
	r.Review.MaxPlayer = r.MaxPlyer
	r.Unlock()
	RoomManager.UpdateRoomList()
}

func (r *Room) EnterRoom(client *Session) bool {
	if _, ok := r.Client[client]; ok {
		return false
	}
	if r.IsFull {
		return false
	}
	r.Client[client] = struct{}{}
	r.PlayerInRoom += 1
	if r.PlayerInRoom == r.MaxPlyer {
		r.IsFull = true
	}
	r.UpdateReview()
	r.UpdateRoomContent()
	return true
}

func (r *Room) KickOut(master *Session, client *Session) bool {
	if master != r.Master {
		return false
	}
	if _, ok := r.Client[client]; ok {
		client.Room = nil
		delete(r.Client, client)
		r.PlayerInRoom -= 1
		r.IsFull = false
		r.UpdateReview()
		r.UpdateRoomContent()
		return true
	}
	return false
}

func (r *Room) DeleteRoom(master *Session) bool {
	if master != r.Master {
		return false
	}
	return true
}

func (r *Room) LeaveRoom(s *Session) bool {
	if s == r.Master {
		r.DeleteRoom(s)
		return true
	} else {
		if _, ok := r.Client[s]; ok {
			s.Room = nil
			delete(r.Client, s)
			r.PlayerInRoom -= 1
			r.IsFull = false
			r.UpdateReview()
			r.UpdateRoomContent()
			return true
		}
	}
	return false
}

func (r *Room) UpdateRoomContent() {
	content := &RoomContent{
		Uuid:    r.Uuid,
		Players: make(map[string]*PlayerInfo),
	}
	pInfo := r.Master.GetPlayerInfo()
	content.Players[pInfo.UserName] = pInfo
	for s, _ := range r.Client {
		pInfo = s.GetPlayerInfo()
		content.Players[pInfo.UserName] = pInfo
	}

	ch := r.Master.GetMsgChan("RoomContent")
	ch.DataCh <- content

	for s, _ := range r.Client {
		ch = s.GetMsgChan("RoomContent")
		ch.DataCh <- content
	}
}

func (r *Room) CheckReady() bool {
	r.UpdateRoomContent()
	if !r.Master.IsReady {
		return false
	}
	for s, _ := range r.Client {
		if !s.IsReady {
			return false
		}
	}
	//r.CreateRoomOnGameServer()
	return true
	//Start Game
}

func (r *Room) CreateRoomOnGameServer() {
	gameCreation := &GameCreation{
		PlayerSessions: make([]*SessionInfo, 0),
		RoomInfo:       &RoomInfo,
	}
	gameCreation.RoomInfo.GameType = r.GameType
	r.Master.SetState(int32(r.Master.State.GetStateCode()) + 1)
	GameCreation.PlayerSessions = append(GameCreation.PlayerSessions, r.Master)
	r.Master.AddMsgChan("ServerInfo", 1)
	for s, _ := range r.Client {
		s.SetState(int32(s.State.GetStateCode()) + 1)
		s.AddMsgChan("ServerInfo", 1)
		GameCreation.PlayerSessions = append(GameCreation.PlayerSessions, s.)
	}
}
