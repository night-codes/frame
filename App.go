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
	App struct {
		count uint
	}
)

// SetDefaultIconFromFile for application windows
func (a *App) SetDefaultIconFromFile(filename string) {
	C.gtk_window_set_default_icon_from_file(C.gcharptr(C.CString(filename)), nil)
}

// SetDefaultIconName for application windows
func (a *App) SetDefaultIconName(name string) {
	C.gtk_window_set_default_icon_name(C.gcharptr(C.CString(name)))
}
