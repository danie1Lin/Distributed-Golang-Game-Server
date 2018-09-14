package ode

// #include <ode/ode.h>
// extern void callMovedCallback(dBodyID body);
import "C"

import (
	"unsafe"
)

var (
	bodyData       = map[Body]interface{}{}
	movedCallbacks = map[Body]MovedCallback{}
)

// Body represents a rigid body.
type Body uintptr

// MovedCallback is called when the body has moved.
type MovedCallback func(b Body)

//export movedCallback
func movedCallback(c C.dBodyID) {
	body := cToBody(c)
	if cb, ok := movedCallbacks[body]; ok {
		cb(body)
	}
}

func cToBody(c C.dBodyID) Body {
	return Body(unsafe.Pointer(c))
}

func (b Body) c() C.dBodyID {
	return C.dBodyID(unsafe.Pointer(b))
}

// Destroy destroys the body.
func (b Body) Destroy() {
	delete(bodyData, b)
	C.dBodyDestroy(b.c())
}

// SetAutoDisableLinearThreshold sets the auto disable linear average threshold.
func (b Body) SetAutoDisableLinearThreshold(thresh float64) {
	C.dBodySetAutoDisableLinearThreshold(b.c(), C.dReal(thresh))
}

// AutoDisableLinearThreshold returns the auto disable linear average threshold.
func (b Body) AutoDisableLinearThreshold() float64 {
	return float64(C.dBodyGetAutoDisableLinearThreshold(b.c()))
}

// SetAutoDisableAngularThreshold sets the auto disable angular average threshold.
func (b Body) SetAutoDisableAngularThreshold(thresh float64) {
	C.dBodySetAutoDisableAngularThreshold(b.c(), C.dReal(thresh))
}

// AutoDisableAngularThreshold returns the auto disable angular average threshold.
func (b Body) AutoDisableAngularThreshold() float64 {
	return float64(C.dBodyGetAutoDisableAngularThreshold(b.c()))
}

// SetAutoAutoDisableAverageSamplesCount sets auto disable average sample count.
func (b Body) SetAutoAutoDisableAverageSamplesCount(count int) {
	C.dBodySetAutoDisableAverageSamplesCount(b.c(), C.uint(count))
}

// AutoDisableAverageSamplesCount returns the auto disable sample count.
func (b Body) AutoDisableAverageSamplesCount() int {
	return int(C.dBodyGetAutoDisableAverageSamplesCount(b.c()))
}

// SetAutoDisableSteps sets the number of auto disable steps.
func (b Body) SetAutoDisableSteps(steps int) {
	C.dBodySetAutoDisableSteps(b.c(), C.int(steps))
}

// AutoDisableSteps returns the number of auto disable steps.
func (b Body) AutoDisableSteps() int {
	return int(C.dBodyGetAutoDisableSteps(b.c()))
}

// SetAutoDisableTime sets the auto disable time.
func (b Body) SetAutoDisableTime(time float64) {
	C.dBodySetAutoDisableTime(b.c(), C.dReal(time))
}

// AutoDisableTime returns the auto disable time.
func (b Body) AutoDisableTime() float64 {
	return float64(C.dBodyGetAutoDisableTime(b.c()))
}

// SetAutoDisable sets wether the body will be auto disabled.
func (b Body) SetAutoDisable(doAutoDisable bool) {
	C.dBodySetAutoDisableFlag(b.c(), C.int(btoi(doAutoDisable)))
}

// AutoDisable returns whether the body will be auto disabled.
func (b Body) AutoDisable() bool {
	return C.dBodyGetAutoDisableFlag(b.c()) != 0
}

// SetAutoDisableDefaults sets auto disable settings to default defaults.
func (b Body) SetAutoDisableDefaults() {
	C.dBodySetAutoDisableDefaults(b.c())
}

// World returns the world which contains the body.
func (b Body) World() World {
	return cToWorld(C.dBodyGetWorld(b.c()))
}

// SetData associates user-specified data with the body.
func (b Body) SetData(data interface{}) {
	bodyData[b] = data
}

// Data returns the user-specified data associated with the body.
func (b Body) Data() interface{} {
	return bodyData[b]
}

// SetPosition sets the position.
func (b Body) SetPosition(pos Vector3) {
	C.dBodySetPosition(b.c(), C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]))
}

// Position returns the position.
func (b Body) Position() Vector3 {
	pos := NewVector3()
	C.dBodyCopyPosition(b.c(), (*C.dReal)(&pos[0]))
	return pos
}

// SetRotation sets the orientation represented by a rotation matrix.
func (b Body) SetRotation(rot Matrix3) {
	C.dBodySetRotation(b.c(), (*C.dReal)(&rot[0][0]))
}

// Rotation returns the orientation represented by a rotation matrix.
func (b Body) Rotation() Matrix3 {
	rot := NewMatrix3()
	C.dBodyCopyRotation(b.c(), (*C.dReal)(&rot[0][0]))
	return rot
}

// SetQuaternion sets the orientation represented by a quaternion.
func (b Body) SetQuaternion(quat Quaternion) {
	C.dBodySetQuaternion(b.c(), (*C.dReal)(&quat[0]))
}

// Quaternion returns the orientation represented by a quaternion.
func (b Body) Quaternion() Quaternion {
	quat := NewQuaternion()
	C.dBodyCopyQuaternion(b.c(), (*C.dReal)(&quat[0]))
	return quat
}

// SetLinearVelocity sets the linear velocity.
func (b Body) SetLinearVelocity(vel Vector3) {
	C.dBodySetLinearVel(b.c(), C.dReal(vel[0]), C.dReal(vel[1]), C.dReal(vel[2]))
}

// LinearVelocity returns the linear velocity.
func (b Body) LinearVelocity() Vector3 {
	return cToVector3(C.dBodyGetLinearVel(b.c()))
}

// SetAngularVelocity sets the angular velocity.
func (b Body) SetAngularVelocity(vel Vector3) {
	C.dBodySetAngularVel(b.c(), C.dReal(vel[0]), C.dReal(vel[1]), C.dReal(vel[2]))
}

// AngularVel returns the angular velocity.
func (b Body) AngularVel() Vector3 {
	return cToVector3(C.dBodyGetAngularVel(b.c()))
}

// SetMass sets the mass.
func (b Body) SetMass(mass *Mass) {
	c := &C.dMass{}
	mass.toC(c)
	C.dBodySetMass(b.c(), c)
}

// Mass returns the mass.
func (b Body) Mass() *Mass {
	mass := NewMass()
	c := &C.dMass{}
	C.dBodyGetMass(b.c(), c)
	mass.fromC(c)
	return mass
}

// SetForce sets the force acting on the center of mass.
func (b Body) SetForce(force Vector3) {
	C.dBodySetForce(b.c(), C.dReal(force[0]), C.dReal(force[1]), C.dReal(force[2]))
}

// Force returns the force acting on the center of mass.
func (b Body) Force() Vector3 {
	return cToVector3(C.dBodyGetForce(b.c()))
}

// SetTorque sets the torque acting on the center of mass.
func (b Body) SetTorque(torque Vector3) {
	C.dBodySetTorque(b.c(), C.dReal(torque[0]), C.dReal(torque[1]), C.dReal(torque[2]))
}

// Torque returns the torque acting on the center of mass.
func (b Body) Torque() Vector3 {
	return cToVector3(C.dBodyGetTorque(b.c()))
}

// RelPointPos returns the position in world coordinates of a point in body
// coordinates.
func (b Body) RelPointPos(pt Vector3) Vector3 {
	pos := NewVector3()
	C.dBodyGetRelPointPos(b.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]), (*C.dReal)(&pos[0]))
	return pos
}

// RelPointVel returns the velocity in world coordinates of a point in body
// coordinates.
func (b Body) RelPointVel(pt Vector3) Vector3 {
	vel := NewVector3()
	C.dBodyGetRelPointVel(b.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]), (*C.dReal)(&vel[0]))
	return vel
}

// PointVel returns the velocity in world coordinates of a point in world
// coordinates.
func (b Body) PointVel(pt Vector3) Vector3 {
	vel := NewVector3()
	C.dBodyGetPointVel(b.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]), (*C.dReal)(&vel[0]))
	return vel
}

// PosRelPoint returns the position in body coordinates of a point in world
// coordinates.
func (b Body) PosRelPoint(pos Vector3) Vector3 {
	pt := NewVector3()
	C.dBodyGetPosRelPoint(b.c(), C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]), (*C.dReal)(&pt[0]))
	return pt
}

// VectorToWorld converts a vector in body coordinates to world coordinates.
func (b Body) VectorToWorld(vec Vector3) Vector3 {
	wld := NewVector3()
	C.dBodyVectorToWorld(b.c(), C.dReal(vec[0]), C.dReal(vec[1]), C.dReal(vec[2]), (*C.dReal)(&wld[0]))
	return wld
}

// VectorFromWorld converts a vector in world coordinates to body coordinates.
func (b Body) VectorFromWorld(wld Vector3) Vector3 {
	vec := NewVector3()
	C.dBodyVectorFromWorld(b.c(), C.dReal(wld[0]), C.dReal(wld[1]), C.dReal(wld[2]), (*C.dReal)(&vec[0]))
	return vec
}

// SetFiniteRotationMode sets whether finite rotation mode is used.
func (b Body) SetFiniteRotationMode(isFinite bool) {
	C.dBodySetFiniteRotationMode(b.c(), C.int(btoi(isFinite)))
}

// FiniteRotationMode returns whether finite rotation mode is used.
func (b Body) FiniteRotationMode() bool {
	return C.dBodyGetFiniteRotationMode(b.c()) != 0
}

// SetFiniteRotationAxis sets the finite rotation axis.
func (b Body) SetFiniteRotationAxis(axis Vector3) {
	C.dBodySetFiniteRotationAxis(b.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// FiniteRotationAxis returns the finite rotation axis.
func (b Body) FiniteRotationAxis() Vector3 {
	axis := NewVector3()
	C.dBodyGetFiniteRotationAxis(b.c(), (*C.dReal)(&axis[0]))
	return axis
}

// NumJoints returns the number of joints attached to the body.
func (b Body) NumJoints() int {
	return int(C.dBodyGetNumJoints(b.c()))
}

// Joint returns the joint attached to the body, specified by index.
func (b Body) Joint(index int) Joint {
	return cToJoint(C.dBodyGetJoint(b.c(), C.int(index)))
}

// SetKinematic sets whether the body is in kinematic mode.
func (b Body) SetKinematic(isKinematic bool) {
	if isKinematic {
		C.dBodySetKinematic(b.c())
	} else {
		C.dBodySetDynamic(b.c())
	}
}

// Kinematic returns whether the body is in kinematic mode.
func (b Body) Kinematic() bool {
	return bool(C.dBodyIsKinematic(b.c()) != 0)
}

// SetGravityEnabled sets whether gravity affects the body.
func (b Body) SetGravityEnabled(isGravityEnabled bool) {
	C.dBodySetGravityMode(b.c(), C.int(btoi(isGravityEnabled)))
}

// GravityEnabled returns whether gravity affects the body.
func (b Body) GravityEnabled() bool {
	return C.dBodyGetGravityMode(b.c()) != 0
}

// SetMovedCallback sets callback to call when the body has moved.
func (b Body) SetMovedCallback(cb MovedCallback) {
	if cb == nil {
		C.dBodySetMovedCallback(b.c(), (*[0]byte)(nil)) // clear callback
		delete(movedCallbacks, b)
	} else {
		movedCallbacks[b] = cb
		C.dBodySetMovedCallback(b.c(), (*[0]byte)(unsafe.Pointer(C.callMovedCallback)))
	}
}

// FirstGeom returns the first geometry associated with the body.
func (b Body) FirstGeom() Geom {
	return cToGeom(C.dBodyGetFirstGeom(b.c()))
}

// SetDampingDefaults sets damping settings to default values.
func (b Body) SetDampingDefaults() {
	C.dBodySetDampingDefaults(b.c())
}

// SetLinearDamping sets the linear damping scale.
func (b Body) SetLinearDamping(scale float64) {
	C.dBodySetLinearDamping(b.c(), C.dReal(scale))
}

// LinearDamping returns the linear damping scale.
func (b Body) LinearDamping() float64 {
	return float64(C.dBodyGetLinearDamping(b.c()))
}

// SetAngularDamping sets the angular damping scale.
func (b Body) SetAngularDamping(scale float64) {
	C.dBodySetAngularDamping(b.c(), C.dReal(scale))
}

// AngularDamping returns the angular damping scale.
func (b Body) AngularDamping() float64 {
	return float64(C.dBodyGetAngularDamping(b.c()))
}

// SetLinearDampingThreshold sets the linear damping threshold.
func (b Body) SetLinearDampingThreshold(threshold float64) {
	C.dBodySetLinearDampingThreshold(b.c(), C.dReal(threshold))
}

// LinearDampingThreshold returns the linear damping threshold.
func (b Body) LinearDampingThreshold() float64 {
	return float64(C.dBodyGetLinearDampingThreshold(b.c()))
}

// SetAngularDampingThreshold sets the angular damping threshold.
func (b Body) SetAngularDampingThreshold(threshold float64) {
	C.dBodySetAngularDampingThreshold(b.c(), C.dReal(threshold))
}

// AngularDampingThreshold returns the angular damping threshold.
func (b Body) AngularDampingThreshold() float64 {
	return float64(C.dBodyGetAngularDampingThreshold(b.c()))
}

// SetMaxAngularSpeed sets the maximum angular speed.
func (b Body) SetMaxAngularSpeed(maxSpeed float64) {
	C.dBodySetMaxAngularSpeed(b.c(), C.dReal(maxSpeed))
}

// MaxAngularSpeed returns the maximum angular speed.
func (b Body) MaxAngularSpeed() float64 {
	return float64(C.dBodyGetMaxAngularSpeed(b.c()))
}

// SetGyroModeEnabled sets whether gyroscopic mode is enabled.
func (b Body) SetGyroModeEnabled(isEnabled bool) {
	C.dBodySetGyroscopicMode(b.c(), C.int(btoi(isEnabled)))
}

// GyroModeEnabled returns whether gyroscopic mode is enabled.
func (b Body) GyroModeEnabled() bool {
	return C.dBodyGetGyroscopicMode(b.c()) != 0
}

// Connected returns whether the body is connected to the given body by a
// joint.
func (b Body) Connected(other Body) bool {
	return C.dAreConnected(b.c(), other.c()) != 0
}

// ConnectedExcluding returns whether the body is connected to the given body
// by a joint, except for joints of the given class.
func (b Body) ConnectedExcluding(other Body, jointType int) bool {
	return C.dAreConnectedExcluding(b.c(), other.c(), C.int(jointType)) != 0
}

// AddForce adds a force in world coordinates at the center of mass.
func (b Body) AddForce(force Vector3) {
	C.dBodyAddForce(b.c(), C.dReal(force[0]), C.dReal(force[1]), C.dReal(force[2]))
}

// AddRelForce adds a force in body coordinates at the center of mass.
func (b Body) AddRelForce(force Vector3) {
	C.dBodyAddRelForce(b.c(), C.dReal(force[0]), C.dReal(force[1]), C.dReal(force[2]))
}

// AddTorque adds a torque in world coordinates at the center of mass.
func (b Body) AddTorque(torque Vector3) {
	C.dBodyAddTorque(b.c(), C.dReal(torque[0]), C.dReal(torque[1]), C.dReal(torque[2]))
}

// AddRelTorque adds a torque in body coordinates at the center of mass.
func (b Body) AddRelTorque(torque Vector3) {
	C.dBodyAddRelTorque(b.c(), C.dReal(torque[0]), C.dReal(torque[1]), C.dReal(torque[2]))
}

// AddForceAtPos adds a force in world coordinates at the position in world
// coordinates.
func (b Body) AddForceAtPos(force, pos Vector3) {
	C.dBodyAddForceAtPos(b.c(), C.dReal(force[0]), C.dReal(force[1]), C.dReal(force[2]),
		C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]))
}

// AddForceAtRelPos adds a force in world coordinates at the position in body
// coordinates.
func (b Body) AddForceAtRelPos(force, pos Vector3) {
	C.dBodyAddForceAtRelPos(b.c(), C.dReal(force[0]), C.dReal(force[1]), C.dReal(force[2]),
		C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]))
}

// AddRelForceAtPos adds a force in body coordinates at the position in world
// coordinates.
func (b Body) AddRelForceAtPos(force, pos Vector3) {
	C.dBodyAddRelForceAtPos(b.c(), C.dReal(force[0]), C.dReal(force[1]), C.dReal(force[2]),
		C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]))
}

// AddRelForceAtRelPos adds a force in body coordinates at the position in body coordinates.
func (b Body) AddRelForceAtRelPos(force, pos Vector3) {
	C.dBodyAddRelForceAtRelPos(b.c(), C.dReal(force[0]), C.dReal(force[1]), C.dReal(force[2]),
		C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]))
}

// SetEnabled sets whether the body is enabled.
func (b Body) SetEnabled(isEnabled bool) {
	if isEnabled {
		C.dBodyEnable(b.c())
	} else {
		C.dBodyDisable(b.c())
	}
}

// Enabled returns whether the body is enabled.
func (b Body) Enabled() bool {
	return bool(C.dBodyIsEnabled(b.c()) != 0)
}

// ConnectingJoint returns the first joint connecting the body to the specified
// body.
func (b Body) ConnectingJoint(other Body) Joint {
	return cToJoint(C.dConnectingJoint(b.c(), other.c()))
}

// ConnectingJointList returns a list of joints connecting the body to the
// specified body.
func (b Body) ConnectingJointList(other Body) []Joint {
	jointList := []Joint{}
	for i := 0; i < b.NumJoints(); i++ {
		joint := b.Joint(i)
		for j := 0; j < joint.NumBodies(); j++ {
			if body := joint.Body(j); body == other {
				jointList = append(jointList, joint)
				break
			}
		}
	}
	return jointList
}
