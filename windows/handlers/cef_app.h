// Copyright (c) 2014 The cefcapi authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cefcapi

#pragma once

#include "handlers/cef_base.h"
#include "handlers/cef_render_process_handler.h"
#include "include/capi/cef_app_capi.h"

// ----------------------------------------------------------------------------
// cef_app_t
// ----------------------------------------------------------------------------

///
// Return the handler for functionality specific to the render process. This
// function is called on the render process main thread.
///
static struct _cef_render_process_handler_t*
    CEF_CALLBACK
    get_render_process_handler(struct _cef_app_t* self)
{
    return initialize_render_process_handler();
}

static void initialize_cef_app(cef_app_t* app)
{
    app->base.size = sizeof(cef_app_t);
    initialize_cef_base((cef_base_t*)app);

    app->get_render_process_handler = get_render_process_handler;
}