package client

import (
	"errors"
	"regexp"
	"strings"
)

// lexClientMsg scans a line inputted by the user of the client and outputs
// the tokens in a Line struct
func lexClientMsg(message string) (l *Line, err error) {

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

	a := make([]string, 1)
	a[0] = arg
	l = &Line{
		cmd: cmd, raw: message,
		args: a,
	}

	return
}

// parseClientMsg parses a client message inputted by the user and outputs
// a Line struct (to be printed within the client) and an out string to be sent
// to the server
// it outputs <Line> and empty string if the command is local to the client,
// i.e. "/help join"
func parseClientMsg(message, nick, context string) (l *Line,
	out string, err error) {
	l, err = lexClientMsg(message)
	if err != nil {
		return
	}
	// quite ugly way to see if the context has changed.
	// since empty context is allowed, the default value is valid
	pr, cont := "", "$"
	switch l.cmd {
	case "CHAN":
		out, pr = clChan(nick, context, l.args)
	case "ME":
		out, pr = clMe(nick, context, l.args)
	case "JOIN":
		out, pr, cont = clJoin(nick, l.args)
	default:
		err = errors.New("Unknown command")
	}

	if err != nil {
		return
	}
	if cont != "$" {
		context = cont
	}
	l.context = context
	l.output = pr

	return
}

func clChan(nick, context string, params []string) (out, pr string) {
	out = "PRIVMSG " + context + " :" + params[0]
	pr = nick + ": " + params[0]
	return
}

func clJoin(nick string, params []string) (out, pr, context string) {
	context = params[0]
	out = "JOIN " + params[0]
	pr = nick + " joined " + params[0]
	return
}

func clMe(nick, context string, params []string) (out, pr string) {
	// TODO
	return
}
