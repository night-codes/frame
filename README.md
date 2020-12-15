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
	frame.NewWindow("Simple program!", 450, 300).
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

	frame.WaitAllWindowClose() // lock main
}
```
