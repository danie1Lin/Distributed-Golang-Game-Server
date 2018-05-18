package storage

import (
	"fmt"
	"github.com/daniel840829/gameServer/service"
	//"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"testing"
)

func TestMgo(t *testing.T) {
	m := &MongoDb{}
	m.Init("Testing_OnlineGame", "UserInfo", "Character")
	m.Save("UserInfo", &service.UserInfo{Uuid: "123", UserName: "Daniel"})
	iter := m.Find("UserInfo", bson.M{"uuid": "123"})
	r := &service.UserInfo{}
	for iter.Next(r) {
		fmt.Println(r)
	}
	if err := iter.Close(); err != nil {
		return
	}
	m.Update("UserInfo", bson.M{"uuid": "123"}, bson.M{"$set": bson.M{"uuid": "456"}})
	m.Delete("UserInfo", bson.M{"uuid": "456"})
}
