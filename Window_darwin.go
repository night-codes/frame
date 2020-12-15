// +build darwin

package frame

/*
#import  "c_darwin.h"
*/
import "C"
import (
	"sync"
)

type (
	// Window struct
	Window struct {
		id         int64
		window     C.WindowObj
		modal      *Window
		parent     *Window
		StateEvent func(State)
		Invoke     func(string)
		app        *App
		state      State
		resizeble  bool // BACK after modal?
	}

	// WindowType struct
	WindowType int

	// WindowPosition struct
	WindowPosition int

	// StrutPosition struct
	StrutPosition int
)

var (
	goRequestID uint64
	goRequests  sync.Map
)

const (
/*
TypeNormal       = WindowType(C.GDK_WINDOW_TYPE_HINT_NORMAL)        // Normal toplevel window.
TypeDialog       = WindowType(C.GDK_WINDOW_TYPE_HINT_DIALOG)        // Dialog window.
TypeMenu         = WindowType(C.GDK_WINDOW_TYPE_HINT_MENU)          // Window used to implement a menu; GTK+ uses this hint only for torn-off menus, see GtkTearoffMenuItem.
TypeToolbar      = WindowType(C.GDK_WINDOW_TYPE_HINT_TOOLBAR)       // Window used to implement toolbars.
TypeSplashscreen = WindowType(C.GDK_WINDOW_TYPE_HINT_SPLASHSCREEN)  // Window used to display a splash screen during application startup.
TypeUtility      = WindowType(C.GDK_WINDOW_TYPE_HINT_UTILITY)       // Utility windows which are not detached toolbars or dialogs.
TypeDock         = WindowType(C.GDK_WINDOW_TYPE_HINT_DOCK)          // Used for creating dock or panel windows.
TypeDesktop      = WindowType(C.GDK_WINDOW_TYPE_HINT_DESKTOP)       // Used for creating the desktop background window.
TypeDropdownMenu = WindowType(C.GDK_WINDOW_TYPE_HINT_DROPDOWN_MENU) // A menu that belongs to a menubar.
TypePopupMenu    = WindowType(C.GDK_WINDOW_TYPE_HINT_POPUP_MENU)    // A menu that does not belong to a menubar, e.g. a context menu.
TypeTooltip      = WindowType(C.GDK_WINDOW_TYPE_HINT_TOOLTIP)       // A tooltip.
TypeNotification = WindowType(C.GDK_WINDOW_TYPE_HINT_NOTIFICATION)  // A notification - typically a “bubble” that belongs to a status icon.
TypeCombo        = WindowType(C.GDK_WINDOW_TYPE_HINT_COMBO)         // A popup from a combo box.
TypeDnd          = WindowType(C.GDK_WINDOW_TYPE_HINT_DND)           // A window that is used to implement a DND cursor.

StrutTop    = StrutPosition(C.PANEL_WINDOW_POSITION_TOP)
StrutBottom = StrutPosition(C.PANEL_WINDOW_POSITION_BOTTOM)
StrutLeft   = StrutPosition(C.PANEL_WINDOW_POSITION_LEFT)
StrutRight  = StrutPosition(C.PANEL_WINDOW_POSITION_RIGHT) */
)

// Load URL to Window webview
func (f *Window) Load(uri string) *Window {
	C.loadURI(C.WindowObj(f.window), C.CString(uri))
	return f
}

// LoadHTML to Window webview
func (f *Window) LoadHTML(html string, baseURI string) *Window {
	C.loadHTML(C.WindowObj(f.window), C.CString(html), C.CString(baseURI))
	return f
}

// SkipTaskbar of window
func (f *Window) SkipTaskbar(skip bool) *Window {
	// C.gtk_window_set_skip_taskbar_hint(C.WindowObj(f.window), gboolean(skip))
	return f
}

// SkipPager of window
func (f *Window) SkipPager(skip bool) *Window {
	// C.gtk_window_set_skip_pager_hint(C.WindowObj(f.window), gboolean(skip))
	return f
}

// SetResizeble of window
func (f *Window) SetResizeble(resizeble bool) *Window {
	f.resizeble = resizeble
	C.setWindowResizeble(C.WindowObj(f.window), C.bool(resizeble))
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
	C.setTitle(C.WindowObj(f.window), C.CString(title))
	return f
}

// SetSize of the window
func (f *Window) SetSize(width, height int) *Window {
	C.resizeWindow(C.WindowObj(f.window), C.int(width), C.int(height))
	return f
}

// Move the window
func (f *Window) Move(x, y int) *Window {
	/*
		// visible := C.gtk_widget_get_visible(C.WindowObj(f.window)) == 1
		C.gtk_window_move(C.WindowObj(f.window), C.gint(C.int(x)), C.gint(C.int(y)))
	*/
	return f
}

// SetCenter of window
func (f *Window) SetCenter() *Window {
	C.setWindowCenter(C.WindowObj(f.window))
	return f
}

// Eval JS
func (f *Window) Eval(js string) string {
	cRet := cRequest(func(id uint64) {
		C.evalJS(C.WindowObj(f.window), C.CString(js), C.ulonglong(id))
	})
	ret, _ := cRet.(string)
	return ret
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
func goWinRet(reqid C.ulonglong, win C.WindowObj) {
	go func() {
		if chi, ok := goRequests.Load(uint64(reqid)); ok {
			if ch, ok := chi.(chan interface{}); ok {
				ch <- win
			}
		}
	}()
}

// SetModal makes current Window attached as modal window to parent
func (f *Window) SetModal(parent *Window) *Window {
	f.parent = parent
	parent.modal = f
	C.setModal(C.WindowObj(f.window), parent.window)
	return f
}

// UnsetModal unset current Window as modal window from another Frames
func (f *Window) UnsetModal() *Window {
	if f.parent != nil {
		f.parent.modal = nil
		f.parent = nil
	}
	C.unsetModal(C.WindowObj(f.window))
	return f
}

// SetDecorated of window
func (f *Window) SetDecorated(decorated bool) *Window {
	// C.gtk_window_set_decorated(C.WindowObj(f.window), gboolean(decorated))
	return f
}

// SetDeletable of window
func (f *Window) SetDeletable(deletable bool) *Window {
	// C.gtk_window_set_deletable(C.WindowObj(f.window), gboolean(deletable))
	return f
}

// KeepAbove the window
func (f *Window) KeepAbove(above bool) *Window {
	// C.gtk_window_set_keep_above(C.WindowObj(f.window), gboolean(above))
	return f
}

// KeepBelow of window
func (f *Window) KeepBelow(below bool) *Window {
	// C.gtk_window_set_keep_below(C.WindowObj(f.window), gboolean(below))
	return f
}

// Show window
func (f *Window) Show() *Window {
	C.showWindow(C.WindowObj(f.window))
	return f
}

// Hide window
func (f *Window) Hide() *Window {
	C.hideWindow(C.WindowObj(f.window))
	return f
}

// Iconify window
func (f *Window) Iconify(iconify bool) *Window {
	if iconify {
		// C.gtk_window_iconify(C.WindowObj(f.window))
	} else {
		// C.gtk_window_deiconify(C.WindowObj(f.window))
	}
	return f
}

// Stick window
func (f *Window) Stick(stick bool) *Window {
	if stick {
		// C.gtk_window_stick(C.WindowObj(f.window))
	} else {
		// C.gtk_window_unstick(C.WindowObj(f.window))
	}
	return f
}

// Maximize window
func (f *Window) Maximize(maximize bool) *Window {
	if maximize {
		// C.gtk_window_maximize(C.WindowObj(f.window))
	} else {
		// C.gtk_window_unmaximize(C.WindowObj(f.window))
	}
	return f
}

// Fullscreen window
func (f *Window) Fullscreen(fullscreen bool) *Window {
	if fullscreen {
		// C.gtk_window_fullscreen(C.WindowObj(f.window))
	} else {
		// C.gtk_window_unfullscreen(C.WindowObj(f.window))
	}
	return f
}

// SetIconFromFile for Window
func (f *Window) SetIconFromFile(filename string) *Window {
	// C.gtk_window_set_icon_from_file(C.WindowObj(f.window), C.gcharptr(C.CString(filename)), nil)
	return f
}

// SetIconName for Window
func (f *Window) SetIconName(name string) *Window {
	// C.gtk_window_set_icon_name(C.WindowObj(f.window), C.gcharptr(C.CString(name)))
	return f
}

// SetBackgroundColor of Window
func (f *Window) SetBackgroundColor(r, g, b int, alfa float64) *Window {
	C.setBackgroundColor(C.WindowObj(f.window), C.int8_t(r), C.int8_t(g), C.int8_t(b), C.double(alfa), C.bool(true))
	return f
}

// SetMaxSize of window
func (f *Window) SetMaxSize(width, height int) *Window {
	C.setMaxWindowSize(C.WindowObj(f.window), C.int(width), C.int(height))
	return f
}

// SetMinSize of window
func (f *Window) SetMinSize(width, height int) *Window {
	C.setMinWindowSize(C.WindowObj(f.window), C.int(width), C.int(height))
	return f
}

// SetOpacity of window
func (f *Window) SetOpacity(opacity float64) *Window {
	// C.gdk_window_set_opacity(C.gtk_widget_get_window(C.WindowObj(f.window)), C.gdouble(opacity))
	return f
}

// SetType of window
func (f *Window) SetType(hint WindowType) *Window {
	// C.gtk_window_set_type_hint(C.WindowObj(f.window), C.GdkWindowTypeHint(int(hint)))
	return f
}

// GetSize returns width and height of window
func (f *Window) GetSize() (width, height int) {
	var cWidth, cHeight C.int
	// C.gtk_window_get_size(C.WindowObj(f.window), &cWidth, &cHeight)
	width, height = int(cWidth), int(cHeight)
	return
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
