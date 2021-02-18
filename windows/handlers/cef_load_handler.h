// CEF C API example
// Project website: https://github.com/cztomczak/cefcapi

#pragma once

#include "handlers/cef_base.h"
#include "include/capi/cef_app_capi.h"
#include "include/capi/cef_load_handler_capi.h"

extern void goStateChange(cef_browser_t* browser, int status);

static void CEF_CALLBACK on_loading_state_change(struct _cef_load_handler_t* self,
    struct _cef_browser_t* browser, int isLoading, int canGoBack,
    int canGoForward)
{
    goStateChange(browser, isLoading);
};

static cef_load_handler_t* initialize_cef_load_handler()
{
    cef_load_handler_t* loadHandler = (cef_load_handler_t*)calloc(1, sizeof(cef_load_handler_t));
    loadHandler->base.size = sizeof(cef_load_handler_t);
    initialize_cef_base((cef_base_t*)loadHandler);

    loadHandler->on_loading_state_change = on_loading_state_change;
    return loadHandler;
}