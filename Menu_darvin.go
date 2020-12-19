// +build darwin

package frame

import (
	"unsafe"
)

/*
#import  "c_darwin.h"

#ifndef WEBVIEW_COCOA
#define WEBVIEW_COCOA
#endif
*/
import "C"

type (
	// Menu of window
	Menu struct {
		app      *App
		title    string
		key      string
		menu     *C.NSMenu
		menuItem *C.NSMenuItem
		parent   *Menu
	}

	// MenuItem element
	MenuItem struct {
		app      *App
		title    string
		key      string
		action   func()
		menuItem *C.NSMenuItem
		parent   *Menu
	}
)

// AddSubMenu returns a new submenu
func (m *Menu) AddSubMenu(title string) *Menu {
	retM := C.addSubMenu(C.MenuObj{
		title: C.CString(title),
		menu:  m.menu,
	})

	menu := Menu{
		app:      m.app,
		title:    title,
		menu:     retM.menu,
		menuItem: retM.menuItem,
		parent:   m,
	}
	return &menu
}

// AddItem returns a new menu item
func (m *Menu) AddItem(title string, action func(), key ...string) *MenuItem {
	k := ""
	if len(key) > 0 {
		k = key[0]
	}
	retM := C.addItem(C.MenuObj{
		title: C.CString(title),
		key:   C.CString(k),
		menu:  m.menu,
	})

	item := MenuItem{
		app:      m.app,
		title:    title,
		menuItem: retM.menuItem,
		action:   action,
		parent:   m,
	}
	menuItems = append(menuItems, &item)
	return &item
}

//export goMenuFunc
func goMenuFunc(m C.MenuObj) {
	go func() {
		for _, mm := range menuItems {
			if unsafe.Pointer(mm.menuItem) == unsafe.Pointer(m.menuItem) {
				go mm.action()
			}
		}
	}()
}
