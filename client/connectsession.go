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
	IrcChannels map[string]*IRCChannel
}

func NewConnectSession(addr string, nick string, realName string) (*ConnectSession, error) {

	Conn, err := irc.NewConn(addr)
	if err != nil {
		return nil, err
	}
	ircChannels := make(map[string]*IRCChannel)

	// Register the user
	Conn.Write("NICK " + nick + "\r\n")
	Conn.Write("USER " + nick + " 0 * :" + realName + "\r\n")

	cs := &ConnectSession{nick, Conn, ircChannels}
	cs.NewChannel("") // Default server channel

	return cs, nil
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

			value, ok := cs.IrcChannels[line.Context()]
			if !ok {
				cs.IrcChannels[""].Ch <- line
			} else {
				value.Ch <- line
			}

		}
	}()
}

func (cs *ConnectSession) NewChannel(context string) {
	cs.IrcChannels[context] = &IRCChannel{Ch: make(chan *Line)}
	//TODO errorstuff
}

func (cs *ConnectSession) DeleteChannel(context string) {
	delete(cs.IrcChannels, context)
}
