#include <ode/ode.h>
#include <stdio.h>
#include "_cgo_export.h"

void callNearCallback(void *data, dGeomID obj1, dGeomID obj2) {
	nearCallback(data, obj1, obj2);
}

void callMovedCallback(dBodyID body) {
	movedCallback(body);
}

int callTriCallback(dGeomID mesh, dGeomID other, int index) {
	return triCallback(mesh, other, index);
}

int callTriRayCallback(dGeomID mesh, dGeomID ray, int index, dReal u, dReal v) {
	return triRayCallback(mesh, ray, index, u, v);
}
