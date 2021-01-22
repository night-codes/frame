// +build windows

package frame

import (
	"time"
	"unsafe"
)

type (
	// Window struct
	Window struct {
		id        int
		thread    int
		browser   unsafe.Pointer
		window    unsafe.Pointer
		destroyed bool

		StateEvent func(State)
		Invoke     func(string)
		MainMenu   *Menu

		OnHide func()
		OnShow func()

		OnFocus   func()
		OnUnfocus func()

		OnIconize   func()
		OnDeiconize func()

		OnDemaximize func()
		OnMaximize   func()

		OnFullscreen     func()
		OnExitFullscreen func()

		app   *App
		state State
	}

	// WindowType struct
	WindowType int

	// WindowPosition struct
	WindowPosition int

	// StrutPosition struct
	StrutPosition int
)

const (
/*
	TypeNormal       = WindowType(C.NSNormalWindowLevel)      // Normal toplevel window.
	TypeDialog       = WindowType(C.NSNormalWindowLevel)      // Dialog window.
	TypeMenu         = WindowType(C.NSTornOffMenuWindowLevel) // Window used to implement a menu; GTK+ uses this hint only for torn-off menus, see GtkTearoffMenuItem.
	TypeToolbar      = WindowType(C.NSNormalWindowLevel)      // Window used to implement toolbars.
	TypeSplashscreen = WindowType(C.NSStatusWindowLevel)      // Window used to display a splash screen during application startup.
	TypeUtility      = WindowType(C.NSNormalWindowLevel)      // Utility windows which are not detached toolbars or dialogs.
	TypeDock         = WindowType(C.NSDockWindowLevel)        // Used for creating dock or panel windows.
	TypeDesktop      = WindowType(C.kCGDesktopWindowLevelKey) // Used for creating the desktop background window.
	TypeDropdownMenu = WindowType(C.NSTornOffMenuWindowLevel) // A menu that belongs to a menubar.
	TypePopupMenu    = WindowType(C.NSPopUpMenuWindowLevel)   // A menu that does not belong to a menubar, e.g. a context menu.
	TypeTooltip      = WindowType(C.NSPopUpMenuWindowLevel)   // A tooltip.
	TypeNotification = WindowType(C.NSStatusWindowLevel)      // A notification - typically a “bubble” that belongs to a status icon.
	TypeCombo        = WindowType(C.NSPopUpMenuWindowLevel)   // A popup from a combo box.
	TypeDnd          = WindowType(C.NSPopUpMenuWindowLevel)   // A window that is used to implement a DND cursor.


StrutTop    = StrutPosition(C.PANEL_WINDOW_POSITION_TOP)
StrutBottom = StrutPosition(C.PANEL_WINDOW_POSITION_BOTTOM)
StrutLeft   = StrutPosition(C.PANEL_WINDOW_POSITION_LEFT)
StrutRight  = StrutPosition(C.PANEL_WINDOW_POSITION_RIGHT) */
)

// SetType of window
func (f *Window) SetType(hint WindowType) *Window {
	// C.gtk_window_set_type_hint(C.WindowObj(f.window), C.GdkWindowTypeHint(int(hint)))
	return f
}

// GetScreenScaleFactor returns scale factor of window monitor
func (f *Window) GetScreenScaleFactor() float64 {
	return 1 // float64(C.getScreenScale(C.WindowObj(f.window)))
}

// GetScreenSize returns size of window monitor
func (f *Window) GetScreenSize() (width, height int) {
	// size := C.getScreenSize(C.WindowObj(f.window))
	// width, height = int(size.width), int(size.height)
	width, height = 0, 0
	return
}

// GetSize returns width and height of window
func (f *Window) GetSize() (width, height int) {
	// size := C.windowSize(C.WindowObj(f.window))
	// width, height = int(size.width), int(size.height)
	width, height = 0, 0
	return
}

// GetWebviewSize returns width and height of window webview content
func (f *Window) GetWebviewSize() (width, height int) {
	// size := C.contentSize(C.WindowObj(f.window))
	// width, height = int(size.width), int(size.height)
	width, height = 0, 0
	return
}

// GetPosition returns position of window
func (f *Window) GetPosition() (x, y int) {
	// position := C.windowPosition(C.WindowObj(f.window))
	// x, y = int(position.x), int(position.y)
	x, y = 0, 0
	return
}

// SetIconFromFile for Window
func (f *Window) SetIconFromFile(filename string) *Window {
	// C.setWindowIconFromFile(C.WindowObj(f.window), C.CString(filename))
	return f
}

// SetOpacity of window
func (f *Window) SetOpacity(opacity float64) *Window {
	// C.setWindowAlpha(C.WindowObj(f.window), C.double(opacity))
	return f
}

// Maximize window
func (f *Window) Maximize(flag bool) *Window {
	// if (flag && !f.state.Maximized) || (!flag && f.state.Maximized) {
	// 	C.toggleMaximize(C.WindowObj(f.window))
	// }
	return f
}

// KeepAbove the window
func (f *Window) KeepAbove(flag bool) *Window {
	// C.windowKeepAbove(C.WindowObj(f.window), C.bool(flag))
	return f
}

// KeepBelow of window
func (f *Window) KeepBelow(flag bool) *Window {
	// C.windowKeepBelow(C.WindowObj(f.window), C.bool(flag))
	return f
}

// SkipTaskbar of window
func (f *Window) SkipTaskbar(flag bool) *Window {
	// C.setWindowSkipTaskbar(C.WindowObj(f.window), C.bool(flag))
	return f
}

// SkipPager of window
func (f *Window) SkipPager(flag bool) *Window {
	// C.setWindowSkipPager(C.WindowObj(f.window), C.bool(flag))
	return f
}

// Stick window
func (f *Window) Stick(flag bool) *Window {
	// C.stickWindow(C.WindowObj(f.window), C.bool(flag))
	return f
}

// Fullscreen window
func (f *Window) Fullscreen(flag bool) *Window {
	// if (flag && !f.state.Fullscreen) || (!flag && f.state.Fullscreen) {
	// 	C.toggleFullScreen(C.WindowObj(f.window))
	// }
	return f
}

// SetDeletable of window
func (f *Window) SetDeletable(flag bool) *Window {
	// C.setWindowDeletable(C.WindowObj(f.window), C.bool(flag))
	return f
}

// SetDecorated of window
func (f *Window) SetDecorated(flag bool) *Window {
	// C.setWindowDecorated(C.WindowObj(f.window), C.bool(flag))
	return f
}

// Iconify window
func (f *Window) Iconify(flag bool) *Window {
	// C.iconifyWindow(C.WindowObj(f.window), C.bool(flag))
	return f
}

// Load URL to Window webview
func (f *Window) Load(uri string) *Window {
	// C.loadURI(C.WindowObj(f.window), C.CString(uri))
	return f
}

// LoadHTML to Window webview
func (f *Window) LoadHTML(html string, baseURI string) *Window {
	// C.loadHTML(C.WindowObj(f.window), C.CString(html), C.CString(baseURI))
	return f
}

// SetResizeble of window
func (f *Window) SetResizeble(flag bool) *Window {
	// f.resizeble = flag
	// C.setWindowResizeble(C.WindowObj(f.window), C.bool(flag))
	return f
}

// SetStateEvent set handler function for window state event
func (f *Window) SetStateEvent(fn func(State)) *Window {
	// f.StateEvent = fn
	return f
}

// SetInvoke set handler function for window state event
func (f *Window) SetInvoke(fn func(string)) *Window {
	// f.Invoke = fn
	return f
}

// SetTitle of window
func (f *Window) SetTitle(title string) *Window {
	// C.setTitle(C.WindowObj(f.window), C.CString(title))
	return f
}

// SetSize of the window
func (f *Window) SetSize(width, height int) *Window {
	// C.resizeWindow(C.WindowObj(f.window), C.int(width), C.int(height))
	return f
}

// SetWebviewSize sets size of webview (without titlebar)
func (f *Window) SetWebviewSize(width, height int) *Window {
	// C.resizeContent(C.WindowObj(f.window), C.int(width), C.int(height))
	return f
}

// Move the window
func (f *Window) Move(x, y int) *Window {
	// C.moveWindow(C.WindowObj(f.window), C.int(x), C.int(y))
	return f
}

// SetCenter of window
func (f *Window) SetCenter() *Window {
	// C.setWindowCenter(C.WindowObj(f.window))
	return f
}

// Eval JS
func (f *Window) Eval(js string) string {
	// cRet := cRequest(func(id uint64) {
	// 	C.evalJS(C.WindowObj(f.window), C.CString(js), C.ulonglong(id))
	// })
	// ret, _ := cRet.(string)
	// return ret
	return ""
}

//export goEvalRet
// func goEvalRet(reqid C.ulonglong, err *C.char) {
// 	go func() {
// 		if chi, ok := goRequests.Load(uint64(reqid)); ok {
// 			if ch, ok := chi.(chan interface{}); ok {
// 				ch <- C.GoString(err)
// 			}
// 		}
// 	}()
// }

// SetModal makes current Window attached as modal window to parent
func (f *Window) SetModal(parent *Window) *Window {
	// C.setModal(C.WindowObj(f.window), parent.window)
	return f
}

// UnsetModal unset current Window as modal window from another Frames
func (f *Window) UnsetModal() *Window {
	// C.unsetModal(C.WindowObj(f.window))
	return f
}

// Show window
func (f *Window) Show() *Window {
	// C.showWindow(C.int(f.thread))
	time.Sleep(time.Second / 10)
	return f
}

// Hide window
func (f *Window) Hide() *Window {
	// C.hideWindow(C.WindowObj(f.window))
	return f
}

// SetBackgroundColor of Window
func (f *Window) SetBackgroundColor(r, g, b int, alfa float64) *Window {
	// C.setBackgroundColor(C.WindowObj(f.window), C.int8_t(r), C.int8_t(g), C.int8_t(b), C.double(alfa), C.bool(true))
	return f
}

// SetMaxSize of window
func (f *Window) SetMaxSize(width, height int) *Window {
	// C.setMaxWindowSize(C.WindowObj(f.window), C.int(width), C.int(height))
	return f
}

// SetMinSize of window
func (f *Window) SetMinSize(width, height int) *Window {
	// C.setMinWindowSize(C.WindowObj(f.window), C.int(width), C.int(height))
	return f
}

// Strut reserves wind space on the screen
func (f *Window) Strut(position StrutPosition, size int) *Window {
	/* monitorWidth, monitorHeight := f.GetScreen().Size()
	scale := f.GetScreen().ScaleFactor()
	var width, height int

	switch position {
	case StrutBottom, StrutTop:
		width, height = monitorWidth, size
	case StrutLeft, StrutRight:
		width, height = size, monitorHeight
	}
	f.
		SetDecorated(false).
		Resize(width, height).
		Stick(true).
		SetType(TypeDock)

	C.windowStrut(C.gtk_widget_get_window(C.WindowObj(f.window)), C.winPosition(position), C.int(width), C.int(height), C.int(monitorWidth), C.int(monitorHeight), C.int(scale))
	C.gtk_window_set_gravity(C.WindowObj(f.window), C.GDK_GRAVITY_NORTH_WEST)

	switch position {
	case StrutTop, StrutLeft:
		f.Move(0, 0)
	case StrutBottom:
		f.Move(0, monitorHeight-height)
	case StrutRight:
		f.Move(monitorWidth-width, 0)
	}
	f.Stick(true).SetType(TypeDock) */
	return f
}
