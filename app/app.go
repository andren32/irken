package app

import (
	"irken/client"
	"irken/gui"
	"os/user"
)

const DEFAULT_TITLE = "Irken"

type IrkenApp struct {
	gui *gui.GUI
	cs  *client.ConnectSession

	conf *client.Config
}

func NewIrkenApp(width, height int, cfgPath string) *IrkenApp {

	g := gui.NewGUI(DEFAULT_TITLE, width, height)
	g.CreateChannelWindow("", func() {})
	conf, err := loadCfg(cfgPath)
	nick, _ := conf.GetCfgValue("nick")
	realname, _ := conf.GetCfgValue("realname")
	cs := client.NewConnectSession(nick, realname)
	if err != nil {
		// TODO: Write "using default value" to server window
	}
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
	return
}
