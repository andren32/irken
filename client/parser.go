package client

import (
	"errors"
	"strings"
	"time"
)

type Line struct {
	Nick, Ident, Host, Src string
	Cmd, Raw               string
	Args                   []string
	Time                   time.Time

	Context   string
	output    string
	outputErr error
}

// Output returns the string to be outputted by the client.
// It returns an error if the string couldn't be parsed
func (l *Line) Output() (string, error) {
	return "", nil
}

// lexMsg scans a IRC message and outputs its tokens in a Line struct
func lexMsg(message string) (l *Line, err error) {

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

	l = &Line{Nick: nick, Ident: ident, Host: host, Src: src,
		Cmd: command, Raw: message,
		Args: params,
		Time: time.Now()}

	return

}

// ParseServerMsg parses an IRC message from an IRC server and outputs
// a string ready to be printed out from the client.
func ParseServerMsg(message string) (l *Line, err error) {
	l, err = lexMsg(message)
	if err != nil {
		return
	}
	var output string
	var context string
	switch l.Cmd {
	case "PRIVMSG":
		output, context = privMsg(l.Nick, l.Args)
	case "PART":
		output, context = part(l.Nick, l.Args)
	case "JOIN":
		output, context = join(l.Nick, l.Args)
	case "QUIT":
		output, context = quit(l.Nick, l.Args)
	default:
		err = errors.New("Unknown command.")
		return
	}

	l.output = output
	l.Context = context
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

// resolvePrefix returns the token of the IRC message prefix
func resolvePrefix(prefix string) (nick, ident, host, src string, err error) {
	if prefix == "" {
		err = errors.New("Invalid prefix")
		return
	}
	src = prefix

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
