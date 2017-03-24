package frame

/*
#cgo pkg-config: webkit2gtk-4.0
#cgo linux CFLAGS: -DLINUX -Wno-deprecated-declarations
#cgo linux LDFLAGS: -lX11
#include <webkit2/webkit2.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include "my.h"
*/
import "C"
import (
	"runtime"
)

func init() {
	runtime.LockOSThread()
}
