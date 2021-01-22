// Copyright (c) 2014 The cefcapi authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/CzarekTomczak/cefcapi

#include "include/capi/cef_app_capi.h"

extern void goPrint(char* text);
extern void goPrintInt(char* text, int num);
extern void goPrintCef(char* text0, cef_string_t* text);
extern char* cefToString(cef_string_t* source);
extern cef_string_t* cefFromString(char* source);

#pragma once

#include "include/capi/cef_base_capi.h"
#include <stdio.h>

// Set to 1 to check if add_ref() and release()
// are called and to track the total number of calls.
// add_ref will be printed as "+", release as "-".
#define DEBUG_REFERENCE_COUNTING 0

// Print only the first execution of the callback,
// ignore the subsequent.
#define DEBUG_CALLBACK(x)          \
    {                              \
        static int first_call = 1; \
        if (first_call == 1) {     \
            first_call = 0;        \
            goPrint(x);            \
        }                          \
    }

// ----------------------------------------------------------------------------
// cef_base_t
// ----------------------------------------------------------------------------

///
// Structure defining the reference count implementation functions. All
// framework structures must include the cef_base_t structure first.
///

///
// Increment the reference count.
///
static void CEF_CALLBACK add_ref(cef_base_t* self)
{
}

///
// Decrement the reference count.  Delete this object when no references
// remain.
///
static int CEF_CALLBACK release(cef_base_t* self)
{
    return 1;
}

///
// Returns the current number of references.
///
static int CEF_CALLBACK has_one_ref(cef_base_t* self)
{
    return 1;
}

static void initialize_cef_base(cef_base_t* base)
{
    size_t size = base->size;
    if (size <= 0) {
        printf("FATAL: initialize_cef_base failed, size member not set\n");
        return;
    }
    /* base->add_ref = add_ref;
    base->release = release;
    base->has_one_ref = has_one_ref; */
}