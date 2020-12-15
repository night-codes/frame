// +build darwin

package frame

/*
#import  "c_darwin.h"
*/
import "C"

import (
	"fmt"
	"strings"
	// "reflect"
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

func goBool(b C.BOOL) bool {
	if b != 0 {
		return true
	}
	return false
}

//export goWindowEvent
func goWindowEvent(windowID C.int, eventTitle *C.char, x C.int, y C.int, w C.int, h C.int) {
	id := int(windowID)
	title := C.GoString(eventTitle)
	// fmt.Println(id, title)
	if len(frames) > id {
		if frames[id].Invoke != nil && strings.HasPrefix(title, "invoke:") {
			frames[id].Invoke(strings.TrimPrefix(title, "invoke:"))
		}

		if frames[id].StateEvent != nil {
			state := frames[id].state
			switch title {
			case "windowDidDeminiaturize":
				frames[id].state.Iconified = false
			case "windowWillMiniaturize":
				frames[id].state.Iconified = true
			case "windowDidBecomeKey":
				for i := range frames {
					frames[i].state.Focused = false
				}
				frames[id].state.Focused = true
			case "windowWillClose":
				frames[id].state.Hidden = true
			case "windowShouldClose":
				frames[id].state.Hidden = true
			}

			if state.Focused != frames[id].state.Focused || state.Iconified != frames[id].state.Iconified || state.Hidden != frames[id].state.Hidden {
				frames[id].StateEvent(frames[id].state)
			}
			if !state.Hidden && frames[id].state.Hidden {
				if frames[id].modalFor != -1 && frames[id].modal == -1 {
					modalFor := frames[id].modalFor
					frames[id].UnsetModal()
					frames[modalFor].Show()
					if frames[modalFor].modalFor != -1 {
						frames[modalFor].SetModal(frames[frames[modalFor].modalFor])
					}
				}
			}
		}
	}
}

//export goScriptEvent
func goScriptEvent() {
	fmt.Println("js...")
}
