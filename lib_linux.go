// +build freebsd linux netbsd openbsd solaris

package frame

/*
#ifndef WEBVIEW_GTK
#define WEBVIEW_GTK
#endif

#include "c_linux.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type (
	// State struct
	State struct {
		Hidden     bool
		Iconified  bool
		Maximized  bool
		Fullscreen bool
		Focused    bool
	}
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

var (
	goRequestID uint64
	goRequests  sync.Map
	mutexNew    sync.Mutex
	winds       = []*Window{}
	lock        sync.Mutex
	appChan     = make(chan *App)
	idItr       int64
)

//export goAppActivated
func goAppActivated(app *C.GtkApplication) {
	appChan <- &App{
		app:       app,
		openedWns: sync.WaitGroup{},
		shown:     make(chan bool),
	}
}

//export callTest
func callTest(p C.gpointer) C.gboolean {
	fmt.Println(p)
	return gboolean(false)
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
func goWindowState(winObj *C.WindowObj, e C.int) {
	go func() {
		for _, win := range winds {
			if int(win.id) == int(winObj.id) && (uint32(e)&cWithdrawn == 0 || uint32(e)&cFocused == 0) {
				prevState := win.state
				newState := State{
					Hidden:     uint32(e)&cWithdrawn != 0,
					Iconified:  uint32(e)&cIconified != 0,
					Maximized:  uint32(e)&cMaximized != 0,
					Fullscreen: uint32(e)&cFullscreen != 0,
					Focused:    uint32(e)&cFocused != 0,
				}
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
				if newState.Maximized && win.maxLimited {
					win.Maximize(false).SetSize(win.maxWidth, win.maxHeight)
				}
			}
		}
	}()
}

//export goInvokeCallback
func goInvokeCallback(win *C.WindowObj, data *C.char) {
	go func() {
		for i := range winds {
			if winds[i].Invoke != nil && int(winds[i].id) == int(win.id) {
				winds[i].Invoke(C.GoString(data))
			}
		}
	}()
}

//export goEvalRet
func goEvalRet(reqid C.ulonglong, err *C.char) {
	go func() {
		if chi, ok := goRequests.Load(uint64(reqid)); ok {
			if ch, ok := chi.(chan interface{}); ok {
				ch <- C.GoString(err)
			}
		}
	}()
}

//export goWinRet
func goWinRet(reqid C.ulonglong, win *C.WindowObj) {
	go func() {
		if chi, ok := goRequests.Load(uint64(reqid)); ok {
			if ch, ok := chi.(chan interface{}); ok {
				ch <- win
			}
		}
	}()
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

func cRequest(fn func(id uint64)) interface{} {
	id := atomic.AddUint64(&goRequestID, 1)
	ch := make(chan interface{})
	goRequests.Store(id, ch)
	defer goRequests.Delete(id)
	fn(id)
	return <-ch
}
