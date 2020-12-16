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

BOOL isFocused(WindowObj ww);
BOOL isVisible(WindowObj ww);
BOOL isZoomed(WindowObj ww);
BOOL isMiniaturized(WindowObj ww);
BOOL isFullscreen(WindowObj ww);
void makeApp();
void makeWindow(char* title, int width, int height, long long unsigned int req_id, int id);
void evalJS(WindowObj ww, const char* js, long long unsigned int reqid);
void hideWindow(WindowObj ww);
void loadHTML(WindowObj ww, char* content, char* baseUrl);
void loadURI(WindowObj ww, char* uri);
void resizeWindow(WindowObj ww, int width, int height);
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

@interface WindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler>
@property (assign) int goWindowID;
@end

@interface AppDelegate : NSObject <NSApplicationDelegate>
@end

void triggerEvent(int goWindowID, NSWindow* movedWindow, NSString* eventTitle);

#endif // !WEBVIEW_H
#endif /* WEBVIEW_COCOA */