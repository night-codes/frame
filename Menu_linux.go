// +build freebsd linux netbsd openbsd solaris

package frame

/*
#ifndef WEBVIEW_GTK
#define WEBVIEW_GTK
#endif

#include "c_linux.h"
*/
import "C"

import (
	"unsafe"
)

type (
	// Menu struct
	Menu struct {
		title    string
		key      string
		menu     *C.GtkWidget
		menuItem *C.GtkWidget
		parent   *Menu
	}

	// MenuItem element
	MenuItem struct {
		Action   func()
		title    string
		key      string
		menuItem *C.GtkWidget
		parent   *Menu
	}
)

var (
	menuItems = []*MenuItem{}
)

// AddSubMenu returns a new submenu
func (m *Menu) AddSubMenu(title string) *Menu {
	retM := C.addSubMenu(C.MenuObj{
		title: C.CString(title),
		menu:  m.menu,
	})

	menu := Menu{
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
		Action:   action,
		title:    title,
		menuItem: retM.menuItem,
		parent:   m,
	}
	menuItems = append(menuItems, &item)
	return &item
}

// AddSeparatorItem adds separator item to menu
func (m *Menu) AddSeparatorItem() {
	C.addSeparatorItem(C.MenuObj{
		menu: m.menu,
	})
}

//export goMenuFunc
func goMenuFunc(mi *C.GtkWidget) {
	go func() {
		for _, mm := range menuItems {
			if unsafe.Pointer(mm.menuItem) == unsafe.Pointer(mi) {
				go mm.Action()
			}
		}
	}()
}
