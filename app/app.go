package app

import (
	"errors"
	"fmt"
	"irken/client"
	"irken/client/msg"
	"irken/gui"
	"log"
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
	handlers map[string]Handler
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
	g.CreateChannelWindow(DEFAULT_CONT, func() {})
	cs := client.NewConnectSession(nick, realname)

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

	ia := &IrkenApp{
		gui:      g,
		cs:       cs,
		conf:     conf,
		handlers: make(map[string]Handler),
	}
	initHandlers(ia)
	ia.BeginInput()
	return ia
}

func (ia *IrkenApp) BeginInput() {
	go func() {
		// TODO: Parse entry
	}()
}

func initHandlers(ia *IrkenApp) {
	ia.handlers["CCONECT"] = func(l *msg.Line) {
		addr := l.Args()[len(l.Args())-1]
		err := ia.cs.Connect(addr)
		if err != nil {
			errMsg := fmt.Sprintf("Couldn't connect to %s\n, error: %v",
				addr, err)
			err = ia.gui.WriteToChannel(errMsg, "")
			handleFatalErr(err)
		}
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
		return err
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
