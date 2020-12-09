// +build freebsd linux netbsd openbsd solaris

package frame

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.0
#cgo linux CFLAGS: -DLINUX -DWEBVIEW_GTK=1 -Wno-deprecated-declarations
#cgo linux LDFLAGS: -lX11
#include "linux.h"
*/
import "C"
import (
	"runtime"
	"sync"
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
func (app *App) NewFrame(title string, sizes ...int) *Frame {
	mutexNew.Lock()
	defer mutexNew.Unlock()
	width := 400
	height := 300

	if len(sizes) > 0 {
		width = sizes[0]
	}

	if len(sizes) > 1 {
		height = sizes[1]
	}

	window := C.makeWindow(C.CString(title), C.int(width), C.int(height))
	box := C.makeBox(window)
	menubar := C.makeMenubar(box)
	webview := C.makeWebview(box)
	frame := &Frame{
		window:  window,
		box:     box,
		webview: webview,
		menubar: menubar,
	}
	frame.SetPosition(PosCenter)
	frames = append(frames, frame)
	return frame
}
