// +build freebsd linux netbsd openbsd solaris

package frame

/*
#include "c_linux.h"
*/
import "C"

import (
	"fmt"
	"reflect"
)

const (
	// Withdrawn - window is not shown
	cWithdrawn = 1 << iota
	cIconified
	cMaximized
	cSticky
	cFullscreen
	cAbove
	cBelow
	cFocused
	cTiled
)

// State struct
type State struct {
	Hidden     bool
	Iconified  bool
	Maximized  bool
	Sticky     bool
	Fullscreen bool
	Focused    bool
	Tiled      bool
}

//export goAppActivated
func goAppActivated() {
	appChan <- &App{}
}

//export goPrint
func goPrint(text *C.char) {
	fmt.Println(C.GoString(text))
}

//export goPrintInt
func goPrintInt(t C.int) {
	fmt.Println(int(t))
}

//export goWindowState
func goWindowState(c *C.GtkWidget, e C.int) {
	for i := range frames {
		if frames[i].StateEvent != nil && reflect.DeepEqual(frames[i].window, c) && (uint32(e)&cWithdrawn == 0 || uint32(e)&cFocused == 0) {
			frames[i].StateEvent(State{
				Hidden:     uint32(e)&cWithdrawn != 0,
				Iconified:  uint32(e)&cIconified != 0,
				Maximized:  uint32(e)&cMaximized != 0,
				Sticky:     uint32(e)&cSticky != 0,
				Fullscreen: uint32(e)&cFullscreen != 0,
				Focused:    uint32(e)&cFocused != 0,
				Tiled:      uint32(e)&cTiled != 0,
			})
		}
	}
}

//export goScriptEvent
func goScriptEvent() {
	fmt.Println("js...")
}

func gboolean(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func goBool(b C.gboolean) bool {
	if b != 0 {
		return true
	}
	return false
}
