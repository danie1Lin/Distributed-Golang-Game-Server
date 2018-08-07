package game

/*
type ClientToGameServer interface {
	// roomManager
	EnterRoom(context.Context, *Empty) (*Success, error)
	LeaveRoom(context.Context, *Empty) (*Success, error)
	// entityManager
	PlayerInput(ClientToGame_PlayerInputServer) error
	// View
	UpdateGameFrame(*Empty, ClientToGame_UpdateGameFrameServer) error
	Pipe(ClientToGame_PipeServer) error
}
*/

type CTGServer struct {
}

func (ctg *CTGServer) PlayerInput(ClientToGame_PlayerInputServer) error
func (ctg *CTGServer) UpdateGameFrame(*Empty, ClientToGame_UpdateGameFrameServer) error
func (ctg *CTGServer) Pipe(ClientToGame_PipeServer) error
