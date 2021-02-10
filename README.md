# frame - simple golang GUI toolkit (gtk-webkit)
## Install

**Attention! This is an experiment! Do not use on production!**

You will need to install gtk-webkit.

`go get github.com/night-codes/frame`


## Example

```go
package main

import (
	"github.com/night-codes/frame"
)

func main() {
	app := frame.MakeApp("My App") // please, use this row as first in main func
	app.SetIconFromFile(basepath + "/moon.png")
	app.NewWindow("Simple program!", 450, 300).
		KeepAbove(false).
		SkipTaskbar(false).
		SkipPager(false).
		SetSize(500, 360).
		SetMinSize(400, 250).
		SetMaxSize(600, 360).
		LoadHTML(`<body style="color:#dddddd; background: transparent">
					<h1>Hello world</h1>
					<p>Test test test...</p>
				</body>`, "about:blank").
		SetBackgroundColor(50, 50, 50, 0.9).
		Show()

	app.WaitAllWindowClose() // lock main
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
| `Window.SetMaxSize`              |       ✅     |         ✅       |         |
| `Window.SetMinSize`              |       ✅     |         ✅       |         |
| `Window.SetModal`                |       ✅     |         ✅       |    ✅   |
| `Window.SetOpacity`              |       ✅     |         ✅       |    ✅   |
| `Window.SetResizeble`            |       ✅     |         ✅       |    ✅   |
| `Window.SetSize`                 |       ✅     |         ✅       |    ✅   |
| `Window.SetStateEvent`           |       ✅     |         ✅       |         |
| `Window.SetTitle`                |       ✅     |         ✅       |    ✅   |
| `Window.Show`                    |       ✅     |         ✅       |    ✅   |
| `Window.SkipPager`               |       ✅     |         ✅       |    ✅   |
| `Window.SkipTaskbar`             |       ✅     |         ✅       |    ✅   |
| `Window.Stick`                   |       ✅     |         ✅       |         |
| `Window.UnsetModal`              |       ✅     |         ✅       |    ✅   |
| `Window.GetWebviewSize`          |       ✅     |         🆗       |    🆗   |
| `Window.SetWebviewSize`          |       ✅     |         🆗       |    🆗   |
| `Window.SetBackgroundColor`      |       ✅     |         ✅       |    🆗   |
| `Window.Strut`                   |              |         ✅       |         |
| `Window.SetType`                 |              |         ✅       |         |
