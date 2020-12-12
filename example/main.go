package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/night-codes/frame"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func main() {
	ch := make(chan bool)
	app := frame.MakeApp(10) // max webviews count

	app.SetDefaultIconFromFile(basepath + "/moon.png")
	wv := app.NewFrame("Simple program!", 400, 400).
		SetBackgroundColor(50, 50, 50, 0.8).
		LoadHTML(`<body style="color:#dddddd; background: transparent">
      <h1>Hello world</h1>
      <p>Test test test...</p>
      </body>`, "http://localhost:1015/panel/").
		SetStateEvent(func(state frame.State) {
			if state.Hidden {
				fmt.Println("Main window closed")
				ch <- true
			}
		}).
		SetInvoke(func(msg string) {
			fmt.Println(":::", msg)
		}).
		Show()

	app.NewFrame("Modal window", 450, 250).
		SetBackgroundColor(80, 50, 50, 0.9).
		LoadHTML(`<body style="color:#dddddd; background: transparent">
      <h1>Some Dialog</h1>
      <p>Modal window...</p>
      </body>`, "").
		SetModal(wv).
		SetResizeble(false).
		SetStateEvent(func(state frame.State) {
			if state.Hidden {
				fmt.Println("Modal window 1 closed")
			}
		}).
		SetCenter().
		Show()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("~", wv.Eval("document.body.style.background = '#449977'; thisIsError1"), "~")
		fmt.Println("~", wv.Eval("window.external.invoke('Wow! This is external invoke!')"), "~")
		time.Sleep(1 * time.Second)
		wv.SetTitle("New title")
		fmt.Println("~", wv.Eval("thisIsError2"), "~")
		fmt.Println("~", wv.Eval("document.body.style.background = '#994477'"), "~")
		// wv2.Hide()
		wv3 := app.NewFrame("Modal window", 330, 130).
			SetBackgroundColor(40, 80, 50, 0.9).
			LoadHTML(`<body style="color:#dddddd; background: transparent">
      <h1>Some Dialog</h1>
      <p>Modal window...</p>
      </body>`, "").
			SetModal(wv).
			SetResizeble(false).
			SetStateEvent(func(state frame.State) {
				if state.Hidden {
					fmt.Println("Modal window 2 closed")
				}
			})

		go func() {
			wv3.Show()
		}()
	}()
	// w, h := wv.GetScreen().Size()
	// fmt.Println("Screen size:", w, h)

	<-ch
	fmt.Println("Application terminated")
}
