// CEF C API example
// Project website: https://github.com/cztomczak/cefcapi

#ifndef HANDLER_CRPH_H_
#define HANDLER_CRPH_H_

#pragma once

#include "handlers/cef_base.h"
#include "handlers/cef_load_handler.h"
#include "include/capi/cef_app_capi.h"

extern void goRegExtension();
extern void goContextCreate(cef_v8value_t* global);

static void CEF_CALLBACK on_render_thread_created(
    struct _cef_render_process_handler_t* self,
    struct _cef_list_value_t* extra_info) {};

///
// Called after WebKit has been initialized.
///
static void CEF_CALLBACK on_web_kit_initialized(struct _cef_render_process_handler_t* self)
{
    goRegExtension();
}

///
// Called after a browser has been created. When browsing cross-origin a new
// browser will be created before the old browser with the same identifier is
// destroyed.
///
static void CEF_CALLBACK on_browser_created(
    struct _cef_render_process_handler_t* self,
    struct _cef_browser_t* browser) {};

///
// Called before a browser is destroyed.
///
static void CEF_CALLBACK on_browser_destroyed(
    struct _cef_render_process_handler_t* self,
    struct _cef_browser_t* browser)
{
    free(self);
};

///
// Called before browser navigation. Return true (1) to cancel the navigation
// or false (0) to allow the navigation to proceed. The |request| object
// cannot be modified in this callback.
///
static int CEF_CALLBACK on_before_navigation(
    struct _cef_render_process_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    struct _cef_request_t* request, cef_navigation_type_t navigation_type,
    int is_redirect)
{
    return 0;
};

///
// Called immediately after the V8 context for a frame has been created. To
// retrieve the JavaScript 'window' object use the
// cef_v8context_t::get_global() function. V8 handles can only be accessed
// from the thread on which they are created. A task runner for posting tasks
// on the associated thread can be retrieved via the
// cef_v8context_t::get_task_runner() function.
///
static void CEF_CALLBACK on_context_created(
    struct _cef_render_process_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    struct _cef_v8context_t* context)
{
    cef_v8value_t* val = context->get_global(context);
    goContextCreate(val);
};

///
// Called immediately before the V8 context for a frame is released. No
// references to the context should be kept after this function is called.
///
static void CEF_CALLBACK on_context_released(
    struct _cef_render_process_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    struct _cef_v8context_t* context) {};

///
// Return the handler for browser load status events.
///
static struct _cef_load_handler_t* CEF_CALLBACK get_load_rp_handler(
    struct _cef_render_process_handler_t* self)
{
    return NULL; // initialize_cef_load_handler();
}

///
// Called for global uncaught exceptions in a frame. Execution of this
// callback is disabled by default. To enable set
// CefSettings.uncaught_exception_stack_size > 0.
///
static void CEF_CALLBACK on_uncaught_exception(
    struct _cef_render_process_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    struct _cef_v8context_t* context, struct _cef_v8exception_t* exception,
    struct _cef_v8stack_trace_t* stackTrace) {};

///
// Called when a new node in the the browser gets focus. The |node| value may
// be NULL if no specific node has gained focus. The node object passed to
// this function represents a snapshot of the DOM at the time this function is
// executed. DOM objects are only valid for the scope of this function. Do not
// keep references to or attempt to access any DOM objects outside the scope
// of this function.
///
static void CEF_CALLBACK on_focused_node_changed(
    struct _cef_render_process_handler_t* self,
    struct _cef_browser_t* browser, struct _cef_frame_t* frame,
    struct _cef_domnode_t* node) {};

///
// Called when a new message is received from a different process. Return true
// (1) if the message was handled or false (0) otherwise. Do not keep a
// reference to or attempt to access the message outside of this callback.
///
static int CEF_CALLBACK on_render_message_received(struct _cef_render_process_handler_t* self, struct _cef_browser_t* browser, cef_process_id_t source_process, struct _cef_process_message_t* message)
{
    cef_string_userfree_t msg = message->get_name(message);
    if (strcmp(cefToString(msg), "KILL") == 0) {
        free(self);
        // free(browser);
        free(msg);
        return 1;
    }
    free(msg);
    return 0;
}

static cef_render_process_handler_t* initialize_render_process_handler()
{
    cef_render_process_handler_t* handler = (cef_render_process_handler_t*)calloc(1, sizeof(cef_render_process_handler_t));
    handler->base.size = sizeof(cef_render_process_handler_t);
    initialize_cef_base((cef_base_t*)handler);

    handler->on_render_thread_created = on_render_thread_created;
    handler->on_web_kit_initialized = on_web_kit_initialized;
    handler->on_browser_created = on_browser_created;
    handler->on_browser_destroyed = on_browser_destroyed;
    handler->get_load_handler = get_load_rp_handler;
    handler->on_before_navigation = on_before_navigation;
    handler->on_context_created = on_context_created;
    handler->on_context_released = on_context_released;
    handler->on_uncaught_exception = on_uncaught_exception;
    handler->on_focused_node_changed = on_focused_node_changed;
    handler->on_process_message_received = on_render_message_received;

    return handler;
}

#endif // HANDLER_CRPH_H_