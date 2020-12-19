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

	app := frame.MakeApp("My App")

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
| `Window.Eval`                    |      [x]     |        [x]       |   [ ]   |
| `Window.Fullscreen`              |      [x]     |        [x]       |   [ ]   |
| `Window.GetScreenScaleFactor`    |      [ ]     |        [x]       |   [ ]   |
| `Window.GetScreenSize)`          |      [ ]     |        [x]       |   [ ]   |
| `Window.GetSize)`                |      [ ]     |        [x]       |   [ ]   |
| `Window.Hide`                    |      [x]     |        [x]       |   [ ]   |
| `Window.Iconify`                 |      [x]     |        [x]       |   [ ]   |
| `Window.KeepAbove`               |      [ ]     |        [x]       |   [ ]   |
| `Window.KeepBelow`               |      [ ]     |        [x]       |   [ ]   |
| `Window.Load`                    |      [x]     |        [x]       |   [ ]   |
| `Window.LoadHTML`                |      [x]     |        [x]       |   [ ]   |
| `Window.Maximize`                |      [x]     |        [x]       |   [ ]   |
| `Window.Move`                    |      [x]     |        [x]       |   [ ]   |
| `Window.SetBackgroundColor`      |      [x]     |        [x]       |   [ ]   |
| `Window.SetCenter`               |      [x]     |        [x]       |   [ ]   |
| `Window.SetDecorated`            |      [x]     |        [x]       |   [ ]   |
| `Window.SetDeletable`            |      [x]     |        [x]       |   [ ]   |
| `Window.SetIconFromFile`         |      [x]     |        [x]       |   [ ]   |
| `Window.SetInvoke`               |      [x]     |        [x]       |   [ ]   |
| `Window.SetMaxSize`              |      [x]     |        [x]       |   [ ]   |
| `Window.SetMinSize`              |      [x]     |        [x]       |   [ ]   |
| `Window.SetModal`                |      [x]     |        [x]       |   [ ]   |
| `Window.SetOpacity`              |      [x]     |        [x]       |   [ ]   |
| `Window.SetResizeble`            |      [x]     |        [x]       |   [ ]   |
| `Window.SetSize`                 |      [x]     |        [x]       |   [ ]   |
| `Window.SetStateEvent`           |      [x]     |        [x]       |   [ ]   |
| `Window.SetTitle`                |      [x]     |        [x]       |   [ ]   |
| `Window.SetType`                 |      [ ]     |        [x]       |   [ ]   |
| `Window.SetZoom`                 |      [x]     |        [x]       |   [ ]   |
| `Window.Show`                    |      [x]     |        [x]       |   [ ]   |
| `Window.SkipPager`               |      [x]     |        [x]       |   [ ]   |
| `Window.SkipTaskbar`             |      [x]     |        [x]       |   [ ]   |
| `Window.Stick`                   |      [x]     |        [x]       |   [ ]   |
| `Window.Strut`                   |      [x]     |        [x]       |   [ ]   |
| `Window.UnsetModal`              |      [x]     |        [x]       |   [ ]   |
