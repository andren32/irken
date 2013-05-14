package msg

import "time"

type Line struct {
	nick, ident, host, src string
	cmd, raw               string
	args                   []string
	time                   time.Time

	context string
	output  string
}

// NewLine returns a Line struct with a given raw message
func NewLine(raw string) *Line {
	return &Line{raw: raw}
}

func (l *Line) Nick() string {
	return l.nick
}

func (l *Line) SetNick(nick string) {
	l.nick = nick
}

func (l *Line) Ident() string {
	return l.ident
}

func (l *Line) SetIdent(ident string) {
	l.ident = ident
}

func (l *Line) Host() string {
	return l.host
}

func (l *Line) SetHost(host string) {
	l.host = host
}

func (l *Line) Src() string {
	return l.src
}

func (l *Line) SetSrc(src string) {
	l.src = src
}

func (l *Line) Cmd() string {
	return l.cmd
}

func (l *Line) SetCmd(cmd string) {
	l.cmd = cmd
}

func (l *Line) Args() []string {
	return l.args
}

func (l *Line) SetArgs(args []string) {
	l.args = args
}

func (l *Line) Time() time.Time {
	return l.time
}

func (l *Line) SetTime(t time.Time) {
	l.time = t
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

func (l *Line) SetContext(context string) {
	l.context = context
}

func (l *Line) SetOutput(output string) {
	l.output = output
}

func (l *Line) Raw() string {
	return l.raw
}

func (l *Line) OutputMsg() string {
	return l.output
}
