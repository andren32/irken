package main

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"irken/gui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	g := gui.NewGUI("Irken", 860, 640)
	g.CreateChannelWindow("", func() {})

	g.StartMain()
}
