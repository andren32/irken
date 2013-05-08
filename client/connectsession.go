// ConnectSession is a connection with buffers for
// each channel and output ready to be written
// ready to be written to the screen.
// Also stores the nick.
// TODO: Add the methods used in app/app.go
package client

import (
	"fmt"
	"irken/client/msg"
	"irken/client/parser_client"
	"irken/client/parser_server"
	"irken/irc"
	"time"
)

type ConnectSession struct {
	// user specific
	nick     string
	realName string
	addr     string
	// etc
	Conn        *irc.Conn
	IrcChannels map[string]*IRCChannel

	pingFreq    time.Duration
	pingResetCh chan struct{}
	connected   bool
	debug       bool
}

func NewConnectSession(nick string, realName string, debug bool) *ConnectSession {
	cs := &ConnectSession{
		nick:        nick,
		realName:    realName,
		IrcChannels: make(map[string]*IRCChannel),
		pingFreq:    time.Minute,
		connected:   false,
		debug:       debug,
	}
	cs.NewChannel("")
	return cs
}

func (cs *ConnectSession) Connect(addr string) error {
	Conn, err := irc.NewConn(addr)
	if err != nil {
		return err
	}
	// Register the user
	err = Conn.Write("NICK " + cs.nick + "\r\n")
	if err != nil {
		return err
	}
	err = Conn.Write("USER " + cs.nick + " 0 * :" + cs.realName + "\r\n")
	if err != nil {
		return err
	}

	cs.Conn = Conn
	cs.connected = true
	cs.addr = addr
	cs.readToChannels()
	cs.sendPings()
	return nil
}

func (cs *ConnectSession) sendPings() {
	cs.pingResetCh = make(chan struct{})
	go func() {
		for {
			select {
			case _, ok := <-cs.pingResetCh:
				if !ok {
					return
				}
			case <-time.After(cs.pingFreq):
				cs.Send("/silentping "+cs.addr, "")
			}
		}
	}()
}

func (cs *ConnectSession) ResetPing() {
	var dummy struct{}
	cs.pingResetCh <- dummy
}

func (cs *ConnectSession) stopPings() {
	close(cs.pingResetCh)
}

func (cs *ConnectSession) Send(s, context string) error {
	line, output, err := parser_client.Parse(s, cs.nick, context)
	if err != nil {
		return err
	}

	if cs.IsConnected() && output != "" {
		err = cs.Conn.Write(output)
		if err != nil {
			return err
		}
		cs.debugPrint("[" + time.Now().Format("15:04:05") + "]" +
			" <-- " + output)

	}
	if line.OutputMsg() != "" {
		cs.IrcChannels[context].Ch <- line
	}
	return nil
}

func (cs *ConnectSession) readToChannels() {
	go func() {
		for cs.connected {
			s, err := cs.Conn.Read()
			if err != nil {
				// HANDLE ERROR...
			}
			cs.debugPrint("[" + time.Now().Format("15:04:05") + "]" +
				" --> " + s)
			line, err := parser_server.Parse(s)

			if err != nil {
				// HANDLE ERROR...
			}

			value, ok := cs.IrcChannels[line.Context()]
			if !ok && line.OutputMsg() != "" {
				cs.IrcChannels[""].Ch <- line
			} else if line.OutputMsg() != "" {
				value.Ch <- line
			}
		}
	}()
}

func (cs *ConnectSession) NewChannel(context string) {
	cs.IrcChannels[context] = &IRCChannel{Ch: make(chan *msg.Line)}
	//TODO errorstuff
}

func (cs *ConnectSession) DeleteChannel(context string) {
	delete(cs.IrcChannels, context)
}

func (cs *ConnectSession) ChannelExist(context string) bool {
	_, ok := cs.IrcChannels[context]
	return ok
}

func (cs *ConnectSession) IsConnected() bool {
	return cs.connected
}

func (cs *ConnectSession) CloseConnection() {
	cs.connected = false
	cs.stopPings()
	cs.Conn.Close()
}

func (cs *ConnectSession) debugPrint(s string) {
	if cs.debug {
		fmt.Println(s)
	}
}
