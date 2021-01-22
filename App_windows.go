// +build windows

package frame

/*
	#cgo CFLAGS: -I./windows
	#include <stdlib.h>
	#include "string.h"
	#include "windows.h"
	#include "winuser.h"
	#include "include/capi/cef_app_capi.h"
	#include "handlers/cef_app.h"
	#include "handlers/cef_client.h"
	#include "handlers/cef_base.h"
	#include "include/capi/cef_client_capi.h"
	#include "include/capi/cef_browser_capi.h"
	#include "include/capi/cef_urlrequest_capi.h"
	#include "include/capi/cef_v8_capi.h"

	static void resizeWebview(HWND hwnd) {
		RECT* rect = (RECT*)malloc(sizeof(RECT));
		GetClientRect(hwnd, rect);
		HWND cefHwnd = GetWindow(hwnd, GW_CHILD);
		if (cefHwnd != NULL) {
			HDWP hdwp = BeginDeferWindowPos(1);
			hdwp = DeferWindowPos(hdwp, cefHwnd, NULL, rect->left, rect->top, rect->right - rect->left, rect->bottom - rect->top, SWP_NOZORDER);
			EndDeferWindowPos(hdwp);
		}
	}

	static void ExecuteJavaScript(cef_browser_t* browser, const char* code, const char* script_url, int start_line) {
		cef_frame_t * frame = browser->get_main_frame(browser);
		cef_string_t * codeCef = cef_string_userfree_utf16_alloc();
		cef_string_from_utf8(code, strlen(code), codeCef);
		cef_string_t * urlVal = cef_string_userfree_utf16_alloc();
		cef_string_from_utf8(script_url, strlen(script_url), urlVal);

		frame->execute_java_script(frame, codeCef, urlVal, start_line);

		cef_string_userfree_utf16_free(urlVal);
		cef_string_userfree_utf16_free(codeCef);
	}

	static void LoadURL(cef_browser_t* browser, cef_string_t* url) {
		cef_frame_t * frame = browser->get_main_frame(browser);
		frame->load_url(frame, url);
	}

	static void BrowserWasResized(cef_browser_t* browser) {
		cef_browser_host_t * host = browser->get_host(browser);
		host->was_resized(host);
	}

	static cef_window_handle_t GetWindowHandle(cef_browser_t* browser) {
		cef_browser_host_t * host = browser->get_host(browser);
		return host->get_window_handle(host);
	}

	static cef_client_t* GetClient(cef_browser_t* browser) {
		cef_browser_host_t * host = browser->get_host(browser);
		return host->get_client(host);
	}

	static cef_frame_t* GetMainFrame(cef_browser_t* browser) {
		cef_frame_t * frame = browser->get_main_frame(browser);
		return frame;
	}

	// Force close the browser
	static void CloseBrowser(cef_browser_t* browser) {
		cef_browser_host_t * host = browser->get_host(browser);
		host->close_browser(host, 1);
	}

	static int SendProcessMessage(cef_browser_t* browser, cef_process_message_t* message) {
		return browser->send_process_message(browser, PID_BROWSER, message);
		// return browser->send_process_message(browser, PID_RENDERER, message);
	}

	static cef_string_utf8_t * cefStringToUtf8(cef_string_t * source) {
		cef_string_utf8_t * output = cef_string_userfree_utf8_alloc();
		if (source == 0) {
			return output;
		}
		cef_string_to_utf8(source->str, source->length, output);
		return output;
	}

	static cef_string_t * GetURL(cef_browser_t* browser) {
		cef_frame_t * frame = browser->get_main_frame(browser);
		return frame->get_url(frame);
	}
*/
import "C"

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/night-codes/cefresources"
	"golang.org/x/sys/windows"
)

type (
	// App is main application object
	App struct {
		Test5      bool
		app        interface{} // *C.GtkApplication
		openedWns  sync.WaitGroup
		shown      chan bool
		mainArgs   C.cef_main_args_t
		appHandler C.cef_app_t
	}
)

var (
	cef2destroy = false
	mutexNew    sync.Mutex
	winds       = []*Window{}
	lock        sync.Mutex
	appChan     = make(chan *App)
	idItr       int64

	resDir   = cefresources.Extract()
	libcef   = windows.NewLazyDLL(filepath.Join(resDir, "libcef.dll"))
	kernel32 = windows.NewLazySystemDLL("kernel32")
	ole32    = windows.NewLazySystemDLL("ole32")
	user32   = windows.NewLazySystemDLL("user32")

	cefCreateBrowser        = libcef.NewProc("cef_browser_host_create_browser")
	cefInitialize           = libcef.NewProc("cef_initialize")
	cefExecuteProcess       = libcef.NewProc("cef_execute_process")
	cefEnableHDPI           = libcef.NewProc("cef_enable_highdpi_support")
	cefStringFromUTF8       = libcef.NewProc("cef_string_utf8_to_utf16")
	cefStringToUTF8         = libcef.NewProc("cef_string_utf16_to_utf8")
	cefAllocUTF8            = libcef.NewProc("cef_string_userfree_utf8_alloc")
	cefFreeUTF8             = libcef.NewProc("cef_string_userfree_utf8_free")
	cefQuitMessageLoop      = libcef.NewProc("cef_quit_message_loop")
	cefShutdown             = libcef.NewProc("cef_shutdown")
	cefRunMessageLoop       = libcef.NewProc("cef_run_message_loop")
	cefGetGlobalCtx         = libcef.NewProc("cef_request_context_get_global_context")
	cefProcessMessageCreate = libcef.NewProc("cef_process_message_create")

	winCoInitializeEx     = ole32.NewProc("CoInitializeEx")
	winGetProcessHeap     = kernel32.NewProc("GetProcessHeap")
	winHeapAlloc          = kernel32.NewProc("HeapAlloc")
	winHeapFree           = kernel32.NewProc("HeapFree")
	winGetCurrentThreadID = kernel32.NewProc("GetCurrentThreadId")
	winLoadImageW         = user32.NewProc("LoadImageW")
	winGetSystemMetrics   = user32.NewProc("GetSystemMetrics")
	winGetDpiForWindow    = user32.NewProc("GetDpiForWindow")
	winRegisterClassExW   = user32.NewProc("RegisterClassExW")
	winCreateWindowExW    = user32.NewProc("CreateWindowExW")
	winDestroyWindow      = user32.NewProc("DestroyWindow")
	winShowWindow         = user32.NewProc("ShowWindow")
	winUpdateWindow       = user32.NewProc("UpdateWindow")
	winSwitchToThisWindow = user32.NewProc("SwitchToThisWindow")
	winSetFocus           = user32.NewProc("SetFocus")
	winGetMessageW        = user32.NewProc("GetMessageW")
	winTranslateMessage   = user32.NewProc("TranslateMessage")
	winDispatchMessageW   = user32.NewProc("DispatchMessageW")
	winDefWindowProcW     = user32.NewProc("DefWindowProcW")
	winGetClientRect      = user32.NewProc("GetClientRect")
	winPostQuitMessage    = user32.NewProc("PostQuitMessage")
	winSetWindowTextW     = user32.NewProc("SetWindowTextW")
	winPostThreadMessageW = user32.NewProc("PostThreadMessageW")
	winGetWindowLongPtrW  = user32.NewProc("GetWindowLongPtrW")
	winSetWindowLongPtrW  = user32.NewProc("SetWindowLongPtrW")
	winAdjustWindowRect   = user32.NewProc("AdjustWindowRect")
	winSetWindowPos       = user32.NewProc("SetWindowPos")
	winRedrawWindow       = user32.NewProc("RedrawWindow")
	winMonitorFromWindow  = user32.NewProc("MonitorFromWindow")
	winGetMonitorInfo     = user32.NewProc("GetMonitorInfoA")

	lifeHandlers = map[uintptr]unsafe.Pointer{}
	cliReqs      = map[uintptr]uint64{}
)

const (
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	SWP_NOREDRAW       = 0x0008
	SWP_NOACTIVATE     = 0x0010
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOSENDCHANGING = 0x0400
	SWP_DRAWFRAME      = SWP_FRAMECHANGED
	SWP_NOREPOSITION   = SWP_NOOWNERZORDER
	SWP_DEFERERASE     = 0x2000
	SWP_ASYNCWINDOWPOS = 0x4000
	RDW_INVALIDATE     = 0x0001
	RDW_UPDATENOW      = 0x0100
)

// ShowWindow constants
const (
	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_MAXIMIZE        = 3
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
)

const (
	MONITOR_DEFAULTTONULL    = 0x00000000
	MONITOR_DEFAULTTOPRIMARY = 0x00000001
	MONITOR_DEFAULTTONEAREST = 0x00000002
)

// MakeApp is make and run one instance of application (At the moment, it is possible to create only one instance)
func MakeApp(appName string) *App {
	lock.Lock()

	app := App{
		mainArgs: C.cef_main_args_t{
			instance: C.GetModuleHandle(nil),
		},
		appHandler: C.cef_app_t{},
	}

	C.initialize_cef_app(&app.appHandler)
	cefEnableHDPI.Call()

	code, _, _ := cefExecuteProcess.Call(
		uintptr(unsafe.Pointer(&app.mainArgs)),
		uintptr(unsafe.Pointer(&app.appHandler)),
		0,
	)
	if int32(code) >= 0 {
		os.Exit(int(code))
	}

	var cefSettings *C.cef_settings_t
	cefSettings = (*C.cef_settings_t)(C.calloc(1, C.sizeof_cef_settings_t))
	cefSettings.size = C.sizeof_cef_settings_t
	// cefSettings.pack_loading_disabled = C.int(1)
	cefSettings.user_agent = *cefString("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	cefSettings.no_sandbox = C.int(1)
	cefSettings.single_process = C.int(1)
	cefSettings.multi_threaded_message_loop = C.int(1)
	cefSettings.context_safety_implementation = C.int(-1)
	cefSettings.log_severity = (C.cef_log_severity_t)(C.int(C.LOGSEVERITY_VERBOSE))
	cefSettings.cache_path = *cefString("")
	cefSettings.log_file = *cefString(resDir + "/log.txt")
	cefSettings.resources_dir_path = *cefString(resDir)
	cefSettings.locales_dir_path = *cefString(resDir)

	cefInitialize.Call(
		uintptr(unsafe.Pointer(&app.mainArgs)),
		uintptr(unsafe.Pointer(cefSettings)),
		uintptr(unsafe.Pointer(&app.appHandler)),
		uintptr(unsafe.Pointer(C.NULL)),
	)

	return &app
}

// SetIconFromFile  sets application icon
func (a *App) SetIconFromFile(filename string) {
	//C.gtk_window_set_default_icon_from_file(C.gcharptr(C.CString(filename)), nil)
}

// WaitAllWindowClose locker
func (a *App) WaitAllWindowClose() {
	fmt.Println("WaitAllWindowClose wait...")
	<-a.shown
	fmt.Println("WaitAllWindowClose - OK")
	// a.openedWns.Wait()
}

// WaitWindowClose locker
func (a *App) WaitWindowClose(win *Window) {
	<-a.shown
	shown := false
	for {
		if !win.state.Hidden {
			shown = true
		}
		if win.state.Hidden && shown {
			break
		}
	}
}

// NewWindow returns window with webview
func (a *App) NewWindow(title string, sizes ...int) *Window {
	mutexNew.Lock()
	defer mutexNew.Unlock()
	id := atomic.AddInt64(&idItr, 1)

	width := 500
	height := 400

	if len(sizes) > 0 {
		width = sizes[0]
	}

	if len(sizes) > 1 {
		height = sizes[1]
	}

	windowInfo := (*C.cef_window_info_t)(C.calloc(1, C.sizeof_cef_window_info_t))
	var thread uintptr
	cRet := cRequest(func(reqid uint64) {
		rect := C.RECT{}
		C.GetClientRect(C.GetDesktopWindow(), &rect)
		windowInfo.style = C.WS_OVERLAPPEDWINDOW | C.WS_TABSTOP // | C.WS_VISIBLE
		windowInfo.transparent_painting_enabled = C.int(1)
		windowInfo.height = C.int(width)
		windowInfo.width = C.int(height)
		windowInfo.window_name = *cefString(title)
		windowInfo.x = C.int(rect.right/2) - C.int(windowInfo.width/2)
		windowInfo.y = C.int(rect.bottom/2) - C.int(windowInfo.height/2)

		var client C.cef_client_t
		C.initialize_cef_client(&client)
		cliReqs[uintptr(unsafe.Pointer(&client))] = reqid

		var settings *C.cef_browser_settings_t
		settings = (*C.cef_browser_settings_t)(C.calloc(1, C.sizeof_cef_browser_settings_t))
		settings.size = C.sizeof_cef_browser_settings_t
		settings.javascript_access_clipboard = C.STATE_ENABLED
		settings.application_cache = C.STATE_ENABLED
		settings.text_area_resize = C.STATE_DISABLED
		settings.plugins = C.STATE_DISABLED
		settings.webgl = C.STATE_ENABLED
		// settings.background_color = "transparent"
		// settings.default_encoding = "UTF-8"

		cefCreateBrowser.Call(
			uintptr(unsafe.Pointer(windowInfo)),
			uintptr(unsafe.Pointer(&client)),
			uintptr(unsafe.Pointer(cefString("about:blank"))),
			uintptr(unsafe.Pointer(settings)),
			uintptr(unsafe.Pointer(C.NULL)),
		)

		thread, _, _ = winGetCurrentThreadID.Call()
	})

	if browser, ok := cRet.(*C.cef_browser_t); ok {
		window := C.GetWindowHandle(browser)
		dpi, _, _ := winGetDpiForWindow.Call(uintptr(unsafe.Pointer(window)))
		monitor, _, _ := winMonitorFromWindow.Call(uintptr(unsafe.Pointer(window)), uintptr(MONITOR_DEFAULTTOPRIMARY))
		info := C.MONITORINFO{cbSize: C.sizeof_MONITORINFO}
		winGetMonitorInfo.Call(monitor, uintptr(unsafe.Pointer(&info)))
		fmt.Println(info.rcWork)
		fmt.Println("DPI:", uint64(dpi))
		// HMONITOR
		windowInfo.width = C.int(uint64(float64(width) * (float64(uint64(dpi)) / 96.0)))
		windowInfo.height = C.int(uint64(float64(height) * (float64(uint64(dpi)) / 96.0)))
		windowInfo.x = C.int(info.rcWork.right/2) - C.int(windowInfo.width/2)
		windowInfo.y = C.int(info.rcWork.bottom/2) - C.int(windowInfo.height/2)

		C.LoadURL(browser, cefString("data:text/html,"+urlEncode("<!DOCTYPE html><html><body><h1>Hello, World!!!</h1>Testtesttest</body>")))
		winSetWindowPos.Call(
			uintptr(unsafe.Pointer(window)),
			uintptr(0),
			uintptr(windowInfo.x),
			uintptr(windowInfo.y),
			uintptr(windowInfo.width),
			uintptr(windowInfo.height),
			uintptr(SWP_SHOWWINDOW|SWP_NOOWNERZORDER|SWP_NOZORDER|SWP_ASYNCWINDOWPOS),
		)

		winShowWindow.Call(uintptr(unsafe.Pointer(window)), uintptr(SW_SHOW))
		winUpdateWindow.Call(uintptr(unsafe.Pointer(window)))
		winSwitchToThisWindow.Call(uintptr(unsafe.Pointer(window)), uintptr(C.TRUE))

		wind := &Window{
			id:      int(id),
			thread:  int(thread),
			browser: unsafe.Pointer(browser),
			window:  unsafe.Pointer(window),
			state:   State{Hidden: true},
			MainMenu: &Menu{
				menu: nil, //ret.menubar,
			},
			app: a,
		}
		winds = append(winds, wind)
		return wind
	}
	/*


		fmt.Println("[::::: goBrowserCreate]",
			uintptr(unsafe.Pointer(browser)),
			uintptr(unsafe.Pointer(frame)),
			uintptr(unsafe.Pointer(win)),
		)

		winSetWindowPos.Call(
			uintptr(unsafe.Pointer(win)),
			uintptr(0),
			uintptr(100),
			uintptr(100),
			uintptr(700),
			uintptr(700),
			uintptr(SWP_SHOWWINDOW|SWP_NOOWNERZORDER|SWP_NOZORDER|SWP_ASYNCWINDOWPOS),
		)

		winRedrawWindow.Call(uintptr(unsafe.Pointer(win)),
			uintptr(unsafe.Pointer(C.NULL)),
			uintptr(unsafe.Pointer(C.NULL)),
			uintptr(RDW_INVALIDATE),
		)

		go func() {
			time.Sleep(time.Second * 10)
			winSetWindowPos.Call(
				uintptr(unsafe.Pointer(win)),
				uintptr(0),
				uintptr(100),
				uintptr(100),
				uintptr(1000),
				uintptr(1000),
				uintptr(SWP_SHOWWINDOW|SWP_NOOWNERZORDER|SWP_NOZORDER|SWP_ASYNCWINDOWPOS),
			)

			winRedrawWindow.Call(uintptr(unsafe.Pointer(win)),
				uintptr(unsafe.Pointer(C.NULL)),
				uintptr(unsafe.Pointer(C.NULL)),
				uintptr(RDW_INVALIDATE),
			)
		}()

		go func() {
			time.Sleep(time.Second / 10)
			C.LoadURL(browser, cefString("data:text/html,%3C!DOCTYPE%20html%3E%3Chtml%3E%3Chead%3E%3Cmeta%20http-equiv%3D%22refresh%22%20content%3D%220%3B%20url%3Dhttp%3A%2F%2Flocalhost%3A8080%2F%22%20%2F%3E%3C%2Fhead%3E%3Cbody%3E%3Ch1%3EHello%2C%20World!!!%3C%2Fh1%3E Test test test%3C%2Fbody%3E"))
			// time.Sleep(time.Second)
			//C.LoadURL(browser, cefString("http://localhost:8080/test2"))
			fmt.Println("[::::: LoadURL]")
		}()
	*/

	return nil
}

func urlEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

func cefToGoString(source *C.cef_string_t) string {
	// output := C.cef_string_userfree_utf8_alloc()
	outputU64, _, _ := cefAllocUTF8.Call()
	output := (*C.cef_string_utf8_t)(unsafe.Pointer(uintptr(outputU64)))
	if source == nil || source.length == 0 {
		return ""
	}

	cefStringToUTF8.Call(
		uintptr(unsafe.Pointer(source.str)),
		uintptr(uint64(source.length)),
		uintptr(unsafe.Pointer(output)),
	)

	defer cefFreeUTF8.Call(
		uintptr(unsafe.Pointer(output)),
	)
	return C.GoString(output.str)
}

//export cefString
func cefString(s string) *C.cef_string_t {
	ret := (*C.cef_string_t)(C.calloc(1, C.sizeof_cef_string_t))
	if len(s) > 0 {
		schar := C.CString(s)
		defer C.free(unsafe.Pointer(schar))
		cefStringFromUTF8.Call(
			uintptr(unsafe.Pointer(schar)),
			uintptr(uint64(C.strlen(schar))),
			uintptr(unsafe.Pointer(ret)),
		)
	}
	return ret
}

//export cefToString
func cefToString(source *C.cef_string_t) *C.char {
	return C.CString(cefToGoString(source))
}

//export cefFromString
func cefFromString(source *C.char) *C.cef_string_t {
	return cefString(C.GoString(source))
}

//export goPrint
func goPrint(text *C.char) {
	fmt.Println(C.GoString(text))
}

//export goPrintInt
func goPrintInt(text *C.char, t C.int) {
	fmt.Println(C.GoString(text), int(t))
}

//export goPrintCef
func goPrintCef(text0 *C.char, text *C.cef_string_t) {
	fmt.Println(C.GoString(text0), cefToGoString(text))
}

//export goGetLifeSpan
func goGetLifeSpan(client *C.cef_client_t) unsafe.Pointer {
	fmt.Println("[::::: goGetLifeSpan]", int(uintptr(unsafe.Pointer(client))))

	if lifeHandler, ok := lifeHandlers[uintptr(unsafe.Pointer(client))]; ok {
		return unsafe.Pointer(lifeHandler)
	}
	lifeHandler := C.initialize_cef_life_span_handler()
	lifeHandlers[uintptr(unsafe.Pointer(client))] = lifeHandler
	return lifeHandler
}

//export goBrowserCreate
func goBrowserCreate(browser *C.cef_browser_t) {
	if reqid, ok := cliReqs[uintptr(unsafe.Pointer(C.GetClient(browser)))]; ok {
		thread, _, _ := winGetCurrentThreadID.Call()
		fmt.Println("[::::::::: PROCESS:2]", thread)
		cRequestRet(reqid, browser)
	}
}

//export goBrowserDoClose
func goBrowserDoClose(browser *C.cef_browser_t) C.int {
	window := C.GetWindowHandle(browser)
	winShowWindow.Call(uintptr(unsafe.Pointer(window)), uintptr(SW_HIDE))
	thread, _, _ := winGetCurrentThreadID.Call()
	fmt.Println("[::::::::: PROCESS:3]", thread)

	// message, _, _ := cefProcessMessageCreate.Call(uintptr(unsafe.Pointer(cefString("KILL"))))
	// C.SendProcessMessage(browser, (*C.cef_process_message_t)(unsafe.Pointer(message)))

	if cef2destroy {
		fmt.Println("[>] GOBROWSERDOCLOSE - 0[<]")
		defer os.Exit(1)
		return C.int(0) // 1
	}
	fmt.Println("[>] GOBROWSERDOCLOSE - 1[<]")
	go closeCef()
	return C.int(0)
}

func closeCef() {
	cef2destroy = true
	var win *Window
	for _, win = range winds {
		go func() {
			C.CloseBrowser((*C.cef_browser_t)(win.browser))
		}()
	}

	fmt.Println("win = range winds")
	// cefShutdown.Call()
	// fmt.Println("cefShutdown.Call()")
	time.Sleep(time.Second * 5)
	fmt.Println("win.app.shown <- true 1")
	win.app.shown <- true
	fmt.Println("win.app.shown <- true 2")
}

//export goBrowserBeforeClose
func goBrowserBeforeClose(browser *C.cef_browser_t) {
	frame := C.GetMainFrame(browser)
	win := C.GetWindowHandle(browser)

	fmt.Println("[::::: goBrowserBeforeClose]",
		uintptr(unsafe.Pointer(browser)),
		uintptr(unsafe.Pointer(frame)),
		uintptr(unsafe.Pointer(win)),
	)
}
