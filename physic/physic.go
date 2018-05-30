package physic

// #include <ode/ode.h>
import "C"

import (
	. "github.com/daniel840829/gameServer/data"
	. "github.com/daniel840829/gameServer/msg"
	"github.com/gazed/vu/math/lin"
	"github.com/golang/protobuf/proto"
	"github.com/ianremmler/ode"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
)

type World struct {
	sync.RWMutex
	World ode.World
	Space ode.Space
	CtGrp ode.JointGroup
	Bodys Bodys
	Cb    func(data interface{}, obj1, obj2 ode.Geom)
}

func (w *World) Init(roomId int64) {
	w.World = ode.NewWorld()
	w.Space = ode.NilSpace().NewHashSpace()
	w.Space.NewPlane(ode.V4(0, 0, 1, 0))
	w.World.SetGravity(ode.V3(0, 0, -0.5))
	ctGrp := ode.NewJointGroup(100)
	w.World.SetData(ctGrp)
	w.Cb = GetCollideHandler(w)
	w.CtGrp = ctGrp
}

func GetCollideHandler(w *World) func(data interface{}, obj1, obj2 ode.Geom) {
	return func(data interface{}, obj1, obj2 ode.Geom) {
		contact := ode.NewContact()
		body1, body2 := obj1.Body(), obj2.Body()

		if body1 != 0 && body2 != 0 && body1.Connected(body2) {
			return
		}
		world := body1.World()
		ctGrp := world.Data().(ode.JointGroup)
		//log.Debug("Body1:", body1.Data(), "\n\rBody2:", body2.Data())
		contact.Surface.Mode = 0
		contact.Surface.Mu = 0.1
		contact.Surface.Mu2 = 0
		cts := obj1.Collide(obj2, 1, 0)
		if len(cts) > 0 {
			contact.Geom = cts[0]
			ct := world.NewContactJoint(ctGrp, contact)
			ct.Attach(body1, body2)
		}
	}
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
		v := value.(ode.Body)
		p := V3_OdeToMsg(v.Position())
		q := Q_OdeToMsg(v.Quaternion())
		t := &TransForm{proto.Clone(p).(*Vector3), proto.Clone(q).(*Quaternion)}
		pos.PosMap[k] = t
		return true
	}
	w.Lock()
	w.Bodys.Range(f)
	w.Unlock()
	return
}
func (w *World) AddBody(id int64, body ode.Body) {
	log.Debug("{physic}[Addbody] Id:", id)
	w.Bodys.Store(id, body)
}

func (w *World) CreateEnitity(objName string, id int64, pos lin.V3, rot lin.Q) {

	obj := ObjData[objName]
	body := w.World.NewBody()
	body.SetPosition(V3_LinToOde(&pos))
	body.SetQuaternion(Q_LinToOde(&rot))
	body.SetData(BodyData{Uuid: id, Name: objName, Type: obj.Type})
	mass := ode.NewMass()
	switch obj.Shape {
	case "Box":
		mass.SetBox(1, ode.NewVector3(obj.Lens[0], obj.Lens[1], obj.Lens[2]))
		box := w.Space.NewBox(ode.NewVector3(obj.Lens[0], obj.Lens[1], obj.Lens[2]))
		box.SetBody(body)
	case "Capsule":
		mass.SetCapsule(1, 1, obj.Lens[0], obj.Lens[1])
		capsule := w.Space.NewCapsule(obj.Lens[0], obj.Lens[1])
		capsule.SetBody(body)
	}
	w.AddBody(id, body)

}

func (w *World) SetTranform(id int64, t *TransForm) {
	q := t.Rotation
	v3 := t.Position
	body, ok := w.Bodys.Get(id)
	if !ok {
		log.Warn("{physic}[SetTranform]Not Found:", id)
	}
	body.SetPosition(V3_MsgToOde(v3))
	body.SetQuaternion(Q_MsgToOde(q))
}

func (w *World) Move(id int64, v float64, omega float64) {
	body, ok := w.Bodys.Get(id)
	if !ok {
		log.Warn("{physic}[Move] No such body ", id)
	}
	w.Lock()
	body.SetAngularVelocity(ode.NewVector3(0, 0, -50.0*omega))
	body.SetLinearVelocity(body.VectorToWorld(ode.NewVector3(0, 50.0*v, 0)))
	log.Debug("{physic}[Move] vel:", body.LinearVelocity(), "ang vel:", body.AngularVel())
	w.Unlock()
	//log.Debug("{physic}[Move] position: ", body.Position(), " Rotation: ", body.Quaternion())
}

type Bodys struct {
	sync.Map
}

func (b *Bodys) Get(id int64) (ode.Body, bool) {
	v, ok := b.Load(id)
	body := v.(ode.Body)
	return body, ok
}

type BodyData struct {
	Uuid int64
	Type string
	Name string
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
	return ode.NewQuaternion(x, y, z, w)
}

func V3_LinToOde(v3 *lin.V3) ode.Vector3 {
	x, y, z := v3.GetS()
	return ode.NewVector3(x, y, z)
}

func Q_OdeToLin(q ode.Quaternion) *lin.Q {
	x := q[0]
	y := q[1]
	z := q[2]
	w := q[3]
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
	msg = &Quaternion{q[0], q[1], q[2], q[3]}
	return
}

func V3_MsgToOde(msg *Vector3) (v3 ode.Vector3) {
	v3 = ode.NewVector3(msg.X, msg.Y, msg.Z)
	return
}
func Q_MsgToOde(msg *Quaternion) (q ode.Quaternion) {
	q = ode.NewQuaternion(msg.X, msg.Y, msg.Z, msg.W)
	return
}
