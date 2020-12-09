#if defined(WEBVIEW_GTK)

#include <string.h>
#include <X11/Xlib.h>
#include <webkit2/webkit2.h>
#include <stdlib.h>

extern void goAppActivated();
extern void goPrint(char *text);
extern void goPrintInt(int num);
extern void goScriptEvent();
extern void goWindowState(GtkWidget *c, int e);

typedef enum
{
	PANEL_WINDOW_POSITION_TOP,
	PANEL_WINDOW_POSITION_BOTTOM,
	PANEL_WINDOW_POSITION_LEFT,
	PANEL_WINDOW_POSITION_RIGHT
} winPosition;

static inline gchar *gcharptr(const char *s) { return (gchar *)s; }
static gint to_gint(int num) { return (gint)num; }
static GtkMenu *to_GtkMenu(GtkWidget *m) { return GTK_MENU(m); }
static GtkMenuShell *to_GtkMenuShell(GtkWidget *m) { return GTK_MENU_SHELL(m); }
static GtkWindow *to_GtkWindow(GtkWidget *w) { return GTK_WINDOW(w); }
static GtkContainer *to_GtkContainer(GtkWidget *w) { return GTK_CONTAINER(w); }
static GtkBox *to_GtkBox(GtkWidget *w) { return GTK_BOX(w); }
static WebKitWebView *to_WebKitWebView(GtkWidget *w) { return WEBKIT_WEB_VIEW(w); }

static int webCount = 1;
static GtkApplication *app;
static GtkWidget **webviews;
static bool *webviewUsed;

static void stateEvent(GtkWidget *c, GdkEventWindowState *event);
static void scriptEvent(GtkWidget *ww, char *n);
static gboolean contextMenuEvent(WebKitWebView *web_view, WebKitContextMenu *context_menu, GdkEvent *event, WebKitHitTestResult *hit_test_result, gpointer user_data);
static GtkWidget *newWebkit();
static void started(GtkApplication *app, gpointer user_data);
static void makeApp(int count);
static void updateVisual(GtkWidget *window);
static GtkWidget *makeWindow(char *name, int width, int height);
static GtkWidget *makeBox(GtkWidget *window);
static GtkWidget *makeWebview(GtkWidget *box);
static GtkWidget *makeMenubar(GtkWidget *box);
static void loadUri(GtkWidget *widget, gchar *uri);
static void loadHTML(GtkWidget *widget, gchar *content, gchar *base_uri);
static void setZoom(GtkWidget *widget, gdouble zoom);
static void setMaxSize(GtkWidget *window, gint width, gint height);
static void setMinSize(GtkWidget *window, gint width, gint height);
static void windowStrut(GdkWindow *window, winPosition position, int width, int height, int monitorWidth, int monitorHeight, int scale);
static void setBackgroundColor(GtkWidget *window, GtkWidget *webview, gint r, gint g, gint b, gdouble alfa);

#endif