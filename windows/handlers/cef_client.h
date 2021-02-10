// Copyright (c) 2014 The cefcapi authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cefcapi

#pragma once

#include "handlers/cef_base.h"
#include "handlers/cef_life_span_handler.h"
#include "handlers/cef_load_handler.h"
#include "handlers/cef_request_handler.h"
#include "include/capi/cef_client_capi.h"

extern void* goGetLifeSpan(cef_client_t* client);

// ----------------------------------------------------------------------------
// struct _cef_client_t
// ----------------------------------------------------------------------------

///
// Implement this structure to provide handler implementations.
///

///
// Return the handler for context menus. If no handler is provided the default
// implementation will be used.
///

static struct _cef_context_menu_handler_t* CEF_CALLBACK get_context_menu_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for dialogs. If no handler is provided the default
// implementation will be used.
///
static struct _cef_dialog_handler_t* CEF_CALLBACK get_dialog_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for browser display state events.
///
static struct _cef_display_handler_t* CEF_CALLBACK get_display_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for download events. If no handler is returned downloads
// will not be allowed.
///
static struct _cef_download_handler_t* CEF_CALLBACK get_download_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for drag events.
///
static struct _cef_drag_handler_t* CEF_CALLBACK get_drag_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for focus events.
///
static struct _cef_focus_handler_t* CEF_CALLBACK get_focus_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for geolocation permissions requests. If no handler is
// provided geolocation access will be denied by default.
///
static struct _cef_geolocation_handler_t* CEF_CALLBACK get_geolocation_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for JavaScript dialogs. If no handler is provided the
// default implementation will be used.
///
static struct _cef_jsdialog_handler_t* CEF_CALLBACK get_jsdialog_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for keyboard events.
///
static struct _cef_keyboard_handler_t* CEF_CALLBACK get_keyboard_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for browser life span events.
///
static cef_life_span_handler_t* CEF_CALLBACK get_life_span_handler(struct _cef_client_t* self)
{
    return (cef_life_span_handler_t*)goGetLifeSpan(self);
}

///
// Return the handler for browser load status events.
///
static struct _cef_load_handler_t* CEF_CALLBACK get_load_handler(struct _cef_client_t* self)
{
    return initialize_cef_load_handler(); //NULL;
}

///
// Return the handler for off-screen rendering events.
///
static struct _cef_render_handler_t* CEF_CALLBACK get_render_handler(struct _cef_client_t* self)
{
    return NULL;
}

///
// Return the handler for browser request events.
///
static struct _cef_request_handler_t* CEF_CALLBACK get_request_handler(struct _cef_client_t* self)
{
    return NULL; // initialize_request_handler();
}

///
// Called when a new message is received from a different process. Return true
// (1) if the message was handled or false (0) otherwise. Do not keep a
// reference to or attempt to access the message outside of this callback.
///
static int CEF_CALLBACK on_process_message_received(struct _cef_client_t* self, struct _cef_browser_t* browser, cef_process_id_t source_process, struct _cef_process_message_t* message)
{
    return 0;
}

static void initialize_cef_client(struct _cef_client_t* client)
{
    client->base.size = sizeof(cef_client_t);
    initialize_cef_base((cef_base_t*)client);
    // callbacks
    client->get_context_menu_handler = get_context_menu_handler;
    client->get_dialog_handler = get_dialog_handler;
    client->get_display_handler = get_display_handler;
    client->get_download_handler = get_download_handler;
    client->get_drag_handler = get_drag_handler;
    client->get_focus_handler = get_focus_handler;
    client->get_geolocation_handler = get_geolocation_handler;
    client->get_jsdialog_handler = get_jsdialog_handler;
    client->get_keyboard_handler = get_keyboard_handler;
    client->get_life_span_handler = get_life_span_handler;
    client->get_load_handler = get_load_handler;
    client->get_render_handler = get_render_handler;
    client->get_request_handler = get_request_handler;
    // client->on_process_message_received = on_process_message_received;
}