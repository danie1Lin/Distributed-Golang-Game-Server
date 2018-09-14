// Package ode provide a Go interface to the Open Dynamics Engine library.
// See the ODE documentation for more information.
package ode

// #cgo LDFLAGS: -lode
// #include <ode/ode.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Initialization flags
const (
	ManualThreadCleanupIFlag = C.dInitFlagManualThreadCleanup
)

// Allocation flags
const (
	BasicDataAFlag     = C.dAllocateFlagBasicData
	CollisionDataAFlag = C.dAllocateFlagCollisionData
	AllAFlag           = C.dAllocateMaskAll
)

// Short constructor aliases for convenience
var (
	V3 = NewVector3
	V4 = NewVector4
	M3 = NewMatrix3
	M4 = NewMatrix4
	Q  = NewQuaternion
	BB = NewAABB
)

// NearCallback is a callback type for handling potential object collisions.
type NearCallback func(data interface{}, obj1, obj2 Geom)

type nearCallbackData struct {
	data interface{}
	fn   NearCallback
}

//export nearCallback
func nearCallback(data unsafe.Pointer, obj1, obj2 C.dGeomID) {
	cbData := (*nearCallbackData)(data)
	cbData.fn(cbData.data, cToGeom(obj1), cToGeom(obj2))
}

// round num up to nearest multiple of align
func alignNum(num, align int) int {
	return (num + (align - 1)) &^ (align - 1)
}

// Vector represents a double precision vector.
type Vector []float64

// NewVector returns a new Vector instance.
func NewVector(size, align int, vals ...float64) Vector {
	alignSize := alignNum(size, align)
	v := make(Vector, size, alignSize)
	copy(v, vals)
	return v
}

func (v Vector) convertC(c *C.dReal, toC bool) {
	for i := range v {
		if toC {
			*c = C.dReal(v[i])
		} else {
			v[i] = float64(*c)
		}
		c = (*C.dReal)(unsafe.Pointer(uintptr(unsafe.Pointer(c)) + unsafe.Sizeof(*c)))
	}
}

func (v Vector) toC(c *C.dReal) {
	v.convertC(c, true)
}

func (v Vector) fromC(c *C.dReal) {
	v.convertC(c, false)
}

// Vector3 represents a 3 component vector.
type Vector3 Vector

func cToVector3(a *C.dReal) Vector3 {
	vec := NewVector3()
	Vector(vec).fromC(a)
	return vec
}

// NewVector3 returns a new Vector3 instance.
func NewVector3(vals ...float64) Vector3 {
	return Vector3(NewVector(3, 4, vals...))
}

// Vector4 represents a 4 component vector.
type Vector4 Vector

// NewVector4 returns a new Vector4 instance.
func NewVector4(vals ...float64) Vector4 {
	return Vector4(NewVector(4, 4, vals...))
}

// Quaternion represents a quaternion.
type Quaternion Vector

// NewQuaternion returns a new Quaternion instance.
func NewQuaternion(vals ...float64) Quaternion {
	return Quaternion(NewVector(4, 1, vals...))
}

// AABB represents an axis-aligned bounding box.
type AABB Vector

// NewAABB returns a new AABB instance.
func NewAABB(vals ...float64) AABB {
	return AABB(NewVector(6, 1, vals...))
}

// Matrix represents a double precision matrix.
type Matrix [][]float64

// NewVector returns a new Matrix instance.
func NewMatrix(numRows, numCols, align int, vals ...float64) Matrix {
	mat := make(Matrix, numRows)
	numAlignCols := alignNum(numCols, align)
	elts := make([]float64, numAlignCols*numRows)
	for i := range mat {
		mat[i], elts = elts[:numCols:numAlignCols], elts[numAlignCols:]
		n := numCols
		if len(vals) < numCols {
			n = len(vals)
		}
		copy(mat[i], vals[:n])
		vals = vals[n:]
	}
	return mat
}

func (m Matrix) convertC(c *C.dReal, toC bool) {
	for i := range m {
		for j := 0; j < cap(m[i]); j++ {
			if j < len(m[i]) {
				if toC {
					*c = C.dReal(m[i][j])
				} else {
					m[i][j] = float64(*c)
				}
			}
			c = (*C.dReal)(unsafe.Pointer(uintptr(unsafe.Pointer(c)) + unsafe.Sizeof(*c)))
		}
	}
}

func (m Matrix) toC(c *C.dReal) {
	m.convertC(c, true)
}

func (m Matrix) fromC(c *C.dReal) {
	m.convertC(c, false)
}

// Matrix3 represents a 3x3 matrix.
type Matrix3 Matrix

// NewMatrix3 returns a new Matrix3 instance.
func NewMatrix3(vals ...float64) Matrix3 {
	return Matrix3(NewMatrix(3, 3, 4, vals...))
}

// Matrix4 represents a 4x4 matrix.
type Matrix4 Matrix

// NewMatrix4 returns a new Matrix4 instance.
func NewMatrix4(vals ...float64) Matrix4 {
	return Matrix4(NewMatrix(4, 4, 4, vals...))
}

// VertexList represents a list of 3D vertices.
type VertexList Matrix

// NewVertexList returns a new VertexList instance.
func NewVertexList(size int, vals ...float64) VertexList {
	return VertexList(NewMatrix(size, 3, 1, vals...))
}

// PlaneList represents a list of plane definitions.
type PlaneList Matrix

// NewPlaneList returns a new PlaneList instance.
func NewPlaneList(size int, vals ...float64) PlaneList {
	return PlaneList(NewMatrix(size, 4, 1, vals...))
}

// TriVertexIndexList represents a list of triangle vertex indices.
type TriVertexIndexList [][]uint32

// NewTriVertexIndexList returns a new TriVertexIndexList instance.
func NewTriVertexIndexList(size int, indices ...uint32) TriVertexIndexList {
	list := make(TriVertexIndexList, size)
	elts := make([]uint32, 3*size)
	for i := range list {
		list[i], elts = elts[:3], elts[3:]
		n := 3
		if len(indices) < 3 {
			n = len(indices)
		}
		copy(list[i], indices[:n])
		indices = indices[n:]
	}
	return list
}

// PolygonList represents a list of polygon definitions
type PolygonList []C.uint

// Init initializes ODE.
func Init(initFlags, allocFlags int) {
	C.dInitODE2(C.uint(initFlags))
	C.dAllocateODEDataForThread(C.uint(allocFlags))
	fmt.Println("Daniel's Version ODE")
}

// Close releases ODE resources.
func Close() {
	C.dCloseODE()
}

// CleanupAllDataForThread manually releases ODE resources for the current thread.
func CleanupAllDataForThread() {
	C.dCleanupODEAllDataForThread()
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
