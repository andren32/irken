package gui

import (
	"errors"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"regexp"
	"unsafe"
)

type SettingsFunc func()

type GUI struct {
	width    int
	height   int
	window   *gtk.Window
	notebook *gtk.Notebook
	pages    map[string]*Page
	settingsBox *gtk.VBox
	settingspopup *gtk.Window
	sf SettingsFunc 
}

type Page struct {
	textView *gtk.TextView
	nickTV   *gtk.TextView
	entry    *gtk.Entry
}

func NewGUI(title string, width, height int) *GUI {
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(title)
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	vbox := gtk.NewVBox(false, 0)

	notebook := gtk.NewNotebook()

	vbox.Add(notebook)
	window.Add(vbox)
	window.SetSizeRequest(width, height)

	gui := &GUI{window: window, notebook: notebook, pages: make(map[string]*Page),
		width: width, height: height}

	gui.createMenu(vbox)

	return gui
}

func (gui *GUI) StartMain() {
	gui.window.ShowAll()
	gtk.Main()
	gdk.ThreadsLeave()
}

func (gui *GUI) CreateChannelWindow(context string, sendFunc func()) {
	var page *gtk.Frame

	conversationRegex := "^\\w"
	regex := regexp.MustCompile(conversationRegex)
	if context == "" {
		page = gtk.NewFrame("Server")
		gui.notebook.AppendPage(page, gtk.NewLabel("Server"))
	} else {
		page = gtk.NewFrame(context)
		gui.notebook.AppendPage(page, gtk.NewLabel(context))
	}

	vbox := gtk.NewVBox(false, 1)
	hbox1 := gtk.NewHBox(false, 1)

	var nickTV *gtk.TextView
	var textView *gtk.TextView

	if context != "" && !regex.MatchString(context) {
		swin := gtk.NewScrolledWindow(nil, nil)
		swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		swin.SetShadowType(gtk.SHADOW_IN)
		textView = gtk.NewTextView()
		textView.SetEditable(false)
		textView.SetCursorVisible(false)
		textView.SetWrapMode(gtk.WRAP_WORD)
		textView.SetSizeRequest(600, 500)
		swin.Add(textView)
		hbox1.Add(swin)

		swin2 := gtk.NewScrolledWindow(nil, nil)
		swin2.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		swin2.SetShadowType(gtk.SHADOW_IN)
		nickTV = gtk.NewTextView()
		nickTV.SetEditable(false)
		nickTV.SetCursorVisible(false)
		nickTV.SetWrapMode(gtk.WRAP_WORD)
		nickTV.SetSizeRequest(200, 500)
		swin2.Add(nickTV)
		hbox1.Add(swin2)
	} else {
		swin := gtk.NewScrolledWindow(nil, nil)
		swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		swin.SetShadowType(gtk.SHADOW_IN)
		textView = gtk.NewTextView()
		textView.SetEditable(false)
		textView.SetCursorVisible(false)
		textView.SetWrapMode(gtk.WRAP_WORD)
		textView.SetSizeRequest(800, 500)
		swin.Add(textView)
		hbox1.Add(swin)
	}

	vbox.Add(hbox1)
	hbox2 := gtk.NewHBox(false, 1)

	// entry
	entry := gtk.NewEntry()
	entry.SetSizeRequest(700, 40)
	entry.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		kev := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		if kev.Keyval == gdk.KEY_Return {
			sendFunc()
		}
	})
	hbox2.Add(entry)

	button := gtk.NewButtonWithLabel("Send")
	button.Clicked(sendFunc)
	hbox2.Add(button)

	vbox.Add(hbox2)

	page.Add(vbox)

	newPage := &Page{textView: textView, nickTV: nickTV, entry: entry}
	gui.pages[context] = newPage
	gui.window.ShowAll()
}

func (gui *GUI) DeleteCurrentWindow() {
	gui.notebook.RemovePage(nil, gui.notebook.GetCurrentPage())
	gui.window.ShowAll()
}

func (gui *GUI) DeleteChannelWindow(context string) error {
	len := gui.notebook.GetNPages()
	for i := 0; i < len; i++ {
		frame := gui.notebook.GetNthPage(i)
		if gui.notebook.GetTabLabelText(frame) == context {
			gui.notebook.RemovePage(nil, i)
			gui.window.ShowAll()
			return nil
		}
	}
	return errors.New("Context doesn't exist")
}

func (gui *GUI) WriteToChannel(s, context string) error {
	var endIter gtk.TextIter
	page, ok := gui.pages[context]
	if !ok {
		return errors.New("WriteToChannel: No Such Window!")
	}
	textBuffer := page.textView.GetBuffer()
	textBuffer.GetEndIter(&endIter)
	textBuffer.Insert(&endIter, s+"\n")

	gui.AutoScroll(textBuffer, &endIter)
	return nil
}

func (gui *GUI) WriteToCurrentWindow(s string) error {
	var endIter gtk.TextIter
	i := gui.notebook.GetCurrentPage()
	frame := gui.notebook.GetNthPage(i)
	labelText := gui.notebook.GetTabLabelText(frame)

	var context string
	if labelText == "Server" {
		context = ""
	} else {
		context = labelText
	}

	page, ok := gui.pages[context]
	if !ok {
		return errors.New("WriteToChannel: No Such Window!")
	}
	textBuffer := page.textView.GetBuffer()
	textBuffer.GetEndIter(&endIter)
	textBuffer.Insert(&endIter, s+"\n")

	gui.AutoScroll(textBuffer, &endIter)
	return nil
}

func (gui *GUI) WriteToNicks(s, context string) error {
	page, ok := gui.pages[context]
	if !ok {
		return errors.New("WriteToChannel: No Such Window!")
	}
	var endIter gtk.TextIter
	textBuffer := page.nickTV.GetBuffer()
	textBuffer.GetEndIter(&endIter)
	textBuffer.Insert(&endIter, s+"\n")
	return nil
}

func (gui *GUI) EmptyNicks(context string) error {
	page, ok := gui.pages[context]
	if !ok {
		return errors.New("WriteToChannel: No Such Window!")
	}
	textBuffer := page.nickTV.GetBuffer()
	textBuffer.SetText("")
	return nil
}

func (gui *GUI) AutoScroll(textbuffer *gtk.TextBuffer, endIter *gtk.TextIter) {
	// TODO
}

func (gui *GUI) GetEntryText(context string) (string, error) {
	page, ok := gui.pages[context]
	if !ok {
		return "", errors.New("GetEntryBuffer: No such window!")
	}
	return page.entry.GetText(), nil
}

func (gui *GUI) EmptyEntryText(context string) error {
	page, ok := gui.pages[context]
	if !ok {
		return errors.New("EmptyEntryBuffer: No such window!")
	}
	page.entry.SetText("")
	return nil
}

func (gui *GUI) Notebook() *gtk.Notebook {
	return gui.notebook
}

func (gui *GUI) createMenu(vbox *gtk.VBox) {
	menubar := gtk.NewMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	menuitem := gtk.NewMenuItem()
	vbox.PackStart(menuitem, false, false, 0)

	cascademenu := gtk.NewMenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("E_xit")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Tools")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	settings := gtk.NewMenuItemWithMnemonic("_Settings")
	settings.Connect("activate", func() {

		gui.settingspopup = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
		gui.settingspopup.SetPosition(gtk.WIN_POS_CENTER)
		gui.settingspopup.SetTitle("Settings")
		gui.settingspopup.SetKeepAbove(true)
		gui.settingspopup.Connect("destroy", func(ctx *glib.CallbackContext) {
			println("settings window got destroy!", ctx.Data().(string))
			gui.CloseSettingsWindow()
		}, "foo")
		gui.settingsBox = gtk.NewVBox(false, 0)

		gui.sf() // settings function

		gui.settingspopup.Add(gui.settingsBox)
		gui.settingspopup.ShowAll()
	})

	submenu.Append(settings)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("_Info")
	menuitem.Connect("activate", func() {
		dialog := gtk.NewMessageDialog(gui.window, gtk.DIALOG_DESTROY_WITH_PARENT,
			gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s",
			"Irken works like most IRC-clients. All commands start with forward-slash(/).\n\nFor a list of commands type /help")
		dialog.Run()
		dialog.Destroy()
	})

	submenu.Append(menuitem)

	menuitem = gtk.NewMenuItemWithMnemonic("_About")
	menuitem.Connect("activate", func() {
		dialog := gtk.NewAboutDialog()
		dialog.SetName("About")
		dialog.SetProgramName("Irken")
		dialog.SetAuthors([]string{"André Nyström - github.com/andren32", "Axel Riese - github.com/axelri"})
		dialog.Run()
		dialog.Destroy()
	})

	submenu.Append(menuitem)
}

// AddSetting adds a setting to the setting menu in the form of an entry.
// Takes the setting name and the initialtext in the entry and returns
// a pointer to the entry.
func (gui *GUI) AddSetting(settingsName, initialText string) *gtk.Entry {
	vbox := gtk.NewVBox(false, 0)

	hbox1 := gtk.NewHBox(false, 0)
	label := gtk.NewLabel(settingsName)
	hbox1.Add(label)
	vbox.Add(hbox1)

	hbox2 := gtk.NewHBox(false, 0)
	entry := gtk.NewEntry()
	entry.SetEditable(true)
	entry.SetText(initialText)
	hbox2.Add(entry)
	vbox.Add(hbox2)

	gui.settingsBox.Add(vbox)
	return entry
}

// AddSettingButton adds a button to the settings menu with the label 
// of the first argument. Calls on buttonFunc on click
func (gui *GUI) AddSettingButton (text string, buttonFunc func()) {
	vbox := gtk.NewVBox(false, 0)
	button := gtk.NewButtonWithLabel(text)
	button.Clicked(buttonFunc)
	vbox.Add(button)
	gui.settingsBox.Add(vbox)
}

func (gui *GUI) SetSettingsFunc (settingsfunc func()) {
	gui.sf = settingsfunc
}

func (gui *GUI) CloseSettingsWindow() {
	gui.settingspopup.Hide()
}