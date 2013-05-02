// ConnectSession is a connection with buffers for
// each channel and output ready to be written
// ready to be written to the screen.
// Also stores the nick.
package client

import "irken/irc"

type ConnectSession struct {
	// user specific
	nick string

	// etc
	Conn        *irc.Conn
	IrcChannels map[string]irc.IRCChannel
}

func NewConnectSession(addr string, nick string, realName string) (*ConnectSession, error) {
	Conn, err := irc.NewConn(addr)
	if err != nil {
		return nil, err
	}
	ircChannels := make(map[string]irc.IRCChannel)

	// Register the user
	Conn.Write("NICK " + nick + "\r\n")
	Conn.Write("USER " + nick + " 0 * :" + realName + "\r\n")

	return &ConnectSession{nick, Conn, ircChannels}, nil
}

func (cs *ConnectSession) ReadToChannels() {
	go func() {
		for {
			s, err := cs.Conn.Read()
			if err != nil {
				// HANDLE ERROR...	
			}
			line, err := ParseServerMsg(s)

			if err != nil {
				// HANDLE ERROR...
			}

			cs.IrcChannels[line.Context()].Ch <- line.Output()
		}
	}()
}
