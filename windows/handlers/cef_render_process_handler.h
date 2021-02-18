// CEF C API example
// Project website: https://github.com/cztomczak/cefcapi

#ifndef HANDLER_CRPH_H_
#define HANDLER_CRPH_H_

#pragma once

#include "handlers/cef_base.h"
#include "handlers/cef_load_handler.h"
#include "include/capi/cef_app_capi.h"

extern void goContextCreate(cef_v8value_t* global);

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

static cef_render_process_handler_t* initialize_render_process_handler()
{
    cef_render_process_handler_t* handler = (cef_render_process_handler_t*)calloc(1, sizeof(cef_render_process_handler_t));
    handler->base.size = sizeof(cef_render_process_handler_t);
    initialize_cef_base((cef_base_t*)handler);

    handler->on_context_created = on_context_created;

    return handler;
}

#endif // HANDLER_CRPH_H_