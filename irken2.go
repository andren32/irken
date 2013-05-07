package main

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"irken/gui"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)
}

func main() {
	g := gui.NewGUI("Irken", 860, 640)

	g.CreateChannelWindow("", func() {
		g.CreateChannelWindow("fafa", func() {})
	})

	g.CreateChannelWindow("hhhh", func() {})
	g.StartMain()
}
