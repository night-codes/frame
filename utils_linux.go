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
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
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

var (
	goRequestID uint64
	goRequests  sync.Map
)

//export goAppActivated
func goAppActivated() {
	appChan <- &App{}
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
func goWindowState(c *C.GtkWidget, e C.int) {
	for i := range frames {
		if frames[i].StateEvent != nil && reflect.DeepEqual(frames[i].window, c) && (uint32(e)&cWithdrawn == 0 || uint32(e)&cFocused == 0) {
			go frames[i].StateEvent(State{
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

//export goInvokeCallback
func goInvokeCallback(webview *C.GtkWidget, data *C.char) {
	for i := range frames {
		if frames[i].Invoke != nil && reflect.DeepEqual(frames[i].webview, webview) {
			go frames[i].Invoke(C.GoString(data))
		}
	}
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

func threadsAddIdle(function C.GSourceFunc, data unsafe.Pointer) bool {
	var ret C.guint
	ret = C.gdk_threads_add_idle(function, (C.gpointer)(data))
	if uint(ret) != 0 {
		return false
	}
	return true
}

func cRequest(fn func(id uint64)) interface{} {
	id := atomic.AddUint64(&goRequestID, 1)
	ch := make(chan interface{})
	goRequests.Store(id, ch)
	defer goRequests.Delete(id)
	go fn(id)
	return <-ch
}
