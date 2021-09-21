// +build darwin

package frame

/*
#ifndef WEBVIEW_COCOA
#define WEBVIEW_COCOA
#endif

#import  "c_darwin.h"
*/
import "C"

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	// "reflect"
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
	appChan <- &App{
		openedWns: sync.WaitGroup{},
		shown:     make(chan bool),
		mainMenu: &Menu{
			menu: ret.mainMenu,
		},
		appMenu: &Menu{
			menu: ret.appMenu,
		},
	}
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
	return C.bool(true) == b
}

func stateSender(win *Window, newState State) {
	prevState := win.state
	win.state = newState
	if win.StateEvent != nil {
		win.StateEvent(newState)
	}

	if prevState.Hidden && !newState.Hidden {
		win.app.openedWns.Add(1)
		select {
		case win.app.shown <- true:
		default:
		}
	} else if !prevState.Hidden && newState.Hidden {
		win.app.openedWns.Done()
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
		state.Hidden = false
	case "windowDidResignKey":
		state.Focused = false
	}

	if !state.Hidden && (state.Focused != win.state.Focused) {
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
