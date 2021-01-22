#include "cef_app.h"
#include "cef_client.h"
//#include "cef_helpers.h"
#include "include/capi/cef_browser_capi.h"
#include "include/capi/cef_client_capi.h"
#include "include/capi/cef_v8_capi.h"
#include "string.h"
#include <stdlib.h> 

extern int cef_browser_t_get_identifier(cef_browser_t* browser);

extern cef_frame_t* cef_browser_t_get_main_frame(cef_browser_t* browser);

extern cef_frame_t* cef_browser_t_get_focused_frame(cef_browser_t* browser);

extern cef_frame_t* cef_browser_t_get_frame_byident(cef_browser_t* browser, int64 identifier);

extern cef_frame_t* cef_browser_t_get_frame(cef_browser_t* browser, const cef_string_t* name);

extern size_t cef_browser_t_get_frame_count(cef_browser_t* browser);

extern void cef_browser_t_get_frame_identifiers(struct _cef_browser_t* self,
    size_t* identifiersCount, int64* identifiers);

extern void cef_browser_t_get_frame_names(struct _cef_browser_t* self,
    cef_string_list_t names);

extern void ExecuteJavaScript(cef_browser_t* browser, const char* code, const char* script_url, int start_line);

extern void LoadURL(cef_browser_t* browser, const char* url);

extern void _LoadString(cef_browser_t* browser, const char* string_val, const char* url);

extern void BrowserWasResized(cef_browser_t* browser);

extern cef_window_handle_t GetWindowHandle(cef_browser_t* browser);

extern cef_window_handle_t GetRootWindowHandle(cef_browser_t* browser);

// Force close the browser
extern void CloseBrowser(cef_browser_t* browser);

extern cef_string_utf8_t* cefStringToUtf8(cef_string_t* source);

extern cef_string_t* GetURL(cef_browser_t* browser);

extern void GetSource(cef_browser_t* browser, cef_string_visitor_t* visitor);

extern void GetText(cef_browser_t* browser, cef_string_visitor_t* visitor);

extern void VisitDOM(cef_browser_t* browser, cef_domvisitor_t* visitor);

extern int SendProcessMessage(struct _cef_browser_t* self,
    cef_process_id_t target_process,
    struct _cef_process_message_t* message);

extern struct _cef_v8context_t* GetV8Context(cef_browser_t* browser);

extern int V8Eval(struct _cef_v8context_t* self,
    const cef_string_t* code, struct _cef_v8value_t** retval,
    struct _cef_v8exception_t** exception);

extern void initialize_cef_string_visitor(struct _cef_string_visitor_t* visitor);