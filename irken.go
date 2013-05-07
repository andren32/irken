package main

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"irken/app"
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
	_ = app.NewIrkenApp("../config.cfg")

	select {}
}
