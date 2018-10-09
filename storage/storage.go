package storage

import (
	"errors"
	"fmt"
	"os"

	"github.com/globalsign/mgo"
	//"github.com/globalsign/mgo/bson"
	//"os"
	//"regexp"
	log "github.com/sirupsen/logrus"
)

const (
	MGO_DB_NAME            = "gameServer"
	UserInfo_COLLECTION    = "UserInfo"
	RegistInput_COLLECTION = "RegistInput"
)

var (
	_MGO_ADDR string = ""
	_MGO_USER string = ""
	_MGO_PASS string = ""
)

var MgoDb *MongoDb = &MongoDb{}

func init() {
	MgoDb.Init(MGO_DB_NAME, UserInfo_COLLECTION, RegistInput_COLLECTION)
}

type Db interface {
	Init(string, ...string)
	Save(string, interface{}) bool
	Update(string, interface{}, interface{}) bool
	Delete(string, interface{}) bool
	Find(string, interface{}) *mgo.Iter
}

type MongoDb struct {
	Session     *mgo.Session
	Collections map[string]*mgo.Collection
}

type FileStore struct {
	FileName string
	Path     string
}

type MgoError struct {
	Op  string
	Err error
}

func (e *MgoError) Error() string {
	return e.Op + ":" + e.Err.Error()
}

func ReadMgoSettingFromEnv() {
	_MGO_PASS = os.Getenv("MGO_PASS")
	_MGO_USER = os.Getenv("MGO_USER")
	_MGO_ADDR = os.Getenv("MGO_ADDR")
}

func (m *MongoDb) Init(dbName string, collectionNames ...string) {
	ReadMgoSettingFromEnv()
	err := error(nil)
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			fmt.Printf("%T \n\r", err)
			a := &MgoError{"mgo:", errors.New("MgoInit")}
			log.Warn("", a.Error())
			//panic("Please open mangodb server")
		}
	}()
	m.Collections = make(map[string]*mgo.Collection)

	m.Session, err = mgo.Dial(_MGO_ADDR)
	if err != nil {
		fmt.Println("mongodb connecting error :", err)
	}
	if err := m.Session.DB(dbName).Login(_MGO_USER, _MGO_PASS); err != nil {
		log.Fatal("Mgo connect error:", err, "/n Do you set env MGO_USER and MGO_PASS? or Do you create a User?")
	}
	for _, collectionName := range collectionNames {
		m.Collections[collectionName] = m.Session.DB(dbName).C(collectionName)
		fmt.Println(collectionName, " is register.")
	}
}

func (m *MongoDb) Save(collectionName string, data interface{}) bool {
	if _, ok := m.Collections[collectionName]; !ok {
		return false
	}
	m.Collections[collectionName].Insert(data)
	return true
}

func (m *MongoDb) Find(collectionName string, query interface{}) *mgo.Iter {
	if _, ok := m.Collections[collectionName]; !ok {
		log.Warn("No such Collection", collectionName)
		return nil
	}
	iter := m.Collections[collectionName].Find(query).Iter()
	return iter
}

func (m *MongoDb) Update(collectionName string, query interface{}, data interface{}) bool {
	if _, ok := m.Collections[collectionName]; !ok {
		return false
	}
	err := m.Collections[collectionName].Update(query, data)
	if err != nil {
		fmt.Println("db update error:", err)
		return false
	}
	return true

}

func (m *MongoDb) Delete(collectionName string, query interface{}) bool {
	if _, ok := m.Collections[collectionName]; !ok {
		return false
	}
	changelog, err := m.Collections[collectionName].RemoveAll(query)
	fmt.Println(changelog, err)
	if err != nil {
		return false
	}
	return true
}
