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
)

type Player struct {
	Entity
}

func (e *Player) Shoot(f *CallFuncInfo) {
	log.Debug("[entity]{Shoot}", f)
	//create shell
	p, q := e.World.GetTransform(e.EntityInfo.Uuid)
	p.Add(p, lin.NewV3S(1, 0, 1.8))
	//send func to all
	id, _ := Uid.NewId(ENTITY_ID)
	e.World.CreateEntity("Shell", id, *p, *q)
	entityInfo := &Character{}
	entityInfo.Uuid = id
	entityInfo.CharacterType = "Shell"
	log.Debug("{Player}[Shoot]", entityInfo)
	e.Room.CreateShell(entityInfo, physic.V3_LinToMsg(p), physic.Q_LinToMsg(q))
	e.World.Move(id, 10.0, 0.0)
	//set Velocity
}
