package app

import (
	"errors"
	"fmt"
	"irken/client/msg"
	"os"
)

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

	// ia.handlers["CDISCONNECT"] = func(l *msg.Line) {
	// 	if !ia.cs.IsConnected() {
	// 		err := ia.gui.WriteToChannel("You are not connected to any server", "")
	// 		handleFatalErr(err)
	// 	}

	// 	for context, _ := range ia.listeners {
	// 		if context != "" {
	// 			ia.cs.DeleteChannel(context)
	// 			ia.gui.DeleteCurrentWindow
	// 		}
	// 		ia.cs.CloseConnection()
	// 	}

	// }

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
			err := ia.gui.WriteToChannel("Error: Not in any channel", "")
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

	ia.handlers["353"] = func(l *msg.Line) { // nick list
		channel, ok := ia.cs.IrcChannels[l.Context()]
		if !ok {
			handleFatalErr(errors.New("353 Nicklist: Channel, " + l.Context() + ", doesn't exist. Raw: " + l.Raw()))
			return
		}
		channel.Nicks += l.OutputMsg() + " "
	}

	ia.handlers["366"] = func(l *msg.Line) { // end of nick list
		channel, ok := ia.cs.IrcChannels[l.Context()]
		if !ok {
			handleFatalErr(errors.New("366 Nicklist: Channel, " + l.Context() + ", doesn't exist. Raw: " + l.Raw()))
			return
		}
		ia.updateNicks(channel.Nicks, l.Context())
	}

	ia.handlers["JOIN"] = func(l *msg.Line) { // end of nick list
		channel, ok := ia.cs.IrcChannels[l.Context()]
		if !ok {
			handleFatalErr(errors.New("366 Nicklist: Channel, " + l.Context() + ", doesn't exist. Raw: " + l.Raw()))
			return
		}
		channel.Nicks += l.Nick() + " "
		ia.updateNicks(channel.Nicks, l.Context())
		ia.gui.WriteToChannel(l.Output(), l.Context())
	}
}
