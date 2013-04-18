// session describes the current properties associated with the
// execution of irken
package session

import "irken/irc"

type Session struct {
	nick string

	conns []*irc.Conn
	bufs  []*Buffer
}

func (s *Session) Nick() string {
	return s.nick
}
