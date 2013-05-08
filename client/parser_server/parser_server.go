package parser_server

import (
	"errors"
	"irken/client/msg"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// lexServerMsg scans a IRC message and outputs its tokens in a Line struct
func lexServerMsg(message string) (l *msg.Line, err error) {

	// make a timestamp as early as possible
	t := time.Now()

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

	l = msg.NewLine(message)
	l.SetNick(nick)
	l.SetIdent(ident)
	l.SetHost(host)
	l.SetSrc(src)

	l.SetCmd(command)
	l.SetArgs(params)
	l.SetTime(t)
	return

}

// Parse parses an IRC message from an IRC server and outputs
// a string ready to be printed out from the client.
func Parse(message string) (l *msg.Line, err error) {
	l, err = lexServerMsg(message)
	if err != nil {
		return
	}
	var output string
	var context string
	switch l.Cmd() {
	case "NOTICE":
		output, context = notice(l.Nick(), l.Args())
	case "NICK":
		output, context = nick(l.Nick(), l.Args())
	case "MODE":
		output, context = mode(l.Nick(), l.Args())
	case "PRIVMSG":
		output, context = privMsg(l.Nick(), l.Args())
		r := "^\\W"
		regex := regexp.MustCompile(r)
		if !regex.MatchString(context) {
			l.SetCmd("P2PMSG")
		}
	case "PART":
		output, context = part(l.Nick(), l.Args())
	case "PING":
		output, context = ping(l.Args())
	case "PONG":
		// TODO: Handle so that pongs from the server doesn't
		// print, but pongs from other users do
		output, context = "", ""
	case "JOIN":
		output, context = join(l.Nick(), l.Args())
	case "QUIT":
		output, context = quit(l.Nick(), l.Args())
	case "328":
		output, context, err = chanUrl(l.Args())
	case "329":
		output, context, err = chanCreated(l.Args())
	case "332":
		output, context, err = topic(l.Args())
	case "333":
		output, context, err = topicSetBy(l.Args())
	case "353":
		output, context = nickList(l.Args())
	case "366":
		output, context = nickListEnd(l.Args())
	default:
		// check for numeric commands
		r := regexp.MustCompile("^\\d+$")
		if r.MatchString(l.Cmd()) {
			output, context = numeric(l.Nick(), l.Args())
		} else {
			err = errors.New("Unknown command.")
			return
		}
	}
	if err != nil {
		return
	}

	l.SetOutput(output)
	l.SetContext(context)
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
	target := params[0]
	r := "^\\W"
	regex := regexp.MustCompile(r)
	if regex.MatchString(target) {
		context = target
	} else {
		context = ""
	}
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

func topic(params []string) (output, context string, err error) {
	topic := params[len(params)-1]
	// ugly way to get a channel context
	context, err = getChanContext(params)
	if err != nil {
		return
	}
	output = "Topic for " + context + " is \"" + topic + "\""
	return
}

func topicSetBy(params []string) (output, context string, err error) {
	context, err = getChanContext(params)
	if err != nil {
		return
	}

	setBy := params[len(params)-2]
	t, err := timeFromUnixString(params[len(params)-1])
	if err != nil {
		return
	}

	output = "Topic set by " + setBy + " on " + t.Format(time.RFC1123)
	return
}

func chanUrl(params []string) (output, context string, err error) {
	context, err = getChanContext(params)
	if err != nil {
		return
	}
	output = "URL for " + context + ": " + params[len(params)-1]
	return
}

func chanCreated(params []string) (output, context string, err error) {
	context, err = getChanContext(params)
	if err != nil {
		return
	}

	t, err := timeFromUnixString(params[len(params)-1])
	if err != nil {
		return
	}

	output = "Channel created on " + t.Format(time.RFC1123)
	return

}

func numeric(nick string, params []string) (output, context string) {
	context = params[0]
	output = params[len(params)-1]
	return
}

func nickList(params []string) (output, context string) {
	context = params[2]
	output = params[len(params)-1]
	return
}

func nickListEnd(params []string) (output, context string) {
	context = params[1]
	output = params[len(params)-1]
	return
}

func ping(params []string) (output, context string) {
	output = "Pinged: "
	if len(params) > 0 {
		output += params[len(params)-1]
	}
	context = ""
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

// getChanContext searches the list of parameters in order to find a
// channel context. It returns an error when it can't find any.
func getChanContext(params []string) (context string, err error) {
	for _, arg := range params {
		if strings.HasPrefix(arg, "#") {
			context = arg
			return
		}
	}
	err = errors.New("Can't find channel context")
	return
}

func timeFromUnixString(uTime string) (t time.Time, err error) {
	tmp, err := strconv.Atoi(uTime)
	if err != nil {
		return
	}
	unixTime := int64(tmp)
	t = time.Unix(unixTime, 0).UTC()
	return
}
