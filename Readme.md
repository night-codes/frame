# frame - simple golang GUI toolkit (gtk-webkit)
## Install

**Attention! This is an experiment! Do not use on production!** 

You will need to install gtk-webkit.

`go get github.com/night-codes/frame`


## Example

```go
package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/night-codes/frame"
)

func main() {
	// start backend server
	host := "localhost:40333"
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/index.html", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `
			<h1>Hello world</h1>
			<p>Test test test...</p>`)
	})
	go router.Run(host)

	// start frame
	f := frame.New("p", 450, 300, "Simple program!").
		KeepBelow(true).
		SkipTaskbar(true).
		SkipPager(true).
		Resize(2000, 360).
		SetLimitSizes(2000, 360, 2000, 360).
		SetBackgroundColor(25, 25, 25).
		Load("http://" + host + "/index.html").
		Show()
}
```
