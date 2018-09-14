package ode

// #include <ode/ode.h>
import "C"

// Mass represents object mass properties.
type Mass struct {
	Center  Vector3
	Inertia Matrix3
	Mass    float64
}

func (m *Mass) toC(c *C.dMass) {
	c.mass = C.dReal(m.Mass)
	Vector(m.Center).toC((*C.dReal)(&c.c[0]))
	Matrix(m.Inertia).toC((*C.dReal)(&c.I[0]))
}

func (m *Mass) fromC(c *C.dMass) {
	m.Mass = float64(c.mass)
	Vector(m.Center).fromC((*C.dReal)(&c.c[0]))
	Matrix(m.Inertia).fromC((*C.dReal)(&c.I[0]))
}

// NewMass returns a new Mass instance.
func NewMass() *Mass {
	return &Mass{
		Center:  NewVector3(),
		Inertia: NewMatrix3(),
	}
}

// Check returns whether the mass's parameter values are valid.
func (m *Mass) Check() bool {
	c := &C.dMass{}
	m.toC(c)
	return C.dMassCheck(c) != 0
}

// SetZero sets the mass to 0.
func (m *Mass) SetZero() {
	c := &C.dMass{}
	C.dMassSetZero(c)
	m.fromC(c)
}

// SetParams sets the mass parameters.
func (m *Mass) SetParams(mass float64, com Vector3, inert Matrix3) {
	c := &C.dMass{}
	C.dMassSetParameters(c, C.dReal(mass),
		C.dReal(com[0]), C.dReal(com[1]), C.dReal(com[2]),
		C.dReal(inert[0][0]), C.dReal(inert[1][1]), C.dReal(inert[2][2]),
		C.dReal(inert[0][1]), C.dReal(inert[0][2]), C.dReal(inert[1][3]))
	m.fromC(c)
}

// SetSphere sets the mass for a sphere of given properties.
func (m *Mass) SetSphere(density, radius float64) {
	c := &C.dMass{}
	C.dMassSetSphere(c, C.dReal(density), C.dReal(radius))
	m.fromC(c)
}

// SetSphereTotal sets the mass for a sphere of given properties.
func (m *Mass) SetSphereTotal(totalMass, radius float64) {
	c := &C.dMass{}
	C.dMassSetSphereTotal(c, C.dReal(totalMass), C.dReal(radius))
	m.fromC(c)
}

// SetCapsule sets the mass for a capsule of given properties.
func (m *Mass) SetCapsule(density float64, direction int, radius, length float64) {
	c := &C.dMass{}
	C.dMassSetCapsule(c, C.dReal(density), C.int(direction), C.dReal(radius),
		C.dReal(length))
	m.fromC(c)
}

// SetCapsuleTotal sets the mass for a capsule of given properties.
func (m *Mass) SetCapsuleTotal(totalMass float64, direction int, radius, length float64) {
	c := &C.dMass{}
	C.dMassSetCapsuleTotal(c, C.dReal(totalMass), C.int(direction), C.dReal(radius),
		C.dReal(length))
	m.fromC(c)
}

// SetCylinder sets the mass for a cylinder of given properties.
func (m *Mass) SetCylinder(density float64, direction int, radius, length float64) {
	c := &C.dMass{}
	C.dMassSetCylinder(c, C.dReal(density), C.int(direction), C.dReal(radius),
		C.dReal(length))
	m.fromC(c)
}

// SetCylinderTotal sets the mass for a cylinder of given properties.
func (m *Mass) SetCylinderTotal(totalMass float64, direction int, radius, length float64) {
	c := &C.dMass{}
	C.dMassSetCylinderTotal(c, C.dReal(totalMass), C.int(direction), C.dReal(radius),
		C.dReal(length))
	m.fromC(c)
}

// SetBox sets the mass for a box of given properties.
func (m *Mass) SetBox(density float64, lens Vector3) {
	c := &C.dMass{}
	C.dMassSetBox(c, C.dReal(density),
		C.dReal(lens[0]), C.dReal(lens[1]), C.dReal(lens[2]))
	m.fromC(c)
}

// SetBoxTotal sets the mass for a box of given properties.
func (m *Mass) SetBoxTotal(totalMass float64, lens Vector3) {
	c := &C.dMass{}
	C.dMassSetBoxTotal(c, C.dReal(totalMass),
		C.dReal(lens[0]), C.dReal(lens[1]), C.dReal(lens[2]))
	m.fromC(c)
}

// SetTrimesh sets the mass for the given triangle mesh.
func (m *Mass) SetTriMesh(density float64, mesh TriMesh) {
	c := &C.dMass{}
	C.dMassSetTrimesh(c, C.dReal(density), mesh.c())
	m.fromC(c)
}

// SetTrimeshTotal sets the mass for the given triangle mesh.
func (m *Mass) SetTriMeshTotal(totalMass float64, mesh TriMesh) {
	c := &C.dMass{}
	C.dMassSetTrimeshTotal(c, C.dReal(totalMass), mesh.c())
	m.fromC(c)
}

// Adjust sets parameters based on the given total mass.
func (m *Mass) Adjust(mass float64) {
	c := &C.dMass{}
	m.toC(c)
	C.dMassAdjust(c, C.dReal(mass))
	m.fromC(c)
}

// Translate translates the mass by vec.
func (m *Mass) Translate(vec Vector3) {
	c := &C.dMass{}
	m.toC(c)
	C.dMassTranslate(c, C.dReal(vec[0]), C.dReal(vec[1]), C.dReal(vec[2]))
	m.fromC(c)
}

// Rotate rotates the mass by rot.
func (m *Mass) Rotate(rot Matrix3) {
	c := &C.dMass{}
	m.toC(c)
	C.dMassRotate(c, (*C.dReal)(&rot[0][0]))
	m.fromC(c)
}

// Add adds the other mass to this mass.
func (m *Mass) Add(other *Mass) {
	c, oc := &C.dMass{}, &C.dMass{}
	m.toC(c)
	other.toC(oc)
	C.dMassAdd(c, oc)
	m.fromC(c)
}
