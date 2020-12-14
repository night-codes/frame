// +build freebsd linux netbsd openbsd solaris

package frame

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.0
#cgo linux CFLAGS: -DLINUX -DWEBVIEW_GTK=1 -Wno-deprecated-declarations
#cgo linux LDFLAGS: -lX11

#ifndef WEBVIEW_GTK
#define WEBVIEW_GTK
#endif

#include "c_linux.h"
*/
import "C"

type (
	// Menu of window
	Menu struct {
	}

	// MenuItem element
	MenuItem struct {
		menu Menu
	}
)

// NewIem returns window with webview
func (m *Menu) NewIem(title string) MenuItem {
	menuitem := MenuItem{
		//
	}
	return menuitem
}
