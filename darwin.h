#if defined(WEBVIEW_COCOA)
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#import <objc/runtime.h>

extern void goAppActivated();
extern void goPrint(char* text);
extern void goPrintInt(int num);
extern void goScriptEvent();
extern void goEvalRet(long long unsigned int reqid, char* err);
// extern void goWindowState(GtkWidget *c, int e);

static WKWebView** webviews;
static NSWindow** windows;
static int windowsUsed;

void makeApp(int);
void runApp();
int makeWindow(char* title, int width, int height);
void showWindow(int);
void hideWindow(int);
void setTitle(int, char* title);
void setBackgroundColor(int, int8_t r, int8_t g, int8_t b, double a, bool titlebarTransparent);
void resizeWindow(int id, int width, int height);
void setMaxWindowSize(int id, int width, int height);
void setMinWindowSize(int id, int width, int height);
void setWindowResizeble(int id, bool resizeble);
void loadUri(int id, char* uri);
void loadHTML(int id, char* content, char* baseUrl);
BOOL isFocused(int id);
BOOL isVisible(int id);
BOOL isZoomed(int id);
BOOL isMiniaturized(int id);
BOOL isFullscreen(int id);
void setModal(int id, int id2);
void unsetModal(int id);
void setWindowCenter(int id);
void evalJS(int id, const char* js, long long unsigned int reqid);
void lock();

@interface WindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler>
@property (assign) int goWindowID;
@end

@interface AppDelegate : NSObject <NSApplicationDelegate>
@end

void triggerEvent(int goWindowID, NSWindow* movedWindow, NSString* eventTitle);

#endif /* WEBVIEW_COCOA */