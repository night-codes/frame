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
		mainMenu  *Menu
		appMenu   *Menu
		openedWns sync.WaitGroup
		shown     chan bool
	}
)

var (
	mutexNew sync.Mutex
	winds    = []*Window{}
	lock     sync.Mutex
	appChan  = make(chan *App)
	idItr    int64
)

// WaitAllWindowClose locker
func (a *App) WaitAllWindowClose() {
	<-a.shown
	a.openedWns.Wait()
}

// WaitWindowClose locker
func (a *App) WaitWindowClose(win *Window) {
	<-a.shown
	shown := false
	for {
		if !win.state.Hidden {
			shown = true
		}
		if win.state.Hidden && shown {
			break
		}
	}
}

// SetIconFromFile sets dock icon
func (a *App) SetIconFromFile(filename string) {
	C.setAppIconFromFile(C.CString(filename))
}

// MakeApp is make and run one instance of application (At the moment, it is possible to create only one instance)
func MakeApp(appName string) *App {
	lock.Lock()
	go func() {
		runtime.LockOSThread()
		C.makeApp(C.CString(appName))
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
		state:     State{Hidden: true},
		app:       a,
		MainMenu:  a.mainMenu,
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

			if state.Hidden != win.state.Hidden || (!state.Hidden && (state.Maximized != win.state.Maximized ||
				state.Fullscreen != win.state.Fullscreen || state.Iconified != win.state.Iconified)) {

				go stateSender(win, state)
			}
		}
	}()

	winds = append(winds, win)
	return win
}
