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
	app := frame.MakeApp(1) // max webviews count
	app.NewFrame("Simple program!", 450, 300).
		KeepAbove(false).
		SkipTaskbar(false).
		SkipPager(false).
		SetSize(500, 360).
		SetLimitSizes(400, 360, 600, 360).
		LoadHTML(`<body style="color:#dddddd; background: transparent">
					<h1>Hello world</h1>
					<p>Test test test...</p>
				</body>`, "http://localhost").
		SetBackgroundColor(50, 50, 50, 0.8).
		Show()

	select {} //
}

```
