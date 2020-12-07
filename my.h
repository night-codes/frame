#include <string.h>
#include <X11/Xlib.h>
#include <webkit2/webkit2.h>
#include <stdlib.h>

extern void goAppActivated();
extern void goPrint(char* text);
extern void goPrintInt(int num);
extern void goScriptEvent();
extern void goWindowState(GtkWidget* c, int e);

typedef enum {
	PANEL_WINDOW_POSITION_TOP,
	PANEL_WINDOW_POSITION_BOTTOM,
	PANEL_WINDOW_POSITION_LEFT,
	PANEL_WINDOW_POSITION_RIGHT
} winPosition;

static inline gchar* gcharptr(const char* s) { return (gchar*)s; }
static gint to_gint(int num) { return (gint)num; }
static GtkMenu* to_GtkMenu(GtkWidget* m) { return GTK_MENU(m); }
static GtkMenuShell* to_GtkMenuShell(GtkWidget* m) { return GTK_MENU_SHELL(m); }
static GtkWindow* to_GtkWindow(GtkWidget* w) { return GTK_WINDOW(w); }
static GtkContainer* to_GtkContainer(GtkWidget* w) { return GTK_CONTAINER(w); }
static GtkBox* to_GtkBox(GtkWidget* w) { return GTK_BOX(w); }
static WebKitWebView* to_WebKitWebView(GtkWidget* w) { return WEBKIT_WEB_VIEW(w); }

static int webCount = 1;
static GtkApplication *app;
static GtkWidget **webviews;
static bool *webviewUsed;


static void stateEvent(GtkWidget* c, GdkEventWindowState* event) {
	goWindowState(c, event->new_window_state);
}
static void scriptEvent(GtkWidget* ww, char *n) {
	goScriptEvent();
	webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(ww), "GetEvents();", NULL, NULL, NULL);
}
static gboolean contextMenuEvent(WebKitWebView *web_view, WebKitContextMenu *context_menu,
                                 GdkEvent *event, WebKitHitTestResult *hit_test_result, gpointer user_data) {
	return TRUE;
}


static GtkWidget* newWebkit() {
	GtkWidget *webview = webkit_web_view_new_with_context(webkit_web_context_get_default());
	WebKitSettings *settings = webkit_settings_new ();
	webkit_settings_set_allow_modal_dialogs(settings, TRUE);
	// webkit_settings_set_enable_smooth_scrolling(settings, TRUE);
	webkit_settings_set_default_charset(settings, "utf-8");
	webkit_settings_set_enable_webgl(settings, TRUE);
	webkit_settings_set_javascript_can_access_clipboard(settings, TRUE);
	webkit_settings_set_javascript_can_open_windows_automatically(settings, TRUE);
	webkit_settings_set_enable_webaudio(settings, TRUE);
	webkit_settings_set_allow_file_access_from_file_urls(settings, TRUE);
	webkit_settings_set_allow_universal_access_from_file_urls(settings, TRUE);
	webkit_settings_set_enable_java(settings, FALSE);
	webkit_settings_set_enable_resizable_text_areas(settings, FALSE);

	// webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS);
	// webkit_settings_set_enable_developer_extras(settings, TRUE);
	webkit_web_view_set_zoom_level(WEBKIT_WEB_VIEW(webview), 1.0);
	webkit_web_view_set_settings (WEBKIT_WEB_VIEW(webview), settings);
	g_signal_connect(webview, "send-script", G_CALLBACK(scriptEvent), NULL);
	g_signal_connect(webview, "context-menu", G_CALLBACK(contextMenuEvent), NULL);
	return webview;
}


// The application is started.
static void started(GtkApplication* app, gpointer user_data) {
	gtk_application_window_new(app); // default window (without window application will be closed)
	webviews = malloc(sizeof(GtkWidget*)*webCount); // create webviews pool
	webviewUsed = malloc(sizeof(bool) * webCount);
	for (int i = 0; i < webCount; i++) {
		webviews[i] =  newWebkit();
		webviewUsed[i] = FALSE;
	}
	goAppActivated(app); // call back
}

// The application is started.
static void makeApp(int count) {
	XInitThreads();
	g_signal_new("send-script",
	             G_TYPE_OBJECT, G_SIGNAL_RUN_FIRST,
	             0, NULL, NULL,
	             g_cclosure_marshal_VOID__POINTER,
	             G_TYPE_NONE, 1, G_TYPE_POINTER);

	if (count > 0) {
		webCount = count;
	}
	app = gtk_application_new(NULL, 0);
	g_signal_connect(app, "activate", G_CALLBACK(started), NULL);
	g_application_run (G_APPLICATION (app), 0, NULL);
}


static void preventDestroy(GtkWidget* window) {
	g_signal_connect (window, "delete-event", G_CALLBACK(gtk_widget_hide_on_delete), NULL);
	g_signal_connect (window, "window-state-event", G_CALLBACK(stateEvent), NULL);
}

static void updateVisual (GtkWidget *window) {
	GdkScreen *screen = gtk_widget_get_screen (window);
	GdkVisual *visual = gdk_screen_get_rgba_visual (screen);
	if (visual) {
		gtk_widget_set_visual(window, visual);
		gtk_widget_set_app_paintable (window, TRUE);
	}
}

static GtkWidget* makeWindow(char *name, int width, int height) {
	GtkWidget *window = gtk_application_window_new (app);
	gtk_window_set_title (GTK_WINDOW (window), name);
	gtk_window_set_default_size (GTK_WINDOW (window), width, height);
	preventDestroy(window);
	g_signal_connect(window, "screen-changed", G_CALLBACK(updateVisual), NULL);
	updateVisual(window);
	gtk_widget_realize (window);
	return window;
}



static GtkWidget* makeBox(GtkWidget *window) {
	GtkWidget *box = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
	gtk_container_add(GTK_CONTAINER(window), box);
	gtk_widget_realize(box);
	gtk_widget_show(box);
	return box;
}


static GtkWidget* makeWebview(GtkWidget *box) {
	GtkWidget* webview = NULL;
	for (int i = 0; i < webCount; i++) {
		if (!webviewUsed[i]) {
			webview = webviews[i];
			webviewUsed[i] = TRUE;
			gtk_box_pack_start(GTK_BOX(box), webview, 1, 1, 0);
			gtk_widget_show(webview);
			break;
		};
	}
	return webview;
}


static GtkWidget* makeMenubar(GtkWidget *box) {
	GtkWidget *menu = gtk_menu_new();
	GtkWidget *item = gtk_menu_item_new_with_label(gcharptr("Файл"));
	GtkWidget *item3 = gtk_menu_item_new_with_label(gcharptr("Опции"));
	GtkWidget *item4 = gtk_menu_item_new_with_label(gcharptr("Справка"));
	GtkWidget *menubar = gtk_menu_bar_new();
	gtk_box_pack_start(GTK_BOX(box), menubar, 0, 1, 0);
	gtk_menu_shell_append(GTK_MENU_SHELL(menubar), item);
	gtk_menu_shell_append(GTK_MENU_SHELL(menubar), item3);
	gtk_menu_shell_append(GTK_MENU_SHELL(menubar), item4);
	gtk_widget_realize(menubar);
	// gtk_widget_show_all(menubar);
	return menubar;
}






static void loadUri(GtkWidget *widget, gchar* uri) {
	webkit_web_view_load_uri(WEBKIT_WEB_VIEW(widget), uri);
}

static void loadHTML(GtkWidget *widget, gchar* content, gchar* base_uri) {
	webkit_web_view_load_html(WEBKIT_WEB_VIEW(widget), content, base_uri);
}

static void setZoom(GtkWidget *widget, gdouble zoom) {
	webkit_web_view_set_zoom_level(WEBKIT_WEB_VIEW(widget), 2);
}

static void setMaxSize(GtkWidget *window, gint width, gint height) {
	GdkGeometry geometry;
	geometry.max_width = width;
	geometry.max_height = height;
	gtk_window_set_geometry_hints(GTK_WINDOW(window), NULL, &geometry,  GDK_HINT_MIN_SIZE | GDK_HINT_MAX_SIZE);
}

static void setMinSize(GtkWidget *window, gint width, gint height) {
	GdkGeometry geometry;
	geometry.min_width = width;
	geometry.min_height = height;
	gtk_window_set_geometry_hints(GTK_WINDOW(window), NULL, &geometry,  GDK_HINT_MIN_SIZE | GDK_HINT_MAX_SIZE);
}



/* panel_window_reset_strut */
static void windowStrut(GdkWindow * window, winPosition position, int width, int height, int monitorWidth, int monitorHeight, int scale) {
	GdkAtom atom;
	GdkAtom cardinal;
	unsigned long strut[12];
	memset(&strut, 0, sizeof(strut));
	gdk_window_set_group(window, window);

	// strut = [ left, right, top, bottom,
	//           left_start_y, left_end_y, right_start_y, right_end_y,
	//           top_start_x, top_end_x, bottom_start_x, bottom_end_x ]
	switch (position) {
	case PANEL_WINDOW_POSITION_TOP:
		strut[2] = height * scale;
		strut[8] = 0;
		strut[9] = monitorWidth * scale;
		break;
	case PANEL_WINDOW_POSITION_BOTTOM:
		strut[3] = height * scale;
		strut[10] = 0;
		strut[11] = monitorWidth * scale;
		break;
	case PANEL_WINDOW_POSITION_LEFT:
		strut[0] = width * scale;
		strut[4] = 0;
		strut[5] = monitorHeight * scale;
		break;
	case PANEL_WINDOW_POSITION_RIGHT:
		strut[1] = width * scale;
		strut[6] = 0;
		strut[7] = monitorHeight * scale;
		break;
	}
	cardinal = gdk_atom_intern("CARDINAL", FALSE);
	atom = gdk_atom_intern("_NET_WM_STRUT", FALSE);
	gdk_property_change(window, atom, cardinal, 32, GDK_PROP_MODE_REPLACE,
	                    (guchar *)strut, 4);
	atom = gdk_atom_intern("_NET_WM_STRUT_PARTIAL", FALSE);
	gdk_property_change(window, atom, cardinal, 32, GDK_PROP_MODE_REPLACE,
	                    (guchar *)strut, 12);
}


static void setBackgroundColor (GtkWidget *window, GtkWidget *webview, gint r, gint g, gint b, gdouble alfa) {
	GdkRGBA rgba;
	rgba.red = (gdouble)r / 255;
	rgba.green = (gdouble)g / 255;
	rgba.blue = (gdouble)b / 255;
	rgba.alpha = alfa;

	updateVisual(window);
	webkit_web_view_set_background_color (WEBKIT_WEB_VIEW(webview), &rgba);
}
