// +build freebsd linux netbsd openbsd solaris

package frame

/*
#include "linux.h"
*/
import "C"
import "unsafe"

type (
	// Frame struct
	Frame struct {
		window     unsafe.Pointer
		box        *C.GtkWidget
		webview    *C.GtkWidget
		menubar    *C.GtkWidget
		StateEvent func(State)
		deferMove  bool
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

const (
	PosNone           = WindowPosition(C.GTK_WIN_POS_NONE)
	PosCenter         = WindowPosition(C.GTK_WIN_POS_CENTER)
	PosMouse          = WindowPosition(C.GTK_WIN_POS_MOUSE)
	PosCenterAlways   = WindowPosition(C.GTK_WIN_POS_CENTER_ALWAYS)
	PosCenterOnParent = WindowPosition(C.GTK_WIN_POS_CENTER_ON_PARENT)

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

// Load URL to Frame webview
func (f *Frame) Load(uri string) *Frame {
	C.loadUri(f.webview, C.gcharptr(C.CString(uri)))
	// C.webkit_web_inspector_attach(C.webkit_web_view_get_inspector(C.to_WebKitWebView(f.webview)))
	return f
}

// LoadHTML to Frame webview
func (f *Frame) LoadHTML(html string, baseURI string) *Frame {
	C.loadHTML(f.webview, C.gcharptr(C.CString(html)), C.gcharptr(C.CString(baseURI)))
	return f
}

// SkipTaskbar of window
func (f *Frame) SkipTaskbar(skip bool) *Frame {
	C.gtk_window_set_skip_taskbar_hint(C.to_GtkWindow(f.window), gboolean(skip))
	return f
}

// SkipPager of window
func (f *Frame) SkipPager(skip bool) *Frame {
	C.gtk_window_set_skip_pager_hint(C.to_GtkWindow(f.window), gboolean(skip))
	return f
}

// SetResizeble of window
func (f *Frame) SetResizeble(resizeble bool) *Frame {
	C.gtk_window_set_resizable(C.to_GtkWindow(f.window), gboolean(resizeble))
	return f
}

// SetStateEvent set handler function for window state event
func (f *Frame) SetStateEvent(fn func(State)) *Frame {
	f.StateEvent = fn
	return f
}

// SetTitle of window
func (f *Frame) SetTitle(title string) *Frame {
	C.gtk_window_set_title(C.to_GtkWindow(f.window), C.gcharptr(C.CString(title)))
	return f
}

// SetDefaultSize of window
func (f *Frame) SetDefaultSize(width, height int) *Frame {
	C.gtk_window_set_default_size(C.to_GtkWindow(f.window), C.gint(C.int(width)), C.gint(C.int(height)))
	return f
}

// Resize the window
func (f *Frame) Resize(width, height int) *Frame {
	C.gtk_window_resize(C.to_GtkWindow(f.window), C.gint(C.int(width)), C.gint(C.int(height)))
	return f
}

// Move the window
func (f *Frame) Move(x, y int) *Frame {
	visible := C.gtk_widget_get_visible(f.window) == 1
	if !visible {
		f.deferMove = true
		f.deferMoveX = x
		f.deferMoveY = y
		return f
	}
	C.gtk_window_move(C.to_GtkWindow(f.window), C.gint(C.int(x)), C.gint(C.int(y)))
	return f
}

// SetPosition of window
func (f *Frame) SetPosition(position WindowPosition) *Frame {
	C.gtk_window_set_position(C.to_GtkWindow(f.window), C.GtkWindowPosition(position))
	return f
}

// SetModal makes current Frame attached as modal window to parent
func (f *Frame) SetModal(parent *Frame) *Frame {
	C.gtk_window_set_transient_for(C.to_GtkWindow(f.window), C.to_GtkWindow(parent.window))
	C.gtk_window_set_destroy_with_parent(C.to_GtkWindow(f.window), C.TRUE)
	C.gtk_window_set_attached_to(C.to_GtkWindow(f.window), parent.window)
	C.gtk_window_set_modal(C.to_GtkWindow(f.window), C.TRUE)
	return f
}

// UnsetModal unset current Frame as modal window from another Frames
func (f *Frame) UnsetModal() *Frame {
	C.gtk_window_set_transient_for(C.to_GtkWindow(f.window), nil)
	C.gtk_window_set_destroy_with_parent(C.to_GtkWindow(f.window), C.FALSE)
	C.gtk_window_set_attached_to(C.to_GtkWindow(f.window), nil)
	C.gtk_window_set_modal(C.to_GtkWindow(f.window), C.FALSE)
	return f
}

// SetDecorated of window
func (f *Frame) SetDecorated(decorated bool) *Frame {
	C.gtk_window_set_decorated(C.to_GtkWindow(f.window), gboolean(decorated))
	return f
}

// SetDeletable of window
func (f *Frame) SetDeletable(deletable bool) *Frame {
	C.gtk_window_set_deletable(C.to_GtkWindow(f.window), gboolean(deletable))
	return f
}

// KeepAbove the window
func (f *Frame) KeepAbove(above bool) *Frame {
	C.gtk_window_set_keep_above(C.to_GtkWindow(f.window), gboolean(above))
	return f
}

// KeepBelow of window
func (f *Frame) KeepBelow(below bool) *Frame {
	C.gtk_window_set_keep_below(C.to_GtkWindow(f.window), gboolean(below))
	return f
}

// Show window
func (f *Frame) Show() *Frame {
	C.gtk_window_present(C.to_GtkWindow(f.window))
	if f.deferMove {
		f.Move(f.deferMoveX, f.deferMoveY)
	}
	return f
}

// Hide window
func (f *Frame) Hide() *Frame {
	C.gtk_window_close(C.to_GtkWindow(f.window))
	return f
}

// Iconify window
func (f *Frame) Iconify(iconify bool) *Frame {
	if iconify {
		C.gtk_window_iconify(C.to_GtkWindow(f.window))
	} else {
		C.gtk_window_deiconify(C.to_GtkWindow(f.window))
	}
	return f
}

// Stick window
func (f *Frame) Stick(stick bool) *Frame {
	if stick {
		C.gtk_window_stick(C.to_GtkWindow(f.window))
	} else {
		C.gtk_window_unstick(C.to_GtkWindow(f.window))
	}
	return f
}

// Maximize window
func (f *Frame) Maximize(maximize bool) *Frame {
	if maximize {
		C.gtk_window_maximize(C.to_GtkWindow(f.window))
	} else {
		C.gtk_window_unmaximize(C.to_GtkWindow(f.window))
	}
	return f
}

// Fullscreen window
func (f *Frame) Fullscreen(fullscreen bool) *Frame {
	if fullscreen {
		C.gtk_window_fullscreen(C.to_GtkWindow(f.window))
	} else {
		C.gtk_window_unfullscreen(C.to_GtkWindow(f.window))
	}
	return f
}

// SetIconFromFile for Frame
func (f *Frame) SetIconFromFile(filename string) *Frame {
	C.gtk_window_set_icon_from_file(C.to_GtkWindow(f.window), C.gcharptr(C.CString(filename)), nil)
	return f
}

// SetIconName for Frame
func (f *Frame) SetIconName(name string) *Frame {
	C.gtk_window_set_icon_name(C.to_GtkWindow(f.window), C.gcharptr(C.CString(name)))
	return f
}

// SetBackgroundColor of Frame
func (f *Frame) SetBackgroundColor(r, g, b int, alfa float64) *Frame {
	C.setBackgroundColor(f.window, f.webview, C.gint(C.int(r)), C.gint(C.int(g)), C.gint(C.int(b)), C.gdouble(alfa))
	return f
}

// SetMaxSize of window
func (f *Frame) SetMaxSize(width, height int) *Frame {
	C.setMaxSize(f.window, C.gint(C.int(width)), C.gint(C.int(height)))
	return f
}

// SetMinSize of window
func (f *Frame) SetMinSize(width, height int) *Frame {
	C.setMinSize(f.window, C.gint(C.int(width)), C.gint(C.int(height)))
	return f
}

// SetOpacity of window
func (f *Frame) SetOpacity(opacity float64) *Frame {
	C.gdk_window_set_opacity(C.gtk_widget_get_window(f.window), C.gdouble(opacity))
	return f
}

// SetType of window
func (f *Frame) SetType(hint WindowType) *Frame {
	C.gtk_window_set_type_hint(C.to_GtkWindow(f.window), C.GdkWindowTypeHint(int(hint)))
	return f
}

// GetScreen where the window placed
func (f *Frame) GetScreen() *Screen {
	screen := C.gtk_widget_get_screen(f.window)
	display := C.gdk_screen_get_display(screen)
	monitor := C.gdk_display_get_monitor_at_window(display, C.gtk_widget_get_window(f.window))
	return &Screen{
		screen:  screen,
		display: display,
		monitor: monitor,
	}
}

// GetSize returns width and height of window
func (f *Frame) GetSize() (width, height int) {
	var cWidth, cHeight C.gint
	C.gtk_window_get_size(C.to_GtkWindow(f.window), &cWidth, &cHeight)
	width, height = int(cWidth), int(cHeight)
	return
}

// Strut reserves frame space on the screen
func (f *Frame) Strut(position StrutPosition, size int) *Frame {
	monitorWidth, monitorHeight := f.GetScreen().Size()
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
	f.Stick(true).SetType(TypeDock)
	return f
}

//		gtk_window_begin_resize_drag ()
//		gtk_window_begin_move_drag ()
//		gtk_window_set_urgency_hint ()
//		gtk_window_set_accept_focus ()
//		gtk_window_set_focus_on_map ()
//		gtk_window_set_startup_id ()
//		gtk_window_reshow_with_initial_size ()
//		gtk_window_set_focus_visible ()
//		gtk_window_set_has_user_ref_count ()
//		gtk_window_set_auto_startup_notification ()
//		gtk_window_set_titlebar
//		gtk_window_add_accel_group ()
//		gtk_window_remove_accel_group ()
//		gtk_window_set_geometry_hints ()
//		gtk_window_set_gravity ()
//		gtk_window_set_hide_titlebar_when_maximized ()
//		gtk_window_set_screen ()
//		gtk_window_add_mnemonic ()
//		gtk_window_remove_mnemonic ()
//		gtk_window_set_focus ()
//		gtk_window_set_default_icon_list ()
//		gtk_window_set_default_icon ()
//		gtk_window_set_icon ()
//		gtk_window_set_icon_list ()
