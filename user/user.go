package user

import (
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/storage"
	. "github.com/daniel840829/gameServer/uuid"
	"github.com/globalsign/mgo/bson"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
)

func NewUserManager(db storage.Db) *UserManager {
	return &UserManager{
		UserOnline: make(map[int64]*User),
		Db:         db,
	}
}

type UserManager struct {
	sync.RWMutex
	UserOnline map[int64]*User
	Db         storage.Db
}

func (um *UserManager) Login(in *LoginInput) (*UserInfo, error) {
	iter := um.Db.Find(storage.RegistInput_COLLECTION, bson.M{"username": in.UserName})
	userPassword := &RegistInput{}
	userInfo := &UserInfo{}
	if iter.Next(userPassword) {
		if userPassword.Pswd == in.Pswd {
			UserInfoIter := um.Db.Find(storage.UserInfo_COLLECTION, bson.M{"username": in.UserName})
			if UserInfoIter.Next(userInfo) {
				user := &User{}
				user.Login(userInfo)
				um.Lock()
				um.UserOnline[userInfo.Uuid] = user
				um.Unlock()
				log.Debug("[User][Login]", userInfo)
				return userInfo, nil
			}
		}
	}
	return userInfo, nil
}

func (um *UserManager) Logout(userId int64) {
	um.Lock()
	user := um.UserOnline[userId]
	_ = user
	delete(um.UserOnline, userId)
	user = nil
	um.Unlock()
}
func (um *UserManager) Regist(in *RegistInput) (*Error, error) {
	if iter := um.Db.Find(storage.RegistInput_COLLECTION, bson.M{"username": in.UserName}); iter.Next(&RegistInput{}) {
		return &Error{ErrMsg: "Username exists"}, nil
	}
	userInfo := NewUserInfo(in.UserName)
	um.Db.Save(storage.UserInfo_COLLECTION, userInfo)
	um.Db.Save(storage.RegistInput_COLLECTION, in)
	log.Debug(userInfo)
	return &Error{}, nil
}

func (um *UserManager) GetUserInfo(id int64) *UserInfo {
	um.RLock()
	userInfo := um.UserOnline[id].GetInfo()
	um.RUnlock()
	return userInfo
}

type User struct {
	UserInfo *UserInfo
	sync.RWMutex
}

func (u *User) Login(userInfo *UserInfo) {
	u.Lock()
	u.UserInfo = userInfo
	u.Unlock()
}

func (u *User) GetInfo() *UserInfo {
	u.RLock()
	userInfo, _ := proto.Clone(u.UserInfo).(*UserInfo)
	u.RUnlock()
	return userInfo
}

func NewCharacter() (c *Character) {
	c = &Character{}
	uuid, _ := Uid.NewId(CHA_ID)
	c.Uuid = uuid
	c.Color = &Color{int32(rand.Intn(256)), int32(rand.Intn(256)), int32(rand.Intn(256))}
	c.CharacterType = "Player"
	c.MaxHealth = 100.0
	c.Ability = &Ability{}
	c.Ability.SPD = 1.0
	c.Ability.TSPD = 1.0
	return
}

func NewUserInfo(userName string) (u *UserInfo) {
	uuid, _ := Uid.NewId(USER_ID)
	u = &UserInfo{OwnCharacter: make(map[int64]*Character)}
	u.Uuid = uuid
	u.UserName = userName
	c := NewCharacter()
	u.OwnCharacter[c.Uuid] = c
	return
}

var Manager *UserManager

func init() {
	Manager = NewUserManager(storage.MgoDb)
}
