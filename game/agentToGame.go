package game

/*
type AgentToGameServer interface {
	// SessionManager
	AquireGameRoom(context.Context, *GameCreation) (*PemKey, error)
}
*/

func NewATGServer() *ATGServer {
	atg := &ATGServer{}
	return atg
}

type ATGServer struct {
}

func (a *ATGServer) AquireGameRoom(context.Context, *GameCreation) (*PemKey, error) {
	return &PemKey, nil
}
