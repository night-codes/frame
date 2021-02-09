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
	#include "include/capi/views/cef_browser_view_capi.h"

	static void ExecuteJavaScript(cef_browser_t* browser, cef_string_t* code, cef_string_t* url, int start_line) {
		cef_frame_t * frame = browser->get_main_frame(browser);
		frame->execute_java_script(frame, code, url, start_line);
	}

	static void LoadURL(cef_browser_t* browser, cef_string_t* url) {
		cef_frame_t * frame = browser->get_main_frame(browser);
		frame->load_url(frame, url);
	}

	static void LoadHTML(cef_browser_t* browser, cef_string_t* html, cef_string_t* url) {
		cef_frame_t * frame = browser->get_main_frame(browser);
		frame->load_string(frame, html, url);
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
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/night-codes/cefresources"
	"golang.org/x/sys/windows"
)

type (
	// App is main application object
	App struct {
		WindowClose *Window
		AllClose    bool
		Test5       bool
		app         interface{} // *C.GtkApplication
		openedWns   sync.WaitGroup
		shown       chan bool
		mainArgs    C.cef_main_args_t
		appHandler  C.cef_app_t
	}

	ceBrowser     *C.cef_browser_t
	ceString      *C.cef_string_t
	C_MONITORINFO C.MONITORINFO
	C_RECT        C.RECT
)

var (
	cef2destroy = false
	shutdown    = false
	mutexNew    sync.Mutex
	winds       = []*Window{}
	lock        sync.Mutex
	appChan     = make(chan *App)
	idItr       int64
	app         *App

	resDir   = cefresources.Extract()
	libcef   = windows.NewLazyDLL(filepath.Join(resDir, "libcef.dll"))
	kernel32 = windows.NewLazySystemDLL("kernel32")
	ole32    = windows.NewLazySystemDLL("ole32")
	user32   = windows.NewLazySystemDLL("user32")
	gdi32    = windows.NewLazySystemDLL("Gdi32")

	monitorinfoSizeof = C.ulong(C.sizeof_MONITORINFO)

	cefCreateBrowser            = libcef.NewProc("cef_browser_host_create_browser")
	cefInitialize               = libcef.NewProc("cef_initialize")
	cefExecuteProcess           = libcef.NewProc("cef_execute_process")
	cefEnableHDPI               = libcef.NewProc("cef_enable_highdpi_support")
	cefStringFromUTF8           = libcef.NewProc("cef_string_utf8_to_utf16")
	cefStringToUTF8             = libcef.NewProc("cef_string_utf16_to_utf8")
	cefAllocUTF8                = libcef.NewProc("cef_string_userfree_utf8_alloc")
	cefFreeUTF8                 = libcef.NewProc("cef_string_userfree_utf8_free")
	cefQuitMessageLoop          = libcef.NewProc("cef_quit_message_loop")
	cefShutdown                 = libcef.NewProc("cef_shutdown")
	cefRunMessageLoop           = libcef.NewProc("cef_run_message_loop")
	cefGetGlobalCtx             = libcef.NewProc("cef_request_context_get_global_context")
	cefBrowserViewGetForBrowser = libcef.NewProc("cef_browser_view_get_for_browser")
	cefProcessMessageCreate     = libcef.NewProc("cef_process_message_create")

	winCoInitializeEx             = ole32.NewProc("CoInitializeEx")
	winGetProcessHeap             = kernel32.NewProc("GetProcessHeap")
	winHeapAlloc                  = kernel32.NewProc("HeapAlloc")
	winHeapFree                   = kernel32.NewProc("HeapFree")
	winGetCurrentThreadID         = kernel32.NewProc("GetCurrentThreadId")
	winExitThread                 = kernel32.NewProc("ExitThread")
	winSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	winLoadImageW                 = user32.NewProc("LoadImageW")
	winGetSystemMetrics           = user32.NewProc("GetSystemMetrics")
	winGetDpiForWindow            = user32.NewProc("GetDpiForWindow")
	winRegisterClassExW           = user32.NewProc("RegisterClassExW")
	winCreateWindowExW            = user32.NewProc("CreateWindowExW")
	winDestroyWindow              = user32.NewProc("DestroyWindow")
	winShowWindow                 = user32.NewProc("ShowWindow")
	winUpdateWindow               = user32.NewProc("UpdateWindow")
	winSwitchToThisWindow         = user32.NewProc("SwitchToThisWindow")
	winSetFocus                   = user32.NewProc("SetFocus")
	winGetMessageW                = user32.NewProc("GetMessageW")
	winTranslateMessage           = user32.NewProc("TranslateMessage")
	winDispatchMessageW           = user32.NewProc("DispatchMessageW")
	winDefWindowProcW             = user32.NewProc("DefWindowProcW")
	winGetClientRect              = user32.NewProc("GetClientRect")
	winPostQuitMessage            = user32.NewProc("PostQuitMessage")
	winSetWindowTextW             = user32.NewProc("SetWindowTextW")
	winPostThreadMessageW         = user32.NewProc("PostThreadMessageW")
	winGetWindowLongPtrW          = user32.NewProc("GetWindowLongPtrW")
	winGetWindowLong              = user32.NewProc("GetWindowLongA")
	winSetWindowLong              = user32.NewProc("SetWindowLongA")
	winAdjustWindowRect           = user32.NewProc("AdjustWindowRect")
	winSetWindowPos               = user32.NewProc("SetWindowPos")
	winRedrawWindow               = user32.NewProc("RedrawWindow")
	winMonitorFromWindow          = user32.NewProc("MonitorFromWindow")
	winGetMonitorInfo             = user32.NewProc("GetMonitorInfoA")
	winSetClassLongPtr            = user32.NewProc("SetClassLongPtrA")
	winGetWindowRect              = user32.NewProc("GetWindowRect")
	winSetParent                  = user32.NewProc("SetParent")

	gdiCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")

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

	WS_OVERLAPPED       = 0x00000000
	WS_POPUP            = 0x80000000
	WS_CHILD            = 0x40000000
	WS_MINIMIZE         = 0x20000000
	WS_VISIBLE          = 0x10000000
	WS_DISABLED         = 0x08000000
	WS_CLIPSIBLINGS     = 0x04000000
	WS_CLIPCHILDREN     = 0x02000000
	WS_MAXIMIZE         = 0x01000000
	WS_CAPTION          = 0x00C00000 // WS_BORDER | WS_DLGFRAME
	WS_BORDER           = 0x00800000
	WS_DLGFRAME         = 0x00400000
	WS_VSCROLL          = 0x00200000
	WS_HSCROLL          = 0x00100000
	WS_SYSMENU          = 0x00080000
	WS_THICKFRAME       = 0x00040000
	WS_GROUP            = 0x00020000
	WS_TABSTOP          = 0x00010000
	WS_MINIMIZEBOX      = 0x00020000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_TILED            = WS_OVERLAPPED
	WS_ICONIC           = WS_MINIMIZE
	WS_SIZEBOX          = WS_THICKFRAME
	WS_TILEDWINDOW      = WS_OVERLAPPEDWINDOW
	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_POPUPWINDOW      = WS_POPUP | WS_BORDER | WS_SYSMENU
	WS_CHILDWINDOW      = WS_CHILD
	WS_EX_LAYERED       = 0x00080000
	WS_EX_COMPOSITED    = 0x02000000

	MAX_PATH              = 260
	LWA_COLORKEY          = 0x00001
	LWA_ALPHA             = 0x00002
	ENUM_CURRENT_SETTINGS = 0xFFFFFFFF
	GWL_STYLE             = -16
	GWL_EXSTYLE           = -20
	GCLP_HBRBACKGROUND    = -10
	HWND_NOTOPMOST        = -2
	HWND_TOPMOST          = -1
	HWND_TOP              = 0
	HWND_BOTTOM           = 1

	MONITOR_DEFAULTTONULL    = 0x00000000
	MONITOR_DEFAULTTOPRIMARY = 0x00000001
	MONITOR_DEFAULTTONEAREST = 0x00000002
)

func appDir() string {
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home + "\\AppData"
}

// C:\Users\mirrr\AppData\Local\Temp
// MakeApp is make and run one instance of application (At the moment, it is possible to create only one instance)
func MakeApp(appName string) *App {
	lock.Lock()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		for _, win := range winds {
			C.CloseBrowser(ceBrowser(win.browser))
		}
	}()

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

	cefSettings := (*C.cef_settings_t)(C.calloc(1, C.sizeof_cef_settings_t))
	cefSettings.size = C.sizeof_cef_settings_t
	// cefSettings.pack_loading_disabled = C.int(1)
	userdata := filepath.Join(appDir(), appName)
	cache := filepath.Join(userdata, "cache")
	if _, err := os.Stat(cache); err != nil {
		os.MkdirAll(cache, 0755)
	}

	cefSettings.user_agent = *cefString("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")
	cefSettings.no_sandbox = C.int(0)
	cefSettings.single_process = C.int(1)
	// cefSettings.remote_debugging_port = C.int(9090)
	cefSettings.multi_threaded_message_loop = C.int(1)
	cefSettings.context_safety_implementation = C.int(0)
	cefSettings.user_data_path = *cefString(userdata)
	cefSettings.persist_session_cookies = C.int(1)
	cefSettings.persist_user_preferences = C.int(1)
	cefSettings.log_severity = (C.cef_log_severity_t)(C.int(C.LOGSEVERITY_DISABLE))
	cefSettings.cache_path = *cefString(cache)
	cefSettings.log_file = *cefString(filepath.Join(resDir, "log.txt"))
	cefSettings.resources_dir_path = *cefString(resDir)
	cefSettings.locales_dir_path = *cefString(resDir)
	cefSettings.background_color = 0xff999999

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
	a.AllClose = true
	for {
		runtime.Gosched()
		time.Sleep(time.Millisecond * 10)
		if shutdown {
			break
		}
	}
	process, _ := os.FindProcess(int(windows.GetCurrentProcessId()))
	process.Kill()
}

// WaitWindowClose locker
func (a *App) WaitWindowClose(win *Window) {
	a.WindowClose = win
	for {
		runtime.Gosched()
		time.Sleep(time.Millisecond * 10)
		if shutdown {
			break
		}
	}
	process, _ := os.FindProcess(int(windows.GetCurrentProcessId()))
	process.Kill()
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
		windowInfo.style = WS_OVERLAPPEDWINDOW | WS_TABSTOP // | WS_VISIBLE
		windowInfo.transparent_painting_enabled = 1
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
		settings.javascript = C.STATE_ENABLED
		settings.javascript_open_windows = C.STATE_DISABLED
		settings.javascript_access_clipboard = C.STATE_ENABLED
		settings.application_cache = C.STATE_ENABLED
		settings.text_area_resize = C.STATE_DISABLED
		settings.plugins = C.STATE_DISABLED
		settings.webgl = C.STATE_ENABLED
		settings.background_color = 0xff999999
		// settings.default_encoding = "UTF-8"

		cefCreateBrowser.Call(
			uintptr(unsafe.Pointer(windowInfo)),
			uintptr(unsafe.Pointer(&client)),
			uintptr(unsafe.Pointer(cefString("about:blank"))),
			uintptr(unsafe.Pointer(settings)),
			uintptr(unsafe.Pointer(C.NULL)),
		)
	})

	if browser, ok := cRet.(ceBrowser); ok {
		window := C.GetWindowHandle(browser)
		dpi, _, _ := winGetDpiForWindow.Call(uintptr(unsafe.Pointer(window)))
		monitor, _, _ := winMonitorFromWindow.Call(uintptr(unsafe.Pointer(window)), uintptr(MONITOR_DEFAULTTOPRIMARY))
		info := C.MONITORINFO{cbSize: monitorinfoSizeof}
		winGetMonitorInfo.Call(monitor, uintptr(unsafe.Pointer(&info)))
		windowInfo.width = C.int(uint64(float64(width) * (float64(uint64(dpi)) / 96.0)))
		windowInfo.height = C.int(uint64(float64(height) * (float64(uint64(dpi)) / 96.0)))
		windowInfo.x = C.int(info.rcWork.right/2) - C.int(windowInfo.width/2)
		windowInfo.y = C.int(info.rcWork.bottom/2) - C.int(windowInfo.height/2)

		winSetWindowPos.Call(
			uintptr(unsafe.Pointer(window)),
			uintptr(0),
			uintptr(windowInfo.x),
			uintptr(windowInfo.y),
			uintptr(windowInfo.width),
			uintptr(windowInfo.height),
			uintptr(SWP_NOOWNERZORDER|SWP_NOZORDER|SWP_ASYNCWINDOWPOS),
		)

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
	return nil
}

func loadHTML(browser unsafe.Pointer, html, uri string) {
	if uri == "" {
		uri = "about:balnk"
	}
	C.LoadHTML(ceBrowser(browser), cefString(html), cefString(uri))
}

func loadURL(browser unsafe.Pointer, uri string) {
	C.LoadURL(ceBrowser(browser), cefString(uri))
}

func urlEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

func cefToGoString(source ceString) string {
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

	defer cefFreeUTF8.Call(uintptr(unsafe.Pointer(output)))
	return C.GoString(output.str)
}

//export cefString
func cefString(s string) ceString {
	ret := (ceString)(C.calloc(1, C.sizeof_cef_string_t))
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
func cefToString(source ceString) *C.char {
	return C.CString(cefToGoString(source))
}

//export cefFromString
func cefFromString(source *C.char) ceString {
	return cefString(C.GoString(source))
}

//export goPrint
func goPrint(text *C.char) {
	fmt.Println(C.GoString(text))
}

func cChar(text string) *C.char {
	return C.CString(text)
}

//export goPrintInt
func goPrintInt(text *C.char, t C.int) {
	fmt.Println(C.GoString(text), int(t))
}

//export goPrintCef
func goPrintCef(text0 *C.char, text ceString) {
	fmt.Println(C.GoString(text0), cefToGoString(text))
}

//export goGetLifeSpan
func goGetLifeSpan(client *C.cef_client_t) unsafe.Pointer {
	if lifeHandler, ok := lifeHandlers[uintptr(unsafe.Pointer(client))]; ok {
		return unsafe.Pointer(lifeHandler)
	}
	lifeHandler := C.initialize_cef_life_span_handler()
	lifeHandlers[uintptr(unsafe.Pointer(client))] = lifeHandler
	return lifeHandler
}

//export goBrowserCreate
func goBrowserCreate(browser ceBrowser) {
	if reqid, ok := cliReqs[uintptr(unsafe.Pointer(C.GetClient(browser)))]; ok {
		cRequestRet(reqid, browser)
	}
}

//export goBrowserDestroyed
func goBrowserDestroyed(browser ceBrowser) C.int {
	return C.int(0)
}

//export goNop
func goNop() {
	runtime.Gosched()
}

//export goBrowserDoClose
func goBrowserDoClose(browser ceBrowser) C.int {
	window := C.GetWindowHandle(browser)
	winShowWindow.Call(uintptr(unsafe.Pointer(window)), uintptr(windows.SW_HIDE))

	if cef2destroy {
		return C.int(0)
	}
	go closeCef()
	return C.int(1)
}

func closeCef() {
	cef2destroy = true
	var win *Window
	for _, win = range winds {
		C.CloseBrowser(ceBrowser(win.browser))
		// winDestroyWindow.Call(uintptr(unsafe.Pointer(win.window)))
	}
	go func() {
		time.Sleep(time.Second / 2)
		shutdown = true
	}()
}

//export goBrowserBeforeClose
func goBrowserBeforeClose(browser ceBrowser) {
	// frame := C.GetMainFrame(browser)
	// win := C.GetWindowHandle(browser)
	fmt.Println("~~~ goBrowserBeforeClose ~~~")
	// cefQuitMessageLoop.Call()
	// cefShutdown.Call()
}
