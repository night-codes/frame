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
	wv := app.NewFrame("Simple program!", 400, 300).
		SetMinSize(300, 200).
		SetMaxSize(500, 400).
		SetBackgroundColor(50, 50, 50, 0.8).
		LoadHTML(`<body style="color:#dddddd; background: transparent">
      <h1>Hello world</h1>
      <p>Test test test...</p>
      </body>`, "").
		SetStateEvent(func(state frame.State) {
			if state.Hidden {
				fmt.Println("Main window closed")
				ch <- true
			}
		}).
		Show()

	wv2 := app.NewFrame("Modal window", 350, 150).
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
		Show()

	go func() {
		time.Sleep(5 * time.Second)

		wv3 := app.NewFrame("Modal window", 330, 130).
			SetBackgroundColor(40, 80, 50, 0.9).
			LoadHTML(`<body style="color:#dddddd; background: transparent">
      <h1>Some Dialog</h1>
      <p>Modal window...</p>
      </body>`, "").
			SetModal(wv2).
			SetResizeble(false).
			SetStateEvent(func(state frame.State) {
				if state.Hidden {
					fmt.Println("Modal window 2 closed")
				}
			})

		go func() {
			wv3.Show()
			wv.Load("https://github.com/mojbro/gocoa/blob/master/window.m")
		}()
	}()
	// w, h := wv.GetScreen().Size()
	// fmt.Println("Screen size:", w, h)

	<-ch
	fmt.Println("Application terminated")
}

// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/night-codes/frame"
// )

// func main() {
// 	fmt.Println("Start")
// 	ch := make(chan bool)
// 	app := frame.MakeApp(10) // max webviews count

// 	app.SetDefaultIconFromFile("moon.png")
// 	fmt.Println("Runned 1")
// 	wv := app.NewFrame("Simple program!", 400, 300).
// 		SetMinSize(400, 300).
// 		SetBackgroundColor(50, 50, 50, 0.8).
// 		LoadHTML(`<body style="color:#dddddd; background: transparent">
// 	<h1>Hello world</h1>
// 	<p>Test test test...</p>
// 	</body>`, "").
// 		SetStateEvent(func(state frame.State) {
// 			if state.Hidden {
// 				fmt.Println("Main window closed")
// 				ch <- true
// 			}
// 		}).
// 		Show()

// 	fmt.Println("Runned 2")
// 	wv2 := app.NewFrame("Modal window", 350, 150).
// 		SetBackgroundColor(80, 50, 50, 0.9).
// 		LoadHTML(`<body style="color:#dddddd; background: transparent">
// 	<h1>Some Dialog</h1>
// 	<p>Modal window...</p>
// 	</body>`, "").
// 		SetModal(wv).
// 		SetResizeble(false).
// 		SetStateEvent(func(state frame.State) {
// 			if state.Hidden {
// 				fmt.Println("Modal window closed")
// 			}
// 		}).
// 		Show()

// 	fmt.Println("Runned 3")
// 	app.NewFrame("Modal window", 330, 130).
// 		SetBackgroundColor(40, 80, 50, 0.9).
// 		LoadHTML(`<body style="color:#dddddd; background: transparent">
// 			<h1>Some Dialog</h1>
// 			<p>Modal window...</p>
// 			</body>`, "").
// 		SetModal(wv2).
// 		SetResizeble(false).
// 		SetStateEvent(func(state frame.State) {
// 			if state.Hidden {
// 				fmt.Println("Modal window closed")
// 			}
// 		}).
// 		Show()

// 	go func() {
// 		time.Sleep(3 * time.Second)
// 		wv2.Resize(600, 600)
// 	}()

// 	fmt.Println("Runned 4")
// 	// w, h := wv.GetScreen().Size()
// 	// fmt.Println("Screen size:", w, h)

// 	<-ch
// 	fmt.Println("Application terminated")
// }
