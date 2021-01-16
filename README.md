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

## Implementation

| Function                         | MacOS (Cocoa)| Linux (WebKitGTK)| Windows |
| -------------------------------- |:------------:|:----------------:|:-------:|
| `App.NewWindow`                  |       ✅     |         ✅       |    ✅   |
| `App.SetIconFromFile`            |       ✅     |         ✅       |         |
| `App.WaitAllWindowClose`         |       ✅     |         ✅       |         |
| `App.WaitWindowClose`            |       ✅     |         ✅       |         |
| `Window.Eval`                    |       ✅     |         ✅       |         |
| `Window.Fullscreen`              |       ✅     |         ✅       |         |
| `Window.GetScreenScaleFactor`    |       ✅     |         ✅       |         |
| `Window.GetSize`                 |       ✅     |         ✅       |         |
| `Window.GetPosition`             |       ✅     |                  |         |
| `Window.GetWebviewSize`          |       ✅     |                  |         |
| `Window.Hide`                    |       ✅     |         ✅       |         |
| `Window.Iconify`                 |       ✅     |         ✅       |         |
| `Window.KeepAbove`               |       ✅     |         ✅       |         |
| `Window.KeepBelow`               |       ✅     |         ✅       |         |
| `Window.Load`                    |       ✅     |         ✅       |         |
| `Window.LoadHTML`                |       ✅     |         ✅       |         |
| `Window.Maximize`                |       ✅     |         ✅       |         |
| `Window.Move`                    |       ✅     |         ✅       |         |
| `Window.SetBackgroundColor`      |       ✅     |         ✅       |         |
| `Window.SetCenter`               |       ✅     |         ✅       |         |
| `Window.SetDecorated`            |       ✅     |         ✅       |         |
| `Window.SetDeletable`            |       ✅     |         ✅       |         |
| `Window.SetIconFromFile`         |       ✅     |         ✅       |         |
| `Window.SetInvoke`               |       ✅     |         ✅       |         |
| `Window.SetMaxSize`              |       ✅     |         ✅       |         |
| `Window.SetMinSize`              |       ✅     |         ✅       |         |
| `Window.SetModal`                |       ✅     |         ✅       |         |
| `Window.SetOpacity`              |       ✅     |         ✅       |         |
| `Window.SetResizeble`            |       ✅     |         ✅       |         |
| `Window.SetSize`                 |       ✅     |         ✅       |         |
| `Window.SetStateEvent`           |       ✅     |         ✅       |         |
| `Window.SetTitle`                |       ✅     |         ✅       |         |
| `Window.SetType`                 |              |         ✅       |         |
| `Window.SetWebviewSize`          |       ✅     |                  |         |
| `Window.Show`                    |       ✅     |         ✅       |    ✅   |
| `Window.SkipPager`               |       ✅     |         ✅       |         |
| `Window.SkipTaskbar`             |       ✅     |         ✅       |         |
| `Window.Stick`                   |       ✅     |         ✅       |         |
| `Window.Strut`                   |       ✅     |         ✅       |         |
| `Window.UnsetModal`              |       ✅     |         ✅       |         |
| `Menu.AddSubMenu`                |       ✅     |         ✅       |         |
| `Menu.AddItem`                   |       ✅     |         ✅       |         |
| `Menu.AddSeparatorItem`          |       ✅     |         ✅       |         |
