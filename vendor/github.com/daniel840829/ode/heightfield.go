package ode

// #include <ode/ode.h>
import "C"

import (
	"unsafe"
)

// HeightfieldData represents heightfield data.
type HeightfieldData uintptr

func cToHeightfieldData(c C.dHeightfieldDataID) HeightfieldData {
	return HeightfieldData(unsafe.Pointer(c))
}

func (h HeightfieldData) c() C.dHeightfieldDataID {
	return C.dHeightfieldDataID(unsafe.Pointer(h))
}

// NewHeightfieldData returns a new HeightfieldData instance.
func NewHeightfieldData() HeightfieldData {
	return cToHeightfieldData(C.dGeomHeightfieldDataCreate())
}

// Destroy destroys the heightfield data.
func (h *HeightfieldData) Destroy() {
	C.dGeomHeightfieldDataDestroy(h.c())
}

// Heightfield is a geometry representing a heightfield.
type Heightfield struct {
	GeomBase
}

// Build builds a heightfield data set.
func (h Heightfield) Build(data HeightfieldData, heightSamples Matrix,
	width, depth, scale, offset, thickness float64, doWrap bool) {

	numWidthSamp, numDepthSamp := len(heightSamples), 0
	var heightSamplesPtr *C.double
	if numDepthSamp > 0 {
		numWidthSamp = len(heightSamples[0])
		if numWidthSamp > 0 {
			heightSamplesPtr = (*C.double)(&heightSamples[0][0])
		}
	}
	C.dGeomHeightfieldDataBuildDouble(data.c(), heightSamplesPtr, 1,
		C.dReal(width), C.dReal(depth), C.int(numWidthSamp), C.int(numDepthSamp),
		C.dReal(scale), C.dReal(offset), C.dReal(thickness), C.int(btoi(doWrap)))
}

// SetBounds sets the minimum and maximum height.
func (h Heightfield) SetBounds(data HeightfieldData, minHeight, maxHeight float64) {
	C.dGeomHeightfieldDataSetBounds(data.c(), C.dReal(minHeight), C.dReal(maxHeight))
}

// SetHeightfieldData associates a data set to the heightfield.
func (h Heightfield) SetHeightfieldData(data HeightfieldData) {
	C.dGeomHeightfieldSetHeightfieldData(h.c(), data.c())
}

// HeightfieldData returns the data set associated with the heightfield.
func (h Heightfield) HeightfieldData() HeightfieldData {
	return cToHeightfieldData(C.dGeomHeightfieldGetHeightfieldData(h.c()))
}
