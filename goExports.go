package frame

/*
#cgo pkg-config: webkit2gtk-4.0
#include <webkit2/webkit2.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include "my.h"
*/
import "C"
import (
	"reflect"
	"runtime"
	"sync"
)

var (
	lock     sync.Mutex
	mutexNew sync.Mutex
	frames   = []*Frame{}
	appChan  = make(chan *App)
)

//export goAppActivated
func goAppActivated() {
	appChan <- &App{}
}

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

func (app *App) New(title string, sizes ...int) *Frame {
	mutexNew.Lock()
	defer mutexNew.Unlock()
	width := 200
	height := 100

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
	return &Frame{
		Window:  window,
		Box:     box,
		Webview: webview,
		Menubar: menubar,
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

//export goWindowState
func goWindowState(c *C.GtkWidget, e C.int) {
	for i := range frames {
		if reflect.DeepEqual(frames[i].Window, c) {
			frames[i].StateEvent(State{
				Hidden:     uint32(e)&cWithdrawn != 0,
				Iconified:  uint32(e)&cIconified != 0,
				Maximized:  uint32(e)&cMaximized != 0,
				Sticky:     uint32(e)&cSticky != 0,
				Fullscreen: uint32(e)&cFullscreen != 0,
				Above:      uint32(e)&cAbove != 0,
				Below:      uint32(e)&cBelow != 0,
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
