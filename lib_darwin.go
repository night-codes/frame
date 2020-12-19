// +build darwin

package frame

/*
#import  "c_darwin.h"

#ifndef WEBVIEW_COCOA
#define WEBVIEW_COCOA
#endif
*/
import "C"

import (
	"fmt"
	"strings"
	"sync/atomic"
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
	Fullscreen bool
	Focused    bool
}

//export goAppActivated
func goAppActivated(ret C.AppMenu) {
	app := &App{
		MainMenu: &Menu{
			menu: ret.mainMenu,
		},
		AppMenu: &Menu{
			menu: ret.appMenu,
		},
	}
	app.MainMenu.app = app
	app.AppMenu.app = app
	appChan <- app
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

func stateSender(win *Window, newState State) {
	oldState := win.state
	win.state = newState
	if win.StateEvent != nil {
		win.StateEvent(newState)
	}

	if !newState.Hidden && oldState.Hidden {
	}
}

//export goWindowEvent
func goWindowEvent(windowID C.int, eventTitle *C.char, x C.int, y C.int, w C.int, h C.int) {
	id := int(windowID)
	title := C.GoString(eventTitle)

	var win *Window
	for i := range winds {
		if int(winds[i].id) == id {
			win = winds[i]
		}
	}

	if win == nil {
		return
	}

	if win.Invoke != nil && strings.HasPrefix(title, "invoke:") {
		win.Invoke(strings.TrimPrefix(title, "invoke:"))
		return
	}

	state := win.state
	switch title {
	case "windowDidBecomeKey":
		state.Focused = true
	case "windowDidResignKey":
		state.Focused = false
	case "windowWillClose":
		state.Hidden = true
	case "windowDidExpose":
		state.Hidden = false
	}

	if state.Focused != win.state.Focused {
		go stateSender(win, state)
	}
}

func cRequest(fn func(id uint64)) interface{} {
	id := atomic.AddUint64(&goRequestID, 1)
	ch := make(chan interface{})
	goRequests.Store(id, ch)
	defer goRequests.Delete(id)
	fn(id)
	return <-ch
}

//export goScriptEvent
func goScriptEvent() {
	fmt.Println("js...")
}
