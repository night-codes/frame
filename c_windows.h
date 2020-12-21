#if defined(WEBVIEW_WINAPI)
#ifndef WEBVIEW_H
#define WEBVIEW_H

#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#define CINTERFACE
// #include <commctrl.h>
// #include <exdisp.h>
// #include <mshtmhst.h>
// #include <mshtml.h>
// #include <shobjidl.h>
#include <stdio.h>
#include <windows.h>

typedef struct idleData {
    int id;
    int id2;
    HWND window;
    char* content;
    char* uri;
    int width;
    int height;
    int x;
    int y;
    int hint;
    double dbl;
    boolean flag;
    long long unsigned int req_id;
} idleData;

#define IDLE_DISPATCH (WM_APP + 1)

typedef struct WindowObj {
    int id;
    int thread;
    long long unsigned int req_id;
    HWND window;
} WindowObj;

typedef void (*idleFn)(WindowObj* data, void* arg);

extern void goPrint(char* text);
extern void goPrintInt(char* text, int num);
extern void goWinRet(long long unsigned int req_id, WindowObj* win);

static int* procs;
static int procs_count;

static void runIdle(int thread, idleFn fn, void* arg)
{
    PostThreadMessage((DWORD)thread, IDLE_DISPATCH, (WPARAM)fn, (LPARAM)arg);
}

static void showWindowIdle(WindowObj* data, void* arg)
{
    ShowWindow(data->window, SW_SHOW);
    UpdateWindow(data->window);
    SwitchToThisWindow(data->window, TRUE);
}

static void showWindow(int thread)
{
    runIdle(thread, showWindowIdle, NULL);
}

static LRESULT CALLBACK wndproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
    WindowObj* ww = (WindowObj*)GetWindowLongPtr(hwnd, GWLP_USERDATA);
    switch (uMsg) {
    case WM_CREATE:
        ww = (WindowObj*)((CREATESTRUCT*)lParam)->lpCreateParams;
        ww->window = hwnd;

        goWinRet(ww->req_id, ww);
        // runIdle(ww->id, showWindowIdle, NULL);
        // window = hwnd;
        // return EmbedBrowserObject(w);
        return TRUE;
    case WM_DESTROY:
        // UnEmbedBrowserObject(w);
        // PostQuitMessage(0);
        return TRUE;
    case WM_SIZE: {
        /* IWebBrowser2* webBrowser2;
        IOleObject* browser = *w->priv.browser;
        if (browser->lpVtbl->QueryInterface(browser, iid_unref(&IID_IWebBrowser2),
                (void**)&webBrowser2)
            == S_OK) {
            RECT rect;
            GetClientRect(hwnd, &rect);
            webBrowser2->lpVtbl->put_Width(webBrowser2, rect.right);
            webBrowser2->lpVtbl->put_Height(webBrowser2, rect.bottom);
        } */
        return TRUE;
    }
    case IDLE_DISPATCH: {
        idleFn f = (idleFn)wParam;
        void* arg = (void*)lParam;
        (*f)(ww, arg);
        return TRUE;
    }
    }
    return DefWindowProc(hwnd, uMsg, wParam, lParam);
}

static int webview_loop(HWND hwnd, int blocking)
{
    MSG msg;
    if (blocking) {
        GetMessage(&msg, 0, 0, 0);
    } else {
        PeekMessage(&msg, 0, 0, 0, PM_REMOVE);
    }
    switch (msg.message) {
    case WM_QUIT:
        return -1;
    case WM_COMMAND:
    case WM_KEYDOWN:
    case WM_KEYUP: {
        HRESULT r = S_OK;
        /* IWebBrowser2* webBrowser2;
        IOleObject* browser = *w->priv.browser;
        if (browser->lpVtbl->QueryInterface(browser, iid_unref(&IID_IWebBrowser2),
                (void**)&webBrowser2)
            == S_OK) {
            IOleInPlaceActiveObject* pIOIPAO;
            if (browser->lpVtbl->QueryInterface(
                    browser, iid_unref(&IID_IOleInPlaceActiveObject),
                    (void**)&pIOIPAO)
                == S_OK) {
                r = pIOIPAO->lpVtbl->TranslateAccelerator(pIOIPAO, &msg);
                pIOIPAO->lpVtbl->Release(pIOIPAO);
            }
            webBrowser2->lpVtbl->Release(webBrowser2);
        } */
        if (r != S_FALSE) {
            break;
        }
    }
    default:
        msg.hwnd = hwnd;
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }
    return 0;
}

static int makeWindow(char* title, int width, int height, long long unsigned int req_id, int id)
{
    WindowObj* ww = (WindowObj*)malloc(sizeof(WindowObj));

    ww->id = id;
    ww->thread = (int)GetCurrentThreadId();
    ww->req_id = req_id;

    static const TCHAR* classname = "WebView";
    WNDCLASSEX wc;
    HINSTANCE hInstance;
    DWORD style;
    RECT clientRect;
    RECT rect;

    hInstance = GetModuleHandle(NULL);
    if (hInstance == NULL) {
        return -1;
    }
    // if (OleInitialize(NULL) != S_OK) {
    //     return -1;
    // }
    ZeroMemory(&wc, sizeof(WNDCLASSEX));
    wc.cbSize = sizeof(WNDCLASSEX);
    wc.hInstance = hInstance;
    wc.lpfnWndProc = wndproc;
    wc.lpszClassName = classname;
    RegisterClassEx(&wc);

    style = WS_OVERLAPPEDWINDOW;
    if (FALSE) {
        style = WS_OVERLAPPED | WS_CAPTION | WS_MINIMIZEBOX | WS_SYSMENU;
    }

    rect.left = 0;
    rect.top = 0;
    rect.right = width;
    rect.bottom = height;
    AdjustWindowRect(&rect, WS_OVERLAPPEDWINDOW, 0);

    GetClientRect(GetDesktopWindow(), &clientRect);
    int left = (clientRect.right / 2) - ((rect.right - rect.left) / 2);
    int top = (clientRect.bottom / 2) - ((rect.bottom - rect.top) / 2);
    rect.right = rect.right - rect.left + left;
    rect.left = left;
    rect.bottom = rect.bottom - rect.top + top;
    rect.top = top;

    HWND hwnd = CreateWindowEx(0, classname, title, style, rect.left, rect.top,
        rect.right - rect.left, rect.bottom - rect.top,
        HWND_DESKTOP, NULL, hInstance, (void*)ww);

    if (hwnd == 0) {
        // OleUninitialize();
        return -1;
    }

    ww->window = hwnd;

    SetWindowLongPtr(hwnd, GWLP_USERDATA, (LONG_PTR)ww);

    // DisplayHTMLPage(w);

    SetWindowText(hwnd, title);
    while (webview_loop(hwnd, 1) == 0) {
    }
    return 0;
}

#endif // !WEBVIEW_H
#endif // WEBVIEW_WINAPI