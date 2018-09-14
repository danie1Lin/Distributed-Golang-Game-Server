package ode

// #include <ode/ode.h>
import "C"

import (
	"unsafe"
)

// World constants
const (
	WorldStepThreadCountUnlimited = C.dWORLDSTEP_THREADCOUNT_UNLIMITED
)

var (
	worldData = map[World]interface{}{}
)

// World represents a simulation world.
type World uintptr

// NewWorld returns a new World instance.
func NewWorld() World {
	return cToWorld(C.dWorldCreate())
}

func cToWorld(c C.dWorldID) World {
	return World(unsafe.Pointer(c))
}

func (w World) c() C.dWorldID {
	return C.dWorldID(unsafe.Pointer(w))
}

// Destroy destroys the world.
func (w World) Destroy() {
	delete(worldData, w)
	C.dWorldDestroy(w.c())
}

// SetData associates user-specified data with the world.
func (w World) SetData(data interface{}) {
	worldData[w] = data
}

// Data returns the user-specified data associated with the world.
func (w World) Data() interface{} {
	return worldData[w]
}

// NewBody returns a new Body instance.
func (w World) NewBody() Body {
	return cToBody(C.dBodyCreate(w.c()))
}

// SetGravity sets the gravity vector.
func (w World) SetGravity(grav Vector3) {
	C.dWorldSetGravity(w.c(), C.dReal(grav[0]), C.dReal(grav[1]), C.dReal(grav[2]))
}

// Gravity returns the gravity vector.
func (w World) Gravity() Vector3 {
	grav := NewVector3()
	C.dWorldGetGravity(w.c(), (*C.dReal)(&grav[0]))
	return grav
}

// SetERP sets the error reduction parameter.
func (w World) SetERP(erp float64) {
	C.dWorldSetERP(w.c(), C.dReal(erp))
}

// ERP returns the error reduction parameter.
func (w World) ERP() float64 {
	return float64(C.dWorldGetERP(w.c()))
}

// SetCFM sets the constraint force mixing value.
func (w World) SetCFM(cfm float64) {
	C.dWorldSetCFM(w.c(), C.dReal(cfm))
}

// CFM returns the constraint force mixing value.
func (w World) CFM() float64 {
	return float64(C.dWorldGetCFM(w.c()))
}

// SetStepIslandsProcessingMaxThreadCount sets the maximum number of threads to
// use for island stepping.
func (w World) SetStepIslandsProcessingMaxThreadCount(count int) {
	C.dWorldSetStepIslandsProcessingMaxThreadCount(w.c(), C.unsigned(count))
}

// StepIslandsProcessingMaxThreadCount returns the maximum number of threads to
// use for island stepping.
func (w World) StepIslandsProcessingMaxThreadCount() int {
	return int(C.dWorldGetStepIslandsProcessingMaxThreadCount(w.c()))
}

// UseSharedWorkingMemory enables sharing working memory with another world,
// and returns whether the operation succeeded.
func (w World) UseSharedWorkingMemory(from World) bool {
	return C.dWorldUseSharedWorkingMemory(w.c(), from.c()) != 0
}

// CleanupWorkingMemory cleans up the world's working memory.
func (w World) CleanupWorkingMemory() {
	C.dWorldCleanupWorkingMemory(w.c())
}

// Step executes a simulation step, and returns whether the operation
// succeeded.
func (w World) Step(stepSize float64) bool {
	return C.dWorldStep(w.c(), C.dReal(stepSize)) != 0
}

// QuickStep executes a simulation quick step, and returns whether the
// operation succeeded.
func (w World) QuickStep(stepSize float64) bool {
	return C.dWorldQuickStep(w.c(), C.dReal(stepSize)) != 0
}

// ImpulseToForce converts an impulse to a force over a step duration.
func (w World) ImpulseToForce(stepSize float64, impulse Vector3) Vector3 {
	force := NewVector3()
	C.dWorldImpulseToForce(w.c(), C.dReal(stepSize),
		C.dReal(impulse[0]), C.dReal(impulse[1]), C.dReal(impulse[2]),
		(*C.dReal)(&force[0]))
	return force
}

// SetQuickStepNumIterations sets the number of iterations to execute during a
// quick step.
func (w World) SetQuickStepNumIterations(num int) {
	C.dWorldSetQuickStepNumIterations(w.c(), C.int(num))
}

// QuickStepNumIterations returns the number of iterations to execute during a
// quick step.
func (w World) QuickStepNumIterations() int {
	return int(C.dWorldGetQuickStepNumIterations(w.c()))
}

// SetQuickStepW sets the over-relaxation parameter.
func (w World) SetQuickStepW(overRelaxation float64) {
	C.dWorldSetQuickStepW(w.c(), C.dReal(overRelaxation))
}

// QuickStepW returns the over-relaxation parameter.
func (w World) QuickStepW() float64 {
	return float64(C.dWorldGetQuickStepW(w.c()))
}

// SetContactMaxCorrectingVelocity sets the maximum correcting velocity that
// contacts are allowed to generate.
func (w World) SetContactMaxCorrectingVelocity(overRelaxation float64) {
	C.dWorldSetContactMaxCorrectingVel(w.c(), C.dReal(overRelaxation))
}

// ContactMaxCorrectingVelocity returns the maximum correcting velocity that
// contacts are allowed to generate.
func (w World) ContactMaxCorrectingVelocity() float64 {
	return float64(C.dWorldGetContactMaxCorrectingVel(w.c()))
}

// SetContactSurfaceLayer sets the depth of the surface layer around all
// geometry objects.
func (w World) SetContactSurfaceLayer(depth float64) {
	C.dWorldSetContactSurfaceLayer(w.c(), C.dReal(depth))
}

// ContactSurfaceLayer returns the depth of the surface layer around all
// geometry objects.
func (w World) ContactSurfaceLayer() float64 {
	return float64(C.dWorldGetContactSurfaceLayer(w.c()))
}

// SetAutoDisableLinearThreshold sets the auto disable linear average threshold.
func (w World) SetAutoDisableLinearThreshold(linearThreshold float64) {
	C.dWorldSetAutoDisableLinearThreshold(w.c(), C.dReal(linearThreshold))
}

// AutoDisableLinearThreshold returns the auto disable linear average threshold.
func (w World) AutoDisableLinearThreshold() float64 {
	return float64(C.dWorldGetAutoDisableLinearThreshold(w.c()))
}

// SetAutoDisableAngularThreshold sets the auto disable angular average threshold.
func (w World) SetAutoDisableAngularThreshold(angularThreshold float64) {
	C.dWorldSetAutoDisableAngularThreshold(w.c(), C.dReal(angularThreshold))
}

// AutoDisableAngularThreshold returns the auto disable angular average threshold.
func (w World) AutoDisableAngularThreshold() float64 {
	return float64(C.dWorldGetAutoDisableAngularThreshold(w.c()))
}

// SetAutoAutoDisableAverageSamplesCount sets auto disable average sample count.
func (w World) SetAutoAutoDisableAverageSamplesCount(averageSamplesCount bool) {
	C.dWorldSetAutoDisableAverageSamplesCount(w.c(), C.unsigned(btoi(averageSamplesCount)))
}

// AutoDisableAverageSamplesCount returns the auto disable sample count.
func (w World) AutoDisableAverageSamplesCount() bool {
	return C.dWorldGetAutoDisableAverageSamplesCount(w.c()) != 0
}

// SetAutoDisableSteps sets the number of auto disable steps.
func (w World) SetAutoDisableSteps(steps int) {
	C.dWorldSetAutoDisableSteps(w.c(), C.int(steps))
}

// AutoDisableSteps returns the number of auto disable steps.
func (w World) AutoDisableSteps() int {
	return int(C.dWorldGetAutoDisableSteps(w.c()))
}

// SetAutoDisableTime sets the auto disable time.
func (w World) SetAutoDisableTime(time float64) {
	C.dWorldSetAutoDisableTime(w.c(), C.dReal(time))
}

// AutoDisableTime returns the auto disable time.
func (w World) AutoDisableTime() float64 {
	return float64(C.dWorldGetAutoDisableTime(w.c()))
}

// SetAutoDisable sets wether the body will be auto disabled.
func (w World) SetAutoDisable(doAutoDisable bool) {
	C.dWorldSetAutoDisableFlag(w.c(), C.int(btoi(doAutoDisable)))
}

// AutoDisable returns whether the body will be auto disabled.
func (w World) AutoDisable() bool {
	return C.dWorldGetAutoDisableFlag(w.c()) != 0
}

// SetLinearDamping sets the linear damping scale.
func (w World) SetLinearDamping(scale float64) {
	C.dWorldSetLinearDamping(w.c(), C.dReal(scale))
}

// LinearDamping returns the linear damping scale.
func (w World) LinearDamping() float64 {
	return float64(C.dWorldGetLinearDamping(w.c()))
}

// SetAngularDamping sets the angular damping scale.
func (w World) SetAngularDamping(scale float64) {
	C.dWorldSetAngularDamping(w.c(), C.dReal(scale))
}

// AngularDamping returns the angular damping scale.
func (w World) AngularDamping() float64 {
	return float64(C.dWorldGetAngularDamping(w.c()))
}

// SetLinearDampingThreshold sets the linear damping threshold.
func (w World) SetLinearDampingThreshold(threshold float64) {
	C.dWorldSetLinearDampingThreshold(w.c(), C.dReal(threshold))
}

// LinearDampingThreshold returns the linear damping threshold.
func (w World) LinearDampingThreshold() float64 {
	return float64(C.dWorldGetLinearDampingThreshold(w.c()))
}

// SetAngularDampingThreshold sets the angular damping threshold.
func (w World) SetAngularDampingThreshold(threshold float64) {
	C.dWorldSetAngularDampingThreshold(w.c(), C.dReal(threshold))
}

// AngularDampingThreshold returns the angular damping threshold.
func (w World) AngularDampingThreshold() float64 {
	return float64(C.dWorldGetAngularDampingThreshold(w.c()))
}

// SetMaxAngularSpeed sets the maximum angular speed.
func (w World) SetMaxAngularSpeed(maxSpeed float64) {
	C.dWorldSetMaxAngularSpeed(w.c(), C.dReal(maxSpeed))
}

// MaxAngularSpeed returns the maximum angular speed.
func (w World) MaxAngularSpeed() float64 {
	return float64(C.dWorldGetMaxAngularSpeed(w.c()))
}

// NewBallJoint returns a new BallJoint instance
func (w World) NewBallJoint(group JointGroup) BallJoint {
	return cToJoint(C.dJointCreateBall(w.c(), group.c())).(BallJoint)
}

// NewHingeJoint returns a new HingeJoint instance
func (w World) NewHingeJoint(group JointGroup) HingeJoint {
	return cToJoint(C.dJointCreateHinge(w.c(), group.c())).(HingeJoint)
}

// NewSliderJoint returns a new SliderJoint instance
func (w World) NewSliderJoint(group JointGroup) SliderJoint {
	return cToJoint(C.dJointCreateSlider(w.c(), group.c())).(SliderJoint)
}

// NewContactJoint returns a new ContactJoint instance
func (w World) NewContactJoint(group JointGroup, contact *Contact) ContactJoint {
	c := &C.dContact{}
	contact.toC(c)
	return cToJoint(C.dJointCreateContact(w.c(), group.c(), c)).(ContactJoint)
}

// NewUniversalJoint returns a new UniversalJoint instance
func (w World) NewUniversalJoint(group JointGroup) UniversalJoint {
	return cToJoint(C.dJointCreateUniversal(w.c(), group.c())).(UniversalJoint)
}

// NewHinge2Joint returns a new Hinge2Joint instance
func (w World) NewHinge2Joint(group JointGroup) Hinge2Joint {
	return cToJoint(C.dJointCreateHinge2(w.c(), group.c())).(Hinge2Joint)
}

// NewFixedJoint returns a new FixedJoint instance
func (w World) NewFixedJoint(group JointGroup) FixedJoint {
	return cToJoint(C.dJointCreateFixed(w.c(), group.c())).(FixedJoint)
}

// NewNullJoint returns a new NullJoint instance
func (w World) NewNullJoint(group JointGroup) NullJoint {
	return cToJoint(C.dJointCreateNull(w.c(), group.c())).(NullJoint)
}

// NewAMotorJoint returns a new AMotorJoint instance
func (w World) NewAMotorJoint(group JointGroup) AMotorJoint {
	return cToJoint(C.dJointCreateAMotor(w.c(), group.c())).(AMotorJoint)
}

// NewLMotorJoint returns a new LMotorJoint instance
func (w World) NewLMotorJoint(group JointGroup) LMotorJoint {
	return cToJoint(C.dJointCreateLMotor(w.c(), group.c())).(LMotorJoint)
}

// NewPlane2DJoint returns a new Plane2DJoint instance
func (w World) NewPlane2DJoint(group JointGroup) Plane2DJoint {
	return cToJoint(C.dJointCreatePlane2D(w.c(), group.c())).(Plane2DJoint)
}

// NewPRJoint returns a new PRJoint instance
func (w World) NewPRJoint(group JointGroup) PRJoint {
	return cToJoint(C.dJointCreatePR(w.c(), group.c())).(PRJoint)
}

// NewPUJoint returns a new PUJoint instance
func (w World) NewPUJoint(group JointGroup) PUJoint {
	return cToJoint(C.dJointCreatePU(w.c(), group.c())).(PUJoint)
}

// NewPistonJoint returns a new PistonJoint instance
func (w World) NewPistonJoint(group JointGroup) PistonJoint {
	return cToJoint(C.dJointCreatePiston(w.c(), group.c())).(PistonJoint)
}

// NewDBallJoint returns a new DBallJoint instance
func (w World) NewDBallJoint(group JointGroup) DBallJoint {
	return cToJoint(C.dJointCreateDBall(w.c(), group.c())).(DBallJoint)
}

// NewDHingeJoint returns a new DHingeJoint instance
func (w World) NewDHingeJoint(group JointGroup) DHingeJoint {
	return cToJoint(C.dJointCreateDHinge(w.c(), group.c())).(DHingeJoint)
}

// NewTransmissionJoint returns a new TransmissionJoint instance
func (w World) NewTransmissionJoint(group JointGroup) TransmissionJoint {
	return cToJoint(C.dJointCreateTransmission(w.c(), group.c())).(TransmissionJoint)
}
