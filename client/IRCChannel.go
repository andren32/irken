package client

// IRCChannel is a channel on the irc
type IRCChannel struct {
	Ch    chan *Line
	Nicks string
}
