package storage

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	//"github.com/globalsign/mgo/bson"
	//"os"
	//"regexp"
)

const (
	_MGO_ADDR              = "127.0.0.1:27017"
	MGO_DB_NAME            = "gameServer"
	UserInfo_COLLECTION    = "UserInfo"
	RegistInput_COLLECTION = "RegistInput"
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

func (m *MongoDb) Init(dbName string, collectionNames ...string) {
	err := error(nil)
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			fmt.Printf("%T \n\r", err)
			a := &MgoError{"fuck", errors.New("MgoInit")}
			fmt.Println(a.Error())
			//panic("Please open mangodb server")
		}
	}()
	m.Collections = make(map[string]*mgo.Collection)

	m.Session, err = mgo.Dial(_MGO_ADDR)
	if err != nil {
		fmt.Println("mongodb connecting error :", err)
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
		fmt.Println("No such Collection", collectionName)
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
