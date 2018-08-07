package entity

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/physic"
	//p "github.com/golang/protobuf/proto"
	. "github.com/daniel840829/gameServer/uuid"
	"github.com/gazed/vu/math/lin"
	//"github.com/golang/protobuf/proto"
	//"github.com/ianremmler/ode"
	log "github.com/sirupsen/logrus"
	//"os"
	//"reflect"
	//"sync"
	//"time"
	"math"
)

type Player struct {
	Entity
}

func (e *Player) Init(gm *GameManager, room *Room, entityInfo *Character) {
	e.Entity.Init(gm, room, entityInfo)
	e.Health = e.EntityInfo.MaxHealth
	log.Info("player's costumeInit")
}

func (e *Player) Shoot(f *CallFuncInfo) {
	log.Debug("[entity]{Shoot}", f)
	//create shell
	p, q := e.World.GetTransform(e.EntityInfo.Uuid)
	muzzle := lin.NewV3().MultQ(lin.NewV3S(0, 1.1, 1.8), q)
	p.Add(p, muzzle)
	id, _ := Uid.NewId(ENTITY_ID)
	e.World.CreateEntity("Shell", id, *p, *q)
	entityInfo := &Character{}
	entityInfo.Uuid = id
	entityInfo.CharacterType = "Shell"
	log.Debug("{Player}[Shoot]", entityInfo)
	e.World.Move(id, float64(f.Value/10), 0.0)
	entity := e.GM.CreateEntity(e.Room, entityInfo, "Shell")
	e.Room.CreateShell(entity, entityInfo, physic.V3_LinToMsg(p), physic.Q_LinToMsg(q))
	//set Velocity
}

func (e *Player) PhysicUpdate() {
	q := e.Obj.CBody.Quaternion()
	rot := e.Obj.CBody.AngularVel()
	q[1] = 0.0
	q[2] = 0.0
	len := math.Sqrt(q[0]*q[0] + q[3]*q[3])
	q[0] /= len
	q[3] /= len
	rot[0] = 0.0
	rot[1] = 0.0
	e.Obj.CBody.SetQuaternion(q)
	e.Obj.CBody.SetAngularVelocity(rot)
	e.Obj.SyncAOEPos()
}
