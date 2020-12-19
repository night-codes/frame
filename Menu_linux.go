// +build freebsd linux netbsd openbsd solaris

package frame

/*
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
