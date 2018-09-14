ode
===

This is a Go binding for the Open Dynamics Engine 3D physics library.  It
sticks fairly closely to the development version of the ODE C API, with a few
stylistic and idiomatic changes thrown in here and there where it seemed
useful.

Get ODE [here](http://bitbucket.org/odedevs/ode/).

ODE must be compiled as a shared library with double precision support.
Triangle mesh indices are expected to be 32 bit, which is the ODE default.  The
following will configure ODE with these options:

`> cd /path/to/ode-src; ./configure --enable-double-precision --enable-shared`
