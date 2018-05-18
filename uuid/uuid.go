package uuid

import (
	log "github.com/sirupsen/logrus"
	"github.com/zheng-ji/goSnowFlake"
)

const (
	GM_ID     = "GameManager"
	ROOM_ID   = "Room"
	USER_ID   = "User"
	ENTITY_ID = "Entity"
	EQUIP_ID  = "Equipment"
	CHA_ID    = "Character"
)

type UID struct {
	Gen         map[string]*goSnowFlake.IdWorker
	WorkersName map[int64]string
}

func (u *UID) RegisterWorker(workers ...string) {
	u.Gen = make(map[string]*goSnowFlake.IdWorker)
	u.WorkersName = make(map[int64]string)
	err := error(nil)
	for idx, worker := range workers {
		u.Gen[worker], err = goSnowFlake.NewIdWorker(int64(idx + 1))
		u.WorkersName[int64(idx+1)] = worker
		if err != nil {
			log.Fatal(err)
			break
		}
	}
}

func (u *UID) NewId(worker string) (id int64, err error) {
	id, err = u.Gen[worker].NextId()
	return
}

func (u *UID) ParseId(id int64) (worker string, ts int64) {
	_, ts, workerId, _ := goSnowFlake.ParseId(id)
	worker, ok := u.WorkersName[workerId]
	if !ok {
		worker = ""
	}
	return
}

var Uid *UID = &UID{}

func init() {
	Uid.RegisterWorker(GM_ID, ROOM_ID, USER_ID, CHA_ID, ENTITY_ID, EQUIP_ID)
}
