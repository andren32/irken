package client

import (
	"errors"
	"strings"
)

// lexMsg scans a IRC message and outputs its tokens.
func lexMsg(message string) (prefix, command string, params []string, err error) {

	// grab prefix if present
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

	command = cmdAndParams[0]
	params = cmdAndParams[1:]
	if trailing != "" {
		params = append(params, trailing)
	}

	return
}

// ParseServerMsg parses an IRC message from an IRC server and outputs
// a string ready to be printed out from the client.
func ParseServerMsg(message string) (output, context string, err error) {
	prefix, command, params, err := lexMsg(message)
	if err != nil {
		return
	}
	switch command {
	case "JOIN":
		return join(prefix, params)
	case "QUIT":
		return quit(prefix, params)
	default:
		err = errors.New("Unknown command.")
		return
	}

	return // bug in old go
}

func join(prefix string, params []string) (string, string, error) {
	nick, err := resolveNick(prefix)
	if err != nil {
		return "", "", err
	}
	channel := strings.Join(params, " ")
	s := nick + " has joined " + channel
	return s, channel, nil
}

func quit(prefix string, params []string) (output, context string, err error) {
	nick, err := resolveNick(prefix)
	if err != nil {
		return
	}
	output = nick + " has quit"
	if len(params) != 0 {
		output += " (" + params[0] + ")"
	}
	return
}

func privMsg(prefix string, params []string) (output, context string, err error) {
	//nick, err := resolveNick(prefix)
	return
}

// resolveNick returns the nick or hostname associated with the IRC message
// it returns empty string when a nick cannot be resolved
func resolveNick(prefix string) (nick string, err error) {
	ind := strings.Index(prefix, "!")
	if ind == -1 {
		err = errors.New("Cannot resolve nick")
		return
	}
	nick = prefix[:ind]
	return
}
