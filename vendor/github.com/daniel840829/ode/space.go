package ode

// #include <ode/ode.h>
// extern void callNearCallback(void *data, dGeomID obj1, dGeomID obj2);
import "C"

import (
	"unsafe"
)

// Sweep and prune axis orders
const (
	SAPAxesXYZ = C.dSAP_AXES_XYZ
	SAPAxesXZY = C.dSAP_AXES_XZY
	SAPAxesYXZ = C.dSAP_AXES_YXZ
	SAPAxesYZX = C.dSAP_AXES_YZX
	SAPAxesZXY = C.dSAP_AXES_ZXY
	SAPAxesZYX = C.dSAP_AXES_ZYX
)

// Space represents a space containing bodies.
type Space interface {
	c() C.dSpaceID
	Destroy()
	SetCleanup(mode bool)
	Cleanup() bool
	SetManualCleanup(mode bool)
	SetSublevel(sublevel int)
	Sublevel() int
	ManualCleanup() bool
	Clean()
	Class() int
	Add(g Geom)
	Remove(g Geom)
	Query(g Geom) bool
	NumGeoms(g Geom) int
	Geom(index int) Geom
	Collide(data interface{}, cb NearCallback)
	NewSphere(radius float64) Sphere
	NewBox(lens Vector3) Box
	NewPlane(params Vector4) Plane
	NewCapsule(radius, length float64) Capsule
	NewCylinder(radius, length float64) Cylinder
	NewRay(length float64) Ray
	NewHeightfield(data HeightfieldData, placeable bool) Heightfield
	NewSimpleSpace() SimpleSpace
	NewHashSpace() HashSpace
	NewQuadTreeSpace(center, extents Vector3, depth int) QuadTreeSpace
	NewSweepAndPruneSpace(axisOrder int) SweepAndPruneSpace
}

// SpaceBase implements Space, and is embedded by specific Space types.
type SpaceBase uintptr

func cToSpace(c C.dSpaceID) Space {
	base := SpaceBase(unsafe.Pointer(c))
	var s Space
	switch int(C.dSpaceGetClass(c)) {
	case SimpleSpaceClass:
		s = SimpleSpace{base}
	case HashSpaceClass:
		s = HashSpace{base}
	case QuadTreeSpaceClass:
		s = QuadTreeSpace{base}
	case SweepAndPruneSpaceClass:
		s = SweepAndPruneSpace{base}
	default:
		s = base
	}
	return s
}

func CToSpace(c C.dSpaceID) Space {
	base := SpaceBase(unsafe.Pointer(c))
	var s Space
	switch int(C.dSpaceGetClass(c)) {
	case SimpleSpaceClass:
		s = SimpleSpace{base}
	case HashSpaceClass:
		s = HashSpace{base}
	case QuadTreeSpaceClass:
		s = QuadTreeSpace{base}
	case SweepAndPruneSpaceClass:
		s = SweepAndPruneSpace{base}
	default:
		s = base
	}
	return s
}

// NilSpace returns the top level "0" space
func NilSpace() Space {
	return SpaceBase(0)
}

func (s SpaceBase) c() C.dSpaceID {
	return C.dSpaceID(unsafe.Pointer(s))
}

// Destroy destroys the space.
func (s SpaceBase) Destroy() {
	C.dSpaceDestroy(s.c())
}

// SetCleanup sets whether contained objects will be destroyed.
func (s SpaceBase) SetCleanup(mode bool) {
	C.dSpaceSetCleanup(s.c(), C.int(btoi(mode)))
}

// Cleanup returns whether contained objects will be destroyed.
func (s SpaceBase) Cleanup() bool {
	return C.dSpaceGetCleanup(s.c()) != 0
}

// SetManualCleanup sets whether this space is marked for manual cleanup.
func (s SpaceBase) SetManualCleanup(mode bool) {
	C.dSpaceSetManualCleanup(s.c(), C.int(btoi(mode)))
}

// ManualCleanup returns whether this space is marked for manual cleanup.
func (s SpaceBase) ManualCleanup() bool {
	return C.dSpaceGetManualCleanup(s.c()) != 0
}

// SetSublevel sets the sublevel for this space.
func (s SpaceBase) SetSublevel(sublevel int) {
	C.dSpaceSetSublevel(s.c(), C.int(sublevel))
}

// Sublevel returns the sublevel for this space.
func (s SpaceBase) Sublevel() int {
	return int(C.dSpaceGetSublevel(s.c()))
}

// Clean cleans the space.
func (s SpaceBase) Clean() {
	C.dSpaceClean(s.c())
}

// Class returns the space class.
func (s SpaceBase) Class() int {
	return int(C.dSpaceGetClass(s.c()))
}

// Add adds a geometry to the space.
func (s SpaceBase) Add(g Geom) {
	C.dSpaceAdd(s.c(), g.c())
}

// Remove removes a geometry from the space.
func (s SpaceBase) Remove(g Geom) {
	C.dSpaceRemove(s.c(), g.c())
}

// Query returns whether a geometry is contained in the space.
func (s SpaceBase) Query(g Geom) bool {
	return C.dSpaceQuery(s.c(), g.c()) != 0
}

// NumGeoms returns the number of geometries contained in the space.
func (s SpaceBase) NumGeoms(g Geom) int {
	return int(C.dSpaceGetNumGeoms(s.c()))
}

// Geom returns the specified contained geometry.
func (s SpaceBase) Geom(index int) Geom {
	return cToGeom(C.dSpaceGetGeom(s.c(), C.int(index)))
}

// Collide tests for collision between  contained objects.
func (s SpaceBase) Collide(data interface{}, cb NearCallback) {
	cbData := &nearCallbackData{fn: cb, data: data}
	C.dSpaceCollide(s.c(), unsafe.Pointer(cbData),
		(*C.dNearCallback)(C.callNearCallback))
}

// NewSphere returns a new Sphere instance.
func (s SpaceBase) NewSphere(radius float64) Sphere {
	return cToGeom(C.dCreateSphere(s.c(), C.dReal(radius))).(Sphere)
}

// NewBox returns a new Box instance.
func (s SpaceBase) NewBox(lens Vector3) Box {
	return cToGeom(C.dCreateBox(s.c(), C.dReal(lens[0]), C.dReal(lens[1]), C.dReal(lens[2]))).(Box)
}

// NewPlane returns a new Plane instance.
func (s SpaceBase) NewPlane(params Vector4) Plane {
	return cToGeom(C.dCreatePlane(s.c(), C.dReal(params[0]), C.dReal(params[1]),
		C.dReal(params[2]), C.dReal(params[3]))).(Plane)
}

// NewCapsule returns a new Capsule instance.
func (s SpaceBase) NewCapsule(radius, length float64) Capsule {
	return cToGeom(C.dCreateCapsule(s.c(), C.dReal(radius), C.dReal(length))).(Capsule)
}

// NewCylinder returns a new Cylinder instance.
func (s SpaceBase) NewCylinder(radius, length float64) Cylinder {
	return cToGeom(C.dCreateCylinder(s.c(), C.dReal(radius), C.dReal(length))).(Cylinder)
}

// NewRay returns a new Ray instance.
func (s SpaceBase) NewRay(length float64) Ray {
	return cToGeom(C.dCreateRay(s.c(), C.dReal(length))).(Ray)
}

// NewHeightfield returns a new Heightfield instance.
func (s SpaceBase) NewHeightfield(data HeightfieldData, placeable bool) Heightfield {
	return cToGeom(C.dCreateHeightfield(s.c(), data.c(), C.int(btoi(placeable)))).(Heightfield)
}

// NewConvex returns a new Convex instance.
func (s SpaceBase) NewConvex(planes PlaneList, pts VertexList, polyList PolygonList) Convex {
	return cToGeom(C.dCreateConvex(s.c(), (*C.dReal)(&planes[0][0]), C.uint(len(planes)),
		(*C.dReal)(&pts[0][0]), C.uint(len(pts)), &polyList[0])).(Convex)
}

// NewTriMesh returns a new TriMesh instance.
func (s SpaceBase) NewTriMesh(data TriMeshData) TriMesh {
	return cToGeom(C.dCreateTriMesh(s.c(), data.c(), nil, nil, nil)).(TriMesh)
}

// NewSimpleSpace returns a new SimpleSpace instance.
func (s SpaceBase) NewSimpleSpace() SimpleSpace {
	return cToSpace(C.dSimpleSpaceCreate(s.c())).(SimpleSpace)
}

// NewHashSpace returns a new HashSpace instance.
func (s SpaceBase) NewHashSpace() HashSpace {
	return cToSpace(C.dHashSpaceCreate(s.c())).(HashSpace)
}

// NewQuadTreeSpace returns a new QuadTreeSpace instance.
func (s SpaceBase) NewQuadTreeSpace(center, extents Vector3, depth int) QuadTreeSpace {
	return cToSpace(C.dQuadTreeSpaceCreate(s.c(), (*C.dReal)(&center[0]),
		(*C.dReal)(&extents[0]), C.int(depth))).(QuadTreeSpace)
}

// NewSweepAndPruneSpace returns a new SweepAndPruneSpace instance.
func (s SpaceBase) NewSweepAndPruneSpace(axisOrder int) SweepAndPruneSpace {
	return cToSpace(C.dSweepAndPruneSpaceCreate(s.c(), C.int(axisOrder))).(SweepAndPruneSpace)
}

// SimpleSpace represents a simple space.
type SimpleSpace struct {
	SpaceBase
}

// HashSpace represents a hash space.
type HashSpace struct {
	SpaceBase
}

// SetLevels sets the minimum and maximum levels.
func (s HashSpace) SetLevels(min, max int) {
	C.dHashSpaceSetLevels(s.c(), C.int(min), C.int(max))
}

// Levels returns the minimum and maximum levels.
func (s HashSpace) Levels() (int, int) {
	var min, max C.int
	C.dHashSpaceGetLevels(s.c(), &min, &max)
	return int(min), int(max)
}

// QuadTreeSpace represents a quad tree space.
type QuadTreeSpace struct {
	SpaceBase
}

// SweepAndPruneSpace represents a sweep and prune space.
type SweepAndPruneSpace struct {
	SpaceBase
}
