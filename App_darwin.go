// +build darwin

package frame

/*
#cgo CFLAGS:  -DWEBVIEW_COCOA=1 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#import <Cocoa/Cocoa.h>
#import  "c_darwin.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
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
	app := <-appChan
	fmt.Println("App started")
	return app
}

// SetDefaultIconFromFile for application windows
func (a *App) SetDefaultIconFromFile(filename string) {
	// C.gtk_window_set_default_icon_from_file(C.gcharptr(C.CString(filename)), nil)
}

// SetDefaultIconName for application windows
func (a *App) SetDefaultIconName(name string) {
	// C.gtk_window_set_default_icon_name(C.gcharptr(C.CString(name)))
}

// NewFrame returns window with webview
func (a *App) NewFrame(title string, sizes ...int) *Frame {
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
	// box := C.makeBox(window)
	// menubar := C.makeMenubar(box)
	// webview := C.makeWebview(box)
	frame := &Frame{
		resizeble: true,
		modal:     -1,
		modalFor:  -1,
		window:    int(window),
		state: State{
			Hidden:     false,
			Iconified:  false,
			Maximized:  false,
			Sticky:     false,
			Fullscreen: false,
			Focused:    false,
			Tiled:      false,
		},
	}

	go func() {
		time.Sleep(time.Second / 2)
		for {
			time.Sleep(time.Second / 100)
			state := frame.state
			frame.state.Focused = goBool(C.isFocused(C.int(frame.window)))
			frame.state.Iconified = goBool(C.isMiniaturized(C.int(frame.window)))
			frame.state.Maximized = goBool(C.isZoomed(C.int(frame.window))) && frame.resizeble
			frame.state.Hidden = !goBool(C.isVisible(C.int(frame.window))) && !goBool(C.isMiniaturized(C.int(frame.window)))
			frame.state.Fullscreen = goBool(C.isFullscreen(C.int(frame.window)))

			// C.unsetModal(C.int(f.window))
			if state.Focused != frame.state.Focused || state.Iconified != frame.state.Iconified || state.Maximized != frame.state.Maximized || state.Hidden != frame.state.Hidden || state.Fullscreen != frame.state.Fullscreen {
				frame.StateEvent(frame.state)
			}
		}
	}()

	frames = append(frames, frame)
	return frame
}
