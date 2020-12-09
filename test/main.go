package main

import (
	"github.com/alediaferia/gogoa"
)

func main() {
	app := gogoa.SharedApplication()
	window := gogoa.NewWindow(200, 200)
	window.SetTitle("Gogoga!")
	app.Run()
}
