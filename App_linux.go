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
	"C"
	"runtime"
	"sync/atomic"
)

// makeApp is make and run one instance of application (At the moment, it is possible to create only one instance)
func makeApp(appName string) *App {
	lock.Lock()
	go func() {
		runtime.LockOSThread()
		C.makeApp(C.CString(appName))
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
