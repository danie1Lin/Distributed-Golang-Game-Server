package session

import (
	"sync"
	"time"

	. "github.com/daniel840829/gameServer/msg"
	. "github.com/daniel840829/gameServer/uuid"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	GameServers:          make(map[*GameServer]struct{}),
}

type GameServer struct {
	ExtIp      string
	AgentPort  string
	ClientPort string
	Client     AgentToGameClient
	Rooms      map[*Room]struct{}
	id         string
}

func newGameServer(ExtIp, clientToGamePort, agentToGamePort, id string) (gameServer *GameServer) {
	conn, err := grpc.Dial(ExtIp+":"+agentToGamePort, grpc.WithInsecure())
	if err != nil {
		log.Warn("Agent can't connect to GameServer", err)
		return
	}
	gameServer = &GameServer{
		id:         id,
		ExtIp:      ExtIp,
		ClientPort: ":" + clientToGamePort,
		AgentPort:  ":" + agentToGamePort,
		Rooms:      make(map[*Room]struct{}, 0),
	}
	gameServer.Client = NewAgentToGameClient(conn)
	log.Debug("client", gameServer.Client)
	return
}

type roomManager struct {
	UserIdleRoomListChan map[*MsgChannel]struct{}
	Rooms                map[int64]*Room
	GameServers          map[*GameServer]struct{}
	sync.RWMutex
	RoomList *RoomList
}

func (rm *roomManager) DeleteRoom(room *Room) {
	room.DeleteRoom()
	delete(rm.Rooms, room.Uuid)
	rm.UpdateRoomList()
}

func (rm *roomManager) ConnectGameServer(ExtIp, clientToGamePort, agentToGamePort, id string) {
	rm.GameServers[newGameServer(ExtIp, clientToGamePort, agentToGamePort, id)] = struct{}{}
	log.Debug("Connect Game Server")
}

func (rm *roomManager) GetGameServer() (game *GameServer) {

	for gs, _ := range rm.GameServers {
		if len(gs.Rooms) >= MAX_ROOM_IN_POD {
			continue
		}
		game = gs
		break
	}
	if game == nil {
		ClusterManager.CreatePod()
		game = rm.GetGameServer()
	}
	return
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
		PlayerInRoom: 1,
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
	Name                string
	GameType            string
	Master              *Session
	Client              map[*Session]struct{}
	Uuid                int64
	IsFull              bool
	IsCreatOnGameServer bool
	MaxPlyer            int32
	PlayerInRoom        int32
	Review              *RoomReview
	GameServer          *GameServer
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

func (r *Room) DeleteRoom() bool {
	if r.IsCreatOnGameServer {

		_, err := r.GameServer.Client.DeletGameRoom(context.Background(), &RoomInfo{Uuid: r.Uuid})
		if err != nil {
			log.Warn(err)
		}
	}
	if r.Master != nil {
		r.Master.Room = nil
	}
	for s, _ := range r.Client {
		s.Room = nil
	}
	return true
}

func (r *Room) LeaveRoom(s *Session) bool {
	if s == r.Master {
		s.Room = nil
		r.Master = nil
		r.PlayerInRoom -= 1
		r.IsFull = false
		r.UpdateReview()
		r.UpdateRoomContent()

	} else {
		if _, ok := r.Client[s]; ok {
			s.Room = nil
			delete(r.Client, s)
			r.PlayerInRoom -= 1
			r.IsFull = false
			r.UpdateReview()
			r.UpdateRoomContent()

		} else {
			return false
		}
	}
	if len(r.Client) == 0 && r.Master == nil {
		RoomManager.DeleteRoom(r)
	}
	return true
}

func (r *Room) UpdateRoomContent() {
	content := &RoomContent{
		Uuid:    r.Uuid,
		Players: make(map[string]*PlayerInfo),
	}
	if r.Master != nil {
		pInfo := r.Master.GetPlayerInfo()
		content.Players[pInfo.UserName] = pInfo
	}
	for s, _ := range r.Client {
		pInfo := s.GetPlayerInfo()
		content.Players[pInfo.UserName] = pInfo
	}

	if r.Master != nil {
		ch := r.Master.GetMsgChan("RoomContent")
		ch.DataCh <- content
	}
	for s, _ := range r.Client {
		ch := s.GetMsgChan("RoomContent")
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
	r.CreateRoomOnGameServer()
	return true
	//Start Game
}

func (r *Room) CreateRoomOnGameServer() {
	if !r.IsCreatOnGameServer {
		gameCreation := &GameCreation{
			PlayerSessions: make([]*SessionInfo, 0),
			RoomInfo:       &RoomInfo{Uuid: r.Uuid},
		}
		gameCreation.RoomInfo.GameType = r.GameType
		r.Master.SetState(int32(SessionInfo_ConnectingGame))
		gameCreation.PlayerSessions = append(gameCreation.PlayerSessions, r.Master.GetSessionInfo())
		for s, _ := range r.Client {
			s.SetState(int32(SessionInfo_ConnectingGame))
			gameCreation.PlayerSessions = append(gameCreation.PlayerSessions, s.GetSessionInfo())
		}
		gs := RoomManager.GetGameServer()
		gs.Rooms[r] = struct{}{}
		pem := make(chan *PemKey)
		var getKey func(sig chan *PemKey)
		getKey = func(sig chan *PemKey) {
			key, err := gs.Client.AquireGameRoom(context.Background(), gameCreation)
			if err != nil {
				log.Warn("GameServer has some issue: ", err)
				time.Sleep(2 * time.Second)
				go getKey(sig)
				return
			}
			sig <- key
		}
		go getKey(pem)
		key := <-pem
		c := r.Master.GetMsgChan("ServerInfo")
		serverInfo := &ServerInfo{
			Addr:      gs.ExtIp,
			Port:      gs.ClientPort,
			PublicKey: key.SSL,
		}
		c.DataCh <- serverInfo
		r.Master.ServerInfo = serverInfo
		for s, _ := range r.Client {
			s.ServerInfo = serverInfo
			c = s.GetMsgChan("ServerInfo")
			c.DataCh <- serverInfo
		}
		r.GameServer = gs
		r.IsCreatOnGameServer = true
	}
}
