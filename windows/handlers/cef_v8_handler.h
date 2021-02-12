// Copyright (c) 2014 The cefcapi authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cefcapi

#include "include/capi/cef_v8_capi.h"

extern cef_v8value_t* valFromString(char* value);
extern cef_v8value_t* valCreateNull();
extern void goInvokeCallback(cef_string_t* value, cef_string_t* value2);

static int CEF_CALLBACK execute(struct _cef_v8handler_t* self,
    const cef_string_t* name, struct _cef_v8value_t* object,
    size_t argumentsCount, struct _cef_v8value_t* const* arguments,
    struct _cef_v8value_t** retval, cef_string_t* exception)
{
    if (argumentsCount > 1 && arguments[0]->is_string(arguments[0]) && arguments[1]->is_string(arguments[1])) {
        cef_string_userfree_t(CEF_CALLBACK * get_string)(struct _cef_value_t * self);
        goInvokeCallback(arguments[0]->get_string_value(arguments[0]), arguments[1]->get_string_value(arguments[1]));
        return 1;
    }
    return 1;
};

static cef_v8handler_t* initialize_cef_v8handler()
{
    cef_v8handler_t* handler = (cef_v8handler_t*)calloc(1, sizeof(cef_v8handler_t));
    handler->base.size = sizeof(cef_v8handler_t);
    initialize_cef_base((cef_base_t*)handler);
    handler->execute = execute;
    return handler;
}

static cef_v8accessor_t* initialize_cef_v8accessor()
{
    cef_v8accessor_t* handler = (cef_v8accessor_t*)calloc(1, sizeof(cef_v8accessor_t));
    handler->base.size = sizeof(cef_v8accessor_t);
    initialize_cef_base((cef_base_t*)handler);
    return handler;
}
