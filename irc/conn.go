package irc

import (
	"net"
	"time"
)

const DEFAULT_PORT = "6667"

type Conn struct {
	// connection info
	user string

	// network utils
	io   *bufio.ReadWriter
	sock *net.TCPConn

	// server properties
	pingFreq time.Duration
}

func NewConn(addr string) (*Conn, error) {

	addr += ":" + DEFAULT_PORT
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	//io := NewReadWriter(bufio.NewReader(sock), bufio.NewWriter(sock))
}

func (c *IRCConn) Connect() {
	//TODO
}

func (c *IRCConn) Read() {
	//TODO
}

func (c *IRCConn) Write() {
	//TODO
}
