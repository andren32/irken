// irc descibes the basic client functionalities of irken
package irc

import (
	"bufio"
	"net"
	"time"
)

const DEFAULT_PORT = "6667"

type Conn struct {
	// sockection info
	user string

	// network utils
	io   *bufio.ReadWriter
	sock *net.TCPConn

	// server properties
	pingFreq time.Duration
}

// NewConn takes a remote network address and a default nick.
// It then tries to connect to the supplied address with the
// properties of the supplied session
func NewConn(addr string, nick string) (*Conn, error) {

	//TODO: check if addr already contains another port then the default port.

	addr += ":" + DEFAULT_PORT
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	sock, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	io := bufio.NewReadWriter(bufio.NewReader(sock), bufio.NewWriter(sock))

	return &Conn{nick, io, sock, time.Minute}, nil
}

func (c *Conn) Read() {
	//TODO
}

func (c *Conn) Write() {
	//TODO
}

func (c *Conn) Close() {
	//TODO
}
