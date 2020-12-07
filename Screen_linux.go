// +build freebsd linux netbsd openbsd solaris
package frame

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.0
#cgo linux CFLAGS: -DLINUX -Wno-deprecated-declarations
#cgo linux LDFLAGS: -lX11
#include "unix.h"
*/
import "C"

type (
	// Screen struct
	Screen struct {
		screen  *C.GdkScreen
		display *C.GdkDisplay
		monitor *C.GdkMonitor
	}
)

// Size of monitor
func (s *Screen) Size() (width, height int) {
	geometry := C.GdkRectangle{}
	C.gdk_monitor_get_geometry(s.monitor, &geometry)
	width, height = int(geometry.width), int(geometry.height)
	return
}

// ScaleFactor of monitor
func (s *Screen) ScaleFactor() int {
	return int(C.gdk_monitor_get_scale_factor(s.monitor))
}

// gdk_monitor_get_scale_factor
// gdk_screen_get_display
// gdk_monitor_get_geometry
// gdk_display_get_monitor_at_window
