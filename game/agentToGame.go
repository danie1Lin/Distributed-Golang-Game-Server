package game

import (
	"github.com/daniel840829/gameServer2/game/session"
	. "github.com/daniel840829/gameServer2/msg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	//	io"
	"time"
)

/*
type AgentToGameServer interface {
	// SessionManager
	AquireGameRoom(context.Context, *GameCreation) (*PemKey, error)
}
*/

type ATGServer struct {
}

type CTGServer struct {
}

func (g *CTGServer) TimeCalibrate(c context.Context, e *Empty) (*TimeStamp, error) {
	log.Debug("Calibration")
	return &TimeStamp{Value: int64(time.Now().UnixNano() / 1000000)}, nil
}

func (g *ATGServer) AquireGameRoom(c context.Context, gc *GameCreation) (*PemKey, error) {
	session.RoomManager.CreateGame(gc)
	return &PemKey{SSL: "HI"}, nil
}

func (g *CTGServer) PlayerInput(stream ClientToGame_PlayerInputServer) error {
	log.Debug("input")
	session := GetSesionFromContext(stream.Context())
	if session == nil {
		return status.Errorf(codes.NotFound, "Session Not Found!")
	}
LOOP2:
	for {
		select {
		default:
			input, err := stream.Recv()
			if HandleRPCError(session, err) {
				break LOOP2
			}
			session.State.HandleInput(input)
		}
	}

	return nil
}
func (g *CTGServer) UpdateGameFrame(e *Empty, stream ClientToGame_UpdateGameFrameServer) error {

	log.Debug("Frame")
	session := GetSesionFromContext(stream.Context())
	if session == nil {
		return status.Errorf(codes.NotFound, "Session Not Found!")
	}
	msgch := session.GetMsgChan("GameFrame")
	if msgch == nil {
		return status.Errorf(codes.Internal, "GameFrame MsgChan Not Found!")
	}
	session.State.StartGame()
LOOP:
	for {
		select {
		case <-msgch.StopSignal:
			break LOOP
		case msg := <-msgch.DataCh:
			err := stream.Send(msg.(*GameFrame))
			if HandleRPCError(session, err) {
				break LOOP
			}
		}
	}

	return nil
}
func (g *CTGServer) Pipe(ClientToGame_PipeServer) error {
	return nil
}

func GetSesionFromContext(c context.Context) *session.Session {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil
	}
	s := session.Manager.GetSession(md)
	return s
}

func HandleRPCError(s *session.Session, e error) (IfEndStream bool) {
	if e == nil {
		return false
	}
	st, _ := status.FromError(e)
	log.Warn(st.Message())
	switch st.Code() {
	case codes.Canceled:
		IfEndStream = ReconnectError(s)
	case codes.Unavailable:
		IfEndStream = ReconnectError(s)
	default:
		IfEndStream = RecordError(s, e)
	}
	return
}

func IgnoreError(s *session.Session) (IfEndStream bool) {
	return true
}

func RecordError(s *session.Session, e error) (IfEndStream bool) {
	log.Warn("RPCError:", e)
	EndConnection(s)
	return true
}
func ReconnectError(s *session.Session) (IfEndStream bool) {
	log.Warn("Wait to reconnect")
	s.State.WaitReconnect()
	return true
}

func EndConnection(s *session.Session) (IfEndStream bool) {
	s.State.EndConnection()
	return true
}
