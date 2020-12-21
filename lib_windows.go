// +build windows

package frame

/*
#ifndef WEBVIEW_WINAPI
#define WEBVIEW_WINAPI
#endif

#include "c_windows.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"sync/atomic"
)

// State struct
type State struct {
	Hidden     bool
	Iconified  bool
	Maximized  bool
	Fullscreen bool
	Focused    bool
}

var (
	goRequestID uint64
	goRequests  sync.Map
)

func cRequest(fn func(id uint64)) interface{} {
	id := atomic.AddUint64(&goRequestID, 1)
	ch := make(chan interface{})
	goRequests.Store(id, ch)
	defer goRequests.Delete(id)
	fn(id)
	return <-ch
}

//export goWinRet
func goWinRet(reqid C.ulonglong, win *C.WindowObj) {
	// fmt.Printf("%+v\n", win.window)
	go func() {
		if chi, ok := goRequests.Load(uint64(reqid)); ok {
			if ch, ok := chi.(chan interface{}); ok {
				ch <- win
			}
		}
	}()
}

//export goPrint
func goPrint(text *C.char) {
	fmt.Println(C.GoString(text))
}

//export goPrintInt
func goPrintInt(text *C.char, t C.int) {
	fmt.Println(C.GoString(text), int(t))
}
