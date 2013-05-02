package client

import "time"

type Line struct {
	nick, ident, host, src string
	cmd, raw               string
	args                   []string
	time                   time.Time

	context string
	output  string
}

// Output returns the formatted string to be outputted by the client
func (l *Line) Output() string {
	return "[" + l.time.Format("15:04:05") + "] " + l.output
}

// Context returns the context string of the line. Empty context equals
// server and/or client context
func (l *Line) Context() string {
	return l.context
}

func (l *Line) Raw() string {
	return l.raw
}
