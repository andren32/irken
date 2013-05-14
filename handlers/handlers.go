package handlers

import (
	"fmt"
	"irken/app"
	"irken/client/msg"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

func Init(ia *app.IrkenApp) {
	ia.AddHandler("CCONNECT", func(l *msg.Line) {
		if ia.HasConnection() {
			ia.WriteToChatWindow("/connect: Disconnect first to connect to a new server", "")
			return
		}

		if len(l.Args()) == 0 {
			ia.WriteToChatWindow("/connect: Connect to what?", "")
			return
		}

		addr := l.Args()[0]
		err := ia.Connect(addr)
		if err != nil {
			errMsg := fmt.Sprintf("/connect: Couldn't connect to %s\n, error: %v",
				addr, err)
			ia.WriteToChatWindow(errMsg, "")
		}
	})

	ia.AddHandler("CHAN", func(l *msg.Line) {
		if !ia.HasConnection() {
			ia.WriteToChatWindow("Error: Must be connected to chat",
				l.Context())
			return
		}

		ia.WriteToChatWindow(l.Output(), l.Context())

	})

	ia.AddHandler("CDISCONNECT", func(l *msg.Line) {
		if !ia.HasConnection() {
			ia.WriteToChatWindow("/disconnect: You are not connected to any server", "")
			return
		}

		for _, context := range ia.CurrentWindowsContexts() {
			err := ia.DeleteChatWindow(context)
			handleFatalErr(err)
		}

		ia.Disconnect()
		ia.WriteToChatWindow(l.OutputMsg(), "")

	})

	ia.AddHandler("CJOIN", func(l *msg.Line) {
		if !ia.HasConnection() {
			ia.WriteToChatWindow("/join: Not connected to any server", "")
			return
		}

		if ia.HasChannel(l.Context()) {
			ia.WriteToChatWindow("/join: You have already joined "+l.Context(), "")
			return
		}

		ia.AddChatWindow(l.Context())
		ia.WriteToChatWindow(l.Output(), l.Context())
	})

	ia.AddHandler("CPART", func(l *msg.Line) {
		if !ia.HasConnection() {
			ia.WriteToChatWindow("/part: Not in any channel", "")
			return
		}

		if l.Context() == "" {
			ia.WriteToChatWindow("/part: You can't \"part\" from the server", "")
			return
		}

		err := ia.DeleteChatWindow(l.Context())
		handleFatalErr(err)
	})

	ia.AddHandler("CQUIT", func(l *msg.Line) {
		// TODO: Clean up, at least check that the server has disconnected
		os.Exit(0)
	})

	ia.AddHandler("CNICK", func(l *msg.Line) {
		if !ia.HasConnection() {
			ia.GUI().WriteToCurrentWindow("/nick: Connect first!")
		}
	})

	ia.AddHandler("353", func(l *msg.Line) { // nick list
		err := ia.AddNicks(l.Context(), l.OutputMsg())
		handleFatalErr(err)
	})

	ia.AddHandler("366", func(l *msg.Line) { // end of nick list
		err := ia.UpdateNicks(l.Context())
		handleFatalErr(err)
	})

	ia.AddHandler("JOIN", func(l *msg.Line) {
		err := ia.AddNewNick(l.Context(), l.Nick())
		if err != nil {
			handleFatalErr(err)
		}
		err = ia.UpdateNicks(l.Context())
		handleFatalErr(err)
		ia.WriteToChatWindow(l.Output(), l.Context())
	})

	ia.AddHandler("NICK", func(l *msg.Line) {
		prevNick := l.Nick()
		newNick := l.Args()[0]
		if prevNick == ia.CurrentNick() {
			ia.ChangeUserNick(newNick)
		}

		for _, context := range ia.CurrentWindowsContexts() {
			ia.ChangeNick(context, prevNick, newNick)
			err := ia.UpdateNicks(context)
			handleFatalErr(err)
			ia.WriteToChatWindow(l.Output(), context)
		}
	})

	ia.AddHandler("QUIT", func(l *msg.Line) {
		if l.Nick() == ia.CurrentNick() {
			ia.WriteToChatWindow(l.Output(), "")
			return
		}
		for _, context := range ia.CurrentWindowsContexts() {
			err := ia.RemoveNick(context, l.Nick())
			handleFatalErr(err)
			ia.WriteToChatWindow(l.Output(), context)
		}
	})

	ia.AddHandler("PART", func(l *msg.Line) {
		l.SetCmd("QUIT")
		err := ia.Handle(l)
		handleFatalErr(err)
	})

	ia.AddHandler("PING", func(l *msg.Line) {
		// TODO: Handle different for server and IRC users
		if len(l.Args()) > 0 {
			err := ia.SendToCurrentServer("PONG :" + l.Args()[len(l.Args())-1])
			if err != nil {
				ia.WriteToChatWindow("Connection lost", "")
			}
		} else {
			err := ia.SendToCurrentServer("PONG")
			if err != nil {
				ia.WriteToChatWindow("Connection lost", "")
			}
		}
		ia.ResetPing()
	})

	ia.AddHandler("P2PMSG", func(l *msg.Line) {
		if ia.HasChannel(l.Context()) {
			// just write to correct window
			ia.WriteToChatWindow(l.Output(), l.Context())
			return
		}

		err := ia.AddChatWindow(l.Context())
		handleFatalErr(err)
		ia.WriteToChatWindow("Beginning conversation with "+
			l.Context(), l.Context())

		ia.WriteToChatWindow(l.Output(), l.Context())
	})

	ia.AddHandler("CMSG", func(l *msg.Line) {
		if !ia.HasConnection() {
			ia.WriteToChatWindow("/msg: Not connected to any server", "")
			return
		}

		if ia.HasChannel(l.Context()) {
			ia.WriteToChatWindow(l.Output(), l.Context())
			return
		}

		err := ia.AddChatWindow(l.Context())
		handleFatalErr(err)
		ia.WriteToChatWindow("Beginning conversation with "+
			l.Context(), l.Context())

		ia.WriteToChatWindow(l.Output(), l.Context())
	})

	ia.AddHandler("401", func(l *msg.Line) {
		// should not close window automatically on
		// disconnect, just tell that the user has disconnected
		prevCmd, _ := ia.RecentChannelCmd(l.Context())
		if prevCmd == "CHAN" {
			ia.WriteToChatWindow(l.Context()+" has disconnected, either /part or wait for him/her to come online",
				l.Context())
			return
		}

		err := ia.DeleteChatWindow(l.Context())
		handleFatalErr(err)
		ia.WriteToChatWindow(l.Output(), "")
	})

	ia.AddHandler("470", func(l *msg.Line) {
		oldCh := l.Args()[1]
		err := ia.DeleteChatWindow(oldCh)
		handleFatalErr(err)
		err = ia.AddChatWindow(l.Context())
		handleFatalErr(err)
		ia.WriteToChatWindow(l.Output(), l.Context())
	})

	ia.AddHandler("CTCP", func(l *msg.Line) {
		switch query := l.Args()[0]; query {
		case "ACTION":
			ia.WriteToChatWindow(l.Output(), l.Context())
		case "PING":
			// ctcp ping commonly uses Unix time in micro second resolution
			uMicro := time.Now().UnixNano() / 1000
			tmp := strconv.Itoa(int(uMicro))
			ia.SendToCurrentServer("NOTICE " + l.Nick() + " :\001PING " + tmp[:10] + " " + tmp[10:] + "\001")
		}
	})

	ia.AddHandler("CRAW", func(l *msg.Line) {
		if !ia.IsDebugging() {
			ia.WriteToChatWindow("/raw: Only available in debug mode",
				l.Context())
			return
		}

		if !ia.HasConnection() {
			ia.WriteToChatWindow("/raw: Can only send when connected",
				l.Context())
			return

		}

		err := ia.SendToCurrentServer(l.OutputMsg())
		if err != nil {
			ia.WriteToChatWindow("/raw: Connection lost",
				l.Context())

		}
		ia.WriteToChatWindow(l.Output(),
			l.Context())
	})


func handleFatalErr(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatalln(err)
	}
}
