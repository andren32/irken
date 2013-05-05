package client

import "irken/client/msg"

// IRCChannel is a channel on the irc
type IRCChannel struct {
	Ch    chan *msg.Line
	Nicks string
}
