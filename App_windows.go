// +build windows

package frame

/*
#cgo windows CXXFLAGS: -std=c++11
#cgo windows,amd64 LDFLAGS: -L./dll/x64 -lwebview -lWebView2Loader
#cgo windows,386 LDFLAGS: -L./dll/x86 -lwebview -lWebView2Loader

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

type (
	// App is main application object
	App struct {
		app       interface{} // *C.GtkApplication
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

// MakeApp is make and run one instance of application (At the moment, it is possible to create only one instance)
func MakeApp(appName string) *App {
	lock.Lock()
	/* go func() {
		runtime.LockOSThread()
		C.makeApp(C.CString(appName))
		runtime.UnlockOSThread()
	}()
	return <-appChan
	*/

	return &App{}
}

// SetIconFromFile  sets application icon
func (a *App) SetIconFromFile(filename string) {
	//C.gtk_window_set_default_icon_from_file(C.gcharptr(C.CString(filename)), nil)
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
	_ = width
	_ = height

	cRet := cRequest(func(reqid uint64) {
		go C.makeWindow(C.CString(title), C.int(width), C.int(height), C.ulonglong(reqid), C.int(int(id)))
	})
	fmt.Println(":::CREATED")

	ret, ok := cRet.(*C.WindowObj)
	if !ok {
		fmt.Println(":::NOT OK")
	}
	// fmt.Printf("%+v\n", ret.window)
	fmt.Println(":::OK")
	// select {}
	wind := &Window{
		id:     id,
		window: ret.window,
		// box:     ret.box,
		// webview: ret.webview,
		// menubar: ret.menubar,
		state: State{Hidden: true},
		MainMenu: &Menu{
			menu: nil, //ret.menubar,
		},
		app: a,
	}
	winds = append(winds, wind)
	return wind
}
