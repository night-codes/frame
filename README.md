# frame - simple golang GUI toolkit (gtk-webkit)
## Install

`go get github.com/night-codes/frame`


## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/night-codes/frame"
)

func main() {
	app := frame.MakeApp("My App") // please, use this row as first in main func
	app.SetIconFromFile("./moon.png")

	window := app.NewWindow("Simple program!", 450, 300).
		LoadHTML(`<body style="color:#dddddd; background: transparent">
					<h1>Hello world</h1>
					<p>Test test test...</p>
				</body>`, "about:blank").
		SetBackgroundColor(50, 50, 50, 0.9)

	go func() {
		window.Show() // show window asynchronously from another go routine

		// You don't have to worry about high resolution screens,
		// the app will look equally good on all screens.
		fmt.Print("Window size: ")
		fmt.Println(window.GetSize()) // Used DPI-related pixels as in browser
		fmt.Print("Window inner size: ")
		fmt.Println(window.GetInnerSize())
	}()

	go func() { // Yes! You can change everything in another threads!
		time.Sleep(time.Second * 5)

		window.SetCenter().
			KeepAbove(true).
			SetSize(900, 600).
			Load("https://html5test.com/")

		fmt.Print("Screen size: ")
		fmt.Println(window.GetScreenSize()) // Used DPI-related pixels as in browser

		time.Sleep(time.Second * 15)
		window.Hide() // will close window and finish application after 15 second
	}()

	app.WaitAllWindowClose() // lock main to avoid app termination (you can also use your own way)
}

```

## Building example:
To build the app use the following commands:

```bash
# Macos
$ mkdir -p Example.app/Contents/MacOS
$ go build -o Example.app/Contents/MacOS/example
$ open example.app # Or click on the app in Finder

# Linux
sudo apt install build-essential
sudo apt install libwebkit2gtk-4.0-dev

go build -o example
./example # run example

# Windows
# OS requires special linker flags for GUI apps.
# It's also recommended to use TDM-GCC-64 compiler for CGo.
# http://tdm-gcc.tdragon.net/download
go build -ldflags="-H windowsgui -s -w" -o example.exe
```

## Implementation

| Function                         | MacOS (Cocoa)| Linux (WebKitGTK)| Windows |
| -------------------------------- |:------------:|:----------------:|:-------:|
| `App.NewWindow`                  |       ✅     |         ✅       |    ✅   |
| `App.SetIconFromFile`            |       ✅     |         ✅       |         |
| `App.WaitAllWindowClose`         |       ✅     |         ✅       |         |
| `App.WaitWindowClose`            |       ✅     |         ✅       |         |
| `Menu.AddSubMenu`                |       ✅     |         ✅       |         |
| `Menu.AddItem`                   |       ✅     |         ✅       |         |
| `Menu.AddSeparatorItem`          |       ✅     |         ✅       |         |
| `Window.Eval`                    |       ✅     |         ✅       |    ✅   |
| `Window.Fullscreen`              |       ✅     |         ✅       |    ✅   |
| `Window.GetScreenSize`           |       ✅     |         ✅       |    ✅   |
| `Window.GetScreenScaleFactor`    |       ✅     |         ✅       |    ✅   |
| `Window.GetSize`                 |       ✅     |         ✅       |    ✅   |
| `Window.GetPosition`             |       ✅     |         ✅       |    ✅   |
| `Window.Hide`                    |       ✅     |         ✅       |    ✅   |
| `Window.Iconify`                 |       ✅     |         ✅       |    ✅   |
| `Window.KeepAbove`               |       ✅     |         ✅       |    ✅   |
| `Window.KeepBelow`               |       ✅     |         ✅       |         |
| `Window.Load`                    |       ✅     |         ✅       |    ✅   |
| `Window.LoadHTML`                |       ✅     |         ✅       |    ✅   |
| `Window.Maximize`                |       ✅     |         ✅       |    ✅   |
| `Window.Move`                    |       ✅     |         ✅       |    ✅   |
| `Window.SetCenter`               |       ✅     |         ✅       |    ✅   |
| `Window.SetDecorated`            |       ✅     |         ✅       |    ✅   |
| `Window.SetDeletable`            |       ✅     |         ✅       |    ✅   |
| `Window.SetIconFromFile`         |       ✅     |         ✅       |         |
| `Window.SetInvoke`               |       ✅     |         ✅       |    ✅   |
| `Window.SetMaxSize`              |       ✅     |         ✅       |    ✅   |
| `Window.SetMinSize`              |       ✅     |         ✅       |    ✅   |
| `Window.SetModal`                |       ✅     |         ✅       |    ✅   |
| `Window.SetOpacity`              |       ✅     |         ✅       |    ✅   |
| `Window.SetResizeble`            |       ✅     |         ✅       |    ✅   |
| `Window.SetSize`                 |       ✅     |         ✅       |    ✅   |
| `Window.SetStateEvent`           |       ✅     |         ✅       |    ✅   |
| `Window.SetTitle`                |       ✅     |         ✅       |    ✅   |
| `Window.Show`                    |       ✅     |         ✅       |    ✅   |
| `Window.SkipPager`               |       ✅     |         ✅       |    ✅   |
| `Window.SkipTaskbar`             |       ✅     |         ✅       |    ✅   |
| `Window.Stick`                   |       ✅     |         ✅       |         |
| `Window.UnsetModal`              |       ✅     |         ✅       |    ✅   |
| `Window.GetInnerSize`            |       ✅     |         🆗       |    ✅   |
| `Window.SetInnerSize`            |       ✅     |         🆗       |    ✅   |
| `Window.SetBackgroundColor`      |       ✅     |         ✅       |    🆗   |
| `Window.Strut`                   |              |         ✅       |         |
| `Window.SetType`                 |              |         ✅       |         |


# License

Copyright 2019-2021, Oleksiy Chechel (alex.mirrr@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.