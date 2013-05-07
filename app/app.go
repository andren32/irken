package app

import (
	"errors"
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"irken/client"
	"irken/client/msg"
	"irken/gui"
	"log"
	"os"
	"os/user"
	"strconv"
)

const DEFAULT_TITLE = "Irken"
const DEFAULT_CONT = ""

// generic handler func that takes a Line argument
type Handler func(*msg.Line)

type IrkenApp struct {
	gui *gui.GUI
	cs  *client.ConnectSession

	conf *client.Config
	// map from a command string to an action
	handlers  map[string]Handler
	listeners map[string]chan struct{}
}

func NewIrkenApp(cfgPath string) *IrkenApp {

	conf, confErr := loadCfg(cfgPath)
	w, _ := conf.GetCfgValue("window_width")
	wWidth, _ := strconv.Atoi(w)
	h, _ := conf.GetCfgValue("window_height")
	wHeight, _ := strconv.Atoi(h)

	nick, _ := conf.GetCfgValue("nick")
	realname, _ := conf.GetCfgValue("realname")
	g := gui.NewGUI(DEFAULT_TITLE, wWidth, wHeight)
	cs := client.NewConnectSession(nick, realname)
	ia := &IrkenApp{
		gui:       g,
		cs:        cs,
		conf:      conf,
		handlers:  make(map[string]Handler),
		listeners: make(map[string]chan struct{}),
	}
	initHandlers(ia)

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

	err := g.WriteToChannel("Welcome to Irken!", DEFAULT_CONT)
	handleFatalErr(err)
	if confErr != nil {
		err := g.WriteToChannel("Cannot parse config file - using default values",
			DEFAULT_CONT)
		handleFatalErr(err)
	}
	err = g.WriteToChannel("Nick is "+nick, DEFAULT_CONT)
	handleFatalErr(err)
	err = g.WriteToChannel("Real name is "+realname, DEFAULT_CONT)
	handleFatalErr(err)

	return ia
}

func (ia *IrkenApp) BeginInput(context string) {
	ch := make(chan struct{})
	ia.listeners[context] = ch
	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				line := <-ia.cs.IrcChannels[context].Ch
				gdk.ThreadsEnter()
				handlErr := ia.handle(line)
				gdk.ThreadsLeave()
				if handlErr != nil {
					gdk.ThreadsEnter()
					err := ia.gui.WriteToChannel(line.Output(), context)
					gdk.ThreadsLeave()
					handleFatalErr(err)
				}
			}

		}
	}()
	return
}

func (ia *IrkenApp) EndInput(context string) {
	close(ia.listeners[context])
}

func initHandlers(ia *IrkenApp) {
	ia.handlers["CCONNECT"] = func(l *msg.Line) {
		if ia.cs.IsConnected() {
			err := ia.gui.WriteToChannel("Disconnect first to connect to a new server", "")
			handleFatalErr(err)
		}

		addr := l.Args()[len(l.Args())-1]
		err := ia.cs.Connect(addr)
		if err != nil {
			errMsg := fmt.Sprintf("Couldn't connect to %s\n, error: %v",
				addr, err)
			err = ia.gui.WriteToChannel(errMsg, "")
			handleFatalErr(err)
		}
	}

	ia.handlers["CJOIN"] = func(l *msg.Line) {
		if !ia.cs.IsConnected() {
			err := ia.gui.WriteToChannel("Error: Not connected to any server",
				"")
			handleFatalErr(err)
			return
		}

		chanCont := l.Context()
		ia.cs.NewChannel(chanCont)

		ia.gui.CreateChannelWindow(chanCont, func() {
			text, err := ia.gui.GetEntryText(chanCont)
			if err != nil {
				err := ia.gui.WriteToChannel("Couldn't get input", chanCont)
				handleFatalErr(err)
			}
			err = ia.cs.Send(text, chanCont)
			if err != nil {
				err := ia.gui.WriteToChannel("Couldn't parse input", chanCont)
				handleFatalErr(err)
			}
			ia.gui.EmptyEntryText(chanCont)
		})
		ia.BeginInput(chanCont)
		ia.gui.Notebook().NextPage()
	}

	ia.handlers["CPART"] = func(l *msg.Line) {
		if !ia.cs.IsConnected() {
			err := ia.gui.WriteToChannel("Error: Not in any channel",
				"")
			handleFatalErr(err)
			return
		}

		ia.cs.DeleteChannel(l.Context())
		ia.EndInput(l.Context())
		ia.gui.DeleteCurrentWindow()
	}

	ia.handlers["CQUIT"] = func(l *msg.Line) {
		// TODO: Clean up, at least check that the server has disconnected
		os.Exit(0)
	}
}

func (ia *IrkenApp) GUI() *gui.GUI {
	return ia.gui
}

func (ia *IrkenApp) AddHandler(h Handler, cmd string) (err error) {
	_, ok := ia.handlers[cmd]
	if ok {
		return errors.New("Command already has a handler")
	}
	ia.handlers[cmd] = h
	return
}

func (ia *IrkenApp) handle(l *msg.Line) (err error) {
	h, ok := ia.handlers[l.Cmd()]
	if !ok {
		return errors.New("Couldn't find a handler")
	}
	h(l)
	return
}

func loadCfg(filename string) (c *client.Config, err error) {
	c, err = client.NewConfig(filename)
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
	}
	if !c.HasValue("window_height") {
		c.AddCfgValue("window_height", "640")
	}

	return
}

func handleFatalErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
