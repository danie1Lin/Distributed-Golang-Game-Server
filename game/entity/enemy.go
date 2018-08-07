package entity

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/physic"
	"github.com/daniel840829/ode"
	"github.com/gazed/vu/math/lin"
	log "github.com/sirupsen/logrus"
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
	e.FindTargetAndAttack(30.0)
}

func (e *Enemy) FindTargetAndAttack(searchRadius float64) {
	e.CD += 1
	//var targetId int64
	var targetPos ode.Vector3
	var targetDis float64 = searchRadius
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
		directionV3 := lin.NewV3().Sub(physic.V3_OdeToLin(targetPos), physic.V3_OdeToLin(e.Obj.CBody.Position()))
		targetQ := physic.Q_OdeToLin(physic.DirectionV3ToQuaternion(directionV3))
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
	} else {
		input := &Input{}
		e.Move(input)
	}
	e.Obj.ClearAOE()
}
