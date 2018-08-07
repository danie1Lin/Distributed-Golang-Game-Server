package entity

import (
	//. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/physic"
	//p "github.com/golang/protobuf/proto"
	//"github.com/gazed/vu/math/lin"
	//. "github.com/daniel840829/gameServer/uuid"
	//"github.com/golang/protobuf/proto"
	"github.com/daniel840829/ode"
	//"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	//"os"
	//"reflect"
	//"sync"
	//"time"
)

const (
	ExplosionRadius = 5.0
)

type Shell struct {
	Entity
	State int
	//0 : not explosion
	//1 : explode
}

func (s *Shell) PhysicUpdate() {
	objData := s.Obj.GetData()
	//log.Debug("[Shell]CollideTimes", objData.CollideTimes)
	switch s.State {
	case 0:
		if objData.CollideTimes > 0 {
			s.Obj.ResetCollide()
			//Send explosion
			shellGeom := s.Obj.CGeom
			q := shellGeom.Quaternion()
			p := shellGeom.Position()
			shellGeom.Destroy()
			s.Obj.CBody.Destroy()
			CGeom := s.World.Space.NewSphere(ExplosionRadius)
			CGeom.SetData(s.Obj)
			CGeom.SetCategoryBits(0)
			CGeom.SetCollideBits(physic.SetBitExcept(physic.SetAllBits(), physic.Skill_Bit, physic.Terrain_Bit))
			s.Obj.CGeom = CGeom
			body := s.World.World.NewBody()
			body.SetKinematic(true)
			body.SetPosition(p)
			body.SetQuaternion(q)
			body.SetLinearVelocity(ode.NewVector3(0.0, 0.0, 0.0))
			body.SetAngularVelocity(ode.NewVector3(0.0, 0.0, 0.0))
			body.SetGravityEnabled(false)
			s.Obj.CBody = body
			s.Obj.CGeom.SetBody(body)
			s.State = 1
		}
	case 1:
		log.Debug("[Shell]{explode}")
		f := func(obj *physic.Obj, times int64) bool {
			log.Debug(obj.GetData(), " is damage ", times)
			//id := obj.GetData().Uuid
			id := obj.GetData().Uuid
			entity, ok := s.Room.EntityInRoom[id]
			if !ok {
				log.Debug("Something Unkown harm")
				return true
			}
			entity.Harm(10)
			return true
		}
		s.Obj.LoopCollideObj(f)
		s.Obj.ResetCollide()
		s.State = 2
	default:
		s.Destroy()
		break
	}
}
