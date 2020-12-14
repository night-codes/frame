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
import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type (
	// App is main application object
	App struct {
		count  uint
		frames []*Frame
	}
)

var (
	mutexNew sync.Mutex
	frames   = []*Frame{}
	lock     sync.Mutex
	appChan  = make(chan *App)
	idItr    uint64
)

// MakeApp is make and run one instance of application (At the moment, it is possible to create only one instance)
func MakeApp(count ...uint) *App {
	var c uint
	if len(count) > 0 {
		c = count[0]
	}
	lock.Lock()
	go func() {
		runtime.LockOSThread()
		C.makeApp(C.int(c))
		runtime.UnlockOSThread()
	}()
	return <-appChan
}

// SetDefaultIconFromFile for application windows
func (a *App) SetDefaultIconFromFile(filename string) {
	C.gtk_window_set_default_icon_from_file(C.gcharptr(C.CString(filename)), nil)
}

// SetDefaultIconName for application windows
func (a *App) SetDefaultIconName(name string) {
	C.gtk_window_set_default_icon_name(C.gcharptr(C.CString(name)))
}

// NewFrame returns window with webview
func (a *App) NewFrame(title string, sizes ...int) *Frame {
	id := atomic.AddUint64(&idItr, 1)
	mutexNew.Lock()
	defer func() {
		time.Sleep(time.Second / 100)
		fmt.Println(id, "|NewFrame - END|")
		mutexNew.Unlock()
	}()

	fmt.Println(id, "|NewFrame|")

	width := 400
	height := 300

	if len(sizes) > 0 {
		width = sizes[0]
	}

	if len(sizes) > 1 {
		height = sizes[1]
	}

	cRet := cRequest(func(id uint64) {
		C.makeWindow(&C.idleData{
			content: C.gcharptr(C.CString(title)),
			width:   C.int(width),
			height:  C.int(height),
			req_id:  C.ulonglong(id),
		})
	})
	ret, _ := cRet.(*C.WindowObj)
	frame := &Frame{
		id:      id,
		window:  ret.window,
		box:     ret.box,
		webview: ret.webview,
		menubar: ret.menubar,
	}
	frame.SetCenter()
	frames = append(frames, frame)

	return frame
}
