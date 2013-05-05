package parser_client

import (
	"errors"
	"irken/client/msg"
	"regexp"
	"strings"
)

// lexClientMsg scans a line inputted by the user of the client and outputs
// the tokens in a Line struct
func lexClientMsg(message string) (l *msg.Line, err error) {

	var cmd string
	var arg string
	// a slash followed by any non-word char is an invalid command
	r := "^/\\W+"
	regex := regexp.MustCompile(r)
	if regex.MatchString(message) {
		err = errors.New("Invalid command")
		return
	}

	if strings.HasPrefix(message, "/") {
		cmdEnd := strings.Index(message, " ")
		if cmdEnd == -1 {
			cmd = strings.ToUpper(message[1:])
			arg = ""
		} else {
			cmd = strings.ToUpper(message[1:cmdEnd])
			arg = message[cmdEnd+1:]
		}
	} else if strings.HasPrefix(message, "\\/") {
		cmd = "CHAN"
		arg = message[1:]
	} else {
		cmd = "CHAN"
		arg = message
	}

	l = msg.NewLine(message)
	l.SetCmd(cmd)
	a := make([]string, 1)
	a[0] = arg
	l.SetArgs(a)

	return
}

// Parse parses a client message inputted by the user and outputs
// a Line struct (to be printed within the client) and an out string to be sent
// to the server
// it outputs <Line> and empty string if the command is local to the client,
// i.e. "/help join"
func Parse(message, nick, context string) (l *msg.Line,
	out string, err error) {
	l, err = lexClientMsg(message)
	if err != nil {
		return
	}
	var pr string
	var cont interface{}
	switch l.Cmd() {
	case "CHAN":
		out, pr = chanMsg(nick, context, l.Args())
	case "ME":
		out, pr = me(nick, context, l.Args())
	case "JOIN":
		out, pr, cont = join(nick, l.Args())
	default:
		err = errors.New("Unknown command")
	}

	if err != nil {
		return
	}
	// since empty context is allowed, an empty interface and a
	// type assertion test tells us if the context has changed
	if c, ok := cont.(string); ok {
		context = c
	}
	l.SetContext(context)
	l.SetOutput(pr)

	return
}

func chanMsg(nick, context string, params []string) (out, pr string) {
	out = "PRIVMSG " + context + " :" + params[0]
	pr = nick + ": " + params[0]
	return
}

func join(nick string, params []string) (out, pr, context string) {
	context = params[0]
	out = "JOIN " + params[0]
	pr = nick + " joined " + params[0]
	return
}

func me(nick, context string, params []string) (out, pr string) {
	// TODO
	return
}
