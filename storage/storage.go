package storage

import (
	"os"
)

type Db interface {
	Save(interface{}) bool
	Update(interface{}) bool
	Deletie(interface{}) bool
	Find(interface{}) interface{}
}

type MongoDb struct {
	HostAdd string
}

type FileStore struct {
	FileName string
	Path     string
}
