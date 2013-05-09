// irc descibes the basic client functionalities of irken
package irc

import (
	"bufio"
	"net"
	"strings"
	"time"
)

const DEFAULT_PORT = "6667"

type Conn struct {
	// network utils
	io   *bufio.ReadWriter
	sock *net.TCPConn

	// server properties
	pingFreq time.Duration
}

// NewConn takes a remote network address and a default nick.
// It then tries to connect to the supplied address with the
// properties of the supplied session
func NewConn(addr string) (*Conn, error) {

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

	return &Conn{io, sock, time.Minute}, nil
}

// Read attempts to read from the connected server and puts the data into io buffer.
// It will wait until it gets one line of a string.
func (c *Conn) Read() (string, error) {
	data, err := c.io.ReadString('\n')
	return strings.Replace(data, "\r\n", "", -1), err
}

// Write attempts to write a string into the io buffer and
// flushes it, sending the data to the connected server
func (c *Conn) Write(s string) error {
	_, err := c.io.WriteString(s + "\r\n")
	if err != nil {
		return err
	}
	err = c.io.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Close attempts to close the connection to the server
func (c *Conn) Close() {
	c.sock.Close()
}
