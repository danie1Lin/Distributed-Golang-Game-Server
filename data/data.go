package data

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type Obj struct {
	Name  string
	Type  string
	Shape string
	Lens  []float64
}

type Objs map[string]Obj

func ReadObjData() Objs {
	raw, err := ioutil.ReadFile("./data/PhysicObj.json")
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	var c Objs
	json.Unmarshal(raw, &c)
	return c
}

var ObjData Objs

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Info("{data}[init] Loading Data...")
	ObjData = ReadObjData()
	log.Info("{data}[init] ObjData: ", ObjData)
}
