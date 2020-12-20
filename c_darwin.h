#if defined(WEBVIEW_COCOA)
#ifndef WEBVIEW_H
#define WEBVIEW_H
#include <Cocoa/Cocoa.h>
#include <WebKit/WebKit.h>
#include <objc/runtime.h>

typedef struct WindowObj {
    int id;
    long long unsigned int req_id;
    BOOL created;
    NSWindow* window;
    WKWebView* webview;
    // GtkWidget* menubar;
} WindowObj;

typedef struct MenuObj {
    char* title;
    char* key;
    NSMenu* menu;
    NSMenuItem* menuItem;
} MenuObj;

typedef struct AppMenu {
    NSMenu* mainMenu;
    NSMenu* appMenu;
} AppMenu;

BOOL isFocused(WindowObj ww);
BOOL isVisible(WindowObj ww);
BOOL isZoomed(WindowObj ww);
BOOL isMiniaturized(WindowObj ww);
BOOL isFullscreen(WindowObj ww);
void makeApp(char* appName);
void makeWindow(char* title, int width, int height, long long unsigned int req_id, int id);
void evalJS(WindowObj ww, const char* js, long long unsigned int reqid);
void hideWindow(WindowObj ww);
void loadHTML(WindowObj ww, char* content, char* baseUrl);
void loadURI(WindowObj ww, char* uri);
void resizeWindow(WindowObj ww, int width, int height);
void resizeContent(WindowObj ww, int width, int height);
void setBackgroundColor(WindowObj ww, int8_t r, int8_t g, int8_t b, double a, bool titlebarTransparent);
void setMaxWindowSize(WindowObj ww, int width, int height);
void setMinWindowSize(WindowObj ww, int width, int height);
void setModal(WindowObj ww, WindowObj parent);
void setTitle(WindowObj ww, char* title);
void setWindowCenter(WindowObj ww);
void setWindowResizeble(WindowObj ww, bool resizeble);
void showWindow(WindowObj ww);
void unsetModal(WindowObj ww);
void moveWindow(WindowObj ww, int x, int y);
void iconifyWindow(WindowObj ww, bool flag);
void setWindowDecorated(WindowObj ww, bool flag);
void setWindowDeletable(WindowObj ww, bool flag);
void toggleFullScreen(WindowObj ww);
void stickWindow(WindowObj ww, bool flag);
void setWindowSkipPager(WindowObj ww, bool flag);
void setWindowSkipTaskbar(WindowObj ww, bool flag);
void windowKeepAbove(WindowObj ww, bool flag);
void windowKeepBelow(WindowObj ww, bool flag);
void setWindowAlpha(WindowObj ww, double opacity);
void toggleMaximize(WindowObj ww);
void setWindowIconFromFile(WindowObj ww, char* filename);
void setAppIconFromFile(char* filename);
MenuObj addSubMenu(MenuObj mm);
MenuObj addItem(MenuObj mm);
MenuObj addSeparatorItem(MenuObj mm);
CGSize windowSize(WindowObj ww);
CGSize contentSize(WindowObj ww);
CGSize getScreenSize(WindowObj ww);
double getScreenScale(WindowObj ww);
CGPoint windowPosition(WindowObj ww);
void setWindowZoom(WindowObj ww, double zoom);

@interface WindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler>
@property (assign) int goWindowID;
@end

@interface AppDelegate : NSObject <NSApplicationDelegate>
@end

void triggerEvent(int goWindowID, NSWindow* movedWindow, NSString* eventTitle);

#endif // !WEBVIEW_H
#endif /* WEBVIEW_COCOA */