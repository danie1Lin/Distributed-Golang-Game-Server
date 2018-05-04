package entity

type IGameBehavier interface {
	Init()
	Tick()
	Destroy()
}

type Room struct {
	RoomInfo
	EntityInRoom map[string]IEntity
}

func (r *Room) Init() {

}
func (r *Room) Tick() {

}
func (r *Room) Destroy() {

}
