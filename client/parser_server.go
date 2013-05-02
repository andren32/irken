package client

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

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

// lexServerMsg scans a IRC message and outputs its tokens in a Line struct
func lexServerMsg(message string) (l *Line, err error) {

	// grab prefix if present
	var prefix string
	prefixEnd := -1
	if strings.HasPrefix(message, ":") {
		prefixEnd = strings.Index(message, " ")
		if prefixEnd == -1 {
			err = errors.New("Message with only a prefix")
			return
		}
		prefix = message[1:prefixEnd]
	}

	// grab trailing param if present
	var trailing string
	trailingStart := strings.Index(message, " :")
	if trailingStart >= 0 {
		trailing = message[trailingStart+2:]
	} else {
		trailingStart = len(message)
	}

	tmp := message[prefixEnd+1 : trailingStart]
	cmdAndParams := strings.Fields(tmp)
	if len(cmdAndParams) < 1 {
		err = errors.New("Cannot lex command")
		return
	}

	command := cmdAndParams[0]
	params := cmdAndParams[1:]
	if trailing != "" {
		params = append(params, trailing)
	}

	nick, ident, host, src, err := resolvePrefix(prefix)
	if err != nil {
		return
	}

	l = &Line{
		nick: nick, ident: ident, host: host, src: src,
		cmd: command, raw: message,
		args: params,
	}

	return

}

// ParseServerMsg parses an IRC message from an IRC server and outputs
// a string ready to be printed out from the client.
func ParseServerMsg(message string) (l *Line, err error) {
	l, err = lexServerMsg(message)
	l.time = time.Now()
	if err != nil {
		return
	}
	var output string
	var context string
	switch l.cmd {
	case "NOTICE":
		output, context = notice(l.nick, l.args)
	case "NICK":
		output, context = nick(l.nick, l.args)
	case "MODE":
		output, context = mode(l.nick, l.args)
	case "PRIVMSG":
		output, context = privMsg(l.nick, l.args)
	case "PART":
		output, context = part(l.nick, l.args)
	case "JOIN":
		output, context = join(l.nick, l.args)
	case "QUIT":
		output, context = quit(l.nick, l.args)
	default:
		// check for numeric commands
		r := regexp.MustCompile("^\\d+$")
		if r.MatchString(l.cmd) {
			l.output, l.context = numeric(l.nick, l.args)
			return
		}
		err = errors.New("Unknown command.")
		return
	}

	l.output = output
	l.context = context
	return
}

func join(nick string, params []string) (output, context string) {
	channel := params[0]
	output = nick + " has joined " + channel
	context = channel
	return
}

func quit(nick string, params []string) (output, context string) {
	output = nick + " has quit"
	if len(params) != 0 {
		output += " (" + params[0] + ")"
	}
	return
}

func notice(nick string, params []string) (output, context string) {
	return privMsg(nick, params)
}

func mode(nick string, params []string) (output, context string) {
	context = params[0]
	output = nick + " changed mode"
	for i := 1; i < len(params); i++ {
		output += " " + params[i]
	}
	output += " for " + context
	return
}

func privMsg(nick string, params []string) (output, context string) {
	output = nick + ": " + params[len(params)-1]
	context = params[0]
	return
}

func part(nick string, params []string) (output, context string) {
	output = nick + " has left " + params[0]
	context = params[0]
	return
}

func nick(nick string, params []string) (output, context string) {
	output = nick + " changed nick to " + params[0]
	return
}

func numeric(nick string, params []string) (output, context string) {
	context = params[0]
	output = params[len(params)-1]
	return
}

// resolvePrefix returns the token of the IRC message prefix
func resolvePrefix(prefix string) (nick, ident, host, src string, err error) {
	src = prefix
	if prefix == "" {
		nick = "<Server>"
		return
	}

	nickEnd := strings.Index(prefix, "!")
	userEnd := strings.Index(prefix, "@")
	if nickEnd != -1 && userEnd != -1 {
		nick = prefix[0:nickEnd]
		ident = prefix[nickEnd+1 : userEnd]
		host = prefix[userEnd+1:]
	} else {
		nick = src
	}

	return
}
