// CEF C API example
// Project website: https://github.com/cztomczak/cefcapi

#pragma once

#include "handlers/cef_base.h"
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_load_handler_capi.h"

static void CEF_CALLBACK on_loading_state_change(struct _cef_load_handler_t* self,
    struct _cef_browser_t* browser, int isLoading, int canGoBack,
    int canGoForward)
{
    goPrint("******* ON_LOADING_STATE_CHANGE");
};

///
// Called when the browser begins loading a frame. The |frame| value will
// never be NULL -- call the is_main() function to check if this frame is the
// main frame. |transition_type| provides information about the source of the
// navigation and an accurate value is only available in the browser process.
// Multiple frames may be loading at the same time. Sub-frames may start or
// continue loading after the main frame load has ended. This function will
// always be called for all frames irrespective of whether the request
// completes successfully. For notification of overall browser load status use
// OnLoadingStateChange instead.
///
static void CEF_CALLBACK on_load_start(struct _cef_load_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    cef_transition_type_t transition_type)
{
    goPrint("******* ON_LOAD_START");
};

///
// Called when the browser is done loading a frame. The |frame| value will
// never be NULL -- call the is_main() function to check if this frame is the
// main frame. Multiple frames may be loading at the same time. Sub-frames may
// start or continue loading after the main frame load has ended. This
// function will always be called for all frames irrespective of whether the
// request completes successfully. For notification of overall browser load
// status use OnLoadingStateChange instead.
///
static void CEF_CALLBACK on_load_end(struct _cef_load_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    int httpStatusCode)
{
    goPrint("******* ON_LOAD_END");
};

///
// Called when the resource load for a navigation fails or is canceled.
// |errorCode| is the error code number, |errorText| is the error text and
// |failedUrl| is the URL that failed to load. See net\base\net_error_list.h
// for complete descriptions of the error codes.
///
static void CEF_CALLBACK on_load_error(struct _cef_load_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    cef_errorcode_t errorCode, const cef_string_t* errorText,
    const cef_string_t* failedUrl)
{
    goPrintCef("errorText:", errorText);
    goPrintCef("failedUrl:", failedUrl);
    goPrint("******* ON_LOAD_ERROR");
};

static cef_load_handler_t* initialize_cef_load_handler()
{
    cef_load_handler_t* loadHandler = (cef_load_handler_t*)calloc(1, sizeof(cef_load_handler_t));
    loadHandler->base.size = sizeof(cef_load_handler_t);
    initialize_cef_base((cef_base_t*)loadHandler);
    DEBUG_CALLBACK("[+ INITIALIZE_CEF_LOAD_HANDLER +]\n");

    loadHandler->on_loading_state_change = on_loading_state_change;
    loadHandler->on_load_start = on_load_start;
    loadHandler->on_load_end = on_load_end;
    loadHandler->on_load_error = on_load_error;
    return loadHandler;
}