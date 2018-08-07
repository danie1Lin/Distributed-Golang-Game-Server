package entity

import (
	. "github.com/daniel840829/gameServer/msg"
)

type AttackBehavier interface {
	Attack(FromPos Position, Value int64)
	AutoAttack(targetId int64)
}

type AttackBase struct {
	AutoAttackable   bool
	CoolDownLeftTime int32
	MaxValue         float64
	MaxCombo         int32
}

func (atk *AttackBase) Attack(FromPos Position, Value float64) {

}

func (atk *AttackBase) AutoAttack(targetId int64) {

}

func NewShootBehavier(autoAttackable bool, coolDownLeftTime int32, maxValue float64, maxCombo int32) *ShootBehavier {
	return &ShootBehavier{
		AttackBase: AttackBase{
			AutoAttackable:   autoAttackable,
			CoolDownLeftTime: coolDownLeftTime,
			MaxValue:         maxValue,
			MaxCombo:         maxCombo,
		},
	}
}

type ShootBehavier struct {
	AttackBase
}

func (atk *ShootBehavier) Attack(FromPos Position, Value float64) {

}

func (atk *ShootBehavier) AutoAttack(targetId int64) {

}
