package frame

/*
#cgo pkg-config: webkit2gtk-4.0
#include <webkit2/webkit2.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include "my.h"
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

func (s *Screen) Size() (width, height int) {
	geometry := C.GdkRectangle{}
	C.gdk_monitor_get_geometry(s.monitor, &geometry)
	width, height = int(geometry.width), int(geometry.height)
	return
}

func (s *Screen) GetScale() int {
	return int(C.gdk_monitor_get_scale_factor(s.monitor))
}

// gdk_monitor_get_scale_factor
// gdk_screen_get_display
// gdk_monitor_get_geometry
// gdk_display_get_monitor_at_window
