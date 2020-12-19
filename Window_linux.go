// +build freebsd linux netbsd openbsd solaris

package frame

/*
#ifndef WEBVIEW_GTK
#define WEBVIEW_GTK
#endif

#include "c_linux.h"


typedef struct GSR {
	GSourceFunc windowSetModal;
	GSourceFunc windowUnsetModal;
	GSourceFunc windowSetIcon;
	GSourceFunc windowSetOpacity;
	GSourceFunc windowSetType;
	GSourceFunc windowSetDecorated;
	GSourceFunc windowSetDeletable;
	GSourceFunc windowKeepAbove;
	GSourceFunc windowKeepBelow;
	GSourceFunc windowIconify;
	GSourceFunc windowStick;
	GSourceFunc windowMaximize;
	GSourceFunc windowFullscreen;
    GSourceFunc windowSetTitle;
    GSourceFunc windowSetSize;
    GSourceFunc windowMove;
    GSourceFunc windowSetBackgroundColor;
    GSourceFunc windowSkipTaskbar;
    GSourceFunc windowSkipPager;
    GSourceFunc windowSetResizeble;
    GSourceFunc windowSetCenter;
    GSourceFunc windowShow;
    GSourceFunc windowHide;
    GSourceFunc evalJS;
    GSourceFunc loadURI;
    GSourceFunc loadHTML;
    GSourceFunc setZoom;
} GSR;

GSR fn = {
	windowSetModal,
	windowUnsetModal,
	windowSetIcon,
	windowSetOpacity,
	windowSetType,
	windowSetDecorated,
	windowSetDeletable,
	windowKeepAbove,
	windowKeepBelow,
	windowIconify,
	windowStick,
	windowMaximize,
	windowFullscreen,
	windowSetTitle,
	windowSetSize,
	windowMove,
    windowSetBackgroundColor,
    windowSkipTaskbar,
    windowSkipPager,
    windowSetResizeble,
    windowSetCenter,
    windowShow,
    windowHide,
    evalJS,
    loadURI,
	loadHTML,
	setZoom
};
*/
import "C"

type (
	// Window struct
	Window struct {
		id         int64
		window     *C.GtkWidget
		box        *C.GtkWidget
		webview    *C.GtkWidget
		menubar    *C.GtkWidget
		StateEvent func(State)
		Invoke     func(string)
		app        *App
		state      State
		maxLimited bool
		maxWidth   int
		maxHeight  int
		minWidth   int
		minHeight  int
	}

	// WindowType struct
	WindowType int

	// WindowPosition struct
	WindowPosition int

	// StrutPosition struct
	StrutPosition int
)

const (
	// Window types
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
	StrutRight  = StrutPosition(C.PANEL_WINDOW_POSITION_RIGHT)
)

func idle(fn C.GSourceFunc, data C.idleData) {
	C.gdk_threads_add_idle(fn, C.gpointer(&data))
}

// SetBackgroundColor of the window
func (f *Window) SetBackgroundColor(r, g, b int, alfa float64) *Window {
	idle(C.fn.windowSetBackgroundColor, C.idleData{
		window:  f.window,
		webview: f.webview,
		rgba: C.GdkRGBA{
			red:   C.gdouble(float64(r) / 255),
			green: C.gdouble(float64(g) / 255),
			blue:  C.gdouble(float64(b) / 255),
			alpha: C.gdouble(alfa),
		},
	})
	return f
}

// ----------------------------------------------
// SetMinSize of window
func (f *Window) SetMinSize(width, height int) *Window {
	f.minHeight = height
	f.minWidth = width
	C.setSizes(f.window, C.gint(C.int(f.maxWidth)), C.gint(C.int(f.maxHeight)), C.gint(C.int(f.minWidth)), C.gint(C.int(f.minHeight)))
	return f
}

// SetMaxSize of window
func (f *Window) SetMaxSize(width, height int) *Window {
	f.maxLimited = width > 0 || height > 0
	monitorWidth, monitorHeight := f.GetScreenSize()

	if height <= 0 {
		height = monitorHeight
	}
	if width <= 0 {
		width = monitorWidth
	}
	f.maxHeight = height
	f.maxWidth = width
	C.setSizes(f.window, C.gint(C.int(f.maxWidth)), C.gint(C.int(f.maxHeight)), C.gint(C.int(f.minWidth)), C.gint(C.int(f.minHeight)))
	return f
}

// SetModal makes current Window attached as modal window to parent
func (f *Window) SetModal(parent *Window) *Window {
	idle(C.fn.windowSetModal, C.idleData{
		window:       f.window,
		windowParent: parent.window,
	})
	return f
}

// UnsetModal unset current Window as modal window from another Frames
func (f *Window) UnsetModal() *Window {
	idle(C.fn.windowUnsetModal, C.idleData{
		window: f.window,
	})
	return f
}

// SetIconFromFile for Window
func (f *Window) SetIconFromFile(filename string) *Window {
	idle(C.fn.windowSetIcon, C.idleData{
		window:  f.window,
		content: C.gcharptr(C.CString(filename)),
	})
	return f
}

// SetOpacity of window
func (f *Window) SetOpacity(opacity float64) *Window {
	idle(C.fn.windowSetOpacity, C.idleData{
		window: f.window,
		dbl:    C.gdouble(opacity),
	})
	return f
}

// SetZoom of webview
func (f *Window) SetZoom(zoom float64) *Window {
	idle(C.fn.setZoom, C.idleData{
		webview: f.webview,
		dbl:     C.gdouble(zoom),
	})
	return f
}

// SetType of window
func (f *Window) SetType(hint WindowType) *Window {
	idle(C.fn.windowSetType, C.idleData{
		window: f.window,
		hint:   C.int(int(hint)),
	})
	return f
}

// ----------------------------------------------

// SetDecorated of window
func (f *Window) SetDecorated(flag bool) *Window {
	idle(C.fn.windowSetDecorated, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// SetDeletable of window
func (f *Window) SetDeletable(flag bool) *Window {
	idle(C.fn.windowSetDeletable, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// KeepAbove the window
func (f *Window) KeepAbove(flag bool) *Window {
	idle(C.fn.windowKeepAbove, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// KeepBelow of window
func (f *Window) KeepBelow(flag bool) *Window {
	idle(C.fn.windowKeepBelow, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// Iconify window
func (f *Window) Iconify(flag bool) *Window {
	idle(C.fn.windowIconify, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// Stick window
func (f *Window) Stick(flag bool) *Window {
	idle(C.fn.windowStick, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// Maximize window
func (f *Window) Maximize(flag bool) *Window {
	idle(C.fn.windowMaximize, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// Fullscreen window
func (f *Window) Fullscreen(flag bool) *Window {
	idle(C.fn.windowFullscreen, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// SetTitle of window
func (f *Window) SetTitle(title string) *Window {
	idle(C.fn.windowSetTitle, C.idleData{
		window:  f.window,
		content: C.gcharptr(C.CString(title)),
	})
	return f
}

// SetSize of the window
func (f *Window) SetSize(width, height int) *Window {
	idle(C.fn.windowSetSize, C.idleData{
		window: f.window,
		width:  C.int(width),
		height: C.int(height),
	})
	return f
}

// Move the window
func (f *Window) Move(x, y int) *Window {
	idle(C.fn.windowMove, C.idleData{
		window: f.window,
		x:      C.int(x),
		y:      C.int(y),
	})
	return f
}

// SkipTaskbar of the window
func (f *Window) SkipTaskbar(flag bool) *Window {
	idle(C.fn.windowSkipTaskbar, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// SkipPager of the window
func (f *Window) SkipPager(flag bool) *Window {
	idle(C.fn.windowSkipPager, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// SetResizeble of the window
func (f *Window) SetResizeble(flag bool) *Window {
	idle(C.fn.windowSetResizeble, C.idleData{
		window: f.window,
		flag:   gboolean(flag),
	})
	return f
}

// SetCenter of the window
func (f *Window) SetCenter() *Window {
	idle(C.fn.windowSetCenter, C.idleData{
		window: f.window,
	})
	return f
}

// Show window
func (f *Window) Show() *Window {
	idle(C.fn.windowShow, C.idleData{
		window: f.window,
	})
	return f
}

// Hide window
func (f *Window) Hide() *Window {
	idle(C.fn.windowHide, C.idleData{
		window: f.window,
	})
	return f
}

// Eval JS
func (f *Window) Eval(js string) string {
	cRet := cRequest(func(id uint64) {
		idle(C.fn.evalJS, C.idleData{
			webview: f.webview,
			content: C.gcharptr(C.CString(js)),
			req_id:  C.ulonglong(id),
		})
	})
	ret, _ := cRet.(string)
	return ret
}

// Load URL to Window webview
func (f *Window) Load(uri string) *Window {
	idle(C.fn.loadURI, C.idleData{
		webview: f.webview,
		uri:     C.gcharptr(C.CString(uri)),
	})
	return f
}

// LoadHTML to Window webview
func (f *Window) LoadHTML(html string, baseURI string) *Window {
	idle(C.fn.loadHTML, C.idleData{
		webview: f.webview,
		content: C.gcharptr(C.CString(html)),
		uri:     C.gcharptr(C.CString(baseURI)),
	})
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

// GetScreenSize returns size of window monitor
func (f *Window) GetScreenSize() (width, height int) {
	geometry := C.getMonitorSize(f.window)
	width, height = int(geometry.width), int(geometry.height)
	return
}

// GetScreenScaleFactor returns scale factor of window monitor
func (f *Window) GetScreenScaleFactor() int {
	return int(C.getMonitorScaleFactor(f.window))
}

// GetSize returns width and height of window
func (f *Window) GetSize() (width, height int) {
	var cWidth, cHeight C.gint
	C.gtk_window_get_size(C.to_GtkWindow(f.window), &cWidth, &cHeight)
	width, height = int(cWidth), int(cHeight)
	return
}

// Strut reserves wind space on the screen
func (f *Window) Strut(position StrutPosition, size int) *Window {
	monitorWidth, monitorHeight := f.GetScreenSize()
	scale := f.GetScreenScaleFactor()
	var width, height int

	switch position {
	case StrutBottom, StrutTop:
		width, height = monitorWidth, size
	case StrutLeft, StrutRight:
		width, height = size, monitorHeight
	}

	f.UnsetModal()
	C.windowStrut(C.gtk_widget_get_window(f.window), C.winPosition(position), C.int(width), C.int(height), C.int(monitorWidth), C.int(monitorHeight), C.int(scale))

	switch position {
	case StrutTop, StrutLeft:
		f.Move(0, 0)
	case StrutBottom:
		f.Move(0, monitorHeight-height)
	case StrutRight:
		f.Move(monitorWidth-width, 0)
	}
	f.Stick(true).SetType(TypeDock)
	return f
}
