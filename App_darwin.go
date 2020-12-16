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
	"sync/atomic"
	"time"
)

type (
	// App is main application object
	App struct {
		count uint
		winds []*Window
	}
)

var (
	mutexNew   sync.Mutex
	winds      = []*Window{}
	lock       sync.Mutex
	appChan    = make(chan *App)
	defaultApp *App
	idItr      int64
)

// NewWindow returns window with webview
func NewWindow(title string, sizes ...int) *Window {
	if defaultApp == nil {
		defaultApp = makeApp()
	}
	return defaultApp.NewWindow(title, sizes...)
}

// WaitAllWindowClose locker
func WaitAllWindowClose() {
	if defaultApp != nil {
		defaultApp.WaitAllWindowClose()
	}
}

// WaitWindowClose locker
func WaitWindowClose(win *Window) {
	if defaultApp != nil && win != nil {
		defaultApp.WaitWindowClose(win)
	}
}

// WaitAllWindowClose locker
func (a *App) WaitAllWindowClose() {
	select {}
}

// WaitWindowClose locker
func (a *App) WaitWindowClose(win *Window) {
	select {}
}

// makeApp is make and run one instance of application (At the moment, it is possible to create only one instance)
func makeApp() *App {
	lock.Lock()
	go func() {
		runtime.LockOSThread()
		C.makeApp()
		runtime.UnlockOSThread()
	}()
	app := <-appChan
	fmt.Println("App started")
	return app
}

// NewWindow returns window with webview
func (a *App) NewWindow(title string, sizes ...int) *Window {
	mutexNew.Lock()
	defer mutexNew.Unlock()
	id := atomic.AddInt64(&idItr, 1)

	width := 500
	height := 400

	if len(sizes) > 0 {
		width = sizes[0]
	}

	if len(sizes) > 1 {
		height = sizes[1]
	}

	cRet := cRequest(func(reqid uint64) {
		C.makeWindow(C.CString(title), C.int(width), C.int(height), C.ulonglong(reqid), C.int(int(id)))
	})
	ret, ok := cRet.(C.WindowObj)
	if !ok {
		panic("Object is not C.WindowObj!")
	}
	win := &Window{
		id:        id,
		resizeble: true,
		window:    ret,
		state:     State{Hidden: false},
		app:       a,
	}

	go func() {
		time.Sleep(time.Second / 2)
		for {
			time.Sleep(time.Second / 10)
			state := win.state
			state.Maximized = goBool(C.isZoomed(C.WindowObj(win.window))) && win.resizeble
			state.Iconified = goBool(C.isMiniaturized(C.WindowObj(win.window)))
			state.Hidden = !goBool(C.isVisible(C.WindowObj(win.window))) && !goBool(C.isMiniaturized(C.WindowObj(win.window)))
			state.Fullscreen = goBool(C.isFullscreen(C.WindowObj(win.window)))
			if state.Iconified {
				state.Focused = false
			}

			if state.Maximized != win.state.Maximized ||
				state.Fullscreen != win.state.Fullscreen ||
				state.Hidden != win.state.Hidden ||
				state.Iconified != win.state.Iconified {

				go stateSender(win, state)
			}
		}
	}()

	winds = append(winds, win)
	return win
}
