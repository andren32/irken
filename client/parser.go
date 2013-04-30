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
	return "", "", nil
}

func join(s string) (prefix, params) {
	nick := resolveNick(prefix)

}
