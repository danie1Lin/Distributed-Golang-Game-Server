package entity

import (
	"fmt"
	. "github.com/daniel840829/gameServer/msg"
	"github.com/daniel840829/gameServer/physic"
	//p "github.com/golang/protobuf/proto"
	//"github.com/gazed/vu/math/lin"
	//. "github.com/daniel840829/gameServer/uuid"
	"github.com/golang/protobuf/proto"
	//"github.com/ianremmler/ode"
	log "github.com/sirupsen/logrus"
	"os"
	//"reflect"
	"sync"
	//"time"
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
	sync.RWMutex
	EntityInfo *Character
	TypeName   string
	Health     float32
	Alive      bool
	I          IEntity
	GM         *GameManager
	Room       *Room
	World      *physic.World
	Obj        *physic.Obj
	Skill      map[string]AttackBehavier
}
type IEntity interface {
	IGameBehavier
	Hit(int32)
	GetInfo() *Character
	Init(*GameManager, *Room, *Character)
	Move(in *Input)
	GetTransform() *TransForm
	Harm(blood float32)
}

func (e *Entity) GetInfo() *Character {
	e.RLock()
	entityInfo := proto.Clone(e.EntityInfo).(*Character)
	e.RUnlock()
	return entityInfo
}
func (e *Entity) Hit(damage int32) {
	fmt.Println("-", damage)
}

func (e *Entity) Harm(blood float32) {
	e.Lock()
	e.Health -= blood
	if e.Health <= 0 {
		//Dead
		e.Alive = false
		e.Unlock()
		e.Destroy()
		return
	}
	f := &CallFuncInfo{
		Func:     "Health",
		Value:    e.Health,
		TargetId: e.EntityInfo.Uuid,
	}
	e.Room.SendFuncToAll(f)
	e.Unlock()
}

func (e *Entity) Init(gm *GameManager, room *Room, entityInfo *Character) {
	e.GM = gm
	e.EntityInfo = entityInfo
	e.Room = room
	e.World = room.World
	var ok bool
	e.Obj, ok = room.World.Objs.Get(entityInfo.Uuid)
	if !ok {
		log.Fatal("[entity]{init} Get obj ", entityInfo.Uuid, " is not found. ")
	}
	//call All client create enitity at some point
	e.costumeInit()
}

func (e *Entity) costumeInit() {
	log.Warn("Please define your costumeInit")
}
func (e *Entity) Tick() {
}
func (e *Entity) Destroy() {
	e.Lock()
	e.GM.DestroyEntity(e.EntityInfo.Uuid)
	e.Room.DestroyEntity(e.EntityInfo.Uuid)
	e.World.DeleteObj(e.EntityInfo.Uuid)
	e.Obj.Destroy()
	e.Obj = nil
	e.Unlock()
	e = nil
}

func (e *Entity) Run() {
}
func (e *Entity) PhysicUpdate() {
}

func (e *Entity) GetTransform() *TransForm {
	return &TransForm{}
}
func (e *Entity) Move(in *Input) {
	turnSpeed := e.EntityInfo.Ability.TSPD
	moveSpeed := e.EntityInfo.Ability.SPD
	moveValue := in.V_Movement
	turnValue := in.H_Movement
	e.Room.World.Move(e.EntityInfo.Uuid, float64(moveValue*moveSpeed), float64(turnValue*turnSpeed))
}

/*
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

*/
