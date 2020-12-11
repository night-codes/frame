
#if defined(WEBVIEW_GTK)
#ifndef WEBVIEW_H
#define WEBVIEW_H

#include <JavaScriptCore/JavaScript.h>
#include <X11/Xlib.h>
#include <stdlib.h>
#include <string.h>
#include <webkit2/webkit2.h>

typedef struct WindowObj {
    GtkWidget* window;
    GtkWidget* box;
    GtkWidget* webview;
    GtkWidget* menubar;
} WindowObj;

extern void goAppActivated();
extern void goPrint(char* text);
extern void goPrintInt(int num);
extern void goScriptEvent();
extern void goWindowState(GtkWidget* c, int e);
extern void goInvokeCallback(GtkWidget* webview, char* data);
extern void goWinRet(long long unsigned reqid, WindowObj* win);
extern void goEvalRet(long long unsigned reqid, char* err);

typedef enum {
    PANEL_WINDOW_POSITION_TOP,
    PANEL_WINDOW_POSITION_BOTTOM,
    PANEL_WINDOW_POSITION_LEFT,
    PANEL_WINDOW_POSITION_RIGHT
} winPosition;

typedef struct idleData {
    GtkWidget* widget;
    GtkWidget* widget2;
    gchar* content;
    gchar* uri;
    int width;
    int height;
    int x;
    int y;
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

static int webCount = 1;
static GtkApplication* app;
static GtkWidget** webviews;
static bool* webviewUsed;

static void stateEvent(GtkWidget* c, GdkEventWindowState* event)
{
    goWindowState(c, event->new_window_state);
}

static void scriptEvent(GtkWidget* ww, char* n)
{
    goScriptEvent();
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(ww), "GetEvents();", NULL, NULL, NULL);
}

static void external_message_received_cb(WebKitUserContentManager* contentManager,
    WebKitJavascriptResult* r,
    gpointer arg)
{
    (void)contentManager;
    GtkWidget* w = (GtkWidget*)arg;
    JSGlobalContextRef context = webkit_javascript_result_get_global_context(r);
    JSValueRef value = webkit_javascript_result_get_value(r);
    JSStringRef js = JSValueToStringCopy(context, value, NULL);
    size_t n = JSStringGetMaximumUTF8CStringSize(js);
    char* s = g_new(char, n);
    JSStringGetUTF8CString(js, s, n);
    goInvokeCallback(w, s);
    JSStringRelease(js);
    g_free(s);
}

// The application is started.
static void started(GtkApplication* app, gpointer user_data)
{
    gtk_application_window_new(app); // default window (without window application will be closed)
    goAppActivated(app); // call back
}

// The application is started.
static void makeApp(int count)
{
    XInitThreads();
    gtk_init(0, NULL);
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
    g_application_run(G_APPLICATION(app), 0, NULL);
}

static void updateVisual(GtkWidget* window)
{
    GdkScreen* screen = gtk_widget_get_screen(window);
    GdkVisual* visual = gdk_screen_get_rgba_visual(screen);
    if (visual) {
        gtk_widget_set_visual(window, visual);
        gtk_widget_set_app_paintable(window, TRUE);
    }
}

static void setZoom(GtkWidget* widget, gdouble zoom)
{
    webkit_web_view_set_zoom_level(WEBKIT_WEB_VIEW(widget), 2);
}

static void evalJSFinished(GObject* object, GAsyncResult* result, gpointer arg)
{
    idleData* data = (idleData*)arg;
    WebKitJavascriptResult* js_result;
    JSCValue* value;
    GError* error = NULL;
    js_result = webkit_web_view_run_javascript_finish(WEBKIT_WEB_VIEW(data->widget), result, &error);
    if (!js_result && error != NULL) {
        goEvalRet(data->req_id, strdup(error->message));
        //g_warning ("Error running javascript: %s", error->message);
        g_error_free(error);
        return;
    }
    goEvalRet(data->req_id, strdup(""));
}

static gboolean evalJS_idle(gpointer arg)
{
    idleData* data = (idleData*)arg;
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(data->widget), data->content, NULL, evalJSFinished, data);
    return FALSE;
}

static void evalJS(idleData* data)
{
    gdk_threads_add_idle(evalJS_idle, data);
}

static gboolean loadUri_idle(gpointer arg)
{
    idleData* data = (idleData*)arg;
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(data->widget), data->uri);
    return FALSE;
}

static void loadUri(idleData* data)
{
    gdk_threads_add_idle(loadUri_idle, data);
}

static gboolean loadHTML_idle(gpointer arg)
{
    idleData* data = (idleData*)arg;
    webkit_web_view_load_html(WEBKIT_WEB_VIEW(data->widget), data->content, data->uri);
    return FALSE;
}

static void loadHTML(idleData* data)
{
    gdk_threads_add_idle(loadHTML_idle, data);
}

static void setBackgroundColor(GtkWidget* window, GtkWidget* webview, gint r, gint g, gint b, gdouble alfa)
{
    GdkRGBA rgba;
    rgba.red = (gdouble)r / 255;
    rgba.green = (gdouble)g / 255;
    rgba.blue = (gdouble)b / 255;
    rgba.alpha = alfa;

    updateVisual(window);
    webkit_web_view_set_background_color(WEBKIT_WEB_VIEW(webview), &rgba);
}

static void setMaxSize(GtkWidget* window, gint width, gint height)
{
    GdkGeometry geometry;
    geometry.max_width = width;
    geometry.max_height = height;
    gtk_window_set_geometry_hints(GTK_WINDOW(window), NULL, &geometry, GDK_HINT_MAX_SIZE);
}

static void setMinSize(GtkWidget* window, gint width, gint height)
{
    GdkGeometry geometry;
    geometry.min_width = width;
    geometry.min_height = height;
    gtk_window_set_geometry_hints(GTK_WINDOW(window), NULL, &geometry, GDK_HINT_MIN_SIZE);
}

/* panel_window_reset_strut */
static void windowStrut(GdkWindow* window, winPosition position, int width, int height, int monitorWidth, int monitorHeight, int scale)
{
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
    gdk_property_change(window, atom, cardinal, 32, GDK_PROP_MODE_REPLACE,
        (guchar*)strut, 4);
    atom = gdk_atom_intern("_NET_WM_STRUT_PARTIAL", FALSE);
    gdk_property_change(window, atom, cardinal, 32, GDK_PROP_MODE_REPLACE,
        (guchar*)strut, 12);
}

static gboolean contextMenuEvent(WebKitWebView* web_view, WebKitContextMenu* context_menu,
    GdkEvent* event, WebKitHitTestResult* hit_test_result, gpointer user_data)
{
    return TRUE;
}

static gboolean makeWindow_idle(gpointer arg)
{
    idleData* data = (idleData*)arg;

    /** WINDOW  */
    GtkWidget* window = gtk_application_window_new(app);
    gtk_window_set_title(GTK_WINDOW(window), data->content);
    gtk_window_set_default_size(GTK_WINDOW(window), data->width, data->height);
    g_signal_connect(window, "delete-event", G_CALLBACK(gtk_widget_hide_on_delete), window);
    g_signal_connect(window, "window-state-event", G_CALLBACK(stateEvent), window);
    g_signal_connect(window, "screen-changed", G_CALLBACK(updateVisual), window);
    updateVisual(window);
    gtk_widget_realize(window);

    /** BOX  */
    GtkWidget* box = gtk_box_new(GTK_ORIENTATION_VERTICAL, 0);
    gtk_container_add(GTK_CONTAINER(window), box);
    gtk_widget_realize(box);
    gtk_widget_show(box);

    /** MENUBAR  */
    GtkWidget* menu = gtk_menu_new();
    GtkWidget* item = gtk_menu_item_new_with_label(gcharptr("Файл"));
    GtkWidget* item3 = gtk_menu_item_new_with_label(gcharptr("Опции"));
    GtkWidget* item4 = gtk_menu_item_new_with_label(gcharptr("Справка"));
    GtkWidget* item5 = gtk_menu_item_new_with_label(gcharptr("Открыть..."));
    GtkWidget* item6 = gtk_menu_item_new_with_label(gcharptr("Выход"));

    GtkWidget* fileMenu = gtk_menu_new();
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(item), fileMenu);

    /*  GtkWidget* child = gtk_bin_get_child(GTK_BIN(item));
    gtk_label_set_markup(GTK_LABEL(child), "<i>new label</i> with <b>markup</b>");
    gtk_accel_label_set_accel(GTK_ACCEL_LABEL(child), GDK_KEY_1, 0);
 */
    GtkWidget* menubar = gtk_menu_bar_new();
    gtk_box_pack_start(GTK_BOX(box), menubar, 0, 1, 0);
    gtk_menu_shell_append(GTK_MENU_SHELL(menubar), item);
    gtk_menu_shell_append(GTK_MENU_SHELL(menubar), item3);
    gtk_menu_shell_append(GTK_MENU_SHELL(menubar), item4);
    gtk_menu_shell_append(GTK_MENU_SHELL(fileMenu), item5);
    gtk_menu_shell_append(GTK_MENU_SHELL(fileMenu), item6);
    gtk_widget_realize(menubar);
    gtk_widget_show(menubar);

    /** WEBVIEW  */
    WebKitUserContentManager* contentManager = webkit_user_content_manager_new();
    webkit_user_content_manager_register_script_message_handler(contentManager, "external");
    GtkWidget* webview = webkit_web_view_new_with_user_content_manager(contentManager);
    g_signal_connect(contentManager, "script-message-received::external", G_CALLBACK(external_message_received_cb), webview);

    WebKitSettings* settings = webkit_settings_new();
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
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(webview), "window.external={invoke:function(x){window.webkit.messageHandlers.external.postMessage(x);}}", NULL, NULL, NULL);
    // webkit_settings_set_hardware_acceleration_policy(settings, WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS);
    // webkit_settings_set_enable_write_console_messages_to_stdout(settings, TRUE);
    // webkit_settings_set_enable_developer_extras(settings, TRUE);
    webkit_web_view_set_zoom_level(WEBKIT_WEB_VIEW(webview), 1.0);
    webkit_web_view_set_settings(WEBKIT_WEB_VIEW(webview), settings);
    g_signal_connect(webview, "send-script", G_CALLBACK(scriptEvent), NULL);
    g_signal_connect(webview, "context-menu", G_CALLBACK(contextMenuEvent), NULL);
    gtk_box_pack_end(GTK_BOX(box), webview, 1, 1, 0);
    // gtk_box_set_center_widget(GTK_BOX(box), webview);
    gtk_widget_show(webview);

    WindowObj* ret = (WindowObj*)malloc(sizeof(WindowObj));
    ret->window = window;
    ret->box = box;
    ret->webview = webview;
    ret->menubar = menubar;
    goWinRet(data->req_id, ret);
    return FALSE;
}

static void makeWindow(idleData* data)
{
    gdk_threads_add_idle(makeWindow_idle, data);
}

#endif // !WEBVIEW_H
#endif // WEBVIEW_GTK