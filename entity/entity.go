package entity

import (
	"fmt"
	"github.com/daniel840829/gameServer/msg"
	//p "github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

type Entity struct {
	TypeName string
	Id       uuid.UUID
	I        IEntity
}

func (e *Entity) Say(s string) {
	fmt.Println(e.Id, "[", e.TypeName, "]:", s)
}

type IEntity interface {
	Init()
	Tick()
	Say(s string)
}

type Player struct {
	Entity
}

func (e *Player) Hit(damage int32) {
	fmt.Println("-", damage)
}

func (e *Player) Init() {
	//call All client create enitity at some point
	msg.CallAllClient(e.TypeName, e.Id, "CreateEnitity", &msg.Pos{Id: e.Id.String()})
	//
}

func (e *Player) Tick() {
}

type EntityInfo struct {
	//mathod's name map to Mathod's info
	MethodMap map[string]EntityMathod
	Type      reflect.Type
}

type EntityMathod struct {
	Func reflect.Value
	Type reflect.Type
	Args int
}

type EntityManager struct {
	EntityTypeMap map[string]EntityInfo
	EntityIdMap   map[uuid.UUID]reflect.Value
}

var eManager *EntityManager = &EntityManager{
	EntityTypeMap: make(map[string]EntityInfo),
	EntityIdMap:   make(map[uuid.UUID]reflect.Value),
}

func (em *EntityManager) Call(entityTypeName string, id uuid.UUID, fName string, args ...reflect.Value) {
	e, ok := em.EntityIdMap[id]
	eInfo, ok := em.EntityTypeMap[entityTypeName]
	if !ok {
		panic("Id not found")
	}

	f := eInfo.MethodMap[fName]
	fmt.Println("f:", f)
	in := make([]reflect.Value, f.Args)
	in[0] = e
	for i := 1; i < f.Args; i++ {
		in[i] = args[i-1]
	}
	f.Func.Call(in)
}

func (em *EntityManager) CreateEnitity(entityTypeName string, isClient bool) (id uuid.UUID) {
	entityInfo, ok := em.EntityTypeMap[entityTypeName]
	if !ok {
		fmt.Println(entityTypeName, "is not regist.")
	}
	vEntityPtr := reflect.New(entityInfo.Type)
	//check uuid repeat
	err := error(nil)
	id, err = uuid.NewV4()
	fmt.Println(id, err)
	for _, ok := em.EntityIdMap[id]; ok; {
		id, _ = uuid.NewV4()
		fmt.Println(id, err)
	}
	em.EntityIdMap[id] = vEntityPtr
	vEntityPtr.Elem().FieldByName("Id").Set(reflect.ValueOf(id))
	vEntityPtr.Elem().FieldByName("TypeName").Set(reflect.ValueOf(vEntityPtr.Type().Elem().Name()))
	em.Call(entityTypeName, id, "Init")
	return
}

func RegisterEnitity(iEntity IEntity) {
	rEntity := reflect.ValueOf(iEntity)
	tEntity := rEntity.Type()
	entityName := tEntity.Elem().Name()
	rEntity.Elem().FieldByName("TypeName").Set(reflect.ValueOf(entityName))
	fmt.Println("t:", tEntity, "v:", rEntity, "m:", rEntity.NumMethod())
	entityInfo := &EntityInfo{MethodMap: make(map[string]EntityMathod)}
	entityInfo.Type = tEntity.Elem()
	for i := 0; i < rEntity.NumMethod(); i++ {
		m := tEntity.Method(i)
		em := EntityMathod{m.Func, m.Type, m.Type.NumIn()}
		entityInfo.MethodMap[tEntity.Method(i).Name] = em
	}
	fmt.Println(entityInfo)
	eManager.EntityTypeMap[entityName] = *entityInfo
	fmt.Println(eManager)
}
