#if defined(WEBVIEW_COCOA)

#import "darwin.h"
#include "_cgo_export.h"

static int webCount = 1;
static bool appInitialized = false; // false first time function is called
static NSApplication* app;

const int DID_RESIZE_EVENT = 0;
const int DID_MOVE_EVENT = 1;
const int DID_MINIATURIZE_EVENT = 2;
const int DID_DEMINIATURIZE_EVENT = 3;

WindowDelegate* windowDelegate = nil;
AppDelegate* appDelegate = nil;

@implementation WindowDelegate
- (void)dealloc
{
    [super dealloc];
}
- (void)windowDidResize:(NSNotification*)notification
{
    NSWindow* movedWindow = notification.object;
    triggerEvent([self goWindowID], movedWindow, @"windowDidResize", DID_RESIZE_EVENT);
}
- (void)windowDidMove:(NSNotification*)notification
{
    NSWindow* movedWindow = notification.object;
    triggerEvent([self goWindowID], movedWindow, @"windowDidMove", DID_MOVE_EVENT);
}
- (void)windowDidMiniaturize:(NSNotification*)notification
{
    NSWindow* movedWindow = notification.object;
    triggerEvent([self goWindowID], movedWindow, @"windowDidMiniaturize", DID_MINIATURIZE_EVENT);
}
- (void)windowDidDeminiaturize:(NSNotification*)notification
{
    NSWindow* movedWindow = notification.object;
    triggerEvent([self goWindowID], movedWindow, @"windowDidDeminiaturize", DID_DEMINIATURIZE_EVENT);
}
@end

void triggerEvent(int goWindowID, NSWindow* movedWindow, NSString* eventTitle, const int eventId)
{
    if ([movedWindow isKeyWindow]) {
        NSRect rect = movedWindow.frame;
        int x = (int)(rect.origin.x);
        int y = (int)(rect.origin.y);
        int w = (int)(rect.size.width);
        int h = (int)(rect.size.height);
        // NSLog(@"%@ %@", eventTitle, movedWindow);
        onWindowEvent(goWindowID, eventId, x, y, w, h);
    }
}

@implementation AppDelegate
- (void)dealloc
{
    [super dealloc];
}
- (void)applicationDidFinishLaunching:(NSNotification*)aNotification
{
    goAppActivated();
}
@end

// The application is started.
void makeApp(int count)
{
    if (appInitialized)
        return;

    if (count > 0) {
        webCount = count;
    }
    app = [NSApplication sharedApplication];

    @autoreleasepool {
        [app setActivationPolicy:NSApplicationActivationPolicyRegular];
        appDelegate = [[AppDelegate alloc] init];
        [app setDelegate:appDelegate];
        windows = malloc(sizeof(NSWindow*) * webCount); // create windows pool
        webviews = malloc(sizeof(WKWebView*) * webCount); // create webviews pool
        windowsUsed = 0;
        NSWindow.allowsAutomaticWindowTabbing = NO;
        for (int id = 0; id < webCount; id++) {
            NSRect r = NSMakeRect(0, 0, 100, 100);
            NSUInteger mask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
            windows[id] = [[[NSWindow alloc] initWithContentRect:r styleMask:mask backing:NSBackingStoreBuffered defer:NO] autorelease];

            // Window
            NSWindow* window = windows[id];
            windowDelegate = [[WindowDelegate alloc] init];
            [windowDelegate setGoWindowID:id];
            [window setDelegate:windowDelegate];

            // Webwiew
            WKWebViewConfiguration* conf = [[WKWebViewConfiguration alloc] init];
            WKWebView* webview = [[WKWebView alloc] initWithFrame:r configuration:conf];
            webviews[id] = webview;
            [webview setAutoresizesSubviews:YES];
            [webview setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
            // NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:webview_check_url(w->url)]];
            // [webview loadRequest:[NSURLRequest requestWithURL:nsURL]];
            [[window contentView] addSubview:webview];
            // [window orderFrontRegardless];
        }

        [app run];
    }

    appInitialized = true;
}

int makeWindow(char* name, int width, int height)
{
    __block int id = windowsUsed;
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window setTitle:[NSString stringWithUTF8String:name]];
        NSRect frame = [window frame];
        frame.origin.y += frame.size.height;
        frame.origin.y -= height;
        frame.size = NSMakeSize(width, height);
        [window setFrame:frame display:YES];
    });

    windowsUsed++;
    return id;
}

void showWindow(int id)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window makeKeyAndOrderFront:app];
        [window center];
    });
}

void resizeWindow(int id, int width, int height)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        NSRect frame = [window frame];
        frame.origin.y += frame.size.height;
        frame.origin.y -= height;
        frame.size = NSMakeSize(width, height);
        [window setMaxSize:NSMakeSize(width, height)];
        [window setFrame:frame display:YES];
    });
}

void setMaxWindowSize(int id, int width, int height)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window setMaxSize:NSMakeSize(width, height)];
    });
}
void setMinWindowSize(int id, int width, int height)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window setMinSize:NSMakeSize(width, height)];
    });
}

void setWindowResizeble(int id, bool resizeble)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        if (resizeble) {
            window.styleMask |= NSWindowStyleMaskResizable;
        } else {
            window.styleMask &= ~NSWindowStyleMaskResizable;
        }
    });
}

void setBackgroundColor(int id, int8_t r, int8_t g, int8_t b, double a, bool titlebarTransparent)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window setBackgroundColor:[NSColor colorWithRed:(CGFloat)r / 255.0 green:(CGFloat)g / 255.0 blue:(CGFloat)b / 255.0 alpha:(CGFloat)a]];
        if (0.5 >= ((r / 255.0 * 299.0) + (g / 255.0 * 587.0) + (b / 255.0 * 114.0)) / 1000.0) {
            [window setAppearance:[NSAppearance appearanceNamed:NSAppearanceNameVibrantDark]];
        } else {
            [window setAppearance:[NSAppearance appearanceNamed:NSAppearanceNameVibrantLight]];
        }
        [window setOpaque:NO];
        if (titlebarTransparent) {
            [window setTitlebarAppearsTransparent:YES];
        }
        [webviews[id] setValue:[NSNumber numberWithBool:NO] forKey:@"drawsBackground"];
    });
}

#endif /* WEBVIEW_COCOA */