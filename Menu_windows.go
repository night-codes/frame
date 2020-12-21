// +build windows

package frame

/*
#ifndef WEBVIEW_WINAPI
#define WEBVIEW_WINAPI
#endif

#include "c_windows.h"
*/
import "C"

type (
	// Menu of window
	Menu struct {
		title    string
		key      string
		menu     interface{} //*C.NSMenu
		menuItem interface{} //*C.NSMenuItem
		parent   *Menu
	}

	// MenuItem element
	MenuItem struct {
		Action   func()
		title    string
		key      string
		menuItem interface{} //*C.NSMenuItem
		parent   *Menu
	}
)

var (
	menuItems = []*MenuItem{}
)

// AddSubMenu returns a new submenu
func (m *Menu) AddSubMenu(title string) *Menu {
	// retM := C.addSubMenu(C.MenuObj{
	// 	title: C.CString(title),
	// 	menu:  m.menu,
	// })

	menu := Menu{
		title: title,
		// menu:     retM.menu,
		// menuItem: retM.menuItem,
		parent: m,
	}
	return &menu
}

// AddItem returns a new menu item
func (m *Menu) AddItem(title string, action func(), key ...string) *MenuItem {
	/*
		k := ""
		if len(key) > 0 {
			k = key[0]
		}
		 retM := C.addItem(C.MenuObj{
			title: C.CString(title),
			key:   C.CString(k),
			menu:  m.menu,
		})
	*/

	item := MenuItem{
		Action: action,
		title:  title,
		// menuItem: retM.menuItem,
		parent: m,
	}
	menuItems = append(menuItems, &item)
	return &item
}

// AddSeparatorItem adds separator item to menu
func (m *Menu) AddSeparatorItem() {
	/* C.addSeparatorItem(C.MenuObj{
		menu: m.menu,
	}) */
}

//export goMenuFunc
// func goMenuFunc(m C.MenuObj) {
// 	go func() {
// 		for _, mm := range menuItems {
// 			if unsafe.Pointer(mm.menuItem) == unsafe.Pointer(m.menuItem) {
// 				go mm.Action()
// 			}
// 		}
// 	}()
// }
