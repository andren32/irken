package app

import (
	"errors"
	"github.com/mattn/go-gtk/gdk"
	"irken/backend/config"
	"irken/backend/conn"
	"irken/backend/msg"
	"irken/gui"
	"log"
	"os/user"
	"sort"
	"strconv"
)

const DEFAULT_TITLE = "Irken"
const DEFAULT_CONT = ""

// generic handler func that takes a Line argument
type Handler func(*msg.Line)

type IrkenApp struct {
	gui *gui.GUI
	cs  *conn.ConnectSession

	conf *config.Config
	// map from a command string to an action
	handlers    map[string]Handler
	activeChans map[string]chan struct{}
	prevCmd     map[string]string
}

func NewIrkenApp(cfgPath string) *IrkenApp {

	conf, confErr := loadCfg(cfgPath)
	w, _ := conf.GetCfgValue("window_width")
	wWidth, _ := strconv.Atoi(w)
	h, _ := conf.GetCfgValue("window_height")
	wHeight, _ := strconv.Atoi(h)
	d, _ := conf.GetCfgValue("debug")
	debug, _ := strconv.ParseBool(d)

	nick, _ := conf.GetCfgValue("nick")
	realname, _ := conf.GetCfgValue("realname")

	g := gui.NewGUI(DEFAULT_TITLE, wWidth, wHeight)
	cs := conn.NewConnectSession(nick, realname, debug)
	ia := &IrkenApp{
		gui:         g,
		cs:          cs,
		conf:        conf,
		handlers:    make(map[string]Handler),
		activeChans: make(map[string]chan struct{}),
		prevCmd:     make(map[string]string),
	}

	g.CreateChannelWindow(DEFAULT_CONT, func() {
		text, err := g.GetEntryText("")
		if err != nil {
			err := g.WriteToChannel("Couldn't get input", "")
			handleFatalErr(err)
		}
		err = ia.cs.Send(text, DEFAULT_CONT)
		if err != nil {
			err := g.WriteToChannel("Couldn't parse input", "")
			handleFatalErr(err)
		}
		g.EmptyEntryText("")
	})
	ia.BeginInput("")

	// Settingsmenu stuff, for maximizing coolness of code this should maybe
	// be it's own function or something like that
	g.SetSettingsFunc(func() {
		nick, _ := ia.conf.GetCfgValue("nick")
		realname, _ := ia.conf.GetCfgValue("realname")
		nickEntry := g.AddSetting("Default nick", nick)
		realNameEntry := g.AddSetting("Real name", realname)
		g.AddSettingButton("Save", func() {
			ia.conf.AddCfgValue("nick", nickEntry.GetText())
			ia.conf.AddCfgValue("realname", realNameEntry.GetText())
			ia.conf.Save()
			if !ia.HasConnection() {
				cs.ChangeNick(nickEntry.GetText())
				cs.ChangeRealName(realNameEntry.GetText())
			}
			g.CloseSettingsWindow()
		})
	})

	ia.WriteToChatWindow("Welcome to Irken!", DEFAULT_CONT)
	if confErr != nil {
		ia.WriteToChatWindow("Cannot parse config file - using default values",
			DEFAULT_CONT)
	}
	ia.WriteToChatWindow("For help type /help", DEFAULT_CONT)
	ia.WriteToChatWindow("Nick is "+nick, DEFAULT_CONT)
	ia.WriteToChatWindow("Real name is "+realname, DEFAULT_CONT)

	return ia
}

func (ia *IrkenApp) BeginInput(context string) {
	ch := make(chan struct{})
	ia.activeChans[context] = ch
	go func() {
		for {
			select {
			case <-ch:
				delete(ia.activeChans, context)
				return
			default:
				line := <-ia.cs.IrcChannels[context].Ch
				gdk.ThreadsEnter()
				handlErr := ia.Handle(line)
				gdk.ThreadsLeave()
				if handlErr != nil {
					gdk.ThreadsEnter()
					err := ia.gui.WriteToChannel(line.Output(), context)
					gdk.ThreadsLeave()
					handleFatalErr(err)
				}
				// to handle user disconnects
				// when a conversation already has begun
				ia.prevCmd[context] = line.Cmd()
			}
		}
	}()
	return
}

func (ia *IrkenApp) EndInput(context string) {
	close(ia.activeChans[context])
}

func (ia *IrkenApp) UpdateNicks(context string) error {
	channel, ok := ia.cs.IrcChannels[context]
	if !ok {
		return errors.New("UpdateNicks: No such channel " + context)
	}
	nicks := channel.Nicks

	ia.gui.EmptyNicks(context)
	var op []string
	var halfop []string
	var voice []string
	var regular []string
	for nick, perm := range nicks {
		switch perm {
		case "@":
			op = append(op, perm+nick)
		case "%":
			halfop = append(halfop, perm+nick)
		case "+":
			voice = append(voice, perm+nick)
		default:
			regular = append(regular, perm+nick)
		}
	}
	var allSorted []string
	sort.Strings(op)
	sort.Strings(halfop)
	sort.Strings(voice)
	sort.Strings(regular)
	allSorted = append(op, halfop...)
	allSorted = append(allSorted, voice...)
	allSorted = append(allSorted, regular...)
	for _, v := range allSorted {
		ia.gui.WriteToNicks(v, context)
	}

	return nil
}

func (ia *IrkenApp) GUI() *gui.GUI {
	return ia.gui
}

func (ia *IrkenApp) AddHandler(cmd string, h Handler) (err error) {
	_, ok := ia.handlers[cmd]
	if ok {
		return errors.New("Command already has a handler")
	}
	ia.handlers[cmd] = h
	return
}

func (ia *IrkenApp) Handle(l *msg.Line) (err error) {
	h, ok := ia.handlers[l.Cmd()]
	if !ok {
		return errors.New("Couldn't find a handler")
	}
	h(l)
	return
}

func loadCfg(filename string) (c *config.Config, err error) {
	c, err = config.NewConfig(filename)
	if !c.HasValue("nick") {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		c.AddCfgValue("nick", u.Username)
	}
	if !c.HasValue("realname") {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		c.AddCfgValue("realname", u.Name)
	}

	if !c.HasValue("window_width") {
		c.AddCfgValue("window_width", "860")
	} else {
		w, _ := c.GetCfgValue("window_width")
		_, err = strconv.Atoi(w)
		if err != nil {
			c.AddCfgValue("window_width", "860")
		}
	}

	if !c.HasValue("window_height") {
		c.AddCfgValue("window_height", "640")
	} else {
		h, _ := c.GetCfgValue("window_height")
		_, err = strconv.Atoi(h)
		if err != nil {
			c.AddCfgValue("window_height", "640")
		}
	}

	if !c.HasValue("debug") {
		c.AddCfgValue("debug", "false")
	} else {
		d, _ := c.GetCfgValue("debug")
		_, err = strconv.ParseBool(d)
		if err != nil {
			c.AddCfgValue("debug", "false")
		}
	}

	return
}

func (ia *IrkenApp) AddChatWindow(context string) error {
	if ia.HasChannel(context) {
		return errors.New("AddChatWindow: Channel " + context + " already in use")
	}

	ia.cs.NewChannel(context)

	ia.gui.CreateChannelWindow(context, func() {
		text, err := ia.gui.GetEntryText(context)
		if err != nil {
			err := ia.gui.WriteToChannel("Couldn't get input", context)
			handleFatalErr(err)
		}
		err = ia.cs.Send(text, context)
		if err != nil {
			err := ia.gui.WriteToChannel("Couldn't parse input", context)
			handleFatalErr(err)
		}
		ia.gui.EmptyEntryText(context)
	})
	ia.BeginInput(context)
	ia.gui.Notebook().NextPage()
	return nil
}

func (ia *IrkenApp) WriteToCurrentWindow(s string) {
	err := ia.gui.WriteToCurrentWindow(s)
	handleFatalErr(err)
}

func (ia *IrkenApp) DeleteChatWindow(context string) error {
	if !ia.cs.ChannelExist(context) {
		return errors.New("DeleteChatWindow: Channel " + context +
			" doesn't exist")
	}

	ia.cs.DeleteChannel(context)
	ia.EndInput(context)
	ia.gui.DeleteChannelWindow(context)
	return nil
}

func (ia *IrkenApp) HasConnection() bool {
	return ia.cs.IsConnected()
}

func (ia *IrkenApp) Connect(addr string) error {
	return ia.cs.Connect(addr)
}

func (ia *IrkenApp) Disconnect() error {
	if !ia.HasConnection() {
		return errors.New("Already disconnected")
	}
	ia.cs.CloseConnection()
	return nil
}

func (ia *IrkenApp) WriteToChatWindow(s, context string) {
	err := ia.gui.WriteToChannel(s, context)
	handleFatalErr(err)
}

func (ia *IrkenApp) CurrentWindowsContexts() []string {
	a := make([]string, 0)
	for context, _ := range ia.cs.IrcChannels {
		if context != "" {
			a = append(a, context)
		}
	}
	return a
}

func (ia *IrkenApp) AddNicks(channel, nicks string) error {
	ch, ok := ia.cs.IrcChannels[channel]
	if !ok {
		return errors.New("AddNicks: No such channel " + channel + " registered")
	}
	ch.AddNicks(nicks)
	return nil
}

func (ia *IrkenApp) ChangeNick(channel, oldNick, newNick string) {
	ch, ok := ia.cs.IrcChannels[channel]
	if !ok {
		handleFatalErr(errors.New("No such channel " + channel + " registered"))
	}
	ch.ChangeNick(oldNick, newNick)

	if oldNick == ia.cs.GetNick() {
		ia.cs.ChangeNick(newNick)
	}
}

func (ia *IrkenApp) NickExist(channel, nick string) (bool, error) {
	ch, ok := ia.cs.IrcChannels[channel]
	if !ok {
		return false, errors.New("NickExist: No such channel" + channel + " registered.")
	}
	return ch.NickExist(nick), nil
}

func (ia *IrkenApp) RemoveNick(channel, nick string) error {
	ch, ok := ia.cs.IrcChannels[channel]
	if !ok {
		return errors.New("RemoveNick: No such channel " + channel + " registered")
	}
	ch.RemoveNick(nick)
	return nil
}

func (ia *IrkenApp) AddNewNick(channel, nick string) error {
	ch, ok := ia.cs.IrcChannels[channel]
	if !ok {
		return errors.New("AddNewNick: No such channel " + channel + " registered")
	}
	ch.AddNick(nick, "")
	return nil
}

func (ia *IrkenApp) HasChannel(context string) bool {
	_, ok := ia.activeChans[context]
	return ok
}

func (ia *IrkenApp) CurrentNick() string {
	return ia.cs.GetNick()
}

func (ia *IrkenApp) ChangeUserNick(nick string) {
	ia.cs.ChangeNick(nick)
}

func (ia *IrkenApp) SendToCurrentServer(s string) error {
	if !ia.cs.IsConnected() {
		return errors.New("Can't send - not connected")
	}
	ia.cs.SendRaw(s)
	return nil
}

func (ia *IrkenApp) RecentChannelCmd(channel string) (cmd string, err error) {
	cmd, ok := ia.prevCmd[channel]
	if !ok {
		err = errors.New("RecentChannelCmd: No such channel " + channel)
	}
	return
}

func (ia *IrkenApp) IsDebugging() bool {
	return ia.cs.IsDebugging()
}

func (ia *IrkenApp) ResetPing() {
	if !ia.cs.IsConnected() {
		return
	}
	ia.cs.ResetPing()
}

func handleFatalErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
