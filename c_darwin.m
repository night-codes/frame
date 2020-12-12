#if defined(WEBVIEW_COCOA)

#import "_cgo_export.h"
#import "c_darwin.h"

static int webCount = 1;
static bool appInitialized = false; // false first time function is called
static NSApplication* app;

WindowDelegate* windowDelegate = nil;
AppDelegate* appDelegate = nil;

@implementation WindowDelegate
- (void)dealloc
{
    [super dealloc];
}
- (void)windowDidResize:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowDidResize");
}
- (void)windowDidMove:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowDidMove");
}
- (void)windowWillMiniaturize:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowWillMiniaturize");
}
- (void)windowDidMiniaturize:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowDidMiniaturize");
}
- (void)windowDidDeminiaturize:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowDidDeminiaturize");
}
- (BOOL)windowShouldClose:(NSWindow*)window
{
    triggerEvent([self goWindowID], window, @"windowShouldClose");
    return YES;
}
- (void)windowDidBecomeKey:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowDidBecomeKey");
}
- (void)windowDidResignKey:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowDidResignKey");
}
- (void)windowWillClose:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowWillClose");
}
- (void)userContentController:(WKUserContentController*)userContentController didReceiveScriptMessage:(WKScriptMessage*)message
{
    int id = [self goWindowID];
    triggerEvent(id, windows[id], [NSString stringWithFormat:@"invoke:%@", [message body]]);
}
@end

void triggerEvent(int goWindowID, NSWindow* window, NSString* eventTitle)
{
    NSRect rect = window.frame;
    int x = (int)(rect.origin.x);
    int y = (int)(rect.origin.y);
    int w = (int)(rect.size.width);
    int h = (int)(rect.size.height);
    goWindowEvent(goWindowID, strdup([eventTitle UTF8String]), x, y, w, h);
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
        // [app activateIgnoringOtherApps:YES];
        windows = malloc(sizeof(NSWindow*) * webCount); // create windows pool
        webviews = malloc(sizeof(WKWebView*) * webCount); // create webviews pool
        windowsUsed = 0;
        NSWindow.allowsAutomaticWindowTabbing = NO;
        for (int winID = 0; winID < webCount; winID++) {
            NSRect r = NSMakeRect(0, 0, 100, 100);
            NSUInteger mask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
            windows[winID] = [[[NSWindow alloc] initWithContentRect:r styleMask:mask backing:NSBackingStoreBuffered defer:NO] autorelease];

            // Window
            NSWindow* window = windows[winID];
            windowDelegate = [[WindowDelegate alloc] init];
            [windowDelegate setGoWindowID:winID];
            [window setDelegate:windowDelegate];
            [window center];

            // Webwiew
            WKWebViewConfiguration* conf = [[WKWebViewConfiguration alloc] init];
            WKUserContentController* ucc = [[WKUserContentController alloc] init];
            WKUserScript* us = [[WKUserScript alloc] initWithSource:
                                                         @"window.external={invoke:function(v){window.webkit.messageHandlers.invoke.postMessage(v)}};"
                                                      injectionTime:WKUserScriptInjectionTimeAtDocumentStart
                                                   forMainFrameOnly:NO];
            [ucc addUserScript:us];
            [ucc addScriptMessageHandler:windowDelegate name:@"invoke"];
            [conf setUserContentController:ucc];
            [[conf preferences] setValue:[NSNumber numberWithBool:YES] forKey:@"developerExtrasEnabled"];

            WKWebView* webview = [[WKWebView alloc] initWithFrame:r configuration:conf];
            webviews[winID] = webview;
            [webview setAutoresizesSubviews:YES];
            [webview setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
            // NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:webview_check_url(w->url)]];
            [[window contentView] addSubview:webview];
        }

        [app run];
    }

    appInitialized = true;
}

int makeWindow(char* title, int width, int height)
{
    __block int id = windowsUsed;
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window setTitle:[NSString stringWithUTF8String:title]];
        NSRect old = [window frame];
        NSRect r = NSMakeRect(old.origin.x - (width - old.size.width) / 2, old.origin.y - (height - old.size.height) / 2, width, height);
        [window setFrame:r display:YES animate:YES];
    });

    windowsUsed++;
    return id;
}

void evalJS(int winid, const char* js, long long unsigned int reqid)
{
    WKWebView* webview = webviews[winid];
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [webview evaluateJavaScript:[NSString stringWithUTF8String:js]
                  completionHandler:^(id self, NSError* error) {
                      if (error != NULL) {
                          goEvalRet(reqid, strdup([[NSString stringWithFormat:@"%@", error.userInfo] UTF8String]));
                      } else {
                          goEvalRet(reqid, strdup(""));
                      }
                  }];
    });
    // goEvalRet(reqid, strdup("ttttttttttttt"));
}

void showWindow(int id)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window orderFrontRegardless];
        // [window makeKeyAndOrderFront:app];
    });
}

void hideWindow(int id)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window orderOut:window];
    });
}

BOOL isFocused(int id)
{
    NSWindow* window = windows[id];
    return [window isKeyWindow];
}

BOOL isVisible(int id)
{
    NSWindow* window = windows[id];
    return [window isVisible];
}

BOOL isZoomed(int id)
{
    NSWindow* window = windows[id];
    return [window isZoomed];
}

BOOL isMiniaturized(int id)
{
    NSWindow* window = windows[id];
    return [window isMiniaturized];
}

BOOL isFullscreen(int id)
{
    NSWindow* window = windows[id];
    return ([window styleMask] & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen;
}

void resizeWindow(int id, int width, int height)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];

        NSRect old = [window frame];
        NSRect r = NSMakeRect(old.origin.x - (width - old.size.width) / 2, old.origin.y - (height - old.size.height) / 2, width, height);
        [window setFrame:r display:YES animate:YES];
    });
}

void setModal(int id, int id2)
{
    // dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(.5 * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{//TODO
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        NSWindow* parentWindow = windows[id2];
        window.level = NSModalPanelWindowLevel;
		[parentWindow makeMainWindow];
		[parentWindow addChildWindow:window ordered:NSWindowAbove];
        window.styleMask &= ~NSWindowStyleMaskMiniaturizable;

        // boxes[id].isUserInteractionEnabled = NO;
        // [parentWindow setIgnoresMouseEvents:YES];
    });
}

void lock()
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
                   });
}

void unsetModal(int id)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        window.styleMask |= NSWindowStyleMaskMiniaturizable;
        window.level = NSNormalWindowLevel;
		NSWindow* parentWindow = [window parentWindow];
		if (parentWindow != NULL) {
			[parentWindow removeChildWindow:window];
		}
    });
}

void setWindowCenter(int id)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window center];
    });
}

void setTitle(int id, char* title)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        [window setTitle:[NSString stringWithUTF8String:title]];
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

void loadUri(int id, char* uri)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        WKWebView* webview = webviews[id];
        NSURL* url = [NSURL URLWithString:[NSString stringWithUTF8String:uri]];
        NSURLRequest* request = [NSURLRequest requestWithURL:url];
        [webview loadRequest:request];
    });
}

void loadHTML(int id, char* content, char* baseUrl)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        WKWebView* webview = webviews[id];
        NSString* htmlString = [NSString stringWithUTF8String:content];
        NSURL* baseURL = [NSURL URLWithString:[NSString stringWithUTF8String:baseUrl]];
        [webview loadHTMLString:htmlString baseURL:baseURL];
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