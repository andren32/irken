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
		cmd = "PRIVMSG"
		arg = message[1:]
	} else {
		cmd = "PRIVMSG"
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
	// TODO
	return
}
