// session describes the current properties associated with the
// execution of irken
package client

import "irken/irc"

type Session struct {
	nick string

	conns []*irc.Conn
	bufs  []*Buffer
}

func (s *Session) Nick() string {
	return s.nick
}
