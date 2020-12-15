// +build freebsd linux netbsd openbsd solaris

package frame

/*
#cgo pkg-config: gtk+-3.0 webkit2gtk-4.0
#cgo linux CFLAGS: -DWEBVIEW_GTK=1 -Wno-deprecated-declarations
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
)

type (
	// appObject is main application object
	appObject struct {
		app       *C.GtkApplication
		count     uint
		winds     []*Window
		openedWns sync.WaitGroup
		shown     chan bool
	}
	// State struct
	State struct {
		Hidden     bool
		Iconified  bool
		Maximized  bool
		Sticky     bool
		Fullscreen bool
		Focused    bool
		Tiled      bool
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
	appChan     = make(chan *appObject)
	idItr       int64
	defaultApp  = makeApp()
)

// NewWindow returns window with webview
func NewWindow(title string, sizes ...int) *Window {
	return defaultApp.NewWindow(title, sizes...)
}

// WaitAllWindowClose locker
func WaitAllWindowClose() {
	defaultApp.WaitAllWindowClose()
}

// WaitWindowClose locker
func WaitWindowClose(win *Window) {
	defaultApp.WaitWindowClose(win)
}

// makeApp is make and run one instance of application (At the moment, it is possible to create only one instance)
func makeApp() *appObject {
	lock.Lock()
	go func() {
		runtime.LockOSThread()
		C.makeApp()
		runtime.UnlockOSThread()
	}()
	return <-appChan
}

// SetDefaultIconFromFile for application windows
func (a *appObject) SetDefaultIconFromFile(filename string) {
	C.gtk_window_set_default_icon_from_file(C.gcharptr(C.CString(filename)), nil)
}

// SetDefaultIconName for application windows
func (a *appObject) SetDefaultIconName(name string) {
	C.gtk_window_set_default_icon_name(C.gcharptr(C.CString(name)))
}

// WaitAllWindowClose locker
func (a *appObject) WaitAllWindowClose() {
	<-a.shown
	a.openedWns.Wait()
}

// WaitWindowClose locker
func (a *appObject) WaitWindowClose(win *Window) {
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

// NewWindow returns window with webview
func (a *appObject) NewWindow(title string, sizes ...int) *Window {
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
		C.makeWindow(&C.idleData{
			id:      C.int(int(id)),
			app:     a.app,
			content: C.gcharptr(C.CString(title)),
			width:   C.int(width),
			height:  C.int(height),
			req_id:  C.ulonglong(reqid),
		})
	})
	ret, _ := cRet.(*C.WindowObj)
	wind := &Window{
		id:      id,
		window:  ret.window,
		box:     ret.box,
		webview: ret.webview,
		menubar: ret.menubar,
		state:   State{Hidden: true},
		app:     a,
	}
	winds = append(winds, wind)
	return wind
}

//export goAppActivated
func goAppActivated(app *C.GtkApplication) {
	appChan <- &appObject{
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
					Sticky:     uint32(e)&cSticky != 0,
					Fullscreen: uint32(e)&cFullscreen != 0,
					Focused:    uint32(e)&cFocused != 0,
					Tiled:      uint32(e)&cTiled != 0,
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
