package timeCalibration

import (
	log "github.com/sirupsen/logrus"
	"runtime"
	"sync"
	"time"
)

type TimeCalibration struct {
	RunningNoMapUserId sync.Map
	Proccess           sync.Map
}

func (tc *TimeCalibration) FromClientTimeDelay(userId int64, t int64) {

}

func (tc *TimeCalibration) ToClientTimeDelay(userId int64, t int64) {

}

type Proccess struct {
	RunningNo   int64
	RpcFuncName string
	StageTime   []*StageInfo
}

type StageInfo struct {
	Fn string
}

func (p *Proccess) String() {

}
