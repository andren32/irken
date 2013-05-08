package app

import (
	"errors"
	"github.com/mattn/go-gtk/gdk"
	"irken/client"
	"irken/client/msg"
	"irken/gui"
	"log"
	"os/user"
	"strconv"
	"sort"
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
	d, _ := conf.GetCfgValue("debug")
	debug, _ := strconv.ParseBool(d)

	nick, _ := conf.GetCfgValue("nick")
	realname, _ := conf.GetCfgValue("realname")
	g := gui.NewGUI(DEFAULT_TITLE, wWidth, wHeight)
	cs := client.NewConnectSession(nick, realname, debug)
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

func (ia *IrkenApp) updateNicks(nicks map[string]string, context string) {
	ia.gui.EmptyNicks(context)
	var op 		[]string
	var halfop 	[]string
	var voice	[]string
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

func handleFatalErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
