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
	conn *irc.Conn
	bufs map[string]Buffer
}

func NewConnectSession(addr string, nick string, realName string) (*ConnectSession, error) {
	conn, err := irc.NewConn(addr)
	if err != nil {
		return nil, err
	}
	bufs := make(map[string]Buffer)

	// Register the user
	conn.Write("NICK " + nick + "\r\n")
	conn.Write("USER " + nick + "0 *:" + realName + "\r\n")

	return &ConnectSession{nick, conn, bufs}, nil
}
