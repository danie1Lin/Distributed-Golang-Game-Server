package ode

// #include <ode/ode.h>
// extern void callNearCallback(void *data, dGeomID obj1, dGeomID obj2);
import "C"

import (
	"unsafe"
)

// Geometry classes
const (
	SphereClass             = C.dSphereClass
	BoxClass                = C.dBoxClass
	CapsuleClass            = C.dCapsuleClass
	CylinderClass           = C.dCylinderClass
	PlaneClass              = C.dPlaneClass
	RayClass                = C.dRayClass
	ConvexClass             = C.dConvexClass
	TriMeshClass            = C.dTriMeshClass
	HeightfieldClass        = C.dHeightfieldClass
	SimpleSpaceClass        = C.dSimpleSpaceClass
	HashSpaceClass          = C.dHashSpaceClass
	SweepAndPruneSpaceClass = C.dSweepAndPruneSpaceClass
	QuadTreeSpaceClass      = C.dQuadTreeSpaceClass

	NumClasses = C.dGeomNumClasses

	MaxUserClasses = C.dMaxUserClasses
	FirstUserClass = C.dFirstUserClass
	LastUserClass  = C.dLastUserClass

	FirstSpaceClass = C.dFirstSpaceClass
	LastSpaceClass  = C.dLastSpaceClass
)

var (
	geomData = map[Geom]interface{}{}
)

// Geom represents rigid body geometry.
type Geom interface {
	c() C.dGeomID
	ToSpace() Space
	Destroy()
	SetData(data interface{})
	Data() interface{}
	SetBody(body Body)
	Body() Body
	SetPosition(pos Vector3)
	Position() Vector3
	SetRotation(rot Matrix3)
	Rotation() Matrix3
	SetQuaternion(quat Quaternion)
	Quaternion() Quaternion
	AABB() AABB
	IsSpace() bool
	Space() Space
	Class() int
	SetCategoryBits(bits int)
	SetCollideBits(bits int)
	CategoryBits() int
	CollideBits() int
	SetEnabled(isEnabled bool)
	Enabled() bool
	RelPointPos(pt Vector3) Vector3
	PosRelPoint(pos Vector3) Vector3
	VectorToWorld(vec Vector3) Vector3
	VectorFromWorld(wld Vector3) Vector3
	OffsetPosition() Vector3
	SetOffsetPosition(pos Vector3)
	OffsetRotation() Matrix3
	SetOffsetRotation(rot Matrix3)
	OffsetQuaternion() Quaternion
	SetOffsetQuaternion(quat Quaternion)
	SetOffsetWorldPosition(pos Vector3)
	SetOffsetWorldRotation(rot Matrix3)
	SetOffsetWorldQuaternion(quat Quaternion)
	ClearOffset()
	IsOffset() bool
	Collide(other Geom, maxContacts uint16, flags int) []ContactGeom
	Collide2(other Geom, data interface{}, cb NearCallback)
	Next() Geom
}

// GeomBase implements Geom, and is embedded by specific Geom types.
type GeomBase uintptr

func cToGeom(c C.dGeomID) Geom {
	base := GeomBase(unsafe.Pointer(c))
	var g Geom
	switch int(C.dGeomGetClass(c)) {
	case SphereClass:
		g = Sphere{base}
	case BoxClass:
		g = Box{base}
	case CapsuleClass:
		g = Capsule{base}
	case CylinderClass:
		g = Cylinder{base}
	case PlaneClass:
		g = Plane{base}
	case RayClass:
		g = Ray{base}
	case HeightfieldClass:
		g = Heightfield{base}
	case TriMeshClass:
		g = TriMesh{base}
	default:
		g = base
	}
	return g
}

func (g GeomBase) c() C.dGeomID {
	return C.dGeomID(unsafe.Pointer(g))
}

func (g GeomBase) ToSpace() Space {
	var space Space
	space = CToSpace(C.dSpaceID(unsafe.Pointer(g)))
	return space
}

// Destroy destroys the GeomBase.
func (g GeomBase) Destroy() {
	delete(geomData, g)
	C.dGeomDestroy(g.c())
}

// SetData associates user-specified data with the geometry.
func (g GeomBase) SetData(data interface{}) {
	geomData[g] = data
}

// Data returns the user-specified data associated with the geometry.
func (g GeomBase) Data() interface{} {
	return geomData[g]
}

// SetBody sets the associated body.
func (g GeomBase) SetBody(body Body) {
	C.dGeomSetBody(g.c(), body.c())
}

// Body returns the body associated with the geometry.
func (g GeomBase) Body() Body {
	return cToBody(C.dGeomGetBody(g.c()))
}

// SetPosition sets the position.
func (g GeomBase) SetPosition(pos Vector3) {
	C.dGeomSetPosition(g.c(), C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]))
}

// Position returns the position.
func (g GeomBase) Position() Vector3 {
	pos := NewVector3()
	C.dGeomCopyPosition(g.c(), (*C.dReal)(&pos[0]))
	return pos
}

// SetRotation sets the orientation represented by a rotation matrix.
func (g GeomBase) SetRotation(rot Matrix3) {
	C.dGeomSetRotation(g.c(), (*C.dReal)(&rot[0][0]))
}

// Rotation returns the orientation represented by a rotation matrix.
func (g GeomBase) Rotation() Matrix3 {
	rot := NewMatrix3()
	C.dGeomCopyRotation(g.c(), (*C.dReal)(&rot[0][0]))
	return rot
}

// SetQuaternion sets the orientation represented by a quaternion.
func (g GeomBase) SetQuaternion(quat Quaternion) {
	C.dGeomSetQuaternion(g.c(), (*C.dReal)(&quat[0]))
}

// Quaternion returns the orientation represented by a quaternion.
func (g GeomBase) Quaternion() Quaternion {
	quat := NewQuaternion()
	C.dGeomGetQuaternion(g.c(), (*C.dReal)(&quat[0]))
	return quat
}

// AABB returns the axis-aligned bounding box.
func (g GeomBase) AABB() AABB {
	aabb := NewAABB()
	C.dGeomGetAABB(g.c(), (*C.dReal)(&aabb[0]))
	return aabb
}

// IsSpace returns whether the geometry is a space.
func (g GeomBase) IsSpace() bool {
	return C.dGeomIsSpace(g.c()) != 0
}

// Space returns the containing space.
func (g GeomBase) Space() Space {
	return cToSpace(C.dGeomGetSpace(g.c()))
}

// Class returns the geometry class.
func (g GeomBase) Class() int {
	return int(C.dGeomGetClass(g.c()))
}

// SetCategoryBits sets the category bitfield.
func (g GeomBase) SetCategoryBits(bits int) {
	C.dGeomSetCategoryBits(g.c(), C.ulong(bits))
}

// CategoryBits returns the category bitfield.
func (g GeomBase) CategoryBits() int {
	return int(C.dGeomGetCategoryBits(g.c()))
}

// SetCollideBits sets the collide bitfield.
func (g GeomBase) SetCollideBits(bits int) {
	C.dGeomSetCollideBits(g.c(), C.ulong(bits))
}

// CollideBits returns the collide bitfield.
func (g GeomBase) CollideBits() int {
	return int(C.dGeomGetCollideBits(g.c()))
}

// SetEnabled sets whether the geometry is enabled.
func (g GeomBase) SetEnabled(isEnabled bool) {
	if isEnabled {
		C.dGeomEnable(g.c())
	} else {
		C.dGeomDisable(g.c())
	}
}

// Enabled returns whether the geometry is enabled.
func (g GeomBase) Enabled() bool {
	return bool(C.dGeomIsEnabled(g.c()) != 0)
}

// RelPointPos returns the position in world coordinates of a point in geometry
// coordinates.
func (g GeomBase) RelPointPos(pt Vector3) Vector3 {
	pos := NewVector3()
	C.dGeomGetRelPointPos(g.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2]), (*C.dReal)(&pos[0]))
	return pos
}

// PosRelPoint returns the position in geometry coordinates of a point in world
// coordinates.
func (g GeomBase) PosRelPoint(pos Vector3) Vector3 {
	pt := NewVector3()
	C.dGeomGetPosRelPoint(g.c(), C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]), (*C.dReal)(&pt[0]))
	return pt
}

// VectorToWorld converts a vector in geometry coordinates to world
// coordinates.
func (g GeomBase) VectorToWorld(vec Vector3) Vector3 {
	wld := NewVector3()
	C.dGeomVectorToWorld(g.c(), C.dReal(vec[0]), C.dReal(vec[1]), C.dReal(vec[2]), (*C.dReal)(&wld[0]))
	return wld
}

// VectorFromWorld converts a vector in world coordinates to geometry
// coordinates.
func (g GeomBase) VectorFromWorld(wld Vector3) Vector3 {
	vec := NewVector3()
	C.dGeomVectorFromWorld(g.c(), C.dReal(wld[0]), C.dReal(wld[1]), C.dReal(wld[2]), (*C.dReal)(&vec[0]))
	return vec
}

// SetOffsetPosition sets the position offset from the body.
func (g GeomBase) SetOffsetPosition(pos Vector3) {
	C.dGeomSetOffsetPosition(g.c(), C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]))
}

// OffsetPosition returns the position offset from the body.
func (g GeomBase) OffsetPosition() Vector3 {
	pos := NewVector3()
	C.dGeomCopyOffsetPosition(g.c(), (*C.dReal)(&pos[0]))
	return pos
}

// SetOffsetRotation sets the orientation offset from the body represented by a
// rotation matrix.
func (g GeomBase) SetOffsetRotation(rot Matrix3) {
	C.dGeomSetOffsetRotation(g.c(), (*C.dReal)(&rot[0][0]))
}

// OffsetRotation returns the orientation offset from the body represented by a
// rotation matrix.
func (g GeomBase) OffsetRotation() Matrix3 {
	rot := NewMatrix3()
	C.dGeomCopyOffsetRotation(g.c(), (*C.dReal)(&rot[0][0]))
	return rot
}

// SetOffsetQuaternion sets the offset from the body orientation represented by
// a quaternion.
func (g GeomBase) SetOffsetQuaternion(quat Quaternion) {
	C.dGeomSetOffsetQuaternion(g.c(), (*C.dReal)(&quat[0]))
}

// OffsetQuaternion returns the orientation offset from the body represented by
// a quaternion.
func (g GeomBase) OffsetQuaternion() Quaternion {
	quat := NewQuaternion()
	C.dGeomGetOffsetQuaternion(g.c(), (*C.dReal)(&quat[0]))
	return quat
}

// SetOffsetWorldPosition sets the offset to the body position such that the
// geom's world position is pos.
func (g GeomBase) SetOffsetWorldPosition(pos Vector3) {
	C.dGeomSetOffsetWorldPosition(g.c(), C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]))
}

// SetOffsetWorldRotation sets the offset to the body orientation such that the
// geom's world orientation is represented by the matrix rot.
func (g GeomBase) SetOffsetWorldRotation(rot Matrix3) {
	C.dGeomSetOffsetWorldRotation(g.c(), (*C.dReal)(&rot[0][0]))
}

// SetOffsetWorldQuaternion sets the offset to the body orientation such that
// the geom's world orientation is represented by the quaternion quat.
func (g GeomBase) SetOffsetWorldQuaternion(quat Quaternion) {
	C.dGeomSetOffsetWorldQuaternion(g.c(), (*C.dReal)(&quat[0]))
}

// ClearOffset removes the body offset.
func (g GeomBase) ClearOffset() {
	C.dGeomClearOffset(g.c())
}

// IsOffset returns whether a body offset has been created.
func (g GeomBase) IsOffset() bool {
	return C.dGeomIsOffset(g.c()) != 0
}

// Collide tests for collision with the given geometry and returns a list of
// contact points.
func (g GeomBase) Collide(other Geom, maxContacts uint16, flags int) []ContactGeom {
	cts := make([]C.dContactGeom, maxContacts)
	numCts := int(C.dCollide(g.c(), other.c(), C.int(int(maxContacts)|flags), &cts[0],
		C.int(unsafe.Sizeof(cts[0]))))
	contacts := make([]ContactGeom, numCts)
	for i := range contacts {
		contacts[i] = *NewContactGeom()
		contacts[i].fromC(&cts[i])
	}
	return contacts
}

// Collide2 tests for collision with the given geometry, applying cb for each
// contact.
func (g GeomBase) Collide2(other Geom, data interface{}, cb NearCallback) {
	cbData := &nearCallbackData{fn: cb, data: data}
	C.dSpaceCollide2(g.c(), other.c(), unsafe.Pointer(cbData),
		(*C.dNearCallback)(C.callNearCallback))
}

// Next returns the next geometry.
func (g GeomBase) Next() Geom {
	return cToGeom(C.dBodyGetNextGeom(g.c()))
}

// Sphere is a geometry representing a sphere.
type Sphere struct {
	GeomBase
}

// SetRadius sets the radius.
func (s Sphere) SetRadius(radius float64) {
	C.dGeomSphereSetRadius(s.c(), C.dReal(radius))
}

// Radius returns the radius.
func (s Sphere) Radius() float64 {
	return float64(C.dGeomSphereGetRadius(s.c()))
}

// SpherePointDepth returns the depth of the given point.
func (s Sphere) SpherePointDepth(pt Vector3) float64 {
	return float64(C.dGeomSpherePointDepth(s.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2])))
}

// Box is a geometry representing a rectangular box.
type Box struct {
	GeomBase
}

// SetLengths sets the lengths of the sides.
func (b Box) SetLengths(lens Vector3) {
	C.dGeomBoxSetLengths(b.c(), C.dReal(lens[0]), C.dReal(lens[1]), C.dReal(lens[2]))
}

// Lengths returns the lengths of the sides.
func (b Box) Lengths() Vector3 {
	lens := NewVector3()
	C.dGeomBoxGetLengths(b.c(), (*C.dReal)(&lens[0]))
	return lens
}

// PointDepth returns the depth of the given point.
func (b Box) PointDepth(pt Vector3) float64 {
	return float64(C.dGeomBoxPointDepth(b.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2])))
}

// Plane is a geometry that represents a plane.
type Plane struct {
	GeomBase
}

// SetParams sets plane parameters.
func (p Plane) SetParams(params Vector4) {
	C.dGeomPlaneSetParams(p.c(), C.dReal(params[0]), C.dReal(params[1]), C.dReal(params[2]), C.dReal(params[3]))
}

// Params returns plane parameters.
func (p Plane) Params() Vector4 {
	params := NewVector4()
	C.dGeomPlaneGetParams(p.c(), (*C.dReal)(&params[0]))
	return params
}

// PointDepth returns the depth of the given point.
func (p Plane) PointDepth(pt Vector3) float64 {
	return float64(C.dGeomPlanePointDepth(p.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2])))
}

// Capsule is a geometry that represents a capsule (a cylinder with a half
// sphere on each end).
type Capsule struct {
	GeomBase
}

// SetParams sets the radius and length.
func (c Capsule) SetParams(radius, length float64) {
	C.dGeomCapsuleSetParams(c.c(), C.dReal(radius), C.dReal(length))
}

// Params returns the radius and length.
func (c Capsule) Params() (float64, float64) {
	var radius, length float64
	C.dGeomCapsuleGetParams(c.c(), (*C.dReal)(&radius), (*C.dReal)(&length))
	return radius, length
}

// PointDepth returns the depth of the given point.
func (c Capsule) PointDepth(pt Vector3) float64 {
	return float64(C.dGeomCapsulePointDepth(c.c(), C.dReal(pt[0]), C.dReal(pt[1]), C.dReal(pt[2])))
}

// Cylinder is a geometry that represents a cylider.
type Cylinder struct {
	GeomBase
}

// SetParams sets the radius and length.
func (c Cylinder) SetParams(radius, length float64) {
	C.dGeomCylinderSetParams(c.c(), C.dReal(radius), C.dReal(length))
}

// Params returns the radius and length.
func (c Cylinder) Params() (float64, float64) {
	var radius, length float64
	C.dGeomCylinderGetParams(c.c(), (*C.dReal)(&radius), (*C.dReal)(&length))
	return radius, length
}

// Ray is a geometry representing a ray.
type Ray struct {
	GeomBase
}

// SetPosDir sets the position and direction.
func (r Ray) SetPosDir(pos, dir Vector3) {
	C.dGeomRaySet(r.c(), C.dReal(pos[0]), C.dReal(pos[1]), C.dReal(pos[2]),
		C.dReal(dir[0]), C.dReal(dir[1]), C.dReal(dir[2]))
}

// PosDir returns the position and direction.
func (r Ray) PosDir() (Vector3, Vector3) {
	pos, dir := NewVector3(), NewVector3()
	C.dGeomRayGet(r.c(), (*C.dReal)(&pos[0]), (*C.dReal)(&dir[0]))
	return pos, dir
}

// SetLength sets the length.
func (r Ray) SetLength(length float64) {
	C.dGeomRaySetLength(r.c(), C.dReal(length))
}

// Length returns the length.
func (r Ray) Length() float64 {
	return float64(C.dGeomRayGetLength(r.c()))
}

// SetFirstContact sets whether to stop collision detection after finding the
// first contact point.
func (r Ray) SetFirstContact(firstContact bool) {
	C.dGeomRaySetFirstContact(r.c(), C.int(btoi(firstContact)))
}

// FirstContact returns whether collision detection will stop after finding the
// first contact.
func (r Ray) FirstContact() bool {
	return C.dGeomRayGetFirstContact(r.c()) != 0
}

// SetBackfaceCull sets whether backface culling is enabled.
func (r Ray) SetBackfaceCull(backfaceCull bool) {
	C.dGeomRaySetBackfaceCull(r.c(), C.int(btoi(backfaceCull)))
}

// BackfaceCull returns whether backface culling is enabled.
func (r Ray) BackfaceCull() bool {
	return C.dGeomRayGetBackfaceCull(r.c()) != 0
}

// SetClosestHit sets whether to only report the closest hit.
func (r Ray) SetClosestHit(closestHit bool) {
	C.dGeomRaySetClosestHit(r.c(), C.int(btoi(closestHit)))
}

// ClosestHit returns whether only the closest hit will be reported.
func (r Ray) ClosestHit() bool {
	return C.dGeomRayGetClosestHit(r.c()) != 0
}

// Convex is a geometry representing a convex object.
type Convex struct {
	GeomBase
}

// Set sets convex object data
func (c Convex) Set(planes PlaneList, pts VertexList, polyList PolygonList) {
	C.dGeomSetConvex(c.c(), (*C.dReal)(&planes[0][0]), C.uint(len(planes)),
		(*C.dReal)(&pts[0][0]), C.uint(len(pts)), &polyList[0])
}
