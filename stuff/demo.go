package main

import (
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("Irken")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)
	textview := gtk.NewTextView()
	textview.SetSizeRequest(800, 600)
	//buffer := textview.GetBuffer()
	swin.Add(textview)
	vbox.Add(swin)

	hbox := gtk.NewHBox(false, 1)

	// entry
	entry := gtk.NewEntry()
	entry.SetSizeRequest(700, 40)
	hbox.Add(entry)

	button := gtk.NewButtonWithLabel("Send")
	button.Clicked(func (){
		println(entry.GetText())
	})
	hbox.Add(button)

	vbox.Add(hbox)

	window.Add(vbox)
	window.SetSizeRequest(800, 640)
	window.ShowAll()
	gtk.Main()
}