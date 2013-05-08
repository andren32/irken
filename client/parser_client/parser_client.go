package parser_client

import (
	"errors"
	"irken/client/msg"
	"regexp"
	"strings"
	"time"
)

// lexClientMsg scans a line inputted by the user of the client and outputs
// the tokens in a Line struct
func lexClientMsg(message string) (l *msg.Line, err error) {

	var cmd string
	var args string
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
			cmd = "C" + strings.ToUpper(message[1:])
			args = ""
		} else {
			cmd = "C" + strings.ToUpper(message[1:cmdEnd])
			args = message[cmdEnd+1:]
		}
	} else if strings.HasPrefix(message, "\\/") {
		cmd = "CHAN"
		args = message[1:]
	} else {
		cmd = "CHAN"
		args = message
	}

	a := make([]string, 0)
	// no need to split command arguments
	if cmd == "CHAN" {
		a = append(a, args)
	} else {
		for _, s := range strings.Fields(args) {
			a = append(a, s)
		}
	}

	l = msg.NewLine(message)
	l.SetCmd(cmd)
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
	case "CMSG":
		out, pr, cont = privMsg(nick, l.Args())
	case "CME":
		out, pr = me(nick, context, l.Args())
	case "CJOIN":
		out, pr, cont = join(nick, l.Args())
	case "CNICK":
		out, pr = nickChange(nick, l.Args())
	case "CLEAVE":
		out, pr = part(nick, context)
		l.SetCmd("CPART")
	case "CPART":
		out, pr = part(nick, context)
	case "CEXIT":
		out, pr = quit(nick, l.Args())
		l.SetCmd("CQUIT")
	case "CQUIT":
		out, pr = quit(nick, l.Args())
	case "CCONNECT":
		pr = connect(nick, l.Args())
	case "CDISCONNECT":
		out, pr = disconnect(nick, l.Args())
	case "CPING":
		out, pr = ping(nick, l.Args())
	case "CSILENTPING":
		out, pr = silentping(l.Args())
	case "CSILENTPONG":
		out, pr = silentpong(l.Args())
	case "CHELP":
		// only arguments are important
		out, pr = "", ""
	default:
		// TODO: Remove all error handling, should
		// instead just parse an error msg
		// to the current window
		out = ""
		pr = "/" + strings.ToLower(l.Cmd()[1:]) + ": Unknown command"

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
	l.SetTime(time.Now())

	return
}

func chanMsg(nick, context string, params []string) (out, pr string) {
	if context == "" {
		out = ""
		pr = "Error: you can only chat in a channel or conversation window"
		return
	}
	out = "PRIVMSG " + context + " :" + params[0]
	pr = nick + ": " + params[0]
	return
}

func privMsg(nick string, params []string) (out, pr, context string) {
	if len(params) < 1 {
		out = ""
		context = ""
		pr = "/msg: You must supply a target to /msg to"
		return
	}
	if len(params) < 2 {
		out = ""
		context = ""
		pr = "/msg: You must supply a message to send to " + params[0]
		return
	}
	context = params[0]
	msg := concatArgs(params[1:])
	out = "PRIVMSG " + context + " " + ":" + msg
	pr = nick + ": " + msg
	return
}

func join(nick string, params []string) (out, pr, context string) {
	if len(params) == 0 {
		out = ""
		pr = "/join: Join what?"
		context = ""
		return
	}
	context = params[0]
	out = "JOIN " + params[0]
	pr = nick + " joined " + params[0]
	return
}

func part(nick, context string) (out, pr string) {
	r := "^\\w"
	regex := regexp.MustCompile(r)
	if regex.MatchString(context) {
		out = ""
		pr = nick + " has left the conversation with " + context
		return
	}
	if context == "" {
		out = ""
		pr = "You can't part from server!"
		return
	}
	out = "PART " + context
	pr = nick + " has left " + context
	return
}

func quit(nick string, params []string) (out, pr string) {
	out = "QUIT"
	var msg string
	if len(params) > 0 {
		msg = concatArgs(params)
		out += " :" + msg
	}
	pr = nick + " has quit (" + msg + ")"
	return
}

func me(nick, context string, params []string) (out, pr string) {
	if len(params) < 1 {
		out = ""
		pr = "/me: must supply an action"
		return
	}
	msg := concatArgs(params)
	out = "PRIVMSG " + context + " :\001ACTION " + msg + "\001"
	pr = "*" + nick + "* " + msg
	return
}

func nickChange(nick string, params []string) (out, pr string) {
	out = "NICK " + params[0]
	pr = ""
	return
}

func ping(nick string, params []string) (out, pr string) {
	if len(params) == 0 {
		out = ""
		pr = "/ping: You must supply a ping target!"
		return
	}

	target := params[0]
	out = "PING " + target
	pr = nick + " pinged " + target
	return
}

func silentping(params []string) (out, pr string) {
	out = "PING " + params[0]
	pr = ""
	return
}

func silentpong(params []string) (out, pr string) {
	out = "PONG"
	if len(params) > 0 {
		out += " :" + params[len(params)-1]
	}
	pr = ""
	return
}

// -- Client commands --
func connect(nick string, params []string) (pr string) {
	pr = nick + " connected to " + params[len(params)-1]
	return
}

func disconnect(nick string, params []string) (out, pr string) {
	msg := concatArgs(params)
	out = "QUIT :" + msg
	pr = nick + " disconnected"
	return
}

func concatArgs(args []string) string {
	return strings.Join(args, " ")
}
