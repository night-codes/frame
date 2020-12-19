#if defined(WEBVIEW_COCOA)
#ifndef WEBVIEW_N
#define WEBVIEW_N

#import "c_darwin.h"
#import "_cgo_export.h"

static bool appInitialized = false;
static bool menuInitialized = false;
static NSApplication* app;
static NSMenu* mainMenu;
static NSMenu* windowsMenu;
static NSMenu* appMenu;
static char* appName;

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
    [window orderOut:app];
    return NO;
}
- (void)windowDidExpose:(NSNotification*)notification
{
    NSWindow* window = notification.object;
    triggerEvent([self goWindowID], window, @"windowDidExpose");
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
- (void)userContentController:(WKUserContentController*)userContentController
      didReceiveScriptMessage:(WKScriptMessage*)message
{
    triggerEvent([self goWindowID], NULL, [NSString stringWithFormat:@"invoke:%@", [message body]]);
}
@end

void triggerEvent(int goWindowID, NSWindow* window, NSString* eventTitle)
{
    if (window == NULL) {
        goWindowEvent(goWindowID, strdup([eventTitle UTF8String]), 0, 0, 0, 0);
    } else {
        NSRect rect = window.frame;
        int x = (int)(rect.origin.x);
        int y = (int)(rect.origin.y);
        int w = (int)(rect.size.width);
        int h = (int)(rect.size.height);
        goWindowEvent(goWindowID, strdup([eventTitle UTF8String]), x, y, w, h);
    }
}

@implementation AppDelegate
- (void)dealloc
{
    [super dealloc];
}

- (void)menuAction:(NSMenuItem*)menuItem
{
	MenuObj mm;
	mm.menuItem = menuItem;
	goMenuFunc(mm);
}

- (void)applicationDidFinishLaunching:(NSNotification*)notification
{
	mainMenu = [[[NSMenu alloc] initWithTitle:@""] autorelease];
	[app setMainMenu:mainMenu];
	NSMenuItem* appMenuItem = [NSMenuItem new];
    [mainMenu addItem:appMenuItem];
    appMenu = [NSMenu new];
	[appMenuItem setSubmenu:appMenu];

	windowsMenu = [[[NSMenu alloc] initWithTitle:@"Window"] autorelease];
	[app setWindowsMenu:windowsMenu];
	NSMenuItem* windowsMenuItem = [NSMenuItem new];
    [windowsMenu addItem:windowsMenuItem];
	[windowsMenuItem setSubmenu:[NSMenu new]];

	[app setActivationPolicy:NSApplicationActivationPolicyRegular];
	[app activateIgnoringOtherApps:YES];
	AppMenu send = {mainMenu,appMenu};
    goAppActivated(send);
}
@end

// The application is started.
void makeApp(char* aName)
{
    if (appInitialized) {
        return;
    }
    appInitialized = true;

	appName = aName;

    app = [NSApplication sharedApplication];
    @autoreleasepool {
        appDelegate = [[AppDelegate alloc] init];
        [app setDelegate:appDelegate];
        NSWindow.allowsAutomaticWindowTabbing = NO;
        [app run];
    }
}


MenuObj addSubMenu(MenuObj mm)
{
	NSMenuItem* aMenuItem = [NSMenuItem new];
	[mm.menu addItem:aMenuItem];
	NSMenu* aMenu = [[NSMenu alloc] initWithTitle:[NSString stringWithUTF8String:mm.title]];
	[aMenuItem setTitle:[NSString stringWithUTF8String:mm.title]];
	dispatch_async(dispatch_get_main_queue(), ^(void) {
		[aMenuItem setSubmenu:aMenu];
	});

	mm.menu = aMenu;
	mm.menuItem = aMenuItem;
	return mm;
}

MenuObj addItem(MenuObj mm)
{
	NSMenuItem* aMenuItem = [[NSMenuItem alloc] initWithTitle:[NSString stringWithUTF8String:mm.title] action:@selector(menuAction:) keyEquivalent:[NSString stringWithUTF8String:mm.key]];
	[aMenuItem setTitle:[NSString stringWithUTF8String:mm.title]];
	dispatch_async(dispatch_get_main_queue(), ^(void) {
		[mm.menu addItem:aMenuItem];
	});

	mm.menuItem = aMenuItem;
	return mm;
}

void makeWindow(char* title, int width, int height, long long unsigned int req_id, int id)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        WindowObj ret;
        ret.id = id;
        ret.req_id = req_id;
        ret.created = FALSE;

        NSRect r = NSMakeRect(0, 0, width, height);
        NSUInteger mask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;

        // Window
        NSWindow* window = [[NSWindow alloc] initWithContentRect:r styleMask:mask backing:NSBackingStoreBuffered defer:NO];
        ret.window = window;
        windowDelegate = [[WindowDelegate alloc] init];
        [windowDelegate setGoWindowID:id];
        [window setDelegate:windowDelegate];
        [window center];

		/*
        NSDockTile *dockTile = [window dockTile];
		[dockTile display]; */

        // Webwiew
        WKWebViewConfiguration* conf = [[WKWebViewConfiguration alloc] init];
        WKUserContentController* ucc = [[WKUserContentController alloc] init];
        WKUserScript* us = [[WKUserScript alloc] initWithSource:@"window.external={invoke:function(v){window.webkit.messageHandlers.invoke.postMessage(v)}};" injectionTime:WKUserScriptInjectionTimeAtDocumentStart forMainFrameOnly:NO];
        [ucc addUserScript:us];
        [ucc addScriptMessageHandler:windowDelegate name:@"invoke"];
        [conf setUserContentController:ucc];
        [[conf preferences] setValue:[NSNumber numberWithBool:YES] forKey:@"developerExtrasEnabled"];

        WKWebView* webview = [[WKWebView alloc] initWithFrame:r configuration:conf];
        ret.webview = webview;
        [webview setAutoresizesSubviews:YES];
        [webview setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
        [[window contentView] addSubview:webview];
        [window setTitle:[NSString stringWithUTF8String:title]];

        goWinRet(req_id, ret);
    });
}

void evalJS(WindowObj ww, const char* js, long long unsigned int reqid)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.webview evaluateJavaScript:[NSString stringWithUTF8String:js]
                     completionHandler:^(id self, NSError* error) {
                         if (error != NULL) {
                             goEvalRet(reqid, strdup([[NSString stringWithFormat:@"%@", error.userInfo] UTF8String]));
                         } else {
                             goEvalRet(reqid, strdup(""));
                         }
                     }];
    });
}

void showWindow(WindowObj ww)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window makeKeyAndOrderFront:app];

		if (appMenu != NULL && !menuInitialized){
   			menuInitialized = true;
			[appMenu setTitle:[[NSString stringWithUTF8String:appName] stringByAppendingString:@"\x1b"]];
		}
    });
}

void hideWindow(WindowObj ww)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window orderOut:ww.window];
    });
}

void iconifyWindow(WindowObj ww, bool flag)
{
	dispatch_async(dispatch_get_main_queue(), ^(void) {
        if (flag) {
            [ww.window miniaturize:ww.window];
        } else {
            [ww.window deminiaturize:ww.window];
        }
    });
}

BOOL isFocused(WindowObj ww)
{
    return [ww.window isKeyWindow];
}

BOOL isVisible(WindowObj ww)
{
    return [ww.window isVisible];
}

BOOL isZoomed(WindowObj ww)
{
    return [ww.window isZoomed];
}

BOOL isMiniaturized(WindowObj ww)
{
    return [ww.window isMiniaturized];
}

BOOL isFullscreen(WindowObj ww)
{
    return ([ww.window styleMask] & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen;
}

void resizeWindow(WindowObj ww, int width, int height)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSRect old = [ww.window frame];
        NSRect r = NSMakeRect(old.origin.x - (width - old.size.width) / 2,
            old.origin.y - (height - old.size.height) / 2, width,
            height);
        [ww.window setFrame:r display:YES animate:YES];
    });
}

// Put window location from the top left corner of screen.
void moveWindow(WindowObj ww, int x, int y)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSRect old = [ww.window frame];
		NSScreen *screen = ww.window.screen;
		NSRect frame = screen.frame;
		[ww.window setFrameTopLeftPoint:NSMakePoint(x, frame.size.height-y)];
    });
}

void setModal(WindowObj ww, WindowObj parent)
{
    // dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(.5 * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{//TODO
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        ww.window.level = NSModalPanelWindowLevel;
        ww.window.styleMask &= ~NSWindowStyleMaskMiniaturizable;
        // ww.window.styleMask |= NSWindowStyleMaskDocModalWindow;
		[parent.window addChildWindow:ww.window ordered:NSWindowAbove];
    });
}

void unsetModal(WindowObj ww)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        ww.window.styleMask |= NSWindowStyleMaskMiniaturizable;
        // ww.window.styleMask &= ~NSWindowStyleMaskDocModalWindow;
        ww.window.level = NSNormalWindowLevel;
		NSWindow* parentWindow = [ww.window parentWindow];
		if (parentWindow != NULL) {
			[parentWindow removeChildWindow:ww.window];
		}
    });
}

void setWindowCenter(WindowObj ww)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window center];
    });
}

void setTitle(WindowObj ww, char* title)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window setTitle:[NSString stringWithUTF8String:title]];
    });
}

void setWindowIconFromFile(WindowObj ww, char* filename)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
		NSImage* img = [[NSImage alloc] initWithContentsOfFile:[NSString stringWithUTF8String:filename]];
		[ww.window setRepresentedURL:[NSURL URLWithString:[NSString stringWithUTF8String:""]]];
		[[ww.window standardWindowButton:NSWindowDocumentIconButton] setImage:img];
		// if (img != nil) {
		// 	[app setApplicationIconImage:img];
		// }
    });
}

void setWindowAlpha(WindowObj ww, double opacity)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window setAlphaValue:opacity];
    });
}

void setMaxWindowSize(WindowObj ww, int width, int height)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window setMaxSize:NSMakeSize(width, height)];
    });
}

void setMinWindowSize(WindowObj ww, int width, int height)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window setMinSize:NSMakeSize(width, height)];
    });
}

void loadURI(WindowObj ww, char* uri)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSURL* url = [NSURL URLWithString:[NSString stringWithUTF8String:uri]];
        NSURLRequest* request = [NSURLRequest requestWithURL:url];
        [ww.webview loadRequest:request];
    });
}

void loadHTML(WindowObj ww, char* content, char* baseUrl)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        NSString* htmlString = [NSString stringWithUTF8String:content];
        NSURL* baseURL = [NSURL URLWithString:[NSString stringWithUTF8String:baseUrl]];
        [ww.webview loadHTMLString:htmlString baseURL:baseURL];
    });
}

void setWindowResizeble(WindowObj ww, bool flag)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        if (flag) {
            ww.window.styleMask |= NSWindowStyleMaskResizable;
        } else {
            ww.window.styleMask &= ~NSWindowStyleMaskResizable;
        }
    });
}

void setWindowDecorated(WindowObj ww, bool flag)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        if (flag) {
            ww.window.styleMask |= NSWindowStyleMaskTitled;
        } else {
            ww.window.styleMask &= ~NSWindowStyleMaskTitled;
        }
    });
}

void setWindowDeletable(WindowObj ww, bool flag)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        if (flag) {
            ww.window.styleMask |= NSWindowStyleMaskClosable;
        } else {
            ww.window.styleMask &= ~NSWindowStyleMaskClosable;
        }
    });
}

void setWindowSkipPager(WindowObj ww, bool flag)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        if (flag) {
			ww.window.collectionBehavior |= NSWindowCollectionBehaviorIgnoresCycle;
			ww.window.collectionBehavior |= NSWindowCollectionBehaviorTransient;
        } else {
			ww.window.collectionBehavior &= ~NSWindowCollectionBehaviorIgnoresCycle;
			ww.window.collectionBehavior &= ~NSWindowCollectionBehaviorTransient;
        }
    });
}

void setWindowSkipTaskbar(WindowObj ww, bool flag)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        if (flag) {
			ww.window.collectionBehavior |= NSWindowCollectionBehaviorTransient;
        } else {
			ww.window.collectionBehavior &= ~NSWindowCollectionBehaviorTransient;
        }
    });
}

void stickWindow(WindowObj ww, bool flag)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        if (flag) {
			ww.window.collectionBehavior |= NSWindowCollectionBehaviorCanJoinAllSpaces;
        } else {
			ww.window.collectionBehavior &= ~NSWindowCollectionBehaviorCanJoinAllSpaces;
        }
    });
}

void toggleFullScreen(WindowObj ww)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window toggleFullScreen:ww.window];
    });
}

void toggleMaximize(WindowObj ww)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window zoom:ww.window];
    });
}

void setBackgroundColor(WindowObj ww, int8_t r, int8_t g, int8_t b, double a,
    bool titlebarTransparent)
{
    dispatch_async(dispatch_get_main_queue(), ^(void) {
        [ww.window setBackgroundColor:[NSColor colorWithRed:(CGFloat)r / 255.0 green:(CGFloat)g / 255.0 blue:(CGFloat)b / 255.0 alpha:(CGFloat)a]];
        if (0.5 >= ((r / 255.0 * 299.0) + (g / 255.0 * 587.0) + (b / 255.0 * 114.0)) / 1000.0) {
            [ww.window setAppearance:[NSAppearance appearanceNamed:NSAppearanceNameVibrantDark]];
        } else {
            [ww.window setAppearance:[NSAppearance appearanceNamed:NSAppearanceNameVibrantLight]];
        }
        [ww.window setOpaque:NO];
        if (titlebarTransparent) {
            [ww.window setTitlebarAppearsTransparent:YES];
			/* ww.window.styleMask |= NSFullSizeContentViewWindowMask;
			[ww.window setMovableByWindowBackground:YES];

			NSVisualEffectView *vibrant=[[NSVisualEffectView alloc] initWithFrame:[[ww.window contentView] bounds]];
        	[vibrant setState:NSVisualEffectStateActive];
			[vibrant setAutoresizingMask:NSViewWidthSizable|NSViewHeightSizable];
			[vibrant setBlendingMode:NSVisualEffectBlendingModeBehindWindow];
			[vibrant setMaterial:NSVisualEffectMaterialDark];
			[[ww.window contentView] addSubview:vibrant positioned:NSWindowBelow relativeTo:NULL]; */
        }
        [ww.webview setValue:[NSNumber numberWithBool:NO] forKey:@"drawsBackground"];
    });
}

#endif // !WEBVIEW_N
#endif /* WEBVIEW_COCOA */