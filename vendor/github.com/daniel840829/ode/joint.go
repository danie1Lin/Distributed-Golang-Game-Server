package ode

// #include <ode/ode.h>
import "C"

import (
	"unsafe"
)

// Joint types
const (
	BallJointType         = C.dJointTypeBall
	HingeJointType        = C.dJointTypeHinge
	SliderJointType       = C.dJointTypeSlider
	ContactJointType      = C.dJointTypeContact
	UniversalJointType    = C.dJointTypeUniversal
	Hinge2JointType       = C.dJointTypeHinge2
	FixedJointType        = C.dJointTypeFixed
	NullJointType         = C.dJointTypeNull
	AMotorJointType       = C.dJointTypeAMotor
	LMotorJointType       = C.dJointTypeLMotor
	Plane2DJointType      = C.dJointTypePlane2D
	PRJointType           = C.dJointTypePR
	PUJointType           = C.dJointTypePU
	PistonJointType       = C.dJointTypePiston
	DBallJointType        = C.dJointTypeDBall
	DHingeJointType       = C.dJointTypeDHinge
	TransmissionJointType = C.dJointTypeTransmission
)

// Joint parameters
const (
	LoStopJtParam        = C.dParamLoStop
	HiStopJtParam        = C.dParamHiStop
	VelJtParam           = C.dParamVel
	LoVelJtParam         = C.dParamLoVel
	HiVelJtParam         = C.dParamHiVel
	FMaxJtParam          = C.dParamFMax
	FudgeFactorJtParam   = C.dParamFudgeFactor
	BounceJtParam        = C.dParamBounce
	CFMJtParam           = C.dParamCFM
	StopERPJtParam       = C.dParamStopERP
	StopCFMJtParam       = C.dParamStopCFM
	SuspensionERPJtParam = C.dParamSuspensionERP
	SuspensionCFMJtParam = C.dParamSuspensionCFM
	ERPJtParam           = C.dParamERP

	NumJtParams = C.dParamsInGroup

	JtParamGroup1         = C.dParamGroup1
	LoStopJtParam1        = C.dParamLoStop1
	HiStopJtParam1        = C.dParamHiStop1
	VelJtParam1           = C.dParamVel1
	LoVelJtParam1         = C.dParamLoVel1
	HiVelJtParam1         = C.dParamHiVel1
	FMaxJtParam1          = C.dParamFMax1
	FudgeFactorJtParam1   = C.dParamFudgeFactor1
	BounceJtParam1        = C.dParamBounce1
	CFMJtParam1           = C.dParamCFM1
	StopERPJtParam1       = C.dParamStopERP1
	StopCFMJtParam1       = C.dParamStopCFM1
	SuspensionERPJtParam1 = C.dParamSuspensionERP1
	SuspensionCFMJtParam1 = C.dParamSuspensionCFM1
	ERPJtParam1           = C.dParamERP1

	JtParamGroup2         = C.dParamGroup2
	LoStopJtParam2        = C.dParamLoStop2
	HiStopJtParam2        = C.dParamHiStop2
	VelJtParam2           = C.dParamVel2
	LoVelJtParam2         = C.dParamLoVel2
	HiVelJtParam2         = C.dParamHiVel2
	FMaxJtParam2          = C.dParamFMax2
	FudgeFactorJtParam2   = C.dParamFudgeFactor2
	BounceJtParam2        = C.dParamBounce2
	CFMJtParam2           = C.dParamCFM2
	StopERPJtParam2       = C.dParamStopERP2
	StopCFMJtParam2       = C.dParamStopCFM2
	SuspensionERPJtParam2 = C.dParamSuspensionERP2
	SuspensionCFMJtParam2 = C.dParamSuspensionCFM2
	ERPJtParam2           = C.dParamERP2

	JtParamGroup3         = C.dParamGroup3
	LoStopJtParam3        = C.dParamLoStop3
	HiStopJtParam3        = C.dParamHiStop3
	VelJtParam3           = C.dParamVel3
	LoVelJtParam3         = C.dParamLoVel3
	HiVelJtParam3         = C.dParamHiVel3
	FMaxJtParam3          = C.dParamFMax3
	FudgeFactorJtParam3   = C.dParamFudgeFactor3
	BounceJtParam3        = C.dParamBounce3
	CFMJtParam3           = C.dParamCFM3
	StopERPJtParam3       = C.dParamStopERP3
	StopCFMJtParam3       = C.dParamStopCFM3
	SuspensionERPJtParam3 = C.dParamSuspensionERP3
	SuspensionCFMJtParam3 = C.dParamSuspensionCFM3
	ERPJtParam3           = C.dParamERP3
)

// Angular motor parameters
const (
	AMotorUser  = C.dAMotorUser
	AMotorEuler = C.dAMotorEuler
)

// Transmission parameters
const (
	TransmissionParallelAxes     = C.dTransmissionParallelAxes
	TransmissionIntersectingAxes = C.dTransmissionIntersectingAxes
	TransmissionChainDrive       = C.dTransmissionChainDrive
)

var (
	jointData = map[Joint]interface{}{}
)

// JointFeedback represents feedback forces and torques associated with a
// joint.
type JointFeedback struct {
	Force1  Vector3 // force applied to body 1
	Torque1 Vector3 // torque applied to body 1
	Force2  Vector3 // force applied to body 2
	Torque2 Vector3 // torque applied to body 2
}

func (f *JointFeedback) fromC(c *C.dJointFeedback) {
	Vector(f.Force1).fromC(&c.f1[0])
	Vector(f.Torque1).fromC(&c.t1[0])
	Vector(f.Force2).fromC(&c.f2[0])
	Vector(f.Torque2).fromC(&c.t2[0])
}

func (f *JointFeedback) toC(c *C.dJointFeedback) {
	Vector(f.Force1).toC((*C.dReal)(&c.f1[0]))
	Vector(f.Torque1).toC((*C.dReal)(&c.t1[0]))
	Vector(f.Force2).toC((*C.dReal)(&c.f2[0]))
	Vector(f.Torque2).toC((*C.dReal)(&c.t2[0]))
}

// JointGroup represents a group of joints.
type JointGroup uintptr

// NewJointGroup returns a new JointGroup instance.
func NewJointGroup(maxJoints int) JointGroup {
	return cToJointGroup(C.dJointGroupCreate(C.int(maxJoints)))
}

func cToJointGroup(c C.dJointGroupID) JointGroup {
	return JointGroup(unsafe.Pointer(c))
}

func (g JointGroup) c() C.dJointGroupID {
	return C.dJointGroupID(unsafe.Pointer(g))
}

// Destroy destroys the joint group.
func (g JointGroup) Destroy() {
	C.dJointGroupDestroy(g.c())
}

func (g JointGroup) C() C.dJointGroupID {
	return C.dJointGroupID(unsafe.Pointer(g))
}
func CToJointGroup(c C.dJointGroupID) JointGroup {
	return JointGroup(unsafe.Pointer(c))
}

// Empty removes all joints from the group.
func (g JointGroup) Empty() {
	C.dJointGroupEmpty(g.c())
}

// Joint represents a joint.
type Joint interface {
	c() C.dJointID
	Destroy()
	SetData(data interface{})
	Data() interface{}
	NumBodies() int
	Attach(body1, body2 Body)
	SetEnabled(isEnabled bool)
	Enabled() bool
	Type() int
	Body(index int) Body
	SetFeedback(f *JointFeedback)
	Feedback() *JointFeedback
}

// JointBase implements Joint, and is embedded by specific Joint types.
type JointBase uintptr

func cToJoint(c C.dJointID) Joint {
	base := JointBase(unsafe.Pointer(c))
	var j Joint
	switch int(C.dJointGetType(c)) {
	case BallJointType:
		j = BallJoint{base}
	case HingeJointType:
		j = HingeJoint{base}
	case SliderJointType:
		j = SliderJoint{base}
	case ContactJointType:
		j = ContactJoint{base}
	case UniversalJointType:
		j = UniversalJoint{base}
	case Hinge2JointType:
		j = Hinge2Joint{base}
	case FixedJointType:
		j = FixedJoint{base}
	case NullJointType:
		j = NullJoint{base}
	case AMotorJointType:
		j = AMotorJoint{base}
	case LMotorJointType:
		j = LMotorJoint{base}
	case Plane2DJointType:
		j = Plane2DJoint{base}
	case PRJointType:
		j = PRJoint{base}
	case PUJointType:
		j = PUJoint{base}
	case PistonJointType:
		j = PistonJoint{base}
	case DBallJointType:
		j = DBallJoint{base}
	case DHingeJointType:
		j = DHingeJoint{base}
	case TransmissionJointType:
		j = TransmissionJoint{base}
	default:
		j = base
	}
	return j
}

func (j JointBase) c() C.dJointID {
	return C.dJointID(unsafe.Pointer(j))
}

// Destroy destroys the joint base.
func (j JointBase) Destroy() {
	delete(jointData, j)
	C.dJointDestroy(j.c())
}

// SetData associates user-specified data with the joint.
func (j JointBase) SetData(data interface{}) {
	jointData[j] = data
}

// Data returns the user-specified data associated with the joint.
func (j JointBase) Data() interface{} {
	return jointData[j]
}

// NumBodies returns the number of attached bodies.
func (j JointBase) NumBodies() int {
	return int(C.dJointGetNumBodies(j.c()))
}

// Attach attaches two bodies with the joint.
func (j JointBase) Attach(body1, body2 Body) {
	C.dJointAttach(j.c(), body1.c(), body2.c())
}

// SetEnabled sets whether the joint is enabled.
func (j JointBase) SetEnabled(isEnabled bool) {
	if isEnabled {
		C.dJointEnable(j.c())
	} else {
		C.dJointDisable(j.c())
	}
}

// Enabled returns whether the joint is enabled.
func (j JointBase) Enabled() bool {
	return bool(C.dJointIsEnabled(j.c()) != 0)
}

// Type returns the joint type.
func (j JointBase) Type() int {
	return int(C.dJointGetType(j.c()))
}

// Body returns the attached body, specified by index.
func (j JointBase) Body(index int) Body {
	return cToBody(C.dJointGetBody(j.c(), C.int(index)))
}

// SetFeedback sets the feedback forces and torques.
func (j JointBase) SetFeedback(f *JointFeedback) {
	c := &C.dJointFeedback{}
	f.toC(c)
	C.dJointSetFeedback(j.c(), c)
}

// Feedback returns the feedback forces and torques.
func (j JointBase) Feedback() *JointFeedback {
	f := &JointFeedback{}
	f.fromC(C.dJointGetFeedback(j.c()))
	return f
}

// BallJoint implements a ball-and-socket joint.
type BallJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j BallJoint) SetParam(parameter int, value float64) {
	C.dJointSetBallParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j BallJoint) Param(parameter int) float64 {
	return float64(C.dJointGetBallParam(j.c(), C.int(parameter)))
}

// SetAnchor sets the anchor point for the first body.
func (j BallJoint) SetAnchor(pt Vector3) {
	C.dJointSetBallAnchor(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor returns the anchor point for the first body.
func (j BallJoint) Anchor() Vector3 {
	pt := NewVector3()
	C.dJointGetBallAnchor(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAnchor2 sets the anchor point for the second body.
func (j BallJoint) SetAnchor2(pt Vector3) {
	C.dJointSetBallAnchor2(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor2 returns the anchor point for the second body.
func (j BallJoint) Anchor2() Vector3 {
	pt := NewVector3()
	C.dJointGetBallAnchor2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// HingeJoint represents a hinge joint.
type HingeJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j HingeJoint) SetParam(parameter int, value float64) {
	C.dJointSetHingeParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j HingeJoint) Param(parameter int) float64 {
	return float64(C.dJointGetHingeParam(j.c(), C.int(parameter)))
}

// SetAnchor sets the anchor point.
func (j HingeJoint) SetAnchor(pt Vector3) {
	C.dJointSetHingeAnchor(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// SetAnchorDelta sets the hinge anchor delta.
func (j HingeJoint) SetAnchorDelta(pt, delta Vector3) {
	C.dJointSetHingeAnchorDelta(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]),
		C.dReal(delta[0]), C.dReal(delta[1]), C.dReal(delta[2]))
}

// Anchor returns the anchor point for the first body.
func (j HingeJoint) Anchor() Vector3 {
	pt := NewVector3()
	C.dJointGetHingeAnchor(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// Anchor2 returns the anchor point for the second body.
func (j HingeJoint) Anchor2() Vector3 {
	pt := NewVector3()
	C.dJointGetHingeAnchor2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAxis sets the hinge axis.
func (j HingeJoint) SetAxis(axis Vector3) {
	C.dJointSetHingeAxis(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// SetAxisOffset set the hinge axis as if the 2 bodies were already at angle appart.
func (j HingeJoint) SetAxisOffset(axis Vector3, angle float64) {
	C.dJointSetHingeAxisOffset(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]),
		C.dReal(angle))
}

// Axis returns the hinge axis.
func (j HingeJoint) Axis() Vector3 {
	axis := NewVector3()
	C.dJointGetHingeAxis(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// AddTorque adds a torque to the joint.
func (j HingeJoint) AddTorque(torque float64) {
	C.dJointAddHingeTorque(j.c(), C.dReal(torque))
}

// Angle returns the joint angle.
func (j HingeJoint) Angle() float64 {
	return float64(C.dJointGetHingeAngle(j.c()))
}

// AngleRate returns the joint angle's rate of change.
func (j HingeJoint) AngleRate() float64 {
	return float64(C.dJointGetHingeAngleRate(j.c()))
}

// SliderJoint represents a slider joints.
type SliderJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j SliderJoint) SetParam(parameter int, value float64) {
	C.dJointSetSliderParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j SliderJoint) Param(parameter int) float64 {
	return float64(C.dJointGetSliderParam(j.c(), C.int(parameter)))
}

// Position returns the slider position.
func (j SliderJoint) Position() float64 {
	return float64(C.dJointGetSliderPosition(j.c()))
}

// PositionRate returns the slider position's rate of change.
func (j SliderJoint) PositionRate() float64 {
	return float64(C.dJointGetSliderPositionRate(j.c()))
}

// SetAxis sets the slider axis.
func (j SliderJoint) SetAxis(axis Vector3) {
	C.dJointSetSliderAxis(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// SetAxisDelta sets the slider axis delta.
func (j SliderJoint) SetAxisDelta(pt, delta Vector3) {
	C.dJointSetSliderAxisDelta(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]),
		C.dReal(delta[0]), C.dReal(delta[1]), C.dReal(delta[2]))
}

// Axis returns the slider axis.
func (j SliderJoint) Axis() Vector3 {
	axis := NewVector3()
	C.dJointGetSliderAxis(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// AddForce adds a force to the joint.
func (j SliderJoint) AddForce(force float64) {
	C.dJointAddSliderForce(j.c(), C.dReal(force))
}

// ContactJoint represents a contact joint.
type ContactJoint struct {
	JointBase
}

// UniversalJoint represents a universal joint.
type UniversalJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j UniversalJoint) SetParam(parameter int, value float64) {
	C.dJointSetUniversalParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j UniversalJoint) Param(parameter int) float64 {
	return float64(C.dJointGetUniversalParam(j.c(), C.int(parameter)))
}

// SetAnchor sets the anchor point.
func (j UniversalJoint) SetAnchor(pt Vector3) {
	C.dJointSetUniversalAnchor(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor returns the anchor point for the first body.
func (j UniversalJoint) Anchor() Vector3 {
	pt := NewVector3()
	C.dJointGetUniversalAnchor(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// Anchor2 returns the anchor point for the second body.
func (j UniversalJoint) Anchor2() Vector3 {
	pt := NewVector3()
	C.dJointGetUniversalAnchor2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAxis1 sets the first axis.
func (j UniversalJoint) SetAxis1(axis Vector3) {
	C.dJointSetUniversalAxis1(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// SetAxis1Offset sets the first axis as if the 2 bodies were already at
// offset1 and offset2 appart with respect to the first and second axes.
func (j UniversalJoint) SetAxis1Offset(axis Vector3, offset1, offset2 float64) {
	C.dJointSetUniversalAxis1Offset(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]),
		C.dReal(offset1), C.dReal(offset2))
}

// Axis1 returns the first axis.
func (j UniversalJoint) Axis1() Vector3 {
	axis := NewVector3()
	C.dJointGetUniversalAxis1(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAxis2 sets the second axis.
func (j UniversalJoint) SetAxis2(axis Vector3) {
	C.dJointSetUniversalAxis2(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// SetAxis2Offset sets the second axis as if the 2 bodies were already at
// offset1 and offset2 appart with respect to the first and second axes.
func (j UniversalJoint) SetAxis2Offset(axis Vector3, offset1, offset2 float64) {
	C.dJointSetUniversalAxis2Offset(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]),
		C.dReal(offset1), C.dReal(offset1))
}

// Axis2 returns the second axis.
func (j UniversalJoint) Axis2() Vector3 {
	axis := NewVector3()
	C.dJointGetUniversalAxis2(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// Angle1 returns the first angle.
func (j UniversalJoint) Angle1() float64 {
	return float64(C.dJointGetUniversalAngle1(j.c()))
}

// Angle1Rate returns the first angle's rate of change.
func (j UniversalJoint) Angle1Rate() float64 {
	return float64(C.dJointGetUniversalAngle1Rate(j.c()))
}

// Angle2 returns the second angle.
func (j UniversalJoint) Angle2() float64 {
	return float64(C.dJointGetUniversalAngle2(j.c()))
}

// Angle2Rate returns the second angle's rate of change.
func (j UniversalJoint) Angle2Rate() float64 {
	return float64(C.dJointGetUniversalAngle2Rate(j.c()))
}

// Angles returns the two angles.
func (j UniversalJoint) Angles() (float64, float64) {
	var angle1, angle2 float64
	C.dJointGetUniversalAngles(j.c(), (*C.dReal)(&angle1), (*C.dReal)(&angle2))
	return angle1, angle2
}

// AddTorques adds torques to the joint.
func (j UniversalJoint) AddTorques(torque1, torque2 float64) {
	C.dJointAddUniversalTorques(j.c(), C.dReal(torque1), C.dReal(torque2))
}

// Hinge2Joint represents two hinge joints in series.
type Hinge2Joint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j Hinge2Joint) SetParam(parameter int, value float64) {
	C.dJointSetHinge2Param(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j Hinge2Joint) Param(parameter int) float64 {
	return float64(C.dJointGetHinge2Param(j.c(), C.int(parameter)))
}

// SetAnchor sets the anchor point.
func (j Hinge2Joint) SetAnchor(pt Vector3) {
	C.dJointSetHinge2Anchor(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor returns the anchor point for the first body.
func (j Hinge2Joint) Anchor() Vector3 {
	pt := NewVector3()
	C.dJointGetHinge2Anchor(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// Anchor2 returns the anchor point for the second body.
func (j Hinge2Joint) Anchor2() Vector3 {
	pt := NewVector3()
	C.dJointGetHinge2Anchor2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAxis1 sets the first axis.
func (j Hinge2Joint) SetAxis1(axis Vector3) {
	C.dJointSetHinge2Axis1(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis1 returns the first axis.
func (j Hinge2Joint) Axis1() Vector3 {
	axis := NewVector3()
	C.dJointGetHinge2Axis1(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAxis2 sets the second axis.
func (j Hinge2Joint) SetAxis2(axis Vector3) {
	C.dJointSetHinge2Axis2(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis2 returns the second axis.
func (j Hinge2Joint) Axis2() Vector3 {
	axis := NewVector3()
	C.dJointGetHinge2Axis2(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// Angle1 returns the first angle.
func (j Hinge2Joint) Angle1() float64 {
	return float64(C.dJointGetHinge2Angle1(j.c()))
}

// Angle1Rate returns the first angle's rate of change.
func (j Hinge2Joint) Angle1Rate() float64 {
	return float64(C.dJointGetHinge2Angle1Rate(j.c()))
}

// Angle2 returns the second angle.
func (j Hinge2Joint) Angle2() float64 {
	return float64(C.dJointGetHinge2Angle2(j.c()))
}

// Angle2Rate returns the second angle's rate of change.
func (j Hinge2Joint) Angle2Rate() float64 {
	return float64(C.dJointGetHinge2Angle2Rate(j.c()))
}

// AddTorques adds torques to the joint.
func (j Hinge2Joint) AddTorques(torque1, torque2 float64) {
	C.dJointAddHinge2Torques(j.c(), C.dReal(torque1), C.dReal(torque2))
}

// FixedJoint represents a fixed joint.
type FixedJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j FixedJoint) SetParam(parameter int, value float64) {
	C.dJointSetFixedParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j FixedJoint) Param(parameter int) float64 {
	return float64(C.dJointGetFixedParam(j.c(), C.int(parameter)))
}

// Fix fixes the joint in its current state.
func (j FixedJoint) Fix() {
	C.dJointSetFixed(j.c())
}

// NullJoint represents a null joint.
type NullJoint struct {
	JointBase
}

// AMotorJoint represents an angular motor joint.
type AMotorJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j AMotorJoint) SetParam(parameter int, value float64) {
	C.dJointSetAMotorParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j AMotorJoint) Param(parameter int) float64 {
	return float64(C.dJointGetAMotorParam(j.c(), C.int(parameter)))
}

// SetNumAxes sets the number of axes.
func (j AMotorJoint) SetNumAxes(num int) {
	C.dJointSetAMotorNumAxes(j.c(), C.int(num))
}

// NumAxes returns the number of axes.
func (j AMotorJoint) NumAxes() int {
	return int(C.dJointGetAMotorNumAxes(j.c()))
}

// SetAxis sets the given axis relative to body rel (1 or 2) or none (0).
func (j AMotorJoint) SetAxis(num, rel int, axis Vector3) {
	C.dJointSetAMotorAxis(j.c(), C.int(num), C.int(rel),
		C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis returns the given axis.
func (j AMotorJoint) Axis(num int) Vector3 {
	axis := NewVector3()
	C.dJointGetAMotorAxis(j.c(), C.int(num), (*C.dReal)(&axis[0]))
	return axis
}

// AxisRel returns the relative mode for the given axis.
func (j AMotorJoint) AxisRel(num int) int {
	return int(C.dJointGetAMotorAxisRel(j.c(), C.int(num)))
}

// SetAngle sets the angle of the given axis.
func (j AMotorJoint) SetAngle(num int, angle float64) {
	C.dJointSetAMotorAngle(j.c(), C.int(num), C.dReal(angle))
}

// Angle returns the angle of the given axis.
func (j AMotorJoint) Angle(num int) float64 {
	return float64(C.dJointGetAMotorAngle(j.c(), C.int(num)))
}

// AngleRate returns the angle's rate of change for the given axis.
func (j AMotorJoint) AngleRate(num int) float64 {
	return float64(C.dJointGetAMotorAngleRate(j.c(), C.int(num)))
}

// SetMode sets the mode.
func (j AMotorJoint) SetMode(mode int) {
	C.dJointSetAMotorMode(j.c(), C.int(mode))
}

// Mode returns the mode.
func (j AMotorJoint) Mode() int {
	return int(C.dJointGetAMotorMode(j.c()))
}

// AddTorques adds torques to the joint.
func (j AMotorJoint) AddTorques(torque1, torque2, torque3 float64) {
	C.dJointAddAMotorTorques(j.c(), C.dReal(torque1), C.dReal(torque2), C.dReal(torque3))
}

// LMotorJoint represents a linear motor joint.
type LMotorJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j LMotorJoint) SetParam(parameter int, value float64) {
	C.dJointSetLMotorParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j LMotorJoint) Param(parameter int) float64 {
	return float64(C.dJointGetLMotorParam(j.c(), C.int(parameter)))
}

// SetNumAxes sets the number of axes.
func (j LMotorJoint) SetNumAxes(num int) {
	C.dJointSetLMotorNumAxes(j.c(), C.int(num))
}

// NumAxes returns the number of axes.
func (j LMotorJoint) NumAxes() int {
	return int(C.dJointGetLMotorNumAxes(j.c()))
}

// SetAxis sets the given axis relative to a body (1 or 2) or none (0).
func (j LMotorJoint) SetAxis(num, rel int, axis Vector3) {
	C.dJointSetLMotorAxis(j.c(), C.int(num), C.int(rel),
		C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis returns the given axis.
func (j LMotorJoint) Axis(num int) Vector3 {
	axis := NewVector3()
	C.dJointGetLMotorAxis(j.c(), C.int(num), (*C.dReal)(&axis[0]))
	return axis
}

// Plane2DJoint represents a plane joint.
type Plane2DJoint struct {
	JointBase
}

// SetXParam sets a joint parameter.
func (j Plane2DJoint) SetXParam(parameter int, value float64) {
	C.dJointSetPlane2DXParam(j.c(), C.int(parameter), C.dReal(value))
}

// SetYParam sets a joint parameter.
func (j Plane2DJoint) SetYParam(parameter int, value float64) {
	C.dJointSetPlane2DYParam(j.c(), C.int(parameter), C.dReal(value))
}

// SetAngleParam sets a joint parameter.
func (j Plane2DJoint) SetAngleParam(parameter int, value float64) {
	C.dJointSetPlane2DAngleParam(j.c(), C.int(parameter), C.dReal(value))
}

// PRJoint represents a prismatic rotoide joint.
type PRJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j PRJoint) SetParam(parameter int, value float64) {
	C.dJointSetPRParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j PRJoint) Param(parameter int) float64 {
	return float64(C.dJointGetPRParam(j.c(), C.int(parameter)))
}

// SetAnchor sets the anchor point.
func (j PRJoint) SetAnchor(pt Vector3) {
	C.dJointSetPRAnchor(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor returns the anchor point.
func (j PRJoint) Anchor() Vector3 {
	pt := NewVector3()
	C.dJointGetPRAnchor(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAxis1 sets the first axis.
func (j PRJoint) SetAxis1(axis Vector3) {
	C.dJointSetPRAxis1(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis1 returns the first axis.
func (j PRJoint) Axis1() Vector3 {
	axis := NewVector3()
	C.dJointGetPRAxis1(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAxis2 sets the second axis.
func (j PRJoint) SetAxis2(axis Vector3) {
	C.dJointSetPRAxis2(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis2 returns the second axis.
func (j PRJoint) Axis2() Vector3 {
	axis := NewVector3()
	C.dJointGetPRAxis2(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// Position returns the slider position.
func (j PRJoint) Position() float64 {
	return float64(C.dJointGetPRPosition(j.c()))
}

// PositionRate returns the slider position's rate of change.
func (j PRJoint) PositionRate() float64 {
	return float64(C.dJointGetPRPositionRate(j.c()))
}

// Angle returns the joint angle.
func (j PRJoint) Angle() float64 {
	return float64(C.dJointGetPRAngle(j.c()))
}

// AngleRate returns the joint angle's rate of change.
func (j PRJoint) AngleRate() float64 {
	return float64(C.dJointGetPRAngleRate(j.c()))
}

// AddTorque adds a torque to the joint.
func (j PRJoint) AddTorque(torque float64) {
	C.dJointAddPRTorque(j.c(), C.dReal(torque))
}

// PUJoint represents a prismatic universal joint.
type PUJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j PUJoint) SetParam(parameter int, value float64) {
	C.dJointSetPUParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j PUJoint) Param(parameter int) float64 {
	return float64(C.dJointGetPUParam(j.c(), C.int(parameter)))
}

// SetAnchor sets the anchor point.
func (j PUJoint) SetAnchor(pt Vector3) {
	C.dJointSetPUAnchor(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor returns the anchor point.
func (j PUJoint) Anchor() Vector3 {
	pt := NewVector3()
	C.dJointGetPUAnchor(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAnchorOffset sets the anchor as if the 2 bodies were already delta appart.
func (j PUJoint) SetAnchorOffset(pt, delta Vector3) {
	C.dJointSetPUAnchorOffset(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]),
		C.dReal(delta[0]), C.dReal(delta[1]), C.dReal(delta[2]))
}

// SetAxis1 sets the first axis.
func (j PUJoint) SetAxis1(axis Vector3) {
	C.dJointSetPUAxis1(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis1 returns the first axis.
func (j PUJoint) Axis1() Vector3 {
	axis := NewVector3()
	C.dJointGetPUAxis1(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAxis2 sets the second axis.
func (j PUJoint) SetAxis2(axis Vector3) {
	C.dJointSetPUAxis2(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis2 returns the second axis.
func (j PUJoint) Axis2() Vector3 {
	axis := NewVector3()
	C.dJointGetPUAxis2(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAxis3 sets the third (prismatic) axis.
func (j PUJoint) SetAxis3(axis Vector3) {
	C.dJointSetPUAxis3(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis3 returns the third (prismatic) axis.
func (j PUJoint) Axis3() Vector3 {
	axis := NewVector3()
	C.dJointGetPUAxis3(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// Position returns the joint position.
func (j PUJoint) Position() float64 {
	return float64(C.dJointGetPUPosition(j.c()))
}

// PositionRate returns the joint position's rate of change.
func (j PUJoint) PositionRate() float64 {
	return float64(C.dJointGetPUPositionRate(j.c()))
}

// Angle1 returns the first angle.
func (j PUJoint) Angle1() float64 {
	return float64(C.dJointGetPUAngle1(j.c()))
}

// Angle1Rate returns the first angle's rate of change.
func (j PUJoint) Angle1Rate() float64 {
	return float64(C.dJointGetPUAngle1Rate(j.c()))
}

// Angle2 returns the second angle.
func (j PUJoint) Angle2() float64 {
	return float64(C.dJointGetPUAngle2(j.c()))
}

// Angle2Rate returns the second angle's rate of change.
func (j PUJoint) Angle2Rate() float64 {
	return float64(C.dJointGetPUAngle2Rate(j.c()))
}

// Angles returns the two joint angles.
func (j PUJoint) Angles() (float64, float64) {
	var angle1, angle2 float64
	C.dJointGetPUAngles(j.c(), (*C.dReal)(&angle1), (*C.dReal)(&angle2))
	return angle1, angle2
}

// PistonJoint represents a piston joint.
type PistonJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j PistonJoint) SetParam(parameter int, value float64) {
	C.dJointSetPistonParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j PistonJoint) Param(parameter int) float64 {
	return float64(C.dJointGetPistonParam(j.c(), C.int(parameter)))
}

// SetAnchor sets the anchor point.
func (j PistonJoint) SetAnchor(pt Vector3) {
	C.dJointSetPistonAnchor(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// SetAnchorOffset sets the anchor as if the 2 bodies were already delta appart.
func (j PistonJoint) SetAnchorOffset(pt, delta Vector3) {
	C.dJointSetPistonAnchorOffset(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]),
		C.dReal(delta[0]), C.dReal(delta[1]), C.dReal(delta[2]))
}

// Anchor2 returns the anchor point for the second body.
func (j PistonJoint) Anchor2() Vector3 {
	pt := NewVector3()
	C.dJointGetPistonAnchor2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAxis sets the piston axis.
func (j PistonJoint) SetAxis(axis Vector3) {
	C.dJointSetPistonAxis(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis returns the piston axis.
func (j PistonJoint) Axis() Vector3 {
	axis := NewVector3()
	C.dJointGetPistonAxis(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// Position returns the piston position.
func (j PistonJoint) Position() float64 {
	return float64(C.dJointGetPistonPosition(j.c()))
}

// PositionRate returns the piston position's rate of change.
func (j PistonJoint) PositionRate() float64 {
	return float64(C.dJointGetPistonPositionRate(j.c()))
}

// Angle returns the joint angle.
func (j PistonJoint) Angle() float64 {
	return float64(C.dJointGetPistonAngle(j.c()))
}

// AngleRate returns the joint angle's rate of change.
func (j PistonJoint) AngleRate() float64 {
	return float64(C.dJointGetPistonAngleRate(j.c()))
}

// AddForce adds a force to the joint.
func (j PistonJoint) AddForce(force float64) {
	C.dJointAddPistonForce(j.c(), C.dReal(force))
}

// DBallJoint represents a double ball joint.
type DBallJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j DBallJoint) SetParam(parameter int, value float64) {
	C.dJointSetDBallParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j DBallJoint) Param(parameter int) float64 {
	return float64(C.dJointGetDBallParam(j.c(), C.int(parameter)))
}

// SetAnchor1 sets the anchor point for the first body.
func (j DBallJoint) SetAnchor1(pt Vector3) {
	C.dJointSetDBallAnchor1(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor1 returns the anchor point for the first body.
func (j DBallJoint) Anchor1() Vector3 {
	pt := NewVector3()
	C.dJointGetDBallAnchor1(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAnchor2 sets the anchor point for the second body.
func (j DBallJoint) SetAnchor2(pt Vector3) {
	C.dJointSetDBallAnchor2(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor2 returns the anchor point for the second body.
func (j DBallJoint) Anchor2() Vector3 {
	pt := NewVector3()
	C.dJointGetDBallAnchor2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// Distance returns the joint distance.
func (j DBallJoint) Distance() float64 {
	return float64(C.dJointGetDBallDistance(j.c()))
}

// DHingeJoint represents a double hinge joint.
type DHingeJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j DHingeJoint) SetParam(parameter int, value float64) {
	C.dJointSetDHingeParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j DHingeJoint) Param(parameter int) float64 {
	return float64(C.dJointGetDHingeParam(j.c(), C.int(parameter)))
}

// SetAxis sets the joint axis.
func (j DHingeJoint) SetAxis(axis Vector3) {
	C.dJointSetDHingeAxis(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis returns the joint axis.
func (j DHingeJoint) Axis() Vector3 {
	axis := NewVector3()
	C.dJointGetDHingeAxis(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAnchor1 sets the anchor point for the first body.
func (j DHingeJoint) SetAnchor1(pt Vector3) {
	C.dJointSetDHingeAnchor1(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor1 returns the anchor point for the first body.
func (j DHingeJoint) Anchor1() Vector3 {
	pt := NewVector3()
	C.dJointGetDHingeAnchor1(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAnchor2 sets the anchor point for the second body.
func (j DHingeJoint) SetAnchor2(pt Vector3) {
	C.dJointSetDHingeAnchor2(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor2 returns the anchor point for the second body.
func (j DHingeJoint) Anchor2() Vector3 {
	pt := NewVector3()
	C.dJointGetDHingeAnchor2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// Distance returns the joint distance.
func (j DHingeJoint) Distance() float64 {
	return float64(C.dJointGetDHingeDistance(j.c()))
}

// TransmissionJoint represents a transmission joint.
type TransmissionJoint struct {
	JointBase
}

// SetParam sets a joint parameter.
func (j TransmissionJoint) SetParam(parameter int, value float64) {
	C.dJointSetTransmissionParam(j.c(), C.int(parameter), C.dReal(value))
}

// Param returns a joint parameter.
func (j TransmissionJoint) Param(parameter int) float64 {
	return float64(C.dJointGetTransmissionParam(j.c(), C.int(parameter)))
}

// SetAxis sets the common axis.
func (j TransmissionJoint) SetAxis(axis Vector3) {
	C.dJointSetTransmissionAxis(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis returns the common axis.
func (j TransmissionJoint) Axis() Vector3 {
	axis := NewVector3()
	C.dJointGetTransmissionAxis(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAxis1 sets the first axis.
func (j TransmissionJoint) SetAxis1(axis Vector3) {
	C.dJointSetTransmissionAxis1(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis1 returns the first axis.
func (j TransmissionJoint) Axis1() Vector3 {
	axis := NewVector3()
	C.dJointGetTransmissionAxis1(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAxis2 sets the second axis.
func (j TransmissionJoint) SetAxis2(axis Vector3) {
	C.dJointSetTransmissionAxis2(j.c(), C.dReal(axis[0]), C.dReal(axis[1]), C.dReal(axis[2]))
}

// Axis2 returns the second axis.
func (j TransmissionJoint) Axis2() Vector3 {
	axis := NewVector3()
	C.dJointGetTransmissionAxis2(j.c(), (*C.dReal)(&axis[0]))
	return axis
}

// SetAnchor1 sets the anchor point for the first body.
func (j TransmissionJoint) SetAnchor1(pt Vector3) {
	C.dJointSetTransmissionAnchor1(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor1 returns the anchor point for the first body.
func (j TransmissionJoint) Anchor1() Vector3 {
	pt := NewVector3()
	C.dJointGetTransmissionAnchor1(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// SetAnchor2 sets the anchor point for the second body.
func (j TransmissionJoint) SetAnchor2(pt Vector3) {
	C.dJointSetTransmissionAnchor2(j.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]))
}

// Anchor2 returns the anchor point for the second body.
func (j TransmissionJoint) Anchor2() Vector3 {
	pt := NewVector3()
	C.dJointGetTransmissionAnchor2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// ContactPoint1 returns the contact point on the first wheel.
func (j TransmissionJoint) ContactPoint1() Vector3 {
	pt := NewVector3()
	C.dJointGetTransmissionContactPoint1(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// ContactPoint2 returns the contact point on the second wheel.
func (j TransmissionJoint) ContactPoint2() Vector3 {
	pt := NewVector3()
	C.dJointGetTransmissionContactPoint2(j.c(), (*C.dReal)(&pt[0]))
	return pt
}

// Angle1 returns the angle of the first wheel.
func (j TransmissionJoint) Angle1() float64 {
	return float64(C.dJointGetTransmissionAngle1(j.c()))
}

// Angle2 returns the angle of the second wheel.
func (j TransmissionJoint) Angle2() float64 {
	return float64(C.dJointGetTransmissionAngle2(j.c()))
}

// SetRadius1 sets the radius of the first wheel.
func (j TransmissionJoint) SetRadius1(radius float64) {
	C.dJointSetTransmissionRadius1(j.c(), C.dReal(radius))
}

// Radius1 returns the radius of the first wheel.
func (j TransmissionJoint) Radius1() float64 {
	return float64(C.dJointGetTransmissionRadius1(j.c()))
}

// SetRadius2 sets the radius of the second wheel.
func (j TransmissionJoint) SetRadius2(radius float64) {
	C.dJointSetTransmissionRadius2(j.c(), C.dReal(radius))
}

// Radius2 returns the radius of the second wheel.
func (j TransmissionJoint) Radius2() float64 {
	return float64(C.dJointGetTransmissionRadius2(j.c()))
}

// SetMode sets the transmission mode.
func (j TransmissionJoint) SetMode(mode int) {
	C.dJointSetTransmissionMode(j.c(), C.int(mode))
}

// Mode returns the transmission mode.
func (j TransmissionJoint) Mode() int {
	return int(C.dJointGetTransmissionMode(j.c()))
}

// SetRatio sets the gear ratio.
func (j TransmissionJoint) SetRatio(ratio float64) {
	C.dJointSetTransmissionRatio(j.c(), C.dReal(ratio))
}

// Ratio returns the gear ratio.
func (j TransmissionJoint) Ratio() float64 {
	return float64(C.dJointGetTransmissionRatio(j.c()))
}

// SetBacklash set the backlash (gear tooth play distance).
func (j TransmissionJoint) SetBacklash(backlash float64) {
	C.dJointSetTransmissionBacklash(j.c(), C.dReal(backlash))
}

// Backlash returns the backlash (gear tooth play distance).
func (j TransmissionJoint) Backlash() float64 {
	return float64(C.dJointGetTransmissionBacklash(j.c()))
}
