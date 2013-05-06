package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"irken/client"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(10)
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("Irken")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	cs, err := client.NewConnectSession("irc.freenode.net", "testurnstf", "whatever")
	if err != nil {
		fmt.Println(err)
	}

	cs.ReadToChannels()

	notebook := gtk.NewNotebook()

	CreateChannelWindow("", cs, notebook)
	CreateChannelWindow("#freenode", cs, notebook)

	window.Add(notebook)
	window.SetSizeRequest(800, 640)
	window.ShowAll()

	gtk.Main()
}

func CreateChannelWindow(context string, cs *client.ConnectSession, notebook *gtk.Notebook) {
	var page *gtk.Frame

	if context == "" {
		page = gtk.NewFrame("Server")
		notebook.AppendPage(page, gtk.NewLabel("Server"))
	} else {
		cs.NewChannel(context)
		page = gtk.NewFrame(context)
		notebook.AppendPage(page, gtk.NewLabel(context))
	}

	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)
	textview := gtk.NewTextView()
	textview.SetEditable(false)
	textview.SetCursorVisible(false)
	textview.SetWrapMode(gtk.WRAP_WORD)
	textview.SetSizeRequest(800, 500)
	textbuffer := textview.GetBuffer()
	swin.Add(textview)
	vbox.Add(swin)

	hbox := gtk.NewHBox(false, 1)

	// entry
	entry := gtk.NewEntry()
	entry.SetSizeRequest(700, 40)
	hbox.Add(entry)

	button := gtk.NewButtonWithLabel("Send")
	button.Clicked(func() {
		err := cs.Send(entry.GetText(), context)
		if err != nil {
			fmt.Println(err)
		}
		entry.SetText("")
	})
	hbox.Add(button)

	vbox.Add(hbox)

	page.Add(vbox)

	go func() {
		for {
			line := <-cs.IrcChannels[context].Ch
			gdk.ThreadsEnter()
			var endIter gtk.TextIter
			textbuffer.GetEndIter(&endIter)
			textbuffer.Insert(&endIter, line.Output()+"\n")
			textview.ScrollToIter(&endIter, 0.0, false, 0.0, 0.0)
			gdk.ThreadsLeave()
		}
	}()
}
