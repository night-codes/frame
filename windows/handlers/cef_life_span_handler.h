// CEF C API example
// Project website: https://github.com/cztomczak/cefcapi

#pragma once

#include "handlers/cef_base.h"
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_life_span_handler_capi.h"

extern void goBrowserCreate(cef_browser_t* browser);
extern int goBrowserDoClose(cef_browser_t* browser);
extern void goBrowserBeforeClose(cef_browser_t* browser);

// ----------------------------------------------------------------------------
// struct cef_life_span_handler_t
// ----------------------------------------------------------------------------

///
// Implement this structure to handle events related to browser life span. The
// functions of this structure will be called on the UI thread unless otherwise
// indicated.
///

// NOTE: There are many more callbacks in cef_life_span_handler,
//       but only on_before_close is implemented here.

///
// Called just before a browser is destroyed. Release all references to the
// browser object and do not attempt to execute any functions on the browser
// object after this callback returns. This callback will be the last
// notification that references |browser|. See do_close() documentation for
// additional usage information.
///
static void CEF_CALLBACK on_before_close(struct _cef_life_span_handler_t* self, struct _cef_browser_t* browser)
{
    goBrowserBeforeClose(browser);
}

static void CEF_CALLBACK on_after_created(struct _cef_life_span_handler_t* self, struct _cef_browser_t* browser)
{
    goBrowserCreate(browser);
}

static int CEF_CALLBACK do_close(struct _cef_life_span_handler_t* self, struct _cef_browser_t* browser)
{
    return goBrowserDoClose(browser);
};

static void* initialize_cef_life_span_handler()
{
    cef_life_span_handler_t* lifeHandler = (cef_life_span_handler_t*)calloc(1, sizeof(cef_life_span_handler_t));
    lifeHandler->base.size = sizeof(cef_life_span_handler_t);
    initialize_cef_base((cef_base_t*)lifeHandler);
    DEBUG_CALLBACK("[+ INITIALIZE_CEF_LIFE_SPAN_HANDLE +]\n");

    lifeHandler->on_after_created = on_after_created;
    lifeHandler->on_before_close = on_before_close;
    lifeHandler->do_close = do_close;
    lifeHandler->base.size = sizeof(cef_life_span_handler_t);
    return lifeHandler;
}