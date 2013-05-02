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
	Conn *irc.Conn
	Bufs map[string]Buffer
}

func NewConnectSession(addr string, nick string, realName string) (*ConnectSession, error) {
	conn, err := irc.NewConn(addr)
	if err != nil {
		return nil, err
	}
	bufs := make(map[string]Buffer)

	// Register the user
	Conn.Write("NICK " + nick + "\r\n")
	Conn.Write("USER " + nick + " 0 * :" + realName + "\r\n")

	return &ConnectSession{nick, conn, bufs}, nil
}

func (cs *ConnectSession) LoadBuffers() {
	go func() {
		for {
			s, err := conn.Read()
			if err != nil {
				// HANDLE ERROR...	
			}
			line, err := client.ParseServerMessage(s)

			if err != nil {
				// HANDLE ERROR...
			}

			Buf[line.Context].Write(line.Output())
		}
	}()
}
