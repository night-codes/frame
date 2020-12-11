// +build darwin

package frame

/*
#import  "c_darwin.h"
*/
import "C"
import (
	"sync"
	"sync/atomic"
)

type (
	// Frame struct
	Frame struct {
		window     int
		modal      int
		modalFor   int
		StateEvent func(State)
		Invoke     func(string)
		state      State
		deferMove  bool
		resizeble  bool
		deferMoveX int
		deferMoveY int
	}

	// WindowType struct
	WindowType int

	// WindowPosition struct
	WindowPosition int

	// StrutPosition struct
	StrutPosition int
)

var (
	jsRequestID uint64
	jsRequests  sync.Map
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

// Load URL to Frame webview
func (f *Frame) Load(uri string) *Frame {
	C.loadUri(C.int(f.window), C.CString(uri))
	return f
}

// LoadHTML to Frame webview
func (f *Frame) LoadHTML(html string, baseURI string) *Frame {
	C.loadHTML(C.int(f.window), C.CString(html), C.CString(baseURI))
	return f
}

// SkipTaskbar of window
func (f *Frame) SkipTaskbar(skip bool) *Frame {
	// C.gtk_window_set_skip_taskbar_hint(C.to_GtkWindow(f.window), gboolean(skip))
	return f
}

// SkipPager of window
func (f *Frame) SkipPager(skip bool) *Frame {
	// C.gtk_window_set_skip_pager_hint(C.to_GtkWindow(f.window), gboolean(skip))
	return f
}

// SetResizeble of window
func (f *Frame) SetResizeble(resizeble bool) *Frame {
	f.resizeble = resizeble
	C.setWindowResizeble(C.int(f.window), C.bool(resizeble))
	return f
}

// SetStateEvent set handler function for window state event
func (f *Frame) SetStateEvent(fn func(State)) *Frame {
	f.StateEvent = fn
	return f
}

// SetInvoke set handler function for window state event
func (f *Frame) SetInvoke(fn func(string)) *Frame {
	f.Invoke = fn
	return f
}

// SetTitle of window
func (f *Frame) SetTitle(title string) *Frame {
	C.setTitle(C.int(f.window), C.CString(title))
	return f
}

// SetSize of the window
func (f *Frame) SetSize(width, height int) *Frame {
	C.resizeWindow(C.int(f.window), C.int(width), C.int(height))
	return f
}

// Move the window
func (f *Frame) Move(x, y int) *Frame {
	/* visible := C.gtk_widget_get_visible(f.window) == 1
	if !visible {
		f.deferMove = true
		f.deferMoveX = x
		f.deferMoveY = y
		return f
	}
	C.gtk_window_move(C.to_GtkWindow(f.window), C.gint(C.int(x)), C.gint(C.int(y))) */
	return f
}

// SetCenter of window
func (f *Frame) SetCenter() *Frame {
	C.setWindowCenter(C.int(f.window))
	return f
}

// Eval JS
func (f *Frame) Eval(js string) string {
	id := atomic.AddUint64(&jsRequestID, 1)
	ch := make(chan string)

	jsRequests.Store(id, ch)
	C.evalJS(C.int(f.window), C.CString(js), C.ulonglong(id))
	ret := <-ch
	jsRequests.Delete(id)
	return ret
}

//export goEvalRet
func goEvalRet(reqid C.ulonglong, err *C.char) {
	go func() {
		chi, ok := jsRequests.Load(uint64(reqid))
		ch := chi.(chan string)
		if ok {
			ch <- C.GoString(err)
		}
	}()
}

// SetModal makes current Frame attached as modal window to parent
func (f *Frame) SetModal(parent *Frame) *Frame {
	f.modalFor = parent.window
	parent.modal = f.window
	C.setModal(C.int(f.window), C.int(parent.window))
	return f
}

// UnsetModal unset current Frame as modal window from another Frames
func (f *Frame) UnsetModal() *Frame {
	if f.modalFor != -1 {
		frames[f.modalFor].modal = -1
		f.modalFor = -1
	}
	C.unsetModal(C.int(f.window))
	return f
}

// SetDecorated of window
func (f *Frame) SetDecorated(decorated bool) *Frame {
	// C.gtk_window_set_decorated(C.to_GtkWindow(f.window), gboolean(decorated))
	return f
}

// SetDeletable of window
func (f *Frame) SetDeletable(deletable bool) *Frame {
	// C.gtk_window_set_deletable(C.to_GtkWindow(f.window), gboolean(deletable))
	return f
}

// KeepAbove the window
func (f *Frame) KeepAbove(above bool) *Frame {
	// C.gtk_window_set_keep_above(C.to_GtkWindow(f.window), gboolean(above))
	return f
}

// KeepBelow of window
func (f *Frame) KeepBelow(below bool) *Frame {
	// C.gtk_window_set_keep_below(C.to_GtkWindow(f.window), gboolean(below))
	return f
}

// Show window
func (f *Frame) Show() *Frame {
	C.showWindow(C.int(f.window))
	if f.deferMove {
		f.Move(f.deferMoveX, f.deferMoveY)
	}
	return f
}

// Hide window
func (f *Frame) Hide() *Frame {
	C.hideWindow(C.int(f.window))
	return f
}

// Iconify window
func (f *Frame) Iconify(iconify bool) *Frame {
	if iconify {
		// C.gtk_window_iconify(C.to_GtkWindow(f.window))
	} else {
		// C.gtk_window_deiconify(C.to_GtkWindow(f.window))
	}
	return f
}

// Stick window
func (f *Frame) Stick(stick bool) *Frame {
	if stick {
		// C.gtk_window_stick(C.to_GtkWindow(f.window))
	} else {
		// C.gtk_window_unstick(C.to_GtkWindow(f.window))
	}
	return f
}

// Maximize window
func (f *Frame) Maximize(maximize bool) *Frame {
	if maximize {
		// C.gtk_window_maximize(C.to_GtkWindow(f.window))
	} else {
		// C.gtk_window_unmaximize(C.to_GtkWindow(f.window))
	}
	return f
}

// Fullscreen window
func (f *Frame) Fullscreen(fullscreen bool) *Frame {
	if fullscreen {
		// C.gtk_window_fullscreen(C.to_GtkWindow(f.window))
	} else {
		// C.gtk_window_unfullscreen(C.to_GtkWindow(f.window))
	}
	return f
}

// SetIconFromFile for Frame
func (f *Frame) SetIconFromFile(filename string) *Frame {
	// C.gtk_window_set_icon_from_file(C.to_GtkWindow(f.window), C.gcharptr(C.CString(filename)), nil)
	return f
}

// SetIconName for Frame
func (f *Frame) SetIconName(name string) *Frame {
	// C.gtk_window_set_icon_name(C.to_GtkWindow(f.window), C.gcharptr(C.CString(name)))
	return f
}

// SetBackgroundColor of Frame
func (f *Frame) SetBackgroundColor(r, g, b int, alfa float64) *Frame {
	C.setBackgroundColor(C.int(f.window), C.int8_t(r), C.int8_t(g), C.int8_t(b), C.double(alfa), C.bool(false))
	return f
}

// SetMaxSize of window
func (f *Frame) SetMaxSize(width, height int) *Frame {
	C.setMaxWindowSize(C.int(f.window), C.int(width), C.int(height))
	return f
}

// SetMinSize of window
func (f *Frame) SetMinSize(width, height int) *Frame {
	C.setMinWindowSize(C.int(f.window), C.int(width), C.int(height))
	return f
}

// SetOpacity of window
func (f *Frame) SetOpacity(opacity float64) *Frame {
	// C.gdk_window_set_opacity(C.gtk_widget_get_window(f.window), C.gdouble(opacity))
	return f
}

// SetType of window
func (f *Frame) SetType(hint WindowType) *Frame {
	// C.gtk_window_set_type_hint(C.to_GtkWindow(f.window), C.GdkWindowTypeHint(int(hint)))
	return f
}

// GetScreen where the window placed
// func (f *Frame) GetScreen() *Screen {
// 	screen := C.gtk_widget_get_screen(f.window)
// 	display := C.gdk_screen_get_display(screen)
// 	monitor := C.gdk_display_get_monitor_at_window(display, C.gtk_widget_get_window(f.window))
// 	return &Screen{
// 		screen:  screen,
// 		display: display,
// 		monitor: monitor,
// 	}
// }

// GetSize returns width and height of window
func (f *Frame) GetSize() (width, height int) {
	var cWidth, cHeight C.int
	// C.gtk_window_get_size(C.to_GtkWindow(f.window), &cWidth, &cHeight)
	width, height = int(cWidth), int(cHeight)
	return
}

// Strut reserves frame space on the screen
func (f *Frame) Strut(position StrutPosition, size int) *Frame {
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

	C.windowStrut(C.gtk_widget_get_window(f.window), C.winPosition(position), C.int(width), C.int(height), C.int(monitorWidth), C.int(monitorHeight), C.int(scale))
	C.gtk_window_set_gravity(C.to_GtkWindow(f.window), C.GDK_GRAVITY_NORTH_WEST)

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
