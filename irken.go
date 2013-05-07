package main

import (
	"fmt"
	"irken/gui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	g := gui.NewGUI("Irken", 860, 640)
	g.CreateChannelWindow("", func() {})
	g.WriteToChannel("sdfösdfs", "")

	g.StartMain()
}
