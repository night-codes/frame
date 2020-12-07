package frame

/*
#include "my.h"
*/
import "C"

type (
	// Frame struct
	Frame struct {
		Window     *C.GtkWidget
		Box        *C.GtkWidget
		Webview    *C.GtkWidget
		Menubar    *C.GtkWidget
		StateEvent func(State)
		deferMove  bool
		deferMoveX int
		deferMoveY int
	}
)

const (
	PosNone           = C.GTK_WIN_POS_NONE
	PosCenter         = C.GTK_WIN_POS_CENTER
	PosMouse          = C.GTK_WIN_POS_MOUSE
	PosCenterAlways   = C.GTK_WIN_POS_CENTER_ALWAYS
	PosCenterOnParent = C.GTK_WIN_POS_CENTER_ON_PARENT

	TypeNormal       = C.GDK_WINDOW_TYPE_HINT_NORMAL        // Normal toplevel window.
	TypeDialog       = C.GDK_WINDOW_TYPE_HINT_DIALOG        // Dialog window.
	TypeMenu         = C.GDK_WINDOW_TYPE_HINT_MENU          // Window used to implement a menu; GTK+ uses this hint only for torn-off menus, see GtkTearoffMenuItem.
	TypeToolbar      = C.GDK_WINDOW_TYPE_HINT_TOOLBAR       // Window used to implement toolbars.
	TypeSplashscreen = C.GDK_WINDOW_TYPE_HINT_SPLASHSCREEN  // Window used to display a splash screen during application startup.
	TypeUtility      = C.GDK_WINDOW_TYPE_HINT_UTILITY       // Utility windows which are not detached toolbars or dialogs.
	TypeDock         = C.GDK_WINDOW_TYPE_HINT_DOCK          // Used for creating dock or panel windows.
	TypeDesktop      = C.GDK_WINDOW_TYPE_HINT_DESKTOP       // Used for creating the desktop background window.
	TypeDropdownMenu = C.GDK_WINDOW_TYPE_HINT_DROPDOWN_MENU // A menu that belongs to a menubar.
	TypePopupMenu    = C.GDK_WINDOW_TYPE_HINT_POPUP_MENU    // A menu that does not belong to a menubar, e.g. a context menu.
	TypeTooltip      = C.GDK_WINDOW_TYPE_HINT_TOOLTIP       // A tooltip.
	TypeNotification = C.GDK_WINDOW_TYPE_HINT_NOTIFICATION  // A notification - typically a “bubble” that belongs to a status icon.
	TypeCombo        = C.GDK_WINDOW_TYPE_HINT_COMBO         // A popup from a combo box.
	TypeDnd          = C.GDK_WINDOW_TYPE_HINT_DND           // A window that is used to implement a DND cursor.

	GravityNorthWest = C.GDK_GRAVITY_NORTH_WEST // the reference point is at the top left corner.
	GravityNorth     = C.GDK_GRAVITY_NORTH      // the reference point is in the middle of the top edge.
	GravityNorthEast = C.GDK_GRAVITY_NORTH_EAST // the reference point is at the top right corner.
	GravityWest      = C.GDK_GRAVITY_WEST       // the reference point is at the middle of the left edge.
	GravityCenter    = C.GDK_GRAVITY_CENTER     // the reference point is at the center of the window.
	GravityEast      = C.GDK_GRAVITY_EAST       // the reference point is at the middle of the right edge.
	GravitySouthWest = C.GDK_GRAVITY_SOUTH_WEST // the reference point is at the lower left corner.
	GravitySouth     = C.GDK_GRAVITY_SOUTH      // the reference point is at the middle of the lower edge.
	GravitySouthEast = C.GDK_GRAVITY_SOUTH_EAST // the reference point is at the lower right corner.
	GravityStatic    = C.GDK_GRAVITY_STATIC     // the reference point is at the top left corner of the window itself, ignoring window manager decorations.

	StrutTop    = C.PANEL_WINDOW_POSITION_TOP
	StrutBottom = C.PANEL_WINDOW_POSITION_BOTTOM
	StrutLeft   = C.PANEL_WINDOW_POSITION_LEFT
	StrutRight  = C.PANEL_WINDOW_POSITION_RIGHT
)

// Load URL to Frame webview
func (f *Frame) Load(uri string) *Frame {
	C.loadUri(f.Webview, C.gcharptr(C.CString(uri)))
	// C.webkit_web_inspector_attach(C.webkit_web_view_get_inspector(C.to_WebKitWebView(f.Webview)))
	return f
}

// LoadHTML to Frame webview
func (f *Frame) LoadHTML(html string, baseUri string) *Frame {
	C.loadHTML(f.Webview, C.gcharptr(C.CString(html)), C.gcharptr(C.CString(baseUri)))
	return f
}

// SetModal of window
func (f *Frame) SetModal(modal bool) *Frame {
	C.gtk_window_set_modal(C.to_GtkWindow(f.Window), gboolean(modal))
	return f
}

// SkipTaskbar of window
func (f *Frame) SkipTaskbar(skip bool) *Frame {
	C.gtk_window_set_skip_taskbar_hint(C.to_GtkWindow(f.Window), gboolean(skip))
	return f
}

// SkipPager of window
func (f *Frame) SkipPager(skip bool) *Frame {
	C.gtk_window_set_skip_pager_hint(C.to_GtkWindow(f.Window), gboolean(skip))
	return f
}

// SetResizeble of Frame window
func (f *Frame) SetResizeble(resizeble bool) *Frame {
	C.gtk_window_set_resizable(C.to_GtkWindow(f.Window), gboolean(resizeble))
	return f
}

// SetStateEvent set handler function for window state event
func (f *Frame) SetStateEvent(fn func(State)) *Frame {
	f.StateEvent = fn
	return f
}

// SetTitle of Frame window
func (f *Frame) SetTitle(title string) *Frame {
	C.gtk_window_set_title(C.to_GtkWindow(f.Window), C.gcharptr(C.CString(title)))
	return f
}

// SetDefaultSize of Frame window
func (f *Frame) SetDefaultSize(width, height int) *Frame {
	C.gtk_window_set_default_size(C.to_GtkWindow(f.Window), C.gint(C.int(width)), C.gint(C.int(height)))
	return f
}

// Resize the window
func (f *Frame) Resize(width, height int) *Frame {
	C.gtk_window_resize(C.to_GtkWindow(f.Window), C.gint(C.int(width)), C.gint(C.int(height)))
	return f
}

// Move the window
func (f *Frame) Move(x, y int) *Frame {
	visible := C.gtk_widget_get_visible(f.Window) == 1
	if !visible {
		f.deferMove = true
		f.deferMoveX = x
		f.deferMoveY = y
		return f
	}
	C.gtk_window_move(C.to_GtkWindow(f.Window), C.gint(C.int(x)), C.gint(C.int(y)))
	return f
}

// SetPosition of Frame window
func (f *Frame) SetPosition(position C.GtkWindowPosition) *Frame {
	C.gtk_window_set_position(C.to_GtkWindow(f.Window), position)
	return f
}

// SetTransientFor of window. Dialog windows should be set transient for the main
// application window they were spawned from. This allows window managers to e.g.
// keep the dialog on top of the main window, or center the dialog over the main window
func (f *Frame) SetTransientFor(parent *Frame) *Frame {
	C.gtk_window_set_transient_for(C.to_GtkWindow(f.Window), C.to_GtkWindow(parent.Window))
	C.gtk_window_set_destroy_with_parent(C.to_GtkWindow(f.Window), C.TRUE)
	return f
}

// SetAttachedTo of Frame
func (f *Frame) SetAttachedTo(parent *Frame) *Frame {
	C.gtk_window_set_attached_to(C.to_GtkWindow(f.Window), parent.Window)
	return f
}

// AttachTo makes current Frame attached as modal window to parent
func (f *Frame) AttachTo(parent *Frame) *Frame {
	C.gtk_window_set_transient_for(C.to_GtkWindow(f.Window), C.to_GtkWindow(parent.Window))
	// C.gtk_window_set_destroy_with_parent(C.to_GtkWindow(f.Window), C.TRUE)
	// C.gtk_window_set_attached_to(C.to_GtkWindow(f.Window), parent.Window)
	C.gtk_window_set_modal(C.to_GtkWindow(f.Window), C.TRUE)
	return f
}

// Detach current Frame from another Frames
func (f *Frame) Detach() *Frame {
	C.gtk_window_set_transient_for(C.to_GtkWindow(f.Window), nil)
	// C.gtk_window_set_destroy_with_parent(C.to_GtkWindow(f.Window), C.FALSE)
	// C.gtk_window_set_attached_to(C.to_GtkWindow(f.Window), nil)
	C.gtk_window_set_modal(C.to_GtkWindow(f.Window), C.FALSE)
	return f
}

// SetDecorated of Frame window
func (f *Frame) SetDecorated(decorated bool) *Frame {
	C.gtk_window_set_decorated(C.to_GtkWindow(f.Window), gboolean(decorated))
	return f
}

// SetDeletable of Frame window
func (f *Frame) SetDeletable(deletable bool) *Frame {
	C.gtk_window_set_deletable(C.to_GtkWindow(f.Window), gboolean(deletable))
	return f
}

// KeepAbove the window
func (f *Frame) KeepAbove(above bool) *Frame {
	C.gtk_window_set_keep_above(C.to_GtkWindow(f.Window), gboolean(above))
	return f
}

// KeepBelow of window
func (f *Frame) KeepBelow(below bool) *Frame {
	C.gtk_window_set_keep_below(C.to_GtkWindow(f.Window), gboolean(below))
	return f
}

// Show Frame window
func (f *Frame) Show() *Frame {
	C.gtk_window_present(C.to_GtkWindow(f.Window))
	if f.deferMove {
		f.Move(f.deferMoveX, f.deferMoveY)
	}
	return f
}

// Hide Frame window
func (f *Frame) Hide() *Frame {
	C.gtk_window_close(C.to_GtkWindow(f.Window))
	return f
}

// Iconify Frame window
func (f *Frame) Iconify() *Frame {
	C.gtk_window_iconify(C.to_GtkWindow(f.Window))
	return f
}

// Deiconify Frame window
func (f *Frame) Deiconify() *Frame {
	C.gtk_window_deiconify(C.to_GtkWindow(f.Window))
	return f
}

// Stick Frame window
func (f *Frame) Stick() *Frame {
	C.gtk_window_stick(C.to_GtkWindow(f.Window))
	return f
}

// Unstick Frame window
func (f *Frame) Unstick() *Frame {
	C.gtk_window_unstick(C.to_GtkWindow(f.Window))
	return f
}

// Maximize Frame window
func (f *Frame) Maximize() *Frame {
	C.gtk_window_maximize(C.to_GtkWindow(f.Window))
	return f
}

// Unmaximize Frame window
func (f *Frame) Unmaximize() *Frame {
	C.gtk_window_unmaximize(C.to_GtkWindow(f.Window))
	return f
}

// Fullscreen Frame window
func (f *Frame) Fullscreen() *Frame {
	C.gtk_window_fullscreen(C.to_GtkWindow(f.Window))
	return f
}

// Unfullscreen Frame window
func (f *Frame) Unfullscreen() *Frame {
	C.gtk_window_unfullscreen(C.to_GtkWindow(f.Window))
	return f
}

// SetRole for Frame window
func (f *Frame) SetRole(role string) *Frame {
	C.gtk_window_set_role(C.to_GtkWindow(f.Window), C.gcharptr(C.CString(role)))
	return f
}

// SetIconFromFile for Frame
func (f *Frame) SetIconFromFile(filename string) *Frame {
	C.gtk_window_set_icon_from_file(C.to_GtkWindow(f.Window), C.gcharptr(C.CString(filename)), nil)
	return f
}

// SetIconName for Frame
func (f *Frame) SetIconName(name string) *Frame {
	C.gtk_window_set_icon_name(C.to_GtkWindow(f.Window), C.gcharptr(C.CString(name)))
	return f
}

// SetBackgroundColor of Frame
func (f *Frame) SetBackgroundColor(r, g, b int, alfa float64) *Frame {
	C.setBackgroundColor(f.Window, f.Webview, C.gint(C.int(r)), C.gint(C.int(g)), C.gint(C.int(b)), C.gdouble(alfa))
	return f
}

// SetMaxSize of Frame window
func (f *Frame) SetMaxSize(width, height int) *Frame {
	C.setMaxSize(f.Window, C.gint(C.int(width)), C.gint(C.int(height)))
	return f
}

// SetMinSize of Frame window
func (f *Frame) SetMinSize(width, height int) *Frame {
	C.setMinSize(f.Window, C.gint(C.int(width)), C.gint(C.int(height)))
	return f
}

// SetOpacity of Frame window
func (f *Frame) SetOpacity(opacity float64) *Frame {
	C.gdk_window_set_opacity(C.gtk_widget_get_window(f.Window), C.gdouble(opacity))
	return f
}

// SetLimitSizes of Frame window
func (f *Frame) SetLimitSizes(minWidth, minHeight, maxWidth, maxHeight int) *Frame {
	f.SetMaxSize(maxWidth, maxHeight)
	f.SetMinSize(minWidth, minHeight)
	return f
}

// SetLimitSizes of Frame window
func (f *Frame) SetType(hint C.GdkWindowTypeHint) *Frame {
	C.gtk_window_set_type_hint(C.to_GtkWindow(f.Window), hint)
	return f
}

// SetGravity of Frame window
func (f *Frame) SetGravity(gravity C.GdkGravity) *Frame {
	C.gtk_window_set_gravity(C.to_GtkWindow(f.Window), gravity)
	return f
}

// GetScreen where the window placed
func (f *Frame) GetScreen() *Screen {
	screen := C.gtk_widget_get_screen(f.Window)
	display := C.gdk_screen_get_display(screen)
	gdk_window := C.gtk_widget_get_window(f.Window)
	monitor := C.gdk_display_get_monitor_at_window(display, gdk_window)
	return &Screen{
		screen:  screen,
		display: display,
		monitor: monitor,
	}
}

// GetSize where the window placed
func (f *Frame) GetSize() (width, height int) {
	var cWidth, cHeight C.gint
	C.gtk_window_get_size(C.to_GtkWindow(f.Window), &cWidth, &cHeight)
	width, height = int(cWidth), int(cHeight)
	return
}

// Strut reserves frame space on the screen
func (f *Frame) Strut(strutPosition C.winPosition, size int) *Frame {
	monitorWidth, monitorHeight := f.GetScreen().Size()
	scale := f.GetScreen().Scale()
	var width, height int

	switch strutPosition {
	case StrutBottom, StrutTop:
		width, height = monitorWidth, size
	case StrutLeft, StrutRight:
		width, height = size, monitorHeight
	}
	f.
		SetDecorated(false).
		Resize(width, height).
		Stick().
		SetType(TypeDock)

	C.windowStrut(C.gtk_widget_get_window(f.Window), strutPosition, C.int(width), C.int(height), C.int(monitorWidth), C.int(monitorHeight), C.int(scale))
	f.SetGravity(GravityNorthWest)

	switch strutPosition {
	case StrutTop, StrutLeft:
		f.Move(0, 0)
	case StrutBottom:
		f.Move(0, monitorHeight-height)
	case StrutRight:
		f.Move(monitorWidth-width, 0)
	}
	f.Stick().SetType(TypeDock)
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
