#if defined(WEBVIEW_COCOA)

#import "darwin.h"
#include "_cgo_export.h"

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
- (BOOL)windowShouldClose:(NSWindow*)window {
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
@end

void triggerEvent(int goWindowID, NSWindow* movedWindow, NSString* eventTitle)
{
    if ([movedWindow isKeyWindow]) {
        NSRect rect = movedWindow.frame;
        int x = (int)(rect.origin.x);
        int y = (int)(rect.origin.y);
        int w = (int)(rect.size.width);
        int h = (int)(rect.size.height);
        // NSLog(@"%@ %@", eventTitle, movedWindow);
        goWindowEvent(goWindowID, strdup([eventTitle UTF8String]), x, y, w, h);
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
		[app activateIgnoringOtherApps:YES];
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
			[window center];

            // Webwiew
            WKWebViewConfiguration* conf = [[WKWebViewConfiguration alloc] init];
            WKWebView* webview = [[WKWebView alloc] initWithFrame:r configuration:conf];
            webviews[id] = webview;
            [webview setAutoresizesSubviews:YES];
            [webview setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
            // NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:webview_check_url(w->url)]];
            [[window contentView] addSubview:webview];
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
		NSRect old = [window frame];
		NSRect r = NSMakeRect(old.origin.x - (width - old.size.width) / 2, old.origin.y - (height - old.size.height) / 2, width, height);
		[window setFrame:r display:YES animate:YES];
    });

    windowsUsed++;
    return id;
}

void showWindow(int id)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
        // [window makeKeyAndOrderFront:app];
		[window makeKeyAndOrderFront:app];
        // [window center];
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
        NSWindow* modalWindow = windows[id];
        NSWindow* window = windows[id2];
		modalWindow.level= NSModalPanelWindowLevel;
		// [window setIgnoresMouseEvents:YES];
		// [modalWindow setIgnoresMouseEvents:NO];
		// window.level= NSMainMenuWindowLevel;
	});
}

void lock(){
	dispatch_async(dispatch_get_main_queue(), ^(void) {});
}

void unsetModal(int id)
{
    // dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
		window.level= NSNormalWindowLevel;
	// });
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


void loadUri(int id, char* uri){
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
		WKWebView* webview = webviews[id];
		NSURL *url = [NSURL URLWithString:[NSString stringWithUTF8String:uri]];
		NSURLRequest *request = [NSURLRequest requestWithURL:url];
		[webview loadRequest:request];
    });
}

void loadHTML(int id, char* content, char* baseUrl){
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSWindow* window = windows[id];
		WKWebView* webview = webviews[id];
		NSString *htmlString = [NSString stringWithUTF8String:content];
		NSURL *baseURL = [NSURL URLWithString:[NSString stringWithUTF8String:baseUrl]];
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