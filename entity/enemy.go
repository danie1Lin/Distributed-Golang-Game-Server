package entity

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/physic"
	"github.com/gazed/vu/math/lin"
	"github.com/ianremmler/ode"
	log "github.com/sirupsen/logrus"
	"math"
	/*
		"fmt"
		//p "github.com/golang/protobuf/proto"
		//. "github.com/daniel840829/gameServer/uuid"
		"github.com/golang/protobuf/proto"
		"github.com/satori/go.uuid"
		"os"
		"reflect"
		"sync"
		//"time"
	*/)

type Enemy struct {
	Player
	CD int64
}

func (e *Enemy) PhysicUpdate() {
	e.Player.PhysicUpdate()
}

func (e *Enemy) Tick() {
	e.CD += 1
	//
	//var targetId int64
	var targetPos ode.Vector3
	var targetDis float64 = 30
	var isFindTarget bool = false
	var isReadyToAttack bool = false
	//Loop obj in AOE
	for id, pos := range e.Obj.AOEObjs {
		//Check if it is real player
		if _, ok := e.Room.EntityOfUser[id]; ok {
			dis := physic.V3_OdeToLin(e.Obj.CBody.PosRelPoint(pos)).Len()
			if targetDis > dis && dis > 1 {
				targetDis = dis
				//targetId = id
				targetPos = pos
				isFindTarget = true
				if targetDis < 20 && e.CD > 100 {
					isReadyToAttack = true
				}
			}
		}
	}
	if isFindTarget {
		Vel := ode.NewVector3()
		for i, v := range e.Obj.CBody.PosRelPoint(targetPos) {
			Vel[i] = v / targetDis * 2
		}
		Vel[2] = 0
		log.Debug("[enemy]{Tick}FindTarget vel: ", Vel)
		/*
			cos_theta := forwardVector.Unit().Dot(relVel.Unit())
			if cos_theta > 1 {
				cos_theta = 1
				return
			} else if cos_theta < -1 {
				cos_theta = -1
				return
			}
			angle := math.Acos(cos_theta)
			if angle > math.Pi || angle < 0 {
				log.Debug("[Enemy] angle out", angle*180/math.Pi)
			} else if math.IsNaN(angle) {
				log.Debug("angle,cos_theta", angle, cos_theta)
			}
		*/
		/*
			    float m = sqrt(2.f + 2.f * dot(u, v));
				    vec3 w = (1.f / m) * cross(u, v);
					    return quat(0.5f * m, w.x, w.y, w.z);
		*/
		/*
			axis := lin.NewV3().Cross(forwardVector, relVel).Unit()
		*/
		/*
			relVel := physic.V3_OdeToLin(targetPos).Unit()
			forwardVector := lin.NewV3S(0, -1, 0).Unit()
			m := math.Sqrt(2.0 + 2.0*relVel.Dot(forwardVector))
			axis := lin.NewV3().Scale(lin.NewV3().Cross(forwardVector, relVel), 1/m)
			targetQ := lin.NewQ().SetS(axis.X, axis.Y, axis.Z, 0.5*m)
		*/
		relVel := lin.NewV3().Sub(physic.V3_OdeToLin(targetPos), physic.V3_OdeToLin(e.Obj.CBody.Position()))
		angle := -1 * math.Atan2(relVel.X, relVel.Y)
		targetQ := lin.NewQ().SetS(0, 0, math.Sin(angle/2), math.Cos(angle/2))
		NowQ := physic.Q_OdeToLin(e.Obj.CBody.Quaternion())
		if isReadyToAttack {
			e.Obj.CBody.SetQuaternion(physic.Q_LinToOde(targetQ))
			e.Shoot(&CallFuncInfo{Value: 15})
			e.CD = 0
		} else {
			MoveToQ := lin.NewQ().Nlerp(NowQ, targetQ, 0.3)
			log.Debugf("[Enemy]{Tick} m%v,targetQ:%v,NowQ:%v,MoveToQ:%v", targetQ, NowQ, MoveToQ)
			e.Obj.CBody.SetQuaternion(physic.Q_LinToOde(MoveToQ))
			input := &Input{
				V_Movement: float32(targetDis / 25.0),
			}
			e.Move(input)

		}
		//e.Obj.CBody.SetAngularVelocity(ode.NewVector3(0, angle*0.5, 0))
		//e.Obj.CBody.SetLinearVelocity(e.Obj.CBody.VectorToWorld(Vel))
	} else {
		input := &Input{}
		e.Move(input)
	}
	e.Obj.ClearAOE()
}
