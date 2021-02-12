package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/night-codes/frame"
)

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	app := frame.MakeApp("My App")
	app.SetIconFromFile(filepath.Join(dir, "/moon.png"))
	fmt.Println(filepath.Join(dir, "/moon.png"))

	wv := app.NewWindow("Simple program!", 500, 400).
		SetBackgroundColor(50, 50, 50, 0.8).
		Move(20, 100).
		SetMinSize(500, 400).
		SetMaxSize(800, 700).
		// SetDecorated(false).
		LoadHTML(`<body style="color:#dddddd; background: transparent">
      <h1>Hello world</h1>
      <p>Test test test...</p>
      </body>`, "http://localhost:1015/panel/").
		SetStateEvent(func(state frame.State) {
			if state.Hidden {
				fmt.Println("Main window closed")
			}
		}).
		SetInvoke(func(msg string) {
			fmt.Println(":::", msg)
		}).
		Show()

	wv2 := app.NewWindow(dir, 400, 300).
		// SetBackgroundColor(80, 50, 50, 0.9).
		LoadHTML(`
		<html style="color:#cccccc; background: transparent">
			<head><script type="text/javascript">setTimeout(function(){document.body.requestFullscreen();}, 1000);</script></head>
		<body style="color:#cccccc; background: transparent">
      <h1>Some Dialog</h1>
      <p>Modal window...</p>
      </body></html>`, "").
		// KeepBelow(true).
		Move(540, 100).
		// SetModal(wv).
		SetStateEvent(func(state frame.State) {
			if state.Hidden {
				fmt.Println("Modal window 1 closed")
			}
		}).
		SetInvoke(func(msg string) {
			fmt.Println(":::", msg)
		}).
		Show()

	go func() {
		// wv.Eval("document.querySelector('html').style.background = 'rgba(0,0,0,0.2)';")
		// wv2.Eval("document.querySelector('html').style.background = '#0000aa99';")
		wv.Eval("window.external.invoke('Wow! This is external invoke!')")
		wv.SetTitle("Новый заголовок")
		wv.Eval("thisIsError2")
		// wv2.Hide()
		wv3 := app.NewWindow("Modal window", 300, 200).
			SetBackgroundColor(40, 80, 50, 0.5).
			SetDeletable(false).
			LoadHTML(`<body style="color:#ffffff; background: transparent">
      <h1>Some Dialog 2</h1>
      <p>Modal window...</p>
	  </body>`, "").
			// SetModal(wv2).
			SetInvoke(func(msg string) {
				fmt.Println(":::", msg)
			})
		t := false

		wv3.
			// SetDecorated(false).
			Move(960, 100).
			SetStateEvent(func(state frame.State) {
				fmt.Printf("%+v\n", state)
				if state.Hidden {
					wv2.LoadHTML(`
					<html>
					<head><script type="text/javascript">window.webkit.messageHandlers.external.postMessage('postMessage invoke');</script></head>
					<body style="color:#999">
					<h1>Super Dialog 3</h1>
					<p>Super modal window...</p>
					</body></html>`, "")

					wv2.Eval("window.external.invoke('message:Some message');")
					fmt.Println("Modal window 2 closed")
					if !t {
						time.Sleep(time.Second / 2)
						wv3.Show()
						t = true
					}
				}
			}).Show()

		go func() {
			time.Sleep(time.Second * 5)
			// wv3.Hide()
			wv3.Load("https://html5test.com/")
			// wv3.SetSize(1000, 700)
			// wv3.Iconify(true)
			wv3.Fullscreen(true)
			go func() {
				time.Sleep(time.Second * 5)
				wv3.Fullscreen(false)
			}()
			go func() {
				time.Sleep(time.Second * 10)
				wv2.Show()
				wv3.SetDeletable(true)
				wv.SetCenter()
				wv2.SetCenter()
				wv3.SetCenter()
				// wv3.Eval("document.body.style.background = '#ff0000'")
				// wv2.Eval("document.body.style.background = '#00ff00'")
				wv.Eval("myfunction();")
			}()
		}()

		editMenu := wv.MainMenu.AddSubMenu("Edit")
		editMenu.AddItem("Find some items", func() {
			fmt.Println("FIND")
		}, "f")

		editMenu.AddSeparatorItem()

		editMenu.AddItem("Ololo", func() {
			fmt.Println("OLOLO")
		}, "o")

		editMenu.AddItem("Test", func() {
			fmt.Println("TEST")
		})
		helpMenu := wv.MainMenu.AddSubMenu("Help")
		regMenu := helpMenu.AddSubMenu("Register application")
		helpMenu.AddItem("About...", func() {
			wv3.Show()
		}, "a")
		regMenu.AddItem("Register by key", func() {
			fmt.Println("Register by key")
		})
		regMenu.AddItem("Buy key in store...", func() {
			fmt.Println("Buy key")
		})

		go func() {
			fmt.Println("~~~~~~=========<<<<<< + >>>>>=========~~~~~~~")
			fmt.Println(wv.GetSize())
			fmt.Println(wv.GetInnerSize())
			fmt.Println(wv.GetPosition())
			fmt.Println(wv.GetScreenSize())
			fmt.Println(wv.GetScreenScaleFactor())
			wv.Eval("window.external.invoke('Window 1: This is external invoke')")
			wv2.Eval("window.external.invoke('Window 2: This is external invoke')")
			wv3.Eval("window.external.invoke('Window 3: This is external invoke')")

		}()
	}()
	// w, h := wv.GetScreen().Size()
	// fmt.Println("Screen size:", w, h)

	app.WaitAllWindowClose()
	// select {}
	fmt.Println("Application terminated")
}
