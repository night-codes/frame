// CEF C API example
// Project website: https://github.com/cztomczak/cefcapi

#ifndef HANDLER_REQUEST_H_
#define HANDLER_REQUEST_H_

#pragma once

#include "handlers/cef_base.h"
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_request_handler_capi.h"

///
// Called on the UI thread before browser navigation. Return true (1) to
// cancel the navigation or false (0) to allow the navigation to proceed. The
// |request| object cannot be modified in this callback.
// cef_load_handler_t::OnLoadingStateChange will be called twice in all cases.
// If the navigation is allowed cef_load_handler_t::OnLoadStart and
// cef_load_handler_t::OnLoadEnd will be called. If the navigation is canceled
// cef_load_handler_t::OnLoadError will be called with an |errorCode| value of
// ERR_ABORTED.
///
static int CEF_CALLBACK on_before_browse(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    struct _cef_request_t* request, int is_redirect)
{
    goPrint("@@@@@ - ON_BEFORE_BROWSE - @@@@@");
    return 0;
}

///
// Called on the UI thread before OnBeforeBrowse in certain limited cases
// where navigating a new or different browser might be desirable. This
// includes user-initiated navigation that might open in a special way (e.g.
// links clicked via middle-click or ctrl + left-click) and certain types of
// cross-origin navigation initiated from the renderer process (e.g.
// navigating the top-level frame to/from a file URL). The |browser| and
// |frame| values represent the source of the navigation. The
// |target_disposition| value indicates where the user intended to navigate
// the browser based on standard Chromium behaviors (e.g. current tab, new
// tab, etc). The |user_gesture| value will be true (1) if the browser
// navigated via explicit user gesture (e.g. clicking a link) or false (0) if
// it navigated automatically (e.g. via the DomContentLoaded event). Return
// true (1) to cancel the navigation or false (0) to allow the navigation to
// proceed in the source browser's top-level frame.
///
static int CEF_CALLBACK on_open_urlfrom_tab(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    const cef_string_t* target_url,
    cef_window_open_disposition_t target_disposition, int user_gesture)
{
    goPrint("@@@@@ - ON_OPEN_URLFROM_TAB - @@@@@");
    return 0;
}

///
// Called on the IO thread before a resource request is loaded. The |request|
// object may be modified. Return RV_CONTINUE to continue the request
// immediately. Return RV_CONTINUE_ASYNC and call cef_request_tCallback::
// cont() at a later time to continue or cancel the request asynchronously.
// Return RV_CANCEL to cancel the request immediately.
//
///
static cef_return_value_t CEF_CALLBACK on_before_resource_load(
    struct _cef_request_handler_t* self, struct _cef_browser_t* browser,
    struct _cef_frame_t* frame, struct _cef_request_t* request,
    struct _cef_request_callback_t* callback)
{
    goPrintCef("get_url:", request->get_url(request));
    goPrint("@@@@@ - ON_BEFORE_RESOURCE_LOAD - @@@@@");
    return RV_CONTINUE;
}

///
// Called on the IO thread before a resource is loaded. To allow the resource
// to load normally return NULL. To specify a handler for the resource return
// a cef_resource_handler_t object. The |request| object should not be
// modified in this callback.
///
static struct _cef_resource_handler_t* CEF_CALLBACK get_resource_handler(
    struct _cef_request_handler_t* self, struct _cef_browser_t* browser,
    struct _cef_frame_t* frame, struct _cef_request_t* request)
{
    goPrint("@@@@@ - GET_RESOURCE_HANDLER - @@@@@");
    return NULL;
}

///
// Called on the IO thread when a resource load is redirected. The |request|
// parameter will contain the old URL and other request-related information.
// The |new_url| parameter will contain the new URL and can be changed if
// desired. The |request| object cannot be modified in this callback.
///
static void CEF_CALLBACK on_resource_redirect(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    struct _cef_request_t* request, cef_string_t* new_url)
{
    goPrint("@@@@@ - ON_RESOURCE_REDIRECT - @@@@@");
}

///
// Called on the IO thread when a resource response is received. To allow the
// resource to load normally return false (0). To redirect or retry the
// resource modify |request| (url, headers or post body) and return true (1).
// The |response| object cannot be modified in this callback.
///
static int CEF_CALLBACK on_resource_response(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    struct _cef_request_t* request, struct _cef_response_t* response)
{
    goPrintCef("test:", cefFromString("TEST - ОЛОЛО!"));
    goPrintCef("get_url:", request->get_url(request));
    goPrintCef("get_mime_type:", response->get_mime_type(response));
    goPrintCef("get_header:", response->get_header(response, cefFromString("Server")));
    goPrintCef("get_header:", response->get_header(response, cefFromString("Content-Type")));
    goPrintCef("get_status_text:", response->get_status_text(response));
    goPrintInt("@@@@@ - ON_RESOURCE_RESPONSE - @@@@@", response->get_error(response));
    return 0;
}

///
// Called on the IO thread to optionally filter resource response content.
// |request| and |response| represent the request and response respectively
// and cannot be modified in this callback.
///
static struct _cef_response_filter_t* CEF_CALLBACK get_resource_response_filter(
    struct _cef_request_handler_t* self, struct _cef_browser_t* browser,
    struct _cef_frame_t* frame, struct _cef_request_t* request,
    struct _cef_response_t* response)
{
    goPrint("@@@@@ - GET_RESOURCE_RESPONSE_FILTER - @@@@@");
    return NULL;
}

///
// Called on the IO thread when a resource load has completed. |request| and
// |response| represent the request and response respectively and cannot be
// modified in this callback. |status| indicates the load completion status.
// |received_content_length| is the number of response bytes actually read.
///
static void CEF_CALLBACK on_resource_load_complete(
    struct _cef_request_handler_t* self, struct _cef_browser_t* browser,
    struct _cef_frame_t* frame, struct _cef_request_t* request,
    struct _cef_response_t* response, cef_urlrequest_status_t status,
    int64 received_content_length)
{
    goPrintCef("test:", cefFromString("TEST - ПЖПЖПЖ!"));
    goPrintCef("get_url:", request->get_url(request));
    goPrintCef("get_mime_type:", response->get_mime_type(response));
    goPrintCef("get_header:", response->get_header(response, cefFromString("Server")));
    goPrintCef("get_header:", response->get_header(response, cefFromString("Content-Type")));
    goPrintCef("get_status_text:", response->get_status_text(response));
    goPrintInt("@@@@@ - ON_RESOURCE_LOAD_COMPLETE - @@@@@", UR_FAILED);
    goPrintInt("@@@@@ - ON_RESOURCE_LOAD_COMPLETE 2 - @@@@@", status);
}

///
// Called on the IO thread when the browser needs credentials from the user.
// |isProxy| indicates whether the host is a proxy server. |host| contains the
// hostname and |port| contains the port number. |realm| is the realm of the
// challenge and may be NULL. |scheme| is the authentication scheme used, such
// as "basic" or "digest", and will be NULL if the source of the request is an
// FTP server. Return true (1) to continue the request and call
// cef_auth_callback_t::cont() either in this function or at a later time when
// the authentication information is available. Return false (0) to cancel the
// request immediately.
///
static int CEF_CALLBACK get_auth_credentials(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame, int isProxy,
    const cef_string_t* host, int port, const cef_string_t* realm,
    const cef_string_t* scheme, struct _cef_auth_callback_t* callback)
{
    goPrint("@@@@@ - GET_AUTH_CREDENTIALS - @@@@@");
    return 1;
}

///
// Called on the IO thread when JavaScript requests a specific storage quota
// size via the webkitStorageInfo.requestQuota function. |origin_url| is the
// origin of the page making the request. |new_size| is the requested quota
// size in bytes. Return true (1) to continue the request and call
// cef_request_tCallback::cont() either in this function or at a later time to
// grant or deny the request. Return false (0) to cancel the request
// immediately.
///
static int CEF_CALLBACK on_quota_request(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser, const cef_string_t* origin_url,
    int64 new_size, struct _cef_request_callback_t* callback)
{
    goPrint("@@@@@ - ON_QUOTA_REQUEST - @@@@@");
    return 1;
}

///
// Called on the UI thread to handle requests for URLs with an unknown
// protocol component. Set |allow_os_execution| to true (1) to attempt
// execution via the registered OS protocol handler, if any. SECURITY WARNING:
// YOU SHOULD USE THIS METHOD TO ENFORCE RESTRICTIONS BASED ON SCHEME, HOST OR
// OTHER URL ANALYSIS BEFORE ALLOWING OS EXECUTION.
///
static void CEF_CALLBACK on_protocol_execution(
    struct _cef_request_handler_t* self, struct _cef_browser_t* browser,
    const cef_string_t* url, int* allow_os_execution)
{
    goPrint("@@@@@ - ON_PROTOCOL_EXECUTION - @@@@@");
}

///
// Called on the UI thread to handle requests for URLs with an invalid SSL
// certificate. Return true (1) and call cef_request_tCallback::cont() either
// in this function or at a later time to continue or cancel the request.
// Return false (0) to cancel the request immediately. If
// CefSettings.ignore_certificate_errors is set all invalid certificates will
// be accepted without calling this function.
///
static int CEF_CALLBACK on_certificate_error(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser, cef_errorcode_t cert_error,
    const cef_string_t* request_url, struct _cef_sslinfo_t* ssl_info,
    struct _cef_request_callback_t* callback)
{
    goPrint("@@@@@ - ON_CERTIFICATE_ERROR - @@@@@");
    return 1;
}

///
// Called on the browser process UI thread when a plugin has crashed.
// |plugin_path| is the path of the plugin that crashed.
///
static void CEF_CALLBACK on_plugin_crashed(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser, const cef_string_t* plugin_path)
{
    goPrint("@@@@@ - ON_PLUGIN_CRASHED - @@@@@");
}

///
// Called on the browser process UI thread when the render view associated
// with |browser| is ready to receive/handle IPC messages in the render
// process.
///
static void CEF_CALLBACK on_render_view_ready(struct _cef_request_handler_t* self,
    struct _cef_browser_t* browser)
{
    goPrint("@@@@@ - ON_RENDER_VIEW_READY - @@@@@");
}

///
// Called on the browser process UI thread when the render process terminates
// unexpectedly. |status| indicates how the process terminated.
///
static void CEF_CALLBACK on_render_process_terminated(
    struct _cef_request_handler_t* self, struct _cef_browser_t* browser,
    cef_termination_status_t status)
{
    goPrint("@@@@@ - ON_RENDER_PROCESS_TERMINATED - @@@@@");
}

static cef_request_handler_t* initialize_request_handler()
{
    cef_request_handler_t* handler = (cef_request_handler_t*)calloc(1, sizeof(cef_request_handler_t));
    handler->base.size = sizeof(cef_request_handler_t);
    initialize_cef_base((cef_base_t*)handler);
    DEBUG_CALLBACK("[+ INITIALIZE_REQUEST_HANDLER +]\n");

    handler->on_before_browse = on_before_browse;
    handler->on_open_urlfrom_tab = on_open_urlfrom_tab;
    // handler->on_before_resource_load = on_before_resource_load;
    // handler->get_resource_handler = get_resource_handler;
    // handler->on_resource_redirect = on_resource_redirect;
    handler->on_resource_response = on_resource_response;
    // handler->get_resource_response_filter = get_resource_response_filter;
    handler->on_resource_load_complete = on_resource_load_complete;
    // handler->get_auth_credentials = get_auth_credentials;
    // handler->on_quota_request = on_quota_request;
    // handler->on_protocol_execution = on_protocol_execution;
    // handler->on_certificate_error = on_certificate_error;
    handler->on_plugin_crashed = on_plugin_crashed;
    handler->on_render_view_ready = on_render_view_ready;
    handler->on_render_process_terminated = on_render_process_terminated;

    return handler;
}

#endif // HANDLER_REQUEST_H_