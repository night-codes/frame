#if defined(WEBVIEW_GTK)
#ifndef WEBVIEW_H
#define WEBVIEW_H

#include <JavaScriptCore/JavaScript.h>
#include <X11/Xlib.h>
#include <stdlib.h>
#include <glib.h>
#include <glib/gstdio.h>
#include <string.h>
#include <webkit2/webkit2.h>

typedef struct WindowObj {
    int id;
    long long unsigned int req_id;
    gboolean created;
    GtkWidget* window;
    GtkWidget* box;
    GtkWidget* webview;
    GtkWidget* menubar;
} WindowObj;

typedef struct MenuObj {
    char* title;
    char* key;
    GtkWidget* menu;
    GtkWidget* menuItem;
} MenuObj;

extern void goAppActivated(GtkApplication* app);
extern void goPrint(char* text);
extern void goPrintInt(int num);
extern void goScriptEvent();
extern void goWindowState(WindowObj* win, int e);
extern void goMenuFunc(GtkWidget* mm);
extern void goInvokeCallback(WindowObj* win, char* data);
extern void goWinRet(long long unsigned int reqid, WindowObj* win);
extern void goEvalRet(long long unsigned int reqid, char* err);

typedef enum {
    PANEL_WINDOW_POSITION_TOP,
    PANEL_WINDOW_POSITION_BOTTOM,
    PANEL_WINDOW_POSITION_LEFT,
    PANEL_WINDOW_POSITION_RIGHT
} winPosition;

typedef struct idleData {
    GtkApplication* app;
    GtkWidget* window;
    GtkWidget* windowParent;
    GtkWidget* webview;
    gchar* content;
    gchar* uri;
    int width;
    int height;
    int x;
    int y;
    int id;
    int hint;
    gdouble dbl;
    gboolean flag;
    GdkRGBA rgba;
    long long unsigned int req_id;
} idleData;

static inline gchar* gcharptr(const char* s) { return (gchar*)s; }
static gint to_gint(int num) { return (gint)num; }
static GtkMenu* to_GtkMenu(GtkWidget* m) { return GTK_MENU(m); }
static GtkMenuShell* to_GtkMenuShell(GtkWidget* m) { return GTK_MENU_SHELL(m); }
static GtkWindow* to_GtkWindow(GtkWidget* w) { return GTK_WINDOW(w); }
static GtkContainer* to_GtkContainer(GtkWidget* w) { return GTK_CONTAINER(w); }
static GtkBox* to_GtkBox(GtkWidget* w) { return GTK_BOX(w); }
static WebKitWebView* to_WebKitWebView(GtkWidget* w) { return WEBKIT_WEB_VIEW(w); }
static char* applicationName = "";

static void stateEvent(GtkWidget* c, GdkEventWindowState* event, gpointer arg)
{
    WindowObj* win = (WindowObj*)arg;
    goWindowState(win, event->new_window_state);
}

static void scriptEvent(GtkWidget* ww, char* n)
{
    goScriptEvent();
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(ww), "GetEvents();", NULL, NULL, NULL);
}

static void external_message_received_cb(WebKitUserContentManager* contentManager, WebKitJavascriptResult* r, gpointer arg)
{
    (void)contentManager;
    WindowObj* win = (WindowObj*)arg;
    JSGlobalContextRef context = webkit_javascript_result_get_global_context(r);
    JSValueRef value = webkit_javascript_result_get_value(r);
    JSStringRef js = JSValueToStringCopy(context, value, NULL);
    size_t n = JSStringGetMaximumUTF8CStringSize(js);
    char* s = g_new(char, n);
    JSStringGetUTF8CString(js, s, n);
    goInvokeCallback(win, strdup(s));
    JSStringRelease(js);
    g_free(s);
}

static void webview_load_changed_cb(WebKitWebView* webview, WebKitLoadEvent load_event, gpointer arg)
{
    WindowObj* data = (WindowObj*)arg;
    switch (load_event) {
    case WEBKIT_LOAD_FINISHED:
        webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(webview), "window.external={invoke:function(x){window.webkit.messageHandlers.external.postMessage(x);}}", NULL, NULL, NULL);
        if (data->created == FALSE) {
            data->created = TRUE;
            goWinRet(data->req_id, data);
        }
        break;
    }
}

// The application is started.
static void started(GtkApplication* app, gpointer user_data)
{
    gtk_application_window_new(app); // default window (without window application will be closed)
    goAppActivated(app); // call back
}

// The application is started.
static void makeApp(char* appName)
{
    applicationName = appName;
    XInitThreads();
    gtk_init(0, NULL);
    GtkApplication* app = gtk_application_new(NULL, 0);
    g_signal_connect(app, "activate", G_CALLBACK(started), NULL);
    g_application_run(G_APPLICATION(app), 0, NULL);
}

static void updateVisual(GtkWidget* window)
{
    while (gtk_events_pending())
        gtk_main_iteration();
    GdkScreen* screen = gtk_widget_get_screen(window);
    GdkVisual* visual = gdk_screen_get_rgba_visual(screen);
    if (visual) {
        gtk_widget_set_visual(window, visual);
        gtk_widget_set_app_paintable(window, TRUE);
    }
}

static gboolean windowSetModal(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_transient_for(GTK_WINDOW(data->window), GTK_WINDOW(data->windowParent));
    gtk_window_set_destroy_with_parent(GTK_WINDOW(data->window), TRUE);
    gtk_window_set_attached_to(GTK_WINDOW(data->window), data->windowParent);
    gtk_window_set_modal(GTK_WINDOW(data->window), TRUE);
    return FALSE;
}

static gboolean windowUnsetModal(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_transient_for(GTK_WINDOW(data->window), NULL);
    gtk_window_set_destroy_with_parent(GTK_WINDOW(data->window), FALSE);
    gtk_window_set_attached_to(GTK_WINDOW(data->window), NULL);
    gtk_window_set_modal(GTK_WINDOW(data->window), FALSE);
    return FALSE;
}

static gboolean windowSetIcon(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_icon_from_file(GTK_WINDOW(data->window), data->content, NULL);
    return FALSE;
}

static gboolean windowSetOpacity(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gdk_window_set_opacity(gtk_widget_get_window(data->window), data->dbl);
    return FALSE;
}

static gboolean windowSetType(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_type_hint(GTK_WINDOW(data->window), data->hint);
    return FALSE;
}

static gboolean windowSetCenter(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_position(GTK_WINDOW(data->window), GTK_WIN_POS_CENTER);
    return FALSE;
}

static gboolean windowSetBackgroundColor(gpointer arg)
{
    idleData* data = (idleData*)arg;
    updateVisual(data->window);
    webkit_web_view_set_background_color(WEBKIT_WEB_VIEW(data->webview), &data->rgba);
    return FALSE;
}

static gboolean windowSetTitle(gpointer arg) // title string
{
    idleData* data = (idleData*)arg;
    gtk_window_set_title(GTK_WINDOW(data->window), data->content);
    return FALSE;
}

static gboolean windowSetSize(gpointer arg) // width, height int
{
    idleData* data = (idleData*)arg;

    gint x;
    gint y;
    gtk_window_get_position(GTK_WINDOW(data->window), &x, &y);

    gint pWidth;
    gint pHeight;
    gtk_window_get_size(GTK_WINDOW(data->window), &pWidth, &pHeight);

    x = x + (pWidth - data->width) / 2;
    y = y + (pHeight - data->height) / 2;
    if (x < 0) {
        x = 0;
    }
    if (y < 0) {
        y = 0;
    }
    gtk_window_move(GTK_WINDOW(data->window), x, y);
    gtk_window_resize(GTK_WINDOW(data->window), data->width, data->height);

    return FALSE;
}

static gboolean windowMove(gpointer arg) // x, y int
{
    idleData* data = (idleData*)arg;
    gtk_window_move(GTK_WINDOW(data->window), data->x, data->y);
    return FALSE;
}

static gboolean windowSetDecorated(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_decorated(GTK_WINDOW(data->window), data->flag);
    return FALSE;
}

static gboolean windowSetDeletable(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_deletable(GTK_WINDOW(data->window), data->flag);
    return FALSE;
}

static gboolean windowKeepAbove(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_keep_above(GTK_WINDOW(data->window), data->flag);
    return FALSE;
}

static gboolean windowKeepBelow(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_keep_below(GTK_WINDOW(data->window), data->flag);
    return FALSE;
}

static gboolean windowIconify(gpointer arg)
{
    idleData* data = (idleData*)arg;
    if (data->flag) {
        gtk_window_iconify(GTK_WINDOW(data->window));
    } else {
        gtk_window_deiconify(GTK_WINDOW(data->window));
        gtk_window_present(GTK_WINDOW(data->window));
    }
    return FALSE;
}

static gboolean windowStick(gpointer arg)
{
    idleData* data = (idleData*)arg;
    if (data->flag) {
        gtk_window_stick(GTK_WINDOW(data->window));
    } else {
        gtk_window_unstick(GTK_WINDOW(data->window));
    }
    return FALSE;
}

static gboolean windowMaximize(gpointer arg)
{
    idleData* data = (idleData*)arg;
    if (data->flag) {
        gtk_window_maximize(GTK_WINDOW(data->window));
    } else {
        gtk_window_unmaximize(GTK_WINDOW(data->window));
    }
    return FALSE;
}

static gboolean windowFullscreen(gpointer arg)
{
    idleData* data = (idleData*)arg;
    if (data->flag) {
        gtk_window_fullscreen(GTK_WINDOW(data->window));
    } else {
        gtk_window_unfullscreen(GTK_WINDOW(data->window));
    }
    return FALSE;
}

static gboolean windowSkipTaskbar(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_skip_taskbar_hint(GTK_WINDOW(data->window), data->flag);
    return FALSE;
}

static gboolean windowSkipPager(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_skip_pager_hint(GTK_WINDOW(data->window), data->flag);
    return FALSE;
}

static gboolean windowSetResizeble(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_set_resizable(GTK_WINDOW(data->window), data->flag);
    return FALSE;
}

static gboolean windowShow(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_present(GTK_WINDOW(data->window));
    return FALSE;
}

static gboolean windowHide(gpointer arg)
{
    idleData* data = (idleData*)arg;
    gtk_window_close(GTK_WINDOW(data->window));
    return FALSE;
}

/*** EVAL JS ***/
static void evalJSFinished(GObject* object, GAsyncResult* result, gpointer arg)
{
    idleData* data = (idleData*)arg;
    WebKitJavascriptResult* js_result;
    JSCValue* value;
    GError* error = NULL;
    js_result = webkit_web_view_run_javascript_finish(WEBKIT_WEB_VIEW(data->webview), result, &error);
    if (!js_result && error != NULL) {
        goEvalRet(data->req_id, strdup(error->message));
        //g_warning ("Error running javascript: %s", error->message);
        g_error_free(error);
        return;
    }
    goEvalRet(data->req_id, strdup(""));
}

static gboolean evalJS(gpointer arg)
{
    idleData* data = (idleData*)arg;
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(data->webview), data->content, NULL, evalJSFinished, data);
    return FALSE;
}

static gboolean loadURI(gpointer arg)
{
    idleData* data = (idleData*)arg;
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(data->webview), data->uri);
    return FALSE;
}

static gboolean loadHTML(gpointer arg)
{
    idleData* data = (idleData*)arg;
    webkit_web_view_load_html(WEBKIT_WEB_VIEW(data->webview), data->content, data->uri);
    return FALSE;
}

/****************************/
/****************************/

static GdkRectangle getMonitorSize(GtkWidget* window)
{
    GdkScreen* screen = gtk_widget_get_screen(window);
    GdkDisplay* display = gdk_screen_get_display(screen);
    GdkMonitor* monitor = gdk_display_get_monitor_at_window(display, gtk_widget_get_window(window));
    GdkRectangle monitorSize;
    gdk_monitor_get_geometry(monitor, &monitorSize);
    return monitorSize;
}

static int getMonitorScaleFactor(GtkWidget* window)
{
    return gdk_window_get_scale_factor(gtk_widget_get_window(window));
}

static void setSizes(GtkWidget* window, gint max_width, gint max_height, gint min_width, gint min_height)
{
    if (min_width>0 || min_height > 0) {
        gtk_widget_set_size_request(window, min_width, min_height);
    }
    if (max_width>0 || max_height > 0) {
        GdkGeometry geometry;
        geometry.max_width = max_width;
        geometry.max_height = max_height;
        gtk_window_set_geometry_hints(GTK_WINDOW(window), NULL, &geometry, GDK_HINT_MAX_SIZE);
    }
}

/* panel_window_reset_strut */
static void windowStrut(GdkWindow* window, winPosition position, int width, int height, int monitorWidth, int monitorHeight, int scale)
{
    gtk_window_set_decorated(GTK_WINDOW(window), FALSE);
    gtk_window_resize(GTK_WINDOW(window), width, height);
    gtk_window_stick(GTK_WINDOW(window));
    gtk_window_set_type_hint(GTK_WINDOW(window), GDK_WINDOW_TYPE_HINT_DOCK);

    GdkAtom atom;
    GdkAtom cardinal;
    unsigned long strut[12];
    memset(&strut, 0, sizeof(strut));
    gdk_window_set_group(window, window);

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
    gdk_property_change(window, atom, cardinal, 32, GDK_PROP_MODE_REPLACE, (guchar*)strut, 4);
    atom = gdk_atom_intern("_NET_WM_STRUT_PARTIAL", FALSE);
    gdk_property_change(window, atom, cardinal, 32, GDK_PROP_MODE_REPLACE, (guchar*)strut, 12);
    gtk_window_set_gravity(GTK_WINDOW(window), GDK_GRAVITY_NORTH_WEST);
}

static gboolean contextMenuEvent(WebKitWebView* web_view, WebKitContextMenu* context_menu,
    GdkEvent* event, WebKitHitTestResult* hit_test_result, gpointer user_data)
{
    goPrint("contextMenuEvent");
    return TRUE;
}

static MenuObj addSubMenu(MenuObj mm)
{
    GtkWidget* aMenuItem = gtk_menu_item_new_with_label(gcharptr(mm.title));
    GtkWidget* aMenu = gtk_menu_new();
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(aMenuItem), aMenu);
    gtk_menu_shell_append(GTK_MENU_SHELL(mm.menu), aMenuItem);
    gtk_widget_show_all(mm.menu);

    mm.menu = aMenu;
    mm.menuItem = aMenuItem;
    return mm;
}

static MenuObj addItem(MenuObj mm)
{
    GtkWidget* aMenuItem = gtk_menu_item_new_with_label(gcharptr(mm.title));
    gtk_menu_shell_append(GTK_MENU_SHELL(mm.menu), aMenuItem);
    // gtk_accel_label_set_accel(GTK_ACCEL_LABEL(child), GDK_KEY_1, 0);
    gtk_widget_show_all(mm.menu);
    g_signal_connect_swapped(aMenuItem, "activate", G_CALLBACK(goMenuFunc), (gpointer)aMenuItem);
    mm.menuItem = aMenuItem;
    return mm;
}

static MenuObj addSeparatorItem(MenuObj mm)
{
    GtkWidget* aMenuItem = gtk_separator_menu_item_new();
    gtk_menu_shell_append(GTK_MENU_SHELL(mm.menu), aMenuItem);
    gtk_widget_show_all(mm.menu);
    mm.menuItem = aMenuItem;
    return mm;
}

static int initCookieManager(WebKitSettings* webkitSettings)
{
    if (!webkitSettings)
        return 0;

    WebKitCookieManager* cookiemanager = webkit_web_context_get_cookie_manager(webkit_web_context_get_default());
    int error = 0;
    gchar* home = getenv("HOME");
    gchar cookieDatabasePath[2048];
    g_snprintf(cookieDatabasePath, 2048, "%s/.config/%s", home, applicationName);
    if (!g_file_test(cookieDatabasePath, G_FILE_TEST_IS_DIR) || !g_access(cookieDatabasePath, /*S_IWUSR|S_IRUSR*/ 0755)) {
        error = g_mkdir_with_parents(cookieDatabasePath, 0755);
    }
    if (!error) {
        gchar cookieDatabase[2048];
        g_sprintf(cookieDatabase, "%s/cdb", cookieDatabasePath);
        webkit_cookie_manager_set_persistent_storage(cookiemanager, cookieDatabase, WEBKIT_COOKIE_PERSISTENT_STORAGE_SQLITE);
    } else {
        g_printerr("LOG-> Init: Failed to init cookie database\n");
        return 0;
    }

    WebKitCookieAcceptPolicy cookiePolicy = WEBKIT_COOKIE_POLICY_ACCEPT_ALWAYS;
    webkit_cookie_manager_set_accept_policy(cookiemanager, cookiePolicy);
}

static gboolean makeWindow_idle(gpointer arg)
{
    idleData* data = (idleData*)arg;
    WindowObj* ret = (WindowObj*)malloc(sizeof(WindowObj));
    ret->id = data->id;
    ret->req_id = data->req_id;
    ret->created = FALSE;

    /** WINDOW  */
    GtkWidget* window = gtk_application_window_new(data->app);
    ret->window = window;
    gtk_window_set_title(GTK_WINDOW(window), data->content);
    gtk_window_set_default_size(GTK_WINDOW(window), data->width, data->height);
    gtk_window_set_position(GTK_WINDOW(window), GTK_WIN_POS_CENTER);
    g_signal_connect(window, "delete-event", G_CALLBACK(gtk_widget_hide_on_delete), window);
    g_signal_connect(window, "window-state-event", G_CALLBACK(stateEvent), ret);
    g_signal_connect(window, "screen-changed", G_CALLBACK(updateVisual), window);
    updateVisual(window);
    gtk_widget_realize(window);

    /** BOX  */
    GtkWidget* box = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
    ret->box = box;
    gtk_container_add(GTK_CONTAINER(window), box);
    gtk_widget_realize(box);
    gtk_widget_show(box);

    /** MENUBAR  */
    GtkWidget* menubar = gtk_menu_bar_new();
    ret->menubar = menubar;
    gtk_box_pack_start(GTK_BOX(box), menubar, 0, 1, 0);
    gtk_widget_realize(menubar);
    gtk_widget_show(menubar);

    /** WEBVIEW  */
    WebKitUserContentManager* contentManager = webkit_user_content_manager_new();
    webkit_user_content_manager_register_script_message_handler(contentManager, "external");
    GtkWidget* webview = webkit_web_view_new_with_user_content_manager(contentManager);
    ret->webview = webview;
    g_signal_connect(contentManager, "script-message-received::external", G_CALLBACK(external_message_received_cb), ret);

    WebKitSettings* settings = webkit_settings_new();
    webkit_settings_set_allow_modal_dialogs(settings, TRUE);
    webkit_settings_set_default_charset(settings, "utf-8");
    webkit_settings_set_enable_webgl(settings, TRUE);
    webkit_settings_set_javascript_can_access_clipboard(settings, TRUE);
    webkit_settings_set_javascript_can_open_windows_automatically(settings, TRUE);
    webkit_settings_set_enable_webaudio(settings, TRUE);
    webkit_settings_set_allow_file_access_from_file_urls(settings, TRUE);
    webkit_settings_set_allow_universal_access_from_file_urls(settings, TRUE);
    webkit_settings_set_enable_java(settings, FALSE);
    webkit_settings_set_enable_resizable_text_areas(settings, FALSE);
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(webview), "window.external={invoke:function(x){window.webkit.messageHandlers.external.postMessage(x);}}", NULL, NULL, NULL);
    // webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS);
    // webkit_settings_set_enable_writes_console_messages_to_stdout(settings, TRUE);
    // webkit_settings_set_enable_developer_extras(settings, TRUE);
    webkit_web_view_set_zoom_level(WEBKIT_WEB_VIEW(webview), 1.0);
    webkit_web_view_set_settings(WEBKIT_WEB_VIEW(webview), settings);
    initCookieManager(settings);
    g_signal_connect(webview, "context-menu", G_CALLBACK(contextMenuEvent), ret);

    gtk_box_pack_end(GTK_BOX(box), webview, 1, 1, 0);
    // gtk_widget_grab_focus(webview);
    gtk_widget_show(webview);

    g_signal_connect(G_OBJECT(webview), "load-changed", G_CALLBACK(webview_load_changed_cb), ret);
    webkit_web_view_load_html(WEBKIT_WEB_VIEW(webview), "<body></body>", "about:blank");
    return FALSE;
}

static void makeWindow(idleData* data)
{
    gdk_threads_add_idle(makeWindow_idle, data);
}

#endif // !WEBVIEW_H
#endif // WEBVIEW_GTK
