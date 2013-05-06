package app

import (
	"irken/client"
	"irken/gui"
	// "os/user"
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
	cs := client.NewConnectSession(conf.GetValue("nick"),
		conf.GetValue("realname"))
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
	// if !c.HasValue("nick") {
	// 	c.AddCfgValue("nick", user.Current().Username)
	// }
	// if !c.HasValue("realname") {
	// 	c.AddCfgValue("realname", user.Current().Name)
	// }
	return
}
