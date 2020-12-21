// +build windows

package frame

/*
#ifndef WEBVIEW_WINAPI
#define WEBVIEW_WINAPI
#endif

#include "c_windows.h"
*/
import "C"

// State struct
type State struct {
	Hidden     bool
	Iconified  bool
	Maximized  bool
	Fullscreen bool
	Focused    bool
}
