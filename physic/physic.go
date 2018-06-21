package physic

// #include <ode/ode.h>
import "C"

import (
	"github.com/daniel840829/gameServer/data"
	. "github.com/daniel840829/gameServer/msg"
	"github.com/gazed/vu/math/lin"
	"github.com/golang/protobuf/proto"
	"github.com/daniel840829/ode"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"sync"
)

// ObjCategory
const (
	Player_Bit = iota
	Enemy_Bit
	Skill_Bit
	Terrain_Bit
	Specified_Bit
	AOE_Bit
	TREASURE_Bit
	TeamA_Bit
	TeamB_Bit
	TeamC_Bit
	TeamD_Bit
)

const (
	No_Team_Mode = iota
	Team_Mode
)

func SetBits(b int, bits ...uint) int {
	for _, bit := range bits {
		b |= 1 << bit
	}
	return b
}

func SetBitExcept(b int, bits ...uint) int {
	for _, bit := range bits {
		b &= ^(1 << bit)
	}
	return b
}

func SetAllBits() int {
	return 65535
}

type World struct {
	sync.RWMutex
	World ode.World
	Space ode.Space
	Floor ode.Plane
	CtGrp ode.JointGroup
	Objs  Objs
	Cb    func(data interface{}, obj1, obj2 ode.Geom)
}

type WorldData struct {
	Cb    func(data interface{}, obj1, obj2 ode.Geom)
	CtGrp ode.JointGroup
}

func (w *World) Init(roomId int64) {
	w.World = ode.NewWorld()
	ctGrp := ode.NewJointGroup(100000)
	var cb CollideCallback
	cb = GetCollideHandler(w)
	w.Cb = cb
	w.World.SetData(&WorldData{Cb: GetCollideHandler(w), CtGrp: ctGrp})
	w.Space = ode.NilSpace().NewHashSpace()
	w.CreateTerrain()
	w.Space.SetSublevel(0)
	w.World.SetGravity(ode.V3(0, 0, -9.8))
	w.CtGrp = ctGrp
}

func (w *World) Destroy() {
	w.Lock()
	w.World.Destroy()
	w.CtGrp.Destroy()
	w.Cb = nil
	w.Unlock()
	w = nil
}

func (w *World) CreateTerrain() {
	plane := w.Space.NewPlane(ode.V4(0, 0, 1, 0))
	plane.SetCategoryBits(SetBits(0, Terrain_Bit))
	plane.SetCollideBits(SetBitExcept(SetAllBits(), AOE_Bit))
	obj := &Obj{
		CGeom:       plane,
		Data:        &ObjData{Name: "Floor", Type: "Terrain"},
		CollideObjs: make(map[*Obj]int64),
	}
	obj.AddGeom(plane)
	wallVector := [][]float64{
		{1, 0, 0, -40},
		{-1, 0, 0, -40},
		{0, 1, 0, -40},
		{0, -1, 0, -40},
	}
	for i := 0; i < 4; i++ {
		v4 := wallVector[i]
		wall := w.Space.NewPlane(ode.V4(v4...))
		wall.SetCategoryBits(SetBits(0, Terrain_Bit))
		wall.SetCollideBits(SetBitExcept(SetAllBits(), AOE_Bit))
		obj := &Obj{
			CGeom:       wall,
			Data:        &ObjData{Name: "Wall" + strconv.Itoa(i), Type: "Terrain"},
			CollideObjs: make(map[*Obj]int64),
		}
		obj.AddGeom(wall)
	}

}

type CollideCallback func(data interface{}, obj1, obj2 ode.Geom)

func GetCollideHandler(w *World) func(data interface{}, obj1, obj2 ode.Geom) {
	var cb func(data interface{}, obj1, obj2 ode.Geom)
	cb = func(data interface{}, obj1, obj2 ode.Geom) {
		body1, body2 := obj1.Body(), obj2.Body()
		var world ode.World
		//log.Debug("GetCollideHandler", body1, " and ", body2)
		if body1 == 0 && body2 == 0 {
			return
		} else if body1 == 0 {
			world = body2.World()
		} else {
			world = body1.World()
		}
		if (obj1.IsSpace()) || (obj2.IsSpace()) {
			spaceCallback := world.Data().(*WorldData).Cb
			obj1.Collide2(obj2, data, spaceCallback)
			// if need to traverses through all spaces and sub-spaces
			if obj1.IsSpace() {
				obj1.ToSpace().Collide(data, spaceCallback)
			}
			if obj2.IsSpace() {
				obj2.ToSpace().Collide(data, spaceCallback)
			}
		} else {
			cat1, cat2 := obj1.CategoryBits(), obj2.CategoryBits()
			col1, col2 := obj1.CollideBits(), obj2.CollideBits()
			if ((cat1 & col2) | (cat2 & col1)) != 0 {
				if cat1 == SetBits(0, AOE_Bit) || cat2 == SetBits(0, AOE_Bit) {
					cts := obj1.Collide(obj2, 1, 0)
					if len(cts) > 0 {
						d1, d2 := obj1.Data().(*Obj), obj2.Data().(*Obj)
						if d1 == d2 {
							return
						} else {
							var p ode.Vector3
							if cat1&SetBits(0, AOE_Bit) != 0 {
								p = body2.Position()
								d1.InAOE(d2.GetData().Uuid, p)
							} else {
								p = body1.Position()
								d2.InAOE(d1.GetData().Uuid, p)
							}
						}
					}
				} else {
					ctGrp := world.Data().(*WorldData).CtGrp
					//log.Debug("Body1:", body1.Data(), "\n\rBody2:", body2.Data())
					contact := ode.NewContact()
					contact.Surface.Mode = 0
					contact.Surface.Mu = 0.1
					contact.Surface.Mu2 = 0
					cts := obj1.Collide(obj2, 1, 0)
					if len(cts) > 0 {
						d1, d2 := obj1.Data().(*Obj), obj2.Data().(*Obj)
						if d1 == d2 {
							log.Info("CollideHandle right")
						} else {
							d1.Collide(d2)
							d2.Collide(d1)
							contact.Geom = cts[0]
							ct := world.NewContactJoint(ctGrp, contact)
							ct.Attach(body1, body2)

						}
					}
				}
			}
		}
	}
	return cb
}
func (w *World) PhysicUpdate() {
	w.Lock()
	w.Space.Collide(0, w.Cb)
	w.World.Step(0.01)
	w.CtGrp.Empty()
	w.Unlock()
}

func (w *World) GetAllTransform() (pos *Position) {
	pos = &Position{}
	pos.PosMap = make(map[int64]*TransForm)
	f := func(key interface{}, value interface{}) bool {
		k := key.(int64)
		v := value.(*Obj).CBody
		p := V3_OdeToMsg(v.Position())
		q := Q_OdeToMsg(v.Quaternion())
		t := &TransForm{proto.Clone(p).(*Vector3), proto.Clone(q).(*Quaternion)}
		pos.PosMap[k] = t
		return true
	}
	w.Lock()
	w.Objs.Range(f)
	w.Unlock()
	return
}

func (w *World) GetTransform(id int64) (p *lin.V3, q *lin.Q) {
	w.Lock()
	obj, ok := w.Objs.Get(id)
	body := obj.CBody
	if !ok {
		w.Unlock()
		log.Warn("[physic]{GetTransform}id is missed", id)
		return
	}
	p = V3_OdeToLin(body.Position())
	q = Q_OdeToLin(body.Quaternion())
	w.Unlock()
	return
}

func (w *World) AddObj(id int64, obj *Obj) {
	log.Debug("{physic}[Addbody] Id:", id)
	w.Objs.Store(id, obj)
}

func (w *World) DeleteObj(id int64) {
	w.Objs.Delete(id)
}
func (w *World) CreateEntity(objName string, id int64, pos lin.V3, rot lin.Q) {

	obj := data.ObjData[objName]
	body := w.World.NewBody()
	mass := ode.NewMass()
	switch obj.Shape {
	case "Box":
		mass.SetBoxTotal(obj.Mass, ode.NewVector3(obj.Lens[0], obj.Lens[1], obj.Lens[2]))
		box := w.Space.NewBox(ode.NewVector3(obj.Lens[0], obj.Lens[1], obj.Lens[2]))
		box.SetCategoryBits(SetBits(0, Player_Bit))
		box.SetCollideBits(SetBitExcept(SetAllBits(), AOE_Bit))
		box.SetBody(body)
		object := &Obj{
			CBody:       body,
			CGeom:       box,
			Data:        &ObjData{Uuid: id, Name: objName, Type: obj.Type},
			CollideObjs: make(map[*Obj]int64),
		}
		object.CreateAOE(w.Space, 100.0)
		w.Lock()
		object.AddGeom(box)
		object.AddBody(body)
		w.AddObj(id, object)
		//joint := w.World.NewPlane2DJoint(ode.JointGroup(0))
		//joint.Attach(body, ode.Body(0))
	case "Capsule":
		mass.SetCapsuleTotal(obj.Mass, 1, obj.Lens[0], obj.Lens[1])
		capsule := w.Space.NewCapsule(obj.Lens[0], obj.Lens[1])
		capsule.SetCategoryBits(SetBits(0, Skill_Bit))
		capsule.SetCollideBits(SetBits(0, Terrain_Bit, Player_Bit, Enemy_Bit))
		capsule.SetBody(body)
		object := &Obj{
			CBody:       body,
			CGeom:       capsule,
			Data:        &ObjData{Uuid: id, Name: objName, Type: obj.Type},
			CollideObjs: make(map[*Obj]int64),
		}
		w.Lock()
		object.AddGeom(capsule)
		object.AddBody(body)
		w.AddObj(id, object)
	default:
		log.Debug("[World]{CreateEntity} No ", objName)
		return
	}
	body.SetPosition(V3_LinToOde(&pos))
	body.SetQuaternion(Q_LinToOde(&rot))
	w.Unlock()
}

func (w *World) SetTranform(id int64, t *TransForm) {
	q := t.Rotation
	v3 := t.Position
	obj, ok := w.Objs.Get(id)
	body := obj.CBody
	if !ok {
		log.Warn("{physic}[SetTranform]Not Found:", id)
	}
	body.SetPosition(V3_MsgToOde(v3))
	body.SetQuaternion(Q_MsgToOde(q))
}

func (w *World) Move(id int64, v float64, omega float64) {
	w.Lock()
	obj, ok := w.Objs.Get(id)
	if !ok {
		log.Warn("{physic}[Move] No such body ", id)
	}
	body := obj.CBody
	body.SetAngularVelocity(ode.NewVector3(0, 0, -3*omega))
	body.SetLinearVelocity(body.VectorToWorld(ode.NewVector3(0, 10*v, body.LinearVelocity()[2])))
	//log.Debug("{physic}[Move] vel:", body.LinearVelocity(), "ang vel:", body.AngularVel())
	w.Unlock()
	//log.Debug("{physic}[Move] position: ", body.Position(), " Rotation: ", body.Quaternion())
}

type Objs struct {
	sync.Map
}

func (g *Objs) Get(id int64) (*Obj, bool) {
	v, ok := g.Load(id)
	obj := v.(*Obj)
	return obj, ok
}

type ObjData struct {
	Uuid         int64
	Type         string
	Name         string
	CollideTimes int64
}

type Obj struct {
	sync.Mutex
	CBody       ode.Body
	CGeom       ode.Geom
	Space       ode.Space
	AOE         ode.Geom
	OtherBodys  []ode.Body
	OtherGeoms  []ode.Geom
	Data        *ObjData
	CollideObjs map[*Obj]int64
	AOEObjs     map[int64]ode.Vector3
}

func (obj *Obj) CreateAOE(space ode.Space, radius float64) {
	obj.AOEObjs = make(map[int64]ode.Vector3)
	s := space.NewSphere(radius)
	s.SetCategoryBits(SetBits(0, AOE_Bit))
	s.SetCollideBits(SetBits(0, Player_Bit))
	s.SetData(obj)
	obj.AOE = s
}

func (obj *Obj) SyncAOEPos() {
	obj.AOE.SetPosition(obj.CBody.Position())
	obj.AOE.SetQuaternion(obj.CBody.Quaternion())
}

func (obj *Obj) InAOE(entityId int64, p ode.Vector3) {
	if _, ok := obj.AOEObjs[entityId]; ok {
		return
	}
	obj.AOEObjs[entityId] = p
}

func (obj *Obj) ClearAOE() {
	obj.AOEObjs = make(map[int64]ode.Vector3)
}
func (obj *Obj) AddGeom(geom ode.Geom) {
	obj.OtherGeoms = append(obj.OtherGeoms, geom)
	geom.SetData(obj)
}

func (obj *Obj) AddBody(body ode.Body) {
	obj.OtherBodys = append(obj.OtherBodys, body)
}

func (obj *Obj) Collide(obj2 *Obj) {
	obj.Data.CollideTimes += 1
	obj.CollideObjs[obj2] += 1
}

func (obj *Obj) ResetCollide() {
	obj.Data.CollideTimes = 0
	obj.CollideObjs = make(map[*Obj]int64)
}

func (obj *Obj) LoopCollideObj(f func(obj *Obj, times int64) bool) {
	for obj, times := range obj.CollideObjs {
		if !f(obj, times) {
			break
		}
	}
}
func (obj *Obj) GetData() *ObjData {
	return obj.Data
}

func (obj *Obj) SetData(objData ObjData) {
	obj.Data = &objData
}

func (obj *Obj) Destroy() {
	obj.Lock()
	log.Debug("[obj]{Destroy} Space:", obj.Space)
	if obj.Space != nil {
		obj.Space.Destroy()
	} else {
		obj.CGeom.Destroy()
		obj.CBody.Destroy()
	}
	obj.Unlock()
}

var (
	roomIdMapCtGrp map[int64]ode.JointGroup
	tempCtGrp      ode.JointGroup
)

func init() {
	log.Debug("ode Init")
	ode.Init(0, ode.AllAFlag)
}

func EulerToQuaternion(yaw, pitch, roll float64) *lin.Q {
	p := pitch * math.Pi / 180 / 2.0
	ya := yaw * math.Pi / 180 / 2.0
	r := roll * math.Pi / 180 / 2.0

	sinp := math.Sin(p)
	siny := math.Sin(ya)
	sinr := math.Sin(r)
	cosp := math.Cos(p)
	cosy := math.Cos(ya)
	cosr := math.Cos(r)

	x := sinr*cosp*cosy - cosr*sinp*siny
	y := cosr*sinp*cosy + sinr*cosp*siny
	z := cosr*cosp*siny - sinr*sinp*cosy
	w := cosr*cosp*cosy + sinr*sinp*siny
	q := lin.NewQ()
	q.SetS(x, y, z, w)

	return q
}

func Q_LinToOde(q *lin.Q) ode.Quaternion {
	x, y, z, w := q.GetS()
	return ode.NewQuaternion(w, x, y, z)
}

func V3_LinToOde(v3 *lin.V3) ode.Vector3 {
	x, y, z := v3.GetS()
	return ode.NewVector3(x, y, z)
}

func Q_OdeToLin(q ode.Quaternion) *lin.Q {
	x := q[1]
	y := q[2]
	z := q[3]
	w := q[0]
	result := lin.NewQ()
	result.SetS(x, y, z, w)
	return result
}
func V3_OdeToLin(q ode.Vector3) *lin.V3 {
	x := q[0]
	y := q[1]
	z := q[2]
	result := lin.NewV3()
	result.SetS(x, y, z)
	return result
}

func V3_OdeToMsg(v3 ode.Vector3) (msg *Vector3) {
	msg = &Vector3{v3[0], v3[1], v3[2]}
	return
}

func Q_OdeToMsg(q ode.Quaternion) (msg *Quaternion) {
	msg = &Quaternion{q[1], q[2], q[3], q[0]}
	return
}

func V3_MsgToOde(msg *Vector3) (v3 ode.Vector3) {
	v3 = ode.NewVector3(msg.X, msg.Y, msg.Z)
	return
}
func Q_MsgToOde(msg *Quaternion) (q ode.Quaternion) {
	q = ode.NewQuaternion(msg.W, msg.X, msg.Y, msg.Z)
	return
}

func Q_LinToMsg(q *lin.Q) (msg *Quaternion) {
	msg = &Quaternion{}
	msg.X, msg.Y, msg.Z, msg.W = q.GetS()
	return
}

func V3_LinToMsg(p *lin.V3) (msg *Vector3) {
	msg = &Vector3{}
	msg.X, msg.Y, msg.Z = p.GetS()
	return
}
