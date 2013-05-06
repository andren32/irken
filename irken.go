package main

import (
	"runtime"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/gdk"
	"irken/gui"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	g := gui.NewGUI("Irken", 860, 640)
	g.CreateChannelWindow("", func() {})

	g.StartMain()
}

/*
	go func() {
		for {
			line := <-cs.IrcChannels[context].Ch
			gdk.ThreadsEnter()
			var endIter gtk.TextIter
			textBuffer.GetEndIter(&endIter)
			textBuffer.Insert(&endIter, line.Output() + "\n")
			textview.ScrollToIter(&endIter, 0.0, false, 0.0, 0.0)
			gdk.ThreadsLeave()
		}
	}()
*/