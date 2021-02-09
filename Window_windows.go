// +build windows

package frame

/*

 */
import "C"

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type (
	// Window struct
	Window struct {
		id        int
		thread    int
		browser   unsafe.Pointer
		window    unsafe.Pointer
		parent    *Window
		destroyed bool

		StateEvent func(State)
		Invoke     func(string)
		MainMenu   *Menu

		OnShow func()
		OnHide func()

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
	dpi, _, _ := winGetDpiForWindow.Call(uintptr(f.window))
	return float64(uint64(dpi)) / 96.0 // float64(C.getScreenScale(C.WindowObj(f.window)))
}

// GetScreenSize returns size of window monitor
func (f *Window) GetScreenSize() (width, height int) {
	width, height = f.getScreenSize()
	scale := f.GetScreenScaleFactor()
	if scale > 0 {
		width = int(float64(width) / scale)
		height = int(float64(height) / scale)
	}
	return
}

func (f *Window) getScreenSize() (width, height int) {
	monitor, _, _ := winMonitorFromWindow.Call(uintptr(unsafe.Pointer(f.window)), uintptr(MONITOR_DEFAULTTOPRIMARY))
	info := C_MONITORINFO{cbSize: monitorinfoSizeof}
	winGetMonitorInfo.Call(monitor, uintptr(unsafe.Pointer(&info)))
	return int(info.rcWork.right), int(info.rcWork.bottom)
}

// GetSize returns width and height of window
func (f *Window) GetSize() (width, height int) {
	width, height = f.getSize()
	scale := f.GetScreenScaleFactor()
	if scale > 0 {
		width = int(float64(width) / scale)
		height = int(float64(height) / scale)
	}
	return
}

func (f *Window) getSize() (width, height int) {
	rect := C_RECT{}
	winGetWindowRect.Call(uintptr(f.window), uintptr(unsafe.Pointer(&rect)))
	width, height = int(rect.right-rect.left), int(rect.bottom-rect.top)
	return
}

func (f *Window) update() {
	w, h := f.GetSize()
	f.SetSize(w-5, h-5)
	f.SetSize(w+5, h+5)
}

// GetPosition returns position of window
func (f *Window) GetPosition() (x, y int) {
	x, y = f.getPosition()
	scale := f.GetScreenScaleFactor()
	if scale > 0 {
		x = int(float64(x) / scale)
		y = int(float64(y) / scale)
	}
	return
}

func (f *Window) getPosition() (x, y int) {
	rect := C_RECT{}
	winGetWindowRect.Call(uintptr(f.window), uintptr(unsafe.Pointer(&rect)))
	x, y = int(rect.left), int(rect.top)
	return
}

// GetWebviewSize returns width and height of window webview content
func (f *Window) GetWebviewSize() (width, height int) {
	return f.GetSize()
}

// SetIconFromFile for Window
func (f *Window) SetIconFromFile(filename string) *Window {
	// C.setWindowIconFromFile(C.WindowObj(f.window), C.CString(filename))
	return f
}

// Maximize window
func (f *Window) Maximize(flag bool) *Window {
	if flag {
		winShowWindow.Call(uintptr(f.window), uintptr(windows.SW_SHOWMAXIMIZED))
	} else {
		winShowWindow.Call(uintptr(f.window), uintptr(windows.SW_RESTORE))
	}
	return f
}

// KeepAbove the window
func (f *Window) KeepAbove(flag bool) *Window {
	hwnd := HWND_TOPMOST
	if !flag {
		hwnd = HWND_NOTOPMOST
	}
	winSetWindowPos.Call(
		uintptr(f.window),
		uintptr(hwnd),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(SWP_NOSIZE|SWP_NOMOVE|SWP_NOACTIVATE),
	)
	return f
}

// KeepBelow of window
func (f *Window) KeepBelow(flag bool) *Window {
	hwnd := HWND_BOTTOM
	if !flag {
		hwnd = 0
	}
	winSetWindowPos.Call(
		uintptr(f.window),
		uintptr(hwnd),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(SWP_NOSIZE|SWP_NOMOVE|SWP_NOACTIVATE),
	)
	return f
}

// SkipTaskbar of window
func (f *Window) SkipTaskbar(flag bool) *Window {
	hwnd := uintptr(f.app.hiddenWindow.window)
	if !flag {
		hwnd = 0
	}
	gwlStyle := GWL_STYLE
	gwlp := GWLP_HWNDPARENT
	winSetWindowLongPtr.Call(uintptr(f.window), uintptr(gwlp), hwnd)
	style, _, _ := winGetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)))
	s := int64(style) & ^WS_MINIMIZEBOX
	if !flag {
		s = int64(style) | WS_MINIMIZEBOX
	}
	winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)), uintptr(s))
	return f
}

// SkipPager of window
func (f *Window) SkipPager(flag bool) *Window {
	hwnd := uintptr(f.app.altWindow.window)
	if !flag {
		hwnd = 0
	}
	gwlStyle := GWL_STYLE
	gwlp := GWLP_HWNDPARENT
	winSetWindowLongPtr.Call(uintptr(f.window), uintptr(gwlp), hwnd)
	style, _, _ := winGetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)))
	s := int64(style) & ^WS_MINIMIZEBOX
	if !flag {
		s = int64(style) | WS_MINIMIZEBOX
	}
	winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)), uintptr(s))
	return f
}

// Stick window
func (f *Window) Stick(flag bool) *Window {
	// C.stickWindow(C.WindowObj(f.window), C.bool(flag))
	return f
}

// Fullscreen window
func (f *Window) Fullscreen(flag bool) *Window {
	f.Maximize(false)
	f.SetDecorated(!flag)
	f.KeepAbove(flag)
	f.Maximize(flag)
	return f
}

// SetDeletable of window
func (f *Window) SetDeletable(flag bool) *Window {
	// C.setWindowDeletable(C.WindowObj(f.window), C.bool(flag))
	return f
}

// SetDecorated of window
func (f *Window) SetDecorated(flag bool) *Window {
	gwlStyle := GWL_STYLE
	style, _, _ := winGetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)))
	if flag {
		winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)), uintptr(int64(style)|WS_CAPTION|WS_THICKFRAME|WS_MINIMIZEBOX|WS_MAXIMIZEBOX|WS_SYSMENU))
	} else {
		winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)), uintptr(int64(style) & ^(WS_CAPTION|WS_THICKFRAME|WS_MINIMIZEBOX|WS_MAXIMIZEBOX|WS_SYSMENU)))
	}
	f.update()
	return f
}

// Iconify window
func (f *Window) Iconify(flag bool) *Window {
	if flag {
		winShowWindow.Call(uintptr(f.window), uintptr(windows.SW_MINIMIZE))
	} else {
		winShowWindow.Call(uintptr(f.window), uintptr(windows.SW_RESTORE))
	}
	return f
}

// Load URL to Window webview
func (f *Window) Load(uri string) *Window {
	loadURL(f.browser, uri)
	return f
}

// LoadHTML to Window webview
func (f *Window) LoadHTML(html string, baseURI string) *Window {
	loadHTML(f.browser, html, baseURI)
	return f
}

// SetResizeble of window
func (f *Window) SetResizeble(flag bool) *Window {
	gwlStyle := GWL_STYLE
	style, _, _ := winGetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)))
	if flag {
		winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)), uintptr(int64(style)|WS_MAXIMIZEBOX|WS_SIZEBOX|WS_THICKFRAME))
	} else {
		winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)), uintptr(int64(style) & ^(WS_MAXIMIZEBOX|WS_SIZEBOX|WS_THICKFRAME)))
	}
	return f
}

// SetStateEvent set handler function for window state event
func (f *Window) SetStateEvent(fn func(State)) *Window {
	f.StateEvent = fn
	return f
}

// SetInvoke set handler function for window state event
func (f *Window) SetInvoke(fn func(string)) *Window {
	f.Invoke = fn
	return f
}

// SetTitle of window
func (f *Window) SetTitle(title string) *Window {
	winSetWindowTextW.Call(uintptr(f.window), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))))
	return f
}

// SetSize of the window
func (f *Window) SetSize(width, height int) *Window {
	scale := f.GetScreenScaleFactor()
	width = int(uint64(float64(width) * scale))
	height = int(uint64(float64(height) * scale))

	pWidth, pHeight := f.getSize()
	x, y := f.getPosition()
	x = x + (pWidth-width)/2
	y = y + (pHeight-height)/2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	winSetWindowPos.Call(
		uintptr(f.window),
		uintptr(0),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(SWP_NOOWNERZORDER|SWP_NOZORDER|SWP_ASYNCWINDOWPOS), // SWP_NOSIZE SWP_NOMOVE
	)
	winUpdateWindow.Call(uintptr(f.window))
	return f
}

// Move the window
func (f *Window) Move(x, y int) *Window {
	scale := f.GetScreenScaleFactor()
	winSetWindowPos.Call(
		uintptr(f.window),
		uintptr(0),
		uintptr(int(uint64(float64(x)*scale))),
		uintptr(int(uint64(float64(y)*scale))),
		uintptr(0),
		uintptr(0),
		uintptr(SWP_NOSIZE|SWP_NOOWNERZORDER|SWP_NOZORDER|SWP_ASYNCWINDOWPOS),
	)
	return f
}

// SetWebviewSize sets size of webview (without titlebar)
func (f *Window) SetWebviewSize(width, height int) *Window {
	return f.SetSize(width, height)
}

// SetCenter of window
func (f *Window) SetCenter() *Window {
	width, height := f.GetSize()
	sWidth, sHeight := f.GetScreenSize()
	f.Move(sWidth/2-width/2, sHeight/2-height/2)
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
	if f.parent != nil {
		f.UnsetModal()
	}

	f.parent = parent
	gwlStyle := GWL_STYLE
	gwlp := GWLP_HWNDPARENT
	winSetWindowLongPtr.Call(uintptr(f.window), uintptr(gwlp), uintptr(parent.window))
	style, _, _ := winGetWindowLong.Call(uintptr(parent.window), uintptr(uint64(gwlStyle)))
	winSetWindowLong.Call(uintptr(parent.window), uintptr(uint64(gwlStyle)), uintptr(int64(style)|WS_DISABLED))

	style, _, _ = winGetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)))
	winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)), uintptr(int64(style) & ^WS_MINIMIZEBOX))

	hwnd := HWND_TOP
	winSetWindowPos.Call(
		uintptr(f.window),
		uintptr(hwnd),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(SWP_NOSIZE|SWP_NOMOVE),
	)
	return f
}

// UnsetModal unset current Window as modal window from another Frames
func (f *Window) UnsetModal() *Window {
	if f.parent != nil {
		winSetParent.Call(uintptr(f.window), 0)
		gwlStyle := GWL_STYLE
		gwlp := GWLP_HWNDPARENT
		winSetWindowLongPtr.Call(uintptr(f.window), uintptr(gwlp), uintptr(0))
		style, _, _ := winGetWindowLong.Call(uintptr(f.parent.window), uintptr(uint64(gwlStyle)))
		winSetWindowLong.Call(uintptr(f.parent.window), uintptr(uint64(gwlStyle)), uintptr(int64(style) & ^WS_DISABLED))

		style, _, _ = winGetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)))
		winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwlStyle)), uintptr(int64(style)|WS_MINIMIZEBOX))
		f.parent = nil
	}

	return f
}

// Show window
func (f *Window) Show() *Window {
	if f.parent != nil && f.parent.state.Hidden {
		f.parent.Show()
	}
	winShowWindow.Call(uintptr(f.window), uintptr(windows.SW_SHOW))

	f.state.Hidden = false
	winUpdateWindow.Call(uintptr(f.window))
	winSwitchToThisWindow.Call(uintptr(f.window), uintptr(1))
	if f.parent != nil {
		gwlStyle := GWL_STYLE
		style, _, _ := winGetWindowLong.Call(uintptr(f.parent.window), uintptr(uint64(gwlStyle)))
		winSetWindowLong.Call(uintptr(f.parent.window), uintptr(uint64(gwlStyle)), uintptr(int64(style)|WS_DISABLED))
	}
	return f
}

// Hide window
func (f *Window) Hide() *Window {
	goBrowserDoClose(ceBrowser(f.browser))
	return f
}

// SetOpacity of window
func (f *Window) SetOpacity(opacity float64) *Window {
	gwl_exstyle := GWL_EXSTYLE
	t, _, _ := winGetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwl_exstyle)))
	winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwl_exstyle)), uintptr(int64(t)|WS_EX_LAYERED))
	winSetLayeredWindowAttributes.Call(uintptr(f.window), uintptr(0), uintptr(uint64(255*opacity)), uintptr(LWA_ALPHA))
	return f
}

// SetBackgroundColor of Window
func (f *Window) SetBackgroundColor(r, g, b int, alfa float64) *Window {
	// gwl_exstyle := GWL_EXSTYLE
	// t, _, _ := winGetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwl_exstyle)))
	gclp_hbrbackground := GCLP_HBRBACKGROUND
	brush, _, _ := gdiCreateSolidBrush.Call(uintptr(0xff000000 | uint32(r)<<16 | uint32(g)<<8 | uint32(b)))
	winSetClassLongPtr.Call(uintptr(f.window), uintptr(gclp_hbrbackground), brush)

	// winSetWindowLong.Call(uintptr(f.window), uintptr(uint64(gwl_exstyle)), uintptr(int64(t)|WS_EX_LAYERED))
	// winSetLayeredWindowAttributes.Call(uintptr(f.window), uintptr(uint32(r)<<16|uint32(g)<<8|uint32(b)), uintptr(uint64(255*alfa)), LWA_COLORKEY|LWA_ALPHA)
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
