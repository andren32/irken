package app

import (
	"irken/client"
	"irken/gui"
	"os/user"
	"strconv"
)

const DEFAULT_TITLE = "Irken"
const DEFAULT_CONT = ""

type IrkenApp struct {
	gui *gui.GUI
	cs  *client.ConnectSession

	conf *client.Config
}

func NewIrkenApp(cfgPath string) *IrkenApp {

	conf, err := loadCfg(cfgPath)
	w, _ := conf.GetCfgValue("window_width")
	wWidth, _ := strconv.Atoi(w)
	h, _ := conf.GetCfgValue("window_height")
	wHeight, _ := strconv.Atoi(h)

	nick, _ := conf.GetCfgValue("nick")
	realname, _ := conf.GetCfgValue("realname")
	g := gui.NewGUI(DEFAULT_TITLE, wWidth, wHeight)
	g.CreateChannelWindow("", func() {})
	cs := client.NewConnectSession(nick, realname)

	g.StartMain()
	g.WriteToChannel("Welcome to Irken", DEFAULT_CONT)
	if err != nil {
		g.WriteToChannel("Cannot parse config file - using default values", DEFAULT_CONT)
	}
	g.WriteToChannel("Nick is "+nick, DEFAULT_CONT)
	g.WriteToChannel("Real name is "+realname, DEFAULT_CONT)
	return &IrkenApp{
		gui:  g,
		cs:   cs,
		conf: conf,
	}
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
