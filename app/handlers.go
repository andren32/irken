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

	ia.handlers["CHAN"] = func(l *msg.Line) {
		if !ia.cs.IsConnected() {
			err := ia.gui.WriteToChannel("Error: Must be connected to chat",
				l.Context())
			handleFatalErr(err)
			return
		}

		err := ia.gui.WriteToChannel(l.Output(), l.Context())
		handleFatalErr(err)

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
			err := ia.gui.WriteToChannel("/join: Not connected to any server",
				"")
			handleFatalErr(err)
			return
		}

		if l.Context() == "" {
			err := ia.gui.WriteToChannel(l.Output(), l.Context())
			handleFatalErr(err)
			return

		}

		if ia.cs.ChannelExist(l.Context()) {
			err := ia.gui.WriteToChannel("/join: You have already joined "+
				l.Context(), "")
			handleFatalErr(err)
			return
		}

		ia.AddChatWindow(l.Context())
	}

	ia.handlers["CPART"] = func(l *msg.Line) {
		if !ia.cs.IsConnected() {
			err := ia.gui.WriteToChannel("/part: Not in any channel", "")
			handleFatalErr(err)
			return
		}

		if l.Context() == "" {
			err := ia.gui.WriteToChannel("/part: You can't \"part\" from the server",
				l.Context())
			handleFatalErr(err)
			return
		}

		ia.DeleteChatWindow(l.Context())
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
		channel.AddNicks(l.OutputMsg())
	}

	ia.handlers["366"] = func(l *msg.Line) { // end of nick list
		channel, ok := ia.cs.IrcChannels[l.Context()]
		if !ok {
			handleFatalErr(errors.New("366 Nicklist: Channel, " + l.Context() + ", doesn't exist. Raw: " + l.Raw()))
			return
		}
		ia.updateNicks(channel.Nicks, l.Context())
	}

	ia.handlers["JOIN"] = func(l *msg.Line) {
		channel, ok := ia.cs.IrcChannels[l.Context()]
		if !ok {
			handleFatalErr(errors.New("366 Nicklist: Channel, " + l.Context() + ", doesn't exist. Raw: " + l.Raw()))
			return
		}
		channel.AddNick(l.Nick(), "")
		ia.updateNicks(channel.Nicks, l.Context())
		ia.gui.WriteToChannel(l.Output(), l.Context())
	}

	ia.handlers["NICK"] = func(l *msg.Line) {
		if len(l.Args()) == 1 {
			prevNick := ia.cs.GetNick()
			ia.cs.ChangeNick(l.Args()[0])

			for context, channel := range ia.cs.IrcChannels {
				if context != "" {
					channel.ChangeNick(prevNick, ia.cs.GetNick())
					ia.updateNicks(channel.Nicks, context)
				}
				ia.gui.WriteToChannel(l.Output(), context)
			}
		} else {
		}
	}
	ia.handlers["PING"] = func(l *msg.Line) {
		// TODO: Handle different for server and IRC users
		if len(l.Args()) > 0 {
			ia.cs.Send("/silentpong "+l.Args()[len(l.Args())-1], "")
		} else {
			ia.cs.Send("/silentpong", "")
		}
		ia.cs.ResetPing()
	}

	ia.handlers["P2PMSG"] = func(l *msg.Line) {
		actualCont := l.Nick()

		if ia.cs.ChannelExist(actualCont) {
			// just write to correct window
			err := ia.gui.WriteToChannel(l.Output(), actualCont)
			handleFatalErr(err)
			return
		}

		ia.AddChatWindow(actualCont)
		err := ia.gui.WriteToChannel("Beginning conversation with "+
			actualCont, actualCont)
		handleFatalErr(err)

		err = ia.gui.WriteToChannel(l.Output(), actualCont)
		handleFatalErr(err)

	}

	ia.handlers["CMSG"] = func(l *msg.Line) {
		if !ia.cs.IsConnected() {
			err := ia.gui.WriteToChannel("/msg: Not connected to any server",
				"")
			handleFatalErr(err)
			return
		}

		if ia.cs.ChannelExist(l.Context()) {
			err := ia.gui.WriteToChannel(l.Output(), l.Context())
			handleFatalErr(err)
			return
		}

		ia.AddChatWindow(l.Context())
		err := ia.gui.WriteToChannel("Beginning conversation with "+
			l.Context(), l.Context())
		handleFatalErr(err)

		err = ia.gui.WriteToChannel(l.Output(), l.Context())
		handleFatalErr(err)
	}

	ia.handlers["401"] = func(l *msg.Line) {
		// should not close window automatically on
		// disconnect, just tell that the user has disconnected
		prevCmd, ok := ia.prevCmd[l.Context()]
		if ok && prevCmd == "CHAN" {
			err := ia.gui.WriteToChannel(l.Context()+" has disconnected, either /part or wait for him/her to come online",
				l.Context())
			handleFatalErr(err)
			return
		}

		ia.DeleteChatWindow(l.Context())
		err := ia.gui.WriteToChannel(l.Output(), "")
		handleFatalErr(err)
	}
}

func (ia *IrkenApp) AddChatWindow(context string) {
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
}

func (ia *IrkenApp) DeleteChatWindow(context string) {
	ia.cs.DeleteChannel(context)
	ia.EndInput(context)
	ia.gui.DeleteCurrentWindow()
}
